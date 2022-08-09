package server

import (
	"context"

	apiv1 "github.com/daichimukai/x/syakyo/proglog/api/v1"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type Config struct {
	CommitLog  CommitLog
	Authorizer Authorizer
}

const (
	objectWildCard = "*"
	produceAction  = "produce"
	consumeAction  = "consume"
)

type Authorizer interface {
	Authorize(subject, object, action string) error
}

var _ apiv1.LogServer = (*grpcServer)(nil)

type grpcServer struct {
	apiv1.UnimplementedLogServer
	*Config
}

func newgrpcServer(config *Config) (srv *grpcServer, err error) {
	srv = &grpcServer{
		Config: config,
	}
	return srv, nil
}

func (s *grpcServer) Produce(ctx context.Context, req *apiv1.ProduceRequest) (*apiv1.ProduceResponse, error) {
	if err := s.Authorizer.Authorize(subject(ctx), objectWildCard, produceAction); err != nil {
		return nil, err
	}

	offset, err := s.CommitLog.Append(req.Record)
	if err != nil {
		return nil, err
	}
	return &apiv1.ProduceResponse{Offset: offset}, nil
}

func (s *grpcServer) Consume(ctx context.Context, req *apiv1.ConsumeRequest) (*apiv1.ConsumeResponse, error) {
	if err := s.Authorizer.Authorize(subject(ctx), objectWildCard, consumeAction); err != nil {
		return nil, err
	}
	record, err := s.CommitLog.Read(req.Offset)
	if err != nil {
		return nil, err
	}
	return &apiv1.ConsumeResponse{Record: record}, nil
}

func (s *grpcServer) ProduceStream(stream apiv1.Log_ProduceStreamServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		res, err := s.Produce(stream.Context(), req)
		if err != nil {
			return err
		}
		if err = stream.Send(res); err != nil {
			return err
		}
	}
}

func (s *grpcServer) ConsumeStream(req *apiv1.ConsumeRequest, stream apiv1.Log_ConsumeStreamServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			res, err := s.Consume(stream.Context(), req)
			switch err.(type) {
			case nil:
			case apiv1.ErrOffsetOutOfRange:
				continue
			default:
				return err
			}
			if err = stream.Send(res); err != nil {
				return err
			}
			req.Offset++
		}
	}
}

type CommitLog interface {
	Append(*apiv1.Record) (uint64, error)
	Read(uint64) (*apiv1.Record, error)
}

func NewGRPCServer(config *Config, grpcOpts ...grpc.ServerOption) (*grpc.Server, error) {
	grpcOpts = append(
		grpcOpts,
		grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(authenticate)),
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(authenticate)),
	)
	gsrv := grpc.NewServer(grpcOpts...)
	srv, err := newgrpcServer(config)
	if err != nil {
		return nil, err
	}
	apiv1.RegisterLogServer(gsrv, srv)
	return gsrv, nil
}

func authenticate(ctx context.Context) (context.Context, error) {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return ctx, status.New(
			codes.Unknown,
			"couldn't find peer info",
		).Err()
	}
	if peer.AuthInfo == nil {
		return context.WithValue(ctx, subjectContextKey{}, ""), nil
	}

	tlsInfo := peer.AuthInfo.(credentials.TLSInfo)
	subject := tlsInfo.State.VerifiedChains[0][0].Subject.CommonName
	ctx = context.WithValue(ctx, subjectContextKey{}, subject)
	return ctx, nil
}

func subject(ctx context.Context) string {
	return ctx.Value(subjectContextKey{}).(string)
}

type subjectContextKey struct{}
