package grpc

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	escrowpb "github.com/tron-us/go-btfs-common/protos/escrow"
	exchangepb "github.com/tron-us/go-btfs-common/protos/exchange"
	guardpb "github.com/tron-us/go-btfs-common/protos/guard"
	hubpb "github.com/tron-us/go-btfs-common/protos/hub"
	ledgerpb "github.com/tron-us/go-btfs-common/protos/ledger"
	tronpb "github.com/tron-us/go-btfs-common/protos/protocol/api"
	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"
	statuspb "github.com/tron-us/go-btfs-common/protos/status"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	defaultSchema  = "http"
	defaultTimeout = 30 * time.Second
)

func (g *ClientBuilder) doWithContext(ctx context.Context, f interface{}) error {
	newCtx, cancel := context.WithTimeout(ctx, g.timeout)
	if cancel != nil {
		defer cancel()
	}
	conn, err := newGRPCConn(newCtx, g.addr)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		return err
	}
	if conn == nil || conn.GetState() != connectivity.Ready {
		return errors.New("failed to get connection")
	}
	switch v := f.(type) {
	case func(context.Context, statuspb.StatusServiceClient) error:
		return v(ctx, statuspb.NewStatusServiceClient(conn))
	case func(context.Context, hubpb.HubQueryServiceClient) error:
		return v(ctx, hubpb.NewHubQueryServiceClient(conn))
	case func(context.Context, guardpb.GuardServiceClient) error:
		return v(ctx, guardpb.NewGuardServiceClient(conn))
	case func(context.Context, escrowpb.EscrowServiceClient) error:
		return v(ctx, escrowpb.NewEscrowServiceClient(conn))
	case func(ctx context.Context, client sharedpb.RuntimeServiceClient) error:
		return v(ctx, sharedpb.NewRuntimeServiceClient(conn))
	case func(ctx context.Context, client grpc_health_v1.HealthClient) error:
		return v(ctx, grpc_health_v1.NewHealthClient(conn))
	case func(ctx context.Context, client ledgerpb.ChannelsClient) error:
		return v(ctx, ledgerpb.NewChannelsClient(conn))
	case func(ctx context.Context, client exchangepb.ExchangeClient) error:
		return v(ctx, exchangepb.NewExchangeClient(conn))
	case func(ctx context.Context, client tronpb.WalletSolidityClient) error:
		return v(ctx, tronpb.NewWalletSolidityClient(conn))
	default:
		return fmt.Errorf("illegal function: %T", f)
	}
}

type ClientBuilder struct {
	addr    string
	timeout time.Duration
}

func (b *ClientBuilder) Timeout(to time.Duration) {
	b.timeout = to
}

func builder(address string) ClientBuilder {
	return ClientBuilder{
		addr:    address,
		timeout: defaultTimeout,
	}
}

func newGRPCConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	u, err := parse(addr)
	if err != nil {
		return nil, err
	}
	opts := []grpc.DialOption{grpc.WithBlock()}
	if u.schema == "http" {
		opts = append(opts, grpc.WithInsecure())
	} else if u.schema == "https" {
		c := credentials.NewTLS(&tls.Config{})
		opts = append(opts, grpc.WithTransportCredentials(c))
	} else {
		return nil, fmt.Errorf("not supported schema: %v", u.schema)
	}
	return grpc.DialContext(ctx, fmt.Sprintf("%s:%d", u.host, u.port), opts...)
}

func parse(rawU string) (*parsedURL, error) {
	if strings.Index(rawU, "//") == 0 {
		rawU = defaultSchema + ":" + rawU
	}
	if !strings.Contains(rawU, "://") {
		rawU = defaultSchema + "://" + rawU
	}
	u, err := url.Parse(rawU)
	if err != nil {
		return nil, err
	}
	h := u.Hostname()
	result := new(parsedURL)
	result.schema = u.Scheme
	result.host = h
	result.port, err = getPort(u)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getPort(u *url.URL) (int, error) {
	p := u.Port()
	if p == "" {
		switch u.Scheme {
		case "http":
			p = "80"
		case "https":
			p = "443"
		default:
			return -1, fmt.Errorf("not support schema: %v", u.Scheme)
		}
	}
	return strconv.Atoi(p)
}

type parsedURL struct {
	schema string
	host   string
	port   int
}
