package agent_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	apiv1 "github.com/daichimukai/x/syakyo/proglog/api/v1"

	"github.com/daichimukai/x/syakyo/proglog/internal/agent"
	"github.com/daichimukai/x/syakyo/proglog/internal/config"
	"github.com/daichimukai/x/syakyo/proglog/internal/loadbalance"
	"github.com/daichimukai/x/syakyo/proglog/internal/testutils"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func TestAgent(t *testing.T) {
	serverTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		Server:        true,
		ServerAddress: "127.0.0.1",
	})
	require.NoError(t, err)

	peerTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.RootClientCertFile,
		KeyFile:       config.RootClientKeyFile,
		CAFile:        config.CAFile,
		Server:        false,
		ServerAddress: "127.0.0.1",
	})
	require.NoError(t, err)

	var agents []*agent.Agent
	for i := 0; i < 3; i++ {
		port := testutils.GetFreePort()
		bindAddr := fmt.Sprintf("%s:%s", "127.0.0.1", port)
		rpcPort, _ := strconv.Atoi(testutils.GetFreePort())

		dataDir, err := os.MkdirTemp("", "agent-test-log")
		require.NoError(t, err)

		var startJoinAddrs []string
		if i != 0 {
			startJoinAddrs = append(startJoinAddrs, agents[0].Config.BindAddr)
		}

		agent, err := agent.New(agent.Config{
			NodeName:        fmt.Sprintf("%d", i),
			StartJoinAddrs:  startJoinAddrs,
			BindAddr:        bindAddr,
			RPCPort:         rpcPort,
			DataDir:         dataDir,
			ACLModelFile:    config.ACLModelFile,
			ACLPolicyFile:   config.ACLPolicyFile,
			ServerTLSConfig: serverTLSConfig,
			PeerTLSConfig:   peerTLSConfig,
			Bootstrap:       i == 0,
		})
		require.NoError(t, err)

		agents = append(agents, agent)

	}

	defer func() {
		for _, agent := range agents {
			err := agent.Shutdown()
			require.NoError(t, err)
			require.NoError(t, os.RemoveAll(agent.Config.DataDir))
		}
	}()

	// wait for agents construct a cluster
	time.Sleep(3000 * time.Millisecond)

	leaderClient := client(t, agents[0], peerTLSConfig)
	produceResponse, err := leaderClient.Produce(context.Background(), &apiv1.ProduceRequest{
		Record: &apiv1.Record{
			Value: []byte("foo"),
		},
	})
	require.NoError(t, err)

	consumeResponse, err := leaderClient.Consume(
		context.Background(),
		&apiv1.ConsumeRequest{
			Offset: produceResponse.Offset,
		},
	)
	require.NoError(t, err)

	require.Equal(t, consumeResponse.Record.Value, []byte("foo"))

	// wait for replicate
	time.Sleep(500 * time.Millisecond)

	followerClient := client(t, agents[1], peerTLSConfig)
	consumeResponse, err = followerClient.Consume(
		context.Background(),
		&apiv1.ConsumeRequest{
			Offset: produceResponse.Offset,
		},
	)
	require.NoError(t, err)
	require.Equal(t, consumeResponse.Record.Value, []byte("foo"))

	// Check that replication is not "cycling", i.e., the replicated record
	// won't be replicated back to the source node.
	consumeResponse, err = leaderClient.Consume(
		context.Background(),
		&apiv1.ConsumeRequest{
			Offset: produceResponse.Offset + 1,
		},
	)
	require.Nil(t, consumeResponse)
	require.Error(t, err)
	require.Equal(t, status.Code(apiv1.ErrOffsetOutOfRange{}.GRPCStatus().Err()), status.Code(err))
}

func client(t *testing.T, agent *agent.Agent, tlsConfig *tls.Config) apiv1.LogClient {
	tlsCreds := credentials.NewTLS(tlsConfig)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(tlsCreds)}
	rpcAddr, err := agent.Config.RPCAddr()
	require.NoError(t, err)

	conn, err := grpc.Dial(fmt.Sprintf("%s://%s", loadbalance.Name, rpcAddr), opts...)
	require.NoError(t, err)

	client := apiv1.NewLogClient(conn)
	return client
}
