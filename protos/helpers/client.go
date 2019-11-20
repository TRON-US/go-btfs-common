package helpers

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	defaultSchema = "http"
)

var (
	domainRegexp = regexp.MustCompile(`^([a-zA-Z0-9-]{1,63}\.)*([a-zA-Z]{1,63})$`)
	ipv4Regexp   = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)
)

type parsedURL struct {
	schema string
	host   string
	port   int
}

func NewGRPCConn(ctx context.Context, address string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	u, err := parse(address)
	if err != nil {
		return nil, err
	}
	if u.schema == "http" {
		opts = append(opts, grpc.WithInsecure())
	} else if u.schema == "https" {
		c := credentials.NewTLS(&tls.Config{})
		opts = append(opts, grpc.WithTransportCredentials(c))
	} else {
		return nil, errors.New("not supported schema")
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
	h, p, err := splitHostPort(u)
	if err != nil {
		return nil, err
	}
	if err := checkHost(h); err != nil {
		return nil, err
	}
	result := new(parsedURL)
	result.schema = u.Scheme
	result.host = h
	result.port, err = getPort(p, u.Scheme)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getPort(p string, schema string) (int, error) {
	if p == "" {
		switch schema {
		case "http":
			p = "80"
		case "https":
			p = "443"
		default:
			return -1, fmt.Errorf("not support schema: %v", schema)
		}
	}
	return strconv.Atoi(p)
}
func checkHost(host string) error {
	if host == "" {
		return fmt.Errorf("empty host")
	}

	host = strings.ToLower(host)
	if domainRegexp.MatchString(host) {
		return nil
	}

	if ipv4Regexp.MatchString(host) {
		return nil
	}

	return fmt.Errorf("invalid host")
}

func splitHostPort(u *url.URL) (host, port string, err error) {
	if u == nil {
		return "", "", errors.New("empty url")
	}
	host = u.Host

	i := strings.LastIndex(host, ":")
	if i == -1 {
		return host, "", nil
	}

	if i == len(host)-1 {
		return "", "", errors.New("empty port")
	}

	port = host[i+1:]
	host = host[:i]

	if _, err := strconv.Atoi(port); err != nil {
		return "", "", err
	}

	return host, port, nil
}
