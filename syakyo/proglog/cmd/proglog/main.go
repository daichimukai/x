package main

import (
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/daichimukai/x/syakyo/proglog/internal/agent"
	"github.com/daichimukai/x/syakyo/proglog/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	cli := &cli{}

	cmd := &cobra.Command{
		Use:     "proglog",
		PreRunE: cli.setupConfig,
		RunE:    cli.run,
	}

	if err := setupFlags(cmd); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

type cfg struct {
	agent.Config
	ServerTLSConfig config.TLSConfig
	PeerTLSConfig   config.TLSConfig
}

type cli struct {
	cfg cfg
}

func (c *cli) run(cmd *cobra.Command, args []string) error {
	agent, err := agent.New(c.cfg.Config)
	if err != nil {
		return err
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
	return agent.Shutdown()
}

func setupFlags(cmd *cobra.Command) error {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Flags().String("config-file", "", "path to config file")

	dataDir := path.Join(os.TempDir(), "proglog")
	cmd.Flags().String("data-dir", dataDir, "directory to store log and Raft data")
	cmd.Flags().String("node-name", hostname, "unique server ID")

	cmd.Flags().String("bind-addr", "127.0.0.1:8401", "address to bind Serf on")
	cmd.Flags().Int("rpc-port", 8400, "port for RPC clients (and Raft) connections")
	cmd.Flags().StringSlice("start-join-address", nil, "Serf address to join")
	cmd.Flags().Bool("bootstrap", false, "bootstrap the cluster")

	cmd.Flags().String("acl-model-file", "", "path to ACL model file")
	cmd.Flags().String("acl-policy-file", "", "path to ACL policy")
	cmd.Flags().String("server-tls-cert-file", "", "path to server tls cert")
	cmd.Flags().String("server-tls-key-file", "", "path to server tls key")
	cmd.Flags().String("server-tls-ca-file", "", "path to server certificate authority")
	cmd.Flags().String("peer-tls-cert-file", "", "path to peer tls cert")
	cmd.Flags().String("peer-tls-key-file", "", "path to peer tls key")
	cmd.Flags().String("peer-tls-ca-file", "", "path to peer certificate authority")

	return viper.BindPFlags(cmd.Flags())
}

func (c *cli) setupConfig(cmd *cobra.Command, args []string) error {
	configFile, err := cmd.Flags().GetString("config-file")
	if err != nil {
		return err
	}
	viper.SetConfigFile(configFile)

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	c.cfg.DataDir = viper.GetString("data-dir")
	c.cfg.NodeName = viper.GetString("node-name")
	c.cfg.BindAddr = viper.GetString("bind-addr")
	c.cfg.RPCPort = viper.GetInt("rpc-port")
	c.cfg.StartJoinAddrs = viper.GetStringSlice("start-join-addrs")
	c.cfg.Bootstrap = viper.GetBool("bootstrap")
	c.cfg.ACLModelFile = viper.GetString("acl-mode-file")
	c.cfg.ACLPolicyFile = viper.GetString("acl-policy-file")
	c.cfg.ServerTLSConfig.CertFile = viper.GetString("server-tls-cert-file")
	c.cfg.ServerTLSConfig.KeyFile = viper.GetString("server-tls-cert-file")
	c.cfg.ServerTLSConfig.CAFile = viper.GetString("server-tls-ca-file")
	c.cfg.PeerTLSConfig.CertFile = viper.GetString("peer-tls-cert-file")
	c.cfg.PeerTLSConfig.KeyFile = viper.GetString("peer-tls-key-file")
	c.cfg.PeerTLSConfig.CAFile = viper.GetString("peer-tls-ca-file")

	if c.cfg.ServerTLSConfig.CertFile != "" &&
		c.cfg.ServerTLSConfig.KeyFile != "" {
		c.cfg.ServerTLSConfig.Server = true
		c.cfg.Config.ServerTLSConfig, err = config.SetupTLSConfig(c.cfg.ServerTLSConfig)
		if err != nil {
			return err
		}
	}

	if c.cfg.PeerTLSConfig.CertFile != "" &&
		c.cfg.PeerTLSConfig.KeyFile != "" {
		c.cfg.Config.PeerTLSConfig, err = config.SetupTLSConfig(c.cfg.PeerTLSConfig)
		if err != nil {
			return err
		}
	}

	return nil
}
