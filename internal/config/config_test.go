// Package config handles configuration via environment variables
package config

import (
	"os"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	for _, k := range []string{"POSTGRES_URI", "API_URL", "WORKER_CRON", "HEALTHCHECKS_URL"} {
		_ = os.Unsetenv(k)
	}

	tests := []struct {
		name    string
		env     map[string]string
		want    *Config
		wantErr bool
	}{
		{
			name: "all valid",
			env: map[string]string{
				"POSTGRES_URI":     "postgresql://user:pass@localhost:5432/database",
				"API_URL":          "http://localhost:4000",
				"WORKER_CRON":      "*/5 * * * *",
				"HEALTHCHECKS_URL": "https://hc-ping.com/11223344",
			},
			want: &Config{
				PostgresURI:     "postgresql://user:pass@localhost:5432/database",
				APIRootURL:      "http://localhost:4000",
				WorkerCron:      "*/5 * * * *",
				HealthchecksURL: "https://hc-ping.com/11223344",
			},
			wantErr: false,
		},
		{
			name: "required missing",
			env: map[string]string{
				"API_URL":          "http://localhost:4000",
				"WORKER_CRON":      "*/5 * * * *",
				"HEALTHCHECKS_URL": "https://hc-ping.com/11223344",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "required empty",
			env: map[string]string{
				"POSTGRES_URI":     "",
				"API_URL":          "http://localhost:4000",
				"WORKER_CRON":      "*/5 * * * *",
				"HEALTHCHECKS_URL": "https://hc-ping.com/11223344",
			},
			want: &Config{
				PostgresURI:     "",
				APIRootURL:      "http://localhost:4000",
				WorkerCron:      "*/5 * * * *",
				HealthchecksURL: "https://hc-ping.com/11223344",
			},
			wantErr: false,
		},
		{
			name: "optional missing",
			env: map[string]string{
				"POSTGRES_URI": "postgresql://user:pass@localhost:5432/database",
				"API_URL":      "http://localhost:4000",
			},
			want: &Config{
				PostgresURI:     "postgresql://user:pass@localhost:5432/database",
				APIRootURL:      "http://localhost:4000",
				WorkerCron:      "*/5 * * * *",
				HealthchecksURL: "",
			},
			wantErr: false,
		},
		{
			name: "invalid URI", // should not produce an error since we don't validate URIs/URLs when parsing config
			env: map[string]string{
				"POSTGRES_URI": "foobar",
				"API_URL":      "fuzzbuzz",
			},
			want: &Config{
				PostgresURI:     "foobar",
				APIRootURL:      "fuzzbuzz",
				WorkerCron:      "*/5 * * * *",
				HealthchecksURL: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			got, err := Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
