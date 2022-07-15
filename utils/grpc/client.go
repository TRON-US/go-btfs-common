package grpc

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/tron-us/go-btfs-common/protos/online"
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
	InjectedConn   = "injected-connection"
)

func (g *ClientBuilder) doWithContext(ctx context.Context, f interface{}) error {
	newCtx, cancel := context.WithTimeout(ctx, g.timeout)
	if cancel != nil {
		defer cancel()
	}
	var conn *grpc.ClientConn
	if instance := ctx.Value(InjectedConn); instance != nil {
		conn = instance.(*grpc.ClientConn)
	} else if c, err := newGRPCConn(newCtx, g.addr); err == nil {
		if c != nil {
			defer c.Close()
		} else if c.GetState() != connectivity.Ready {
			return errors.New("failed to get connection")
		}
		conn = c
	} else {
		return err
	}
	switch v := f.(type) {
	case func(context.Context, online.OnlineServiceClient) error:
		return wrapError("OnlineClient", v(ctx, online.NewOnlineServiceClient(conn)))
	case func(context.Context, statuspb.StatusServiceClient) error:
		return wrapError("StatusClient", v(ctx, statuspb.NewStatusServiceClient(conn)))
	case func(context.Context, hubpb.HubQueryServiceClient) error:
		return wrapError("HubQueryClient", v(ctx, hubpb.NewHubQueryServiceClient(conn)))
	case func(context.Context, hubpb.HubParseServiceClient) error:
		return wrapError("HubParseClient", v(ctx, hubpb.NewHubParseServiceClient(conn)))
	case func(context.Context, guardpb.GuardServiceClient) error:
		return wrapError("GuardClient", v(ctx, guardpb.NewGuardServiceClient(conn)))
	case func(context.Context, escrowpb.EscrowServiceClient) error:
		return wrapError("EscrowClient", v(ctx, escrowpb.NewEscrowServiceClient(conn)))
	case func(ctx context.Context, client sharedpb.RuntimeServiceClient) error:
		return wrapError("RuntimeClient", v(ctx, sharedpb.NewRuntimeServiceClient(conn)))
	case func(ctx context.Context, client grpc_health_v1.HealthClient) error:
		return wrapError("HealthClient", v(ctx, grpc_health_v1.NewHealthClient(conn)))
	case func(ctx context.Context, client ledgerpb.ChannelsClient) error:
		return wrapError("LedgerClient", v(ctx, ledgerpb.NewChannelsClient(conn)))
	case func(ctx context.Context, client exchangepb.ExchangeClient) error:
		return wrapError("ExchangeClient", v(ctx, exchangepb.NewExchangeClient(conn)))
	case func(ctx context.Context, client tronpb.WalletSolidityClient) error:
		return wrapError("WalletSolidityClient", v(ctx, tronpb.NewWalletSolidityClient(conn)))
	case func(ctx context.Context, client tronpb.WalletClient) error:
		return wrapError("WalletClient", v(ctx, tronpb.NewWalletClient(conn)))
	default:
		return fmt.Errorf("illegal function: %T", f)
	}
}

func wrapError(prefix string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%v: %v", prefix, err)
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
