package server

import (
	"context"
	"net"
	"os"
	"testing"

	apiv1 "github.com/daichimukai/x/syakyo/proglog/api/v1"
	"github.com/daichimukai/x/syakyo/proglog/internal/auth"
	"github.com/daichimukai/x/syakyo/proglog/internal/config"
	"github.com/daichimukai/x/syakyo/proglog/internal/log"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func TestServer(t *testing.T) {
	for scenario, fn := range map[string]func(
		t *testing.T,
		rootClient apiv1.LogClient,
		nobodyClient apiv1.LogClient,
		config *Config,
	){
		"produce/consume a message to/from the log succeeds": testProduceConsume,
		"produce/consume stream succeeds":                    testProduceConsumeStream,
		"produce/consume past log boundary fails":            testConsumePastBoundary,
		"unauthorized client fails":                          testUnauthorized,
	} {
		t.Run(scenario, func(t *testing.T) {
			rootClient, nobodyClient, config, teardown := setupTest(t, nil)
			defer teardown()
			fn(t, rootClient, nobodyClient, config)
		})

	}
}

func testProduceConsume(t *testing.T, client apiv1.LogClient, _ apiv1.LogClient, config *Config) {
	ctx := context.Background()

	want := &apiv1.Record{
		Value: []byte("hello world"),
	}

	produce, err := client.Produce(ctx, &apiv1.ProduceRequest{Record: want})
	require.NoError(t, err)
	want.Offset = produce.Offset

	consume, err := client.Consume(ctx, &apiv1.ConsumeRequest{Offset: produce.Offset})
	require.NoError(t, err)

	require.Equal(t, want.Value, consume.Record.Value)
	require.Equal(t, want.Offset, consume.Record.Offset)
}

func testProduceConsumeStream(t *testing.T, client apiv1.LogClient, _ apiv1.LogClient, config *Config) {
	ctx := context.Background()

	records := []*apiv1.Record{
		{
			Value: []byte("first message"),
		},
		{
			Value: []byte("second message"),
		},
	}

	{
		stream, err := client.ProduceStream(ctx)
		require.NoError(t, err)

		for offset, record := range records {
			err = stream.Send(&apiv1.ProduceRequest{
				Record: record,
			})
			require.NoError(t, err)

			res, err := stream.Recv()
			require.NoError(t, err)

			require.Equal(t, uint64(offset), res.Offset)
		}
	}

	{
		stream, err := client.ConsumeStream(ctx, &apiv1.ConsumeRequest{Offset: 0})
		require.NoError(t, err)

		for i, record := range records {
			res, err := stream.Recv()
			require.NoError(t, err)

			require.Equal(t, &apiv1.Record{
				Value:  record.Value,
				Offset: uint64(i),
			}, res.Record)
		}
	}
}

func testConsumePastBoundary(t *testing.T, client apiv1.LogClient, _ apiv1.LogClient, config *Config) {
	ctx := context.Background()

	produce, err := client.Produce(ctx, &apiv1.ProduceRequest{
		Record: &apiv1.Record{
			Value: []byte("hello world"),
		},
	})
	require.NoError(t, err)

	consume, err := client.Consume(ctx, &apiv1.ConsumeRequest{
		Offset: produce.Offset + 1,
	})
	require.Nil(t, consume)

	got := status.Code(err)
	want := status.Code(apiv1.ErrOffsetOutOfRange{}.GRPCStatus().Err())
	require.Equal(t, want, got)
}

func testUnauthorized(
	t *testing.T,
	_,
	client apiv1.LogClient,
	config *Config,
) {
	ctx := context.Background()
	produce, err := client.Produce(ctx, &apiv1.ProduceRequest{
		Record: &apiv1.Record{
			Value: []byte("hello world"),
		},
	})
	require.Error(t, err)
	require.Nil(t, produce)

	require.Equal(t, codes.PermissionDenied, status.Code(err))

	consume, err := client.Consume(ctx, &apiv1.ConsumeRequest{
		Offset: 0,
	})
	require.Error(t, err)
	require.Nil(t, consume)

	require.Equal(t, codes.PermissionDenied, status.Code(err))
}

func setupTest(t *testing.T, fn func(*Config)) (
	rootClient apiv1.LogClient,
	nobodyClient apiv1.LogClient,
	cfg *Config,
	teardown func(),
) {
	t.Helper()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	newClient := func(crtPath, keyPath string) (*grpc.ClientConn, apiv1.LogClient, []grpc.DialOption) {
		tlsConfig, err := config.SetupTLSConfig(config.TLSConfig{
			CertFile: crtPath,
			KeyFile:  keyPath,
			CAFile:   config.CAFile,
			Server:   false,
		})
		require.NoError(t, err)
		tlsCreds := credentials.NewTLS(tlsConfig)
		opts := []grpc.DialOption{grpc.WithTransportCredentials(tlsCreds)}
		conn, err := grpc.Dial(l.Addr().String(), opts...)
		require.NoError(t, err)
		client := apiv1.NewLogClient(conn)
		return conn, client, opts
	}

	var rootConn *grpc.ClientConn
	rootConn, rootClient, _ = newClient(config.RootClientCertFile, config.RootClientKeyFile)

	var nobodyConn *grpc.ClientConn
	nobodyConn, nobodyClient, _ = newClient(config.NobodyClientCertFile, config.NobodyClientKeyFile)

	serverTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		ServerAddress: l.Addr().String(),
		Server:        true,
	})
	require.NoError(t, err)
	serverCreds := credentials.NewTLS(serverTLSConfig)

	dir, err := os.MkdirTemp("", "server-test")
	require.NoError(t, err)

	clog, err := log.NewLog(dir, log.Config{})
	require.NoError(t, err)

	authorizer := auth.New(config.ACLModelFile, config.ACLPolicyFile)

	cfg = &Config{
		CommitLog:  clog,
		Authorizer: authorizer,
	}
	if fn != nil {
		fn(cfg)
	}

	server, err := NewGRPCServer(cfg, grpc.Creds(serverCreds))
	require.NoError(t, err)

	go func() {
		server.Serve(l)
	}()

	return rootClient, nobodyClient, cfg, func() {
		rootConn.Close()
		nobodyConn.Close()
		server.Stop()
		l.Close()
		clog.Remove()
	}
}
