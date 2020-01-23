package grpc

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/connectivity"
)

func TestNewGRPCConn(t *testing.T) {
	tests := []struct {
		in  string
		out connectivity.State
		err bool
	}{
		{in: "ftp://status-dev.btfs.io", err: true},
		{in: "http://status-dev.btfs.io:443", err: true},
		{in: "https://status-dev.btfs.io", out: connectivity.Ready},
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
		{in: "localhost:8080", out: &parsedURL{schema: "http", host: "localhost", port: 8080}},
		{in: "btfs.io", out: &parsedURL{schema: "http", host: "btfs.io", port: 80}},
		{in: "127.0.0.1", out: &parsedURL{schema: "http", host: "127.0.0.1", port: 80}},
		{in: "127.0.0.1:8080", out: &parsedURL{schema: "http", host: "127.0.0.1", port: 8080}},
		{in: "https://btfs.io", out: &parsedURL{schema: "https", host: "btfs.io", port: 443}},
		{in: "http://btfs.io", out: &parsedURL{schema: "http", host: "btfs.io", port: 80}},
		{in: "https://btfs.io:50051", out: &parsedURL{schema: "https", host: "btfs.io", port: 50051}},
		{in: "http://btfs.io:2333", out: &parsedURL{schema: "http", host: "btfs.io", port: 2333}},
		{in: "https://127.0.0.1", out: &parsedURL{schema: "https", host: "127.0.0.1", port: 443}},
		{in: "http://127.0.0.1", out: &parsedURL{schema: "http", host: "127.0.0.1", port: 80}},
		{in: "hub-service:8080", out: &parsedURL{schema: "http", host: "hub-service", port: 8080}},
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
