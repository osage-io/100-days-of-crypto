package vaulttransit

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestEncryptRotateRewrap(t *testing.T) {
	srv := newFake("root")
	defer srv.Close()
	c := &Client{BaseURL: srv.URL, Token: "root"}
	ctx := context.Background()

	if err := c.CreateKey(ctx, "app"); err != nil {
		t.Fatal(err)
	}
	ver, err := c.LatestVersion(ctx, "app")
	if err != nil {
		t.Fatal(err)
	}
	if ver != 1 {
		t.Fatalf("latest = %d", ver)
	}

	plain := []byte("field cto notebook")
	ct1, err := c.Encrypt(ctx, "app", plain)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(ct1, "vault:v1:") {
		t.Fatalf("ciphertext %q", ct1)
	}
	got, err := c.Decrypt(ctx, "app", ct1)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("decrypt got %q", got)
	}

	if err := c.Rotate(ctx, "app"); err != nil {
		t.Fatal(err)
	}
	ver, err = c.LatestVersion(ctx, "app")
	if err != nil {
		t.Fatal(err)
	}
	if ver != 2 {
		t.Fatalf("after rotate latest = %d", ver)
	}
	// The old version still opens. Rotation did not strand ciphertext.
	got, err = c.Decrypt(ctx, "app", ct1)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("old ciphertext decrypt got %q", got)
	}

	ct2, err := c.Rewrap(ctx, "app", ct1)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(ct2, "vault:v2:") {
		t.Fatalf("rewrap ciphertext %q", ct2)
	}
	if ct2 == ct1 {
		t.Fatal("rewrap returned the original ciphertext")
	}
	got, err = c.Decrypt(ctx, "app", ct2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("rewrap decrypt got %q", got)
	}
}

func TestMissingToken(t *testing.T) {
	srv := newFake("root")
	defer srv.Close()
	c := &Client{BaseURL: srv.URL}
	if err := c.CreateKey(context.Background(), "app"); err == nil {
		t.Fatal("expected a token error")
	}
}

func TestLiveVault(t *testing.T) {
	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")
	if addr == "" || token == "" {
		t.Skip("set VAULT_ADDR and VAULT_TOKEN to run against a dev server")
	}
	c := &Client{BaseURL: addr, Token: token}
	ctx := context.Background()
	name := "days-of-crypto-lab"
	if err := c.CreateKey(ctx, name); err != nil && !strings.Contains(err.Error(), "exists") {
		// A second run against a long-lived dev server may already have the key.
		t.Fatal(err)
	}
	plain := []byte("live transit")
	ct, err := c.Encrypt(ctx, name, plain)
	if err != nil {
		t.Fatal(err)
	}
	got, err := c.Decrypt(ctx, name, ct)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("got %q", got)
	}
}

// fake is a Transit stand-in. The prefix matches Vault's
// vault:v<keyVersion>: template. The blob is nonce || gcm ciphertext,
// which is enough to exercise rotate and rewrap. It is not Vault's blob.
type fake struct {
	token string
	mu    sync.Mutex
	keys  map[string]*ring
}

type ring struct {
	keys [][]byte // index 0 unused
}

func newFake(token string) *httptest.Server {
	f := &fake{token: token, keys: map[string]*ring{}}
	return httptest.NewServer(f)
}

func (f *fake) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Vault-Token") != f.token {
		writeErr(w, http.StatusForbidden, "permission denied")
		return
	}
	body, _ := io.ReadAll(r.Body)
	path := strings.TrimPrefix(r.URL.Path, "/v1/transit/")
	switch {
	case r.Method == http.MethodPost && strings.HasPrefix(path, "keys/") && strings.HasSuffix(path, "/rotate"):
		name := strings.TrimSuffix(strings.TrimPrefix(path, "keys/"), "/rotate")
		f.rotate(w, name)
	case r.Method == http.MethodPost && strings.HasPrefix(path, "keys/"):
		name := strings.TrimPrefix(path, "keys/")
		f.create(w, name)
	case r.Method == http.MethodGet && strings.HasPrefix(path, "keys/"):
		name := strings.TrimPrefix(path, "keys/")
		f.read(w, name)
	case r.Method == http.MethodPost && strings.HasPrefix(path, "encrypt/"):
		f.encrypt(w, strings.TrimPrefix(path, "encrypt/"), body)
	case r.Method == http.MethodPost && strings.HasPrefix(path, "decrypt/"):
		f.decrypt(w, strings.TrimPrefix(path, "decrypt/"), body)
	case r.Method == http.MethodPost && strings.HasPrefix(path, "rewrap/"):
		f.rewrap(w, strings.TrimPrefix(path, "rewrap/"), body)
	default:
		writeErr(w, http.StatusNotFound, "no such path")
	}
}

func (f *fake) create(w http.ResponseWriter, name string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.keys[name]; ok {
		writeErr(w, http.StatusBadRequest, "key already exists")
		return
	}
	k := make([]byte, 32)
	if _, err := rand.Read(k); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	f.keys[name] = &ring{keys: [][]byte{nil, k}}
	w.WriteHeader(http.StatusNoContent)
}

func (f *fake) rotate(w http.ResponseWriter, name string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	ring, ok := f.keys[name]
	if !ok {
		writeErr(w, http.StatusBadRequest, "key not found")
		return
	}
	k := make([]byte, 32)
	if _, err := rand.Read(k); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	ring.keys = append(ring.keys, k)
	w.WriteHeader(http.StatusNoContent)
}

func (f *fake) read(w http.ResponseWriter, name string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	ring, ok := f.keys[name]
	if !ok {
		writeErr(w, http.StatusBadRequest, "key not found")
		return
	}
	writeJSON(w, map[string]any{
		"data": map[string]any{
			"name":           name,
			"type":           "aes256-gcm96",
			"latest_version": len(ring.keys) - 1,
		},
	})
}

func (f *fake) encrypt(w http.ResponseWriter, name string, body []byte) {
	var req struct {
		Plaintext string `json:"plaintext"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	pt, err := base64.StdEncoding.DecodeString(req.Plaintext)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "plaintext must be base64")
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	ring, ok := f.keys[name]
	if !ok {
		writeErr(w, http.StatusBadRequest, "key not found")
		return
	}
	ver := len(ring.keys) - 1
	ct, err := seal(ring.keys[ver], pt)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, map[string]any{"data": map[string]any{"ciphertext": fmt.Sprintf("vault:v%d:%s", ver, ct)}})
}

func (f *fake) decrypt(w http.ResponseWriter, name string, body []byte) {
	var req struct {
		Ciphertext string `json:"ciphertext"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	pt, err := f.open(name, req.Ciphertext)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, map[string]any{"data": map[string]any{
		"plaintext": base64.StdEncoding.EncodeToString(pt),
	}})
}

func (f *fake) rewrap(w http.ResponseWriter, name string, body []byte) {
	var req struct {
		Ciphertext string `json:"ciphertext"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	pt, err := f.open(name, req.Ciphertext)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	f.mu.Lock()
	ring := f.keys[name]
	ver := len(ring.keys) - 1
	key := ring.keys[ver]
	f.mu.Unlock()
	ct, err := seal(key, pt)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, map[string]any{"data": map[string]any{"ciphertext": fmt.Sprintf("vault:v%d:%s", ver, ct)}})
}

func (f *fake) open(name, ciphertext string) ([]byte, error) {
	ver, blob, err := splitCiphertext(ciphertext)
	if err != nil {
		return nil, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	ring, ok := f.keys[name]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}
	if ver < 1 || ver >= len(ring.keys) {
		return nil, fmt.Errorf("ciphertext version %d", ver)
	}
	return open(ring.keys[ver], blob)
}

func splitCiphertext(ct string) (int, string, error) {
	const prefix = "vault:v"
	if !strings.HasPrefix(ct, prefix) {
		return 0, "", fmt.Errorf("missing vault:v prefix")
	}
	rest := strings.TrimPrefix(ct, prefix)
	verText, blob, ok := strings.Cut(rest, ":")
	if !ok {
		return 0, "", fmt.Errorf("missing version separator")
	}
	ver, err := strconv.Atoi(verText)
	if err != nil {
		return 0, "", err
	}
	return ver, blob, nil
}

func seal(key, plaintext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	sealed := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func open(key []byte, blob string) ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(raw) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ct := raw[:gcm.NonceSize()], raw[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ct, nil)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]any{"errors": []string{msg}})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
