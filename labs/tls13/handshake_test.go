package tls13

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestLocalHandshakeCapturesClientHelloAndKeyLog(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	got, err := LocalHandshake(ctx)
	if err != nil {
		t.Fatal(err)
	}
	// RFC 8446: ClientHello.legacy_version is 0x0303. The record version Go
	// sends is 0x0301, the middlebox compatibility value from appendix D.
	if got.ClientVersion != 0x0303 {
		t.Fatalf("ClientHello.legacy_version = %#x, want 0x0303", got.ClientVersion)
	}
	if got.RecordVersion != 0x0301 {
		t.Fatalf("record version = %#x, want 0x0301", got.RecordVersion)
	}
	for _, label := range []string{
		"CLIENT_HANDSHAKE_TRAFFIC_SECRET",
		"SERVER_HANDSHAKE_TRAFFIC_SECRET",
		"CLIENT_TRAFFIC_SECRET_0",
		"SERVER_TRAFFIC_SECRET_0",
	} {
		if !strings.Contains(got.KeyLog, label) {
			t.Fatalf("key log missing %s\n%s", label, got.KeyLog)
		}
	}
}

func TestParseClientHelloRejectsShortInput(t *testing.T) {
	if _, _, err := ParseClientHello([]byte{22, 0x03, 0x01}); err == nil {
		t.Fatal("expected an error")
	}
}
