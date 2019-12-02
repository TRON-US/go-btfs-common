package grpc

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tron-us/go-btfs-common/protos/escrow"
	"github.com/tron-us/go-btfs-common/protos/guard"
	"github.com/tron-us/go-btfs-common/protos/hub"
	"github.com/tron-us/go-btfs-common/protos/status"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
)

const (
	defaultSchema  = "http"
	defaultTimeout = 30 * time.Second
)

var (
	domainRegexp = regexp.MustCompile(`^(localhost)|([a-zA-Z0-9-]{1,63}\.)+([a-zA-Z]{1,63})$`)
	ipv4Regexp   = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)
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
	case func(context.Context, status.StatusClient) error:
		return v(ctx, status.NewStatusClient(conn))
	case func(context.Context, hub.HubQueryClient) error:
		return v(ctx, hub.NewHubQueryClient(conn))
	case func(context.Context, guard.GuardServiceClient) error:
		return v(ctx, guard.NewGuardServiceClient(conn))
	case func(context.Context, escrow.EscrowServiceClient) error:
		return v(ctx, escrow.NewEscrowServiceClient(conn))
	default:
		return fmt.Errorf("illegal function: %T", f)
	}
}

type ClientBuilder struct {
	addr    string
	timeout time.Duration
}

func (b *ClientBuilder) Timeout(to time.Duration) *ClientBuilder {
	b.timeout = to
	return b
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
	if err := checkHost(h); err != nil {
		return nil, err
	}
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

func checkHost(host string) error {
	if host == "" {
		return errors.New("empty host")
	}
	host = strings.ToLower(host)
	if domainRegexp.MatchString(host) {
		return nil
	}
	if ipv4Regexp.MatchString(host) {
		return nil
	}
	return fmt.Errorf("invalid host: %v", host)
}

type parsedURL struct {
	schema string
	host   string
	port   int
}
