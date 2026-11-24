// Package tls13 runs a local TLS 1.3 handshake and keeps the ClientHello
// plus the NSS key log.
//
// The certificate is generated in the test process and trusted only by that
// client. Nothing here listens on a public address on purpose. KeyLogWriter
// is how you feed a handshake to Wireshark. Treat the log as a secret.
package tls13

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"time"
)

// Result is what a local handshake leaves behind.
type Result struct {
	// KeyLog is the NSS key log from the client. It contains traffic secrets.
	KeyLog string
	// ClientHello is the first record the client wrote.
	ClientHello []byte
	// RecordVersion is the legacy record-layer version on that ClientHello.
	RecordVersion uint16
	// ClientVersion is ClientHello.legacy_version.
	ClientVersion uint16
}

// LocalHandshake dials a TLS 1.3 server on 127.0.0.1. The server certificate
// is self-signed and installed only in this client's root pool.
func LocalHandshake(ctx context.Context) (*Result, error) {
	cert, certPEM, err := newCert()
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(certPEM) {
		return nil, errors.New("tls13: failed to install the lab certificate")
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	defer ln.Close()

	errCh := make(chan error, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			errCh <- err
			return
		}
		defer conn.Close()
		server := tls.Server(conn, &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS13,
			MaxVersion:   tls.VersionTLS13,
		})
		errCh <- server.HandshakeContext(ctx)
		server.Close()
	}()

	raw, err := (&net.Dialer{}).DialContext(ctx, "tcp", ln.Addr().String())
	if err != nil {
		return nil, err
	}
	rec := &recordingConn{Conn: raw}
	var keyLog buffer
	client := tls.Client(rec, &tls.Config{
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
		RootCAs:      pool,
		ServerName:   "localhost",
		KeyLogWriter: &keyLog,
	})
	if err := client.HandshakeContext(ctx); err != nil {
		client.Close()
		return nil, err
	}
	client.Close()
	if err := <-errCh; err != nil {
		return nil, err
	}

	recordVersion, clientVersion, err := ParseClientHello(rec.out)
	if err != nil {
		return nil, err
	}
	return &Result{
		KeyLog:        keyLog.String(),
		ClientHello:   append([]byte(nil), rec.out...),
		RecordVersion: recordVersion,
		ClientVersion: clientVersion,
	}, nil
}

// ParseClientHello reads the first TLS record and expects a ClientHello.
// It returns the legacy record version and ClientHello.legacy_version.
func ParseClientHello(record []byte) (recordVersion, clientVersion uint16, err error) {
	if len(record) < 11 {
		return 0, 0, errors.New("tls13: record too short")
	}
	if record[0] != 22 {
		return 0, 0, fmt.Errorf("tls13: content type %d, want handshake (22)", record[0])
	}
	recordVersion = binary.BigEndian.Uint16(record[1:3])
	recLen := int(binary.BigEndian.Uint16(record[3:5]))
	if recLen < 6 || len(record) < 5+recLen {
		return 0, 0, errors.New("tls13: truncated record")
	}
	hs := record[5 : 5+recLen]
	if hs[0] != 1 {
		return 0, 0, fmt.Errorf("tls13: handshake type %d, want ClientHello (1)", hs[0])
	}
	clientVersion = binary.BigEndian.Uint16(hs[4:6])
	return recordVersion, clientVersion, nil
}

func newCert() (tls.Certificate, []byte, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, nil, err
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "localhost"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		DNSNames:     []string{"localhost"},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return tls.Certificate{}, nil, err
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return tls.Certificate{}, nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return tls.Certificate{}, nil, err
	}
	return cert, certPEM, nil
}

type recordingConn struct {
	net.Conn
	out []byte
}

func (c *recordingConn) Write(p []byte) (int, error) {
	n, err := c.Conn.Write(p)
	c.out = append(c.out, p[:n]...)
	return n, err
}

type buffer struct {
	b []byte
}

func (b *buffer) Write(p []byte) (int, error) {
	b.b = append(b.b, p...)
	return len(p), nil
}

func (b *buffer) String() string { return string(b.b) }
