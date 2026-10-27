package main

import "testing"

func TestValidateGameServerKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{name: "short key", key: "short", wantErr: true},
		{name: "documented example key", key: "zneondrive-local-server-key-change-me-0001", wantErr: true},
		{name: "valid configured key", key: "ci-web-stack-test-key-not-for-production", wantErr: false},
		{name: "trimmed example key", key: " zneondrive-local-server-key-change-me-0001 ", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validateGameServerKey(tt.key)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateGameServerKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
