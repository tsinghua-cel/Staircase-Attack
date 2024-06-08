package beaconapi

import (
	"context"
	"encoding/hex"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpcopentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	grpcutil "github.com/prysmaticlabs/prysm/v4/api/grpc"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"

	log "github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"
)

// default grpc port is 4000
func GetValidators(grpcEndpoint string) ([]string, error) {
	client := NewGrpcBeaconChainClient(grpcEndpoint)
	req := &ethpb.ListValidatorsRequest{}
	vals, err := client.ListValidators(context.Background(), req)
	if err != nil {
		return nil, err
	}

	var pubkeys = make([]string, 0, len(vals.ValidatorList))
	for _, val := range vals.ValidatorList {
		pubkeys = append(pubkeys, hex.EncodeToString(val.Validator.PublicKey))
	}
	return pubkeys, nil
}

// ConstructDialOptions constructs a list of grpc dial options
func ConstructDialOptions(
	maxCallRecvMsgSize int,
	withCert string,
	grpcRetries uint,
	grpcRetryDelay time.Duration,
	extraOpts ...grpc.DialOption,
) []grpc.DialOption {
	var transportSecurity grpc.DialOption
	if withCert != "" {
		creds, err := credentials.NewClientTLSFromFile(withCert, "")
		if err != nil {
			log.WithError(err).Error("Could not get valid credentials")
			return nil
		}
		transportSecurity = grpc.WithTransportCredentials(creds)
	} else {
		transportSecurity = grpc.WithInsecure()
		log.Warn("You are using an insecure gRPC connection. If you are running your beacon node and " +
			"validator on the same machines, you can ignore this message. If you want to know " +
			"how to enable secure connections, see: https://docs.prylabs.network/docs/prysm-usage/secure-grpc")
	}

	if maxCallRecvMsgSize == 0 {
		maxCallRecvMsgSize = 10 * 5 << 20 // Default 50Mb
	}

	dialOpts := []grpc.DialOption{
		transportSecurity,
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxCallRecvMsgSize),
			grpcretry.WithMax(grpcRetries),
			grpcretry.WithBackoff(grpcretry.BackoffLinear(grpcRetryDelay)),
		),
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
		grpc.WithUnaryInterceptor(middleware.ChainUnaryClient(
			grpcopentracing.UnaryClientInterceptor(),
			grpcprometheus.UnaryClientInterceptor,
			grpcretry.UnaryClientInterceptor(),
			grpcutil.LogRequests,
		)),
		grpc.WithChainStreamInterceptor(
			grpcutil.LogStream,
			grpcopentracing.StreamClientInterceptor(),
			grpcprometheus.StreamClientInterceptor,
			grpcretry.StreamClientInterceptor(),
		),
		grpc.WithResolvers(&multipleEndpointsGrpcResolverBuilder{}),
	}

	dialOpts = append(dialOpts, extraOpts...)
	return dialOpts
}

func NewGrpcBeaconChainClient(endpoint string) ethpb.BeaconChainClient {

	dialOpts := ConstructDialOptions(
		1024*1024,
		"",
		5,
		time.Second,
	)
	if dialOpts == nil {
		return nil
	}

	grpcConn, err := grpc.DialContext(context.Background(), endpoint, dialOpts...)
	if err != nil {
		return nil
	}
	return ethpb.NewBeaconChainClient(grpcConn)
}
