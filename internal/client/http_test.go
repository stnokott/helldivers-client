package client

import (
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestRateLimitHTTPClientRetryAfter(t *testing.T) {
	c := &rateLimitHTTPClient{
		log: log.Default(),
	}

	type args struct {
		header http.Header
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "regular",
			args: args{
				header: http.Header{"Retry-After": []string{"56"}},
			},
			want: 56 * time.Second,
		},
		{
			name: "not a number",
			args: args{
				header: http.Header{"Retry-After": []string{"Foo"}},
			},
			want: defaultBackoff,
		},
		{
			name: "not present",
			args: args{
				header: http.Header{"Foo": []string{"Bar"}},
			},
			want: defaultBackoff,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.retryAfter(tt.args.header); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rateLimitHTTPClient.retryAfter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRateLimitHTTPClient(t *testing.T) {
	tests := []struct {
		name       string
		serverFunc func() http.HandlerFunc
		wantErr    bool
	}{
		{
			name: "HTTP 200",
			serverFunc: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
					w.Write([]byte("OK"))
				}
			},
			wantErr: false,
		},
		{
			name: "HTTP 500 once",
			serverFunc: func() http.HandlerFunc {
				ok := false
				ptrOK := &ok
				return func(w http.ResponseWriter, r *http.Request) {
					if !*ptrOK {
						w.WriteHeader(500)
						w.Write([]byte("Internal server error"))
						*ptrOK = true
					} else {
						w.WriteHeader(200)
						w.Write([]byte("OK"))
					}
				}
			},
			wantErr: false,
		},
		{
			name: "HTTP 429 once",
			serverFunc: func() http.HandlerFunc {
				ok := false
				ptrOK := &ok
				return func(w http.ResponseWriter, r *http.Request) {
					if !*ptrOK {
						w.Header().Add("Retry-After", "3")
						w.WriteHeader(429)
						w.Write([]byte("Too many requests"))
						*ptrOK = true
					} else {
						w.WriteHeader(200)
						w.Write([]byte("OK"))
					}
				}
			},
			wantErr: false,
		},
		{
			name: "HTTP 500 forever",
			serverFunc: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(500)
					w.Write([]byte("Internal server error"))
				}
			},
			wantErr: true,
		},
		{
			name: "HTTP 429 forever",
			serverFunc: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("Retry-After", "1")
					w.WriteHeader(429)
					w.Write([]byte("Too many requests"))
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.serverFunc())
			defer server.Close()

			c := rateLimitHTTPClient{
				client:   server.Client(),
				maxRetry: 3,
				log:      log.Default(),
			}
			req, err := http.NewRequest("GET", server.URL, nil)
			if err != nil {
				t.Errorf("failed to create request: %v", err)
				return
			}

			resp, err := c.Do(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("rateLimitHTTPClient.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				defer func() {
					_ = resp.Body.Close()
				}()
			}
		})
	}
}
