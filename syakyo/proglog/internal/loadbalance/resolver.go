package loadbalance

import (
	"context"
	"fmt"
	"sync"

	apiv1 "github.com/daichimukai/x/syakyo/proglog/api/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

type Resolver struct {
	mu            sync.Mutex
	clientConn    resolver.ClientConn
	resolverConn  *grpc.ClientConn
	serviceConfig *serviceconfig.ParseResult
	logger        *zap.Logger
}

var _ resolver.Builder = (*Resolver)(nil)

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.logger = zap.L().Named("resolver")
	r.clientConn = cc
	r.serviceConfig = r.clientConn.ParseServiceConfig(
		// seriously?
		fmt.Sprintf(`{"loadBalancingConfig": [{"%s": {}}]}`, Name),
	)
	var dialOpts []grpc.DialOption
	if opts.DialCreds != nil {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(opts.DialCreds))
	}
	var err error
	r.resolverConn, err = grpc.Dial(target.URL.Host, dialOpts...)
	if err != nil {
		return nil, err
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

const Name = "proglog"

func (r *Resolver) Scheme() string {
	return Name
}

var _ resolver.Resolver = (*Resolver)(nil)

// ResolveNow implements resolver.Resolver.ResolveNow.
func (r *Resolver) ResolveNow(resolver.ResolveNowOptions) {
	r.mu.Lock()
	defer r.mu.Unlock()

	client := apiv1.NewLogClient(r.resolverConn)
	ctx := context.Background()
	res, err := client.GetServers(ctx, &apiv1.GetServersRequest{})
	if err != nil {
		r.logger.Error("failed to resolve server", zap.Error(err))
		return
	}
	var addrs []resolver.Address
	for _, server := range res.Servers {
		addrs = append(addrs, resolver.Address{
			Addr: server.RpcAddr,
			Attributes: attributes.New(
				"is_leader",
				server.IsLeader,
			),
		})
	}

	r.clientConn.UpdateState(resolver.State{
		Addresses:     addrs,
		ServiceConfig: r.serviceConfig,
	})
}

// Close implements resolver.Resolver.Close.
func (r *Resolver) Close() {
	if err := r.resolverConn.Close(); err != nil {
		r.logger.Error(
			"failed to close conn",
			zap.Error(err),
		)
	}
}

func init() {
	resolver.Register(&Resolver{})
}
