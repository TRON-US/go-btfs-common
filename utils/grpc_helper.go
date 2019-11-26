package utils

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	defaultSchema = "http"
	defaultName   = "default"
)

var (
	domainRegexp = regexp.MustCompile(`^(localhost)|([a-zA-Z0-9-]{1,63}\.)+([a-zA-Z]{1,63})$`)
	ipv4Regexp   = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)
	dialTimeouts = map[string]time.Duration{
		GUARD:         30 * time.Second,
		ESCROW:        30 * time.Second,
		HUB:           30 * time.Second,
		STATUS_SERVER: 30 * time.Second,
		defaultName:   30 * time.Second,
	}
)

type parsedURL struct {
	schema string
	host   string
	port   int
}

// WithContextTimeout get grpc connection with a default timeout
func GRPCWithContextTimeout(ctx context.Context, address string, f func(ctx context.Context,
	conn *grpc.ClientConn) error) error {
	en := WhoAmI()
	if en == "" {
		en = defaultName
	}
	newCtx, cancel := context.WithTimeout(ctx, dialTimeouts[en])
	conn, err := newGRPCConn(newCtx, address)
	if err != nil {
		return err
	}
	if conn == nil || conn.GetState() != connectivity.Ready {
		return errors.New("failed to get connection")
	}
	defer conn.Close()
	if cancel != nil {
		defer cancel()
	}
	return f(ctx, conn)
}

func newGRPCConn(ctx context.Context, address string) (*grpc.ClientConn, error) {
	u, err := parse(address)
	if err != nil {
		return nil, err
	}
	opts := make([]grpc.DialOption, 0)
	opts = append(opts, grpc.WithBlock())
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
