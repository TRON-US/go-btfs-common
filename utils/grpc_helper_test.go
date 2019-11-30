package utils

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"testing"
	"time"
)

func TestGRPCWithContextTimeout(t *testing.T) {
	f := func(ctx context.Context, conn *grpc.ClientConn) error {
		return nil
	}
	tests := []struct {
		addr string
		f    func(ctx context.Context, conn *grpc.ClientConn) error
		err  bool
	}{
		{addr: "ftp://db-grpc-dev.btfs.io", f: f, err: true},
		{addr: "http://db-grpc-dev.btfs.io:443", f: f, err: true},
		{addr: "https://db-grpc-dev.btfs.io", f: func(ctx context.Context, conn *grpc.ClientConn) error {
			return errors.New("just an error")
		}, err: true},
		{addr: "https://db-grpc-dev.btfs.io", f: f, err: false},
	}
	for _, tt := range tests {
		err := GRPCWithContextTimeout(context.Background(), tt.addr, tt.f)
		if !tt.err && err != nil {
			t.Errorf(`%v, %T: unexpected error "%v"`, tt.addr, tt.f, err)
			continue
		}
		if tt.err && err == nil {
			t.Errorf(`%v, %T: expected error`, tt.addr, tt.f)
		}
	}
}

func TestNewGRPCConn(t *testing.T) {
	tests := []struct {
		in  string
		out connectivity.State
		err bool
	}{
		{in: "ftp://db-grpc-dev.btfs.io", err: true},
		{in: "http://db-grpc-dev.btfs.io:443", err: true},
		{in: "https://db-grpc-dev.btfs.io", out: connectivity.Ready},
	}
	for _, tt := range tests {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
		conn, err := newGRPCConn(ctx, tt.in)
		if cancelFunc != nil {
			defer cancelFunc()
		}
		if !tt.err && err != nil {
			t.Errorf(`%v: unexpected error "%v"`, tt.in, err)
			continue
		}
		if tt.err && err == nil {
			t.Errorf(`%v: expected error`, tt.in)
		}
		if conn != nil {
			defer conn.Close()
			state := conn.GetState()
			if tt.out != state {
				// don't want to break the test
				t.Errorf(`%v: got "%v", want "%v"`, tt.in, state, tt.out)
			}
		}
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		in  string
		out *parsedURL
		err bool
	}{
		{in: "", err: true},
		{in: "/", err: true},
		{in: "//", err: true},
		{in: "http:/www.google.com", err: true},
		{in: "http:///www.google.com", err: true},
		{in: "javascript:void(0)", err: true},
		{in: "<script>", err: true},
		{in: "127.0.0:8080", err: true},
		{in: "127.0.0", err: true},
		{in: "127.0.0.0.1", err: true},
		{in: "ftp://127.0.0.1", err: true},
		{in: "//localhost", out: &parsedURL{schema: "http", host: "localhost", port: 80}},
		{in: "localhost", out: &parsedURL{schema: "http", host: "localhost", port: 80}},
		{in: "localhost:8080", out: &parsedURL{schema: "http", host: "localhost", port: 8080}},
		{in: "btfs.io", out: &parsedURL{schema: "http", host: "btfs.io", port: 80}},
		{in: "btfs.io:8080", out: &parsedURL{schema: "http", host: "btfs.io", port: 8080}},
		{in: "127.0.0.1", out: &parsedURL{schema: "http", host: "127.0.0.1", port: 80}},
		{in: "127.0.0.1:8080", out: &parsedURL{schema: "http", host: "127.0.0.1", port: 8080}},
		{in: "https://btfs.io", out: &parsedURL{schema: "https", host: "btfs.io", port: 443}},
		{in: "http://btfs.io", out: &parsedURL{schema: "http", host: "btfs.io", port: 80}},
		{in: "https://btfs.io:50051", out: &parsedURL{schema: "https", host: "btfs.io", port: 50051}},
		{in: "http://btfs.io:2333", out: &parsedURL{schema: "http", host: "btfs.io", port: 2333}},
		{in: "https://127.0.0.1", out: &parsedURL{schema: "https", host: "127.0.0.1", port: 443}},
		{in: "http://127.0.0.1", out: &parsedURL{schema: "http", host: "127.0.0.1", port: 80}},
	}
	for _, tt := range tests {
		url, err := parse(tt.in)
		if !tt.err && err != nil {
			t.Errorf(`%v: unexpected error "%v"`, tt.in, err)
			continue
		}
		if tt.err && err == nil {
			t.Errorf(`%v: expected error`, tt.in)
		}
		if tt.out != nil {
			if tt.out.schema != url.schema || tt.out.host != url.host || tt.out.port != url.port {
				t.Errorf(`%v: got "%v", want "%v"`, tt.in, url, tt.out)
			}
		}
	}
}
