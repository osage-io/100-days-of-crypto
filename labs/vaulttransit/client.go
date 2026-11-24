// Package vaulttransit is a small client for the Vault Transit encrypt,
// decrypt, rotate, and rewrap calls.
//
// It is not the official Vault SDK and it is not a Vault server. The tests
// run it against an in-process stand-in that follows the HTTP shape:
// base64 plaintext in, a vault:v<keyVersion>: ciphertext out. The bytes
// after that prefix are not Vault's on-disk layout. Point BaseURL at
// `vault server -dev` if you want the real engine. See client_test.go.
package vaulttransit

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client talks to a Transit mount.
type Client struct {
	// BaseURL is the Vault address, with no /v1 suffix. Example: http://127.0.0.1:8200
	BaseURL string
	// Token is sent as X-Vault-Token.
	Token string
	// HTTP is optional. nil means http.DefaultClient.
	HTTP *http.Client
}

// CreateKey creates an aes256-gcm96 key.
func (c *Client) CreateKey(ctx context.Context, name string) error {
	return c.post(ctx, "/v1/transit/keys/"+url.PathEscape(name), map[string]string{
		"type": "aes256-gcm96",
	}, nil)
}

// LatestVersion reads the key and returns latest_version.
func (c *Client) LatestVersion(ctx context.Context, name string) (int, error) {
	var resp struct {
		Data struct {
			LatestVersion int `json:"latest_version"`
		} `json:"data"`
	}
	if err := c.get(ctx, "/v1/transit/keys/"+url.PathEscape(name), &resp); err != nil {
		return 0, err
	}
	return resp.Data.LatestVersion, nil
}

// Encrypt base64-encodes plaintext and returns the ciphertext string.
func (c *Client) Encrypt(ctx context.Context, name string, plaintext []byte) (string, error) {
	var resp struct {
		Data struct {
			Ciphertext string `json:"ciphertext"`
		} `json:"data"`
	}
	err := c.post(ctx, "/v1/transit/encrypt/"+url.PathEscape(name), map[string]string{
		"plaintext": base64.StdEncoding.EncodeToString(plaintext),
	}, &resp)
	if err != nil {
		return "", err
	}
	return resp.Data.Ciphertext, nil
}

// Decrypt returns the plaintext bytes from a Transit ciphertext.
func (c *Client) Decrypt(ctx context.Context, name, ciphertext string) ([]byte, error) {
	var resp struct {
		Data struct {
			Plaintext string `json:"plaintext"`
		} `json:"data"`
	}
	err := c.post(ctx, "/v1/transit/decrypt/"+url.PathEscape(name), map[string]string{
		"ciphertext": ciphertext,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(resp.Data.Plaintext)
}

// Rotate adds a new key version. Older ciphertexts still decrypt.
func (c *Client) Rotate(ctx context.Context, name string) error {
	return c.post(ctx, "/v1/transit/keys/"+url.PathEscape(name)+"/rotate", map[string]any{}, nil)
}

// Rewrap decrypts ciphertext with the version in its prefix and encrypts
// it with the latest version. The plaintext is not in the response.
func (c *Client) Rewrap(ctx context.Context, name, ciphertext string) (string, error) {
	var resp struct {
		Data struct {
			Ciphertext string `json:"ciphertext"`
		} `json:"data"`
	}
	err := c.post(ctx, "/v1/transit/rewrap/"+url.PathEscape(name), map[string]string{
		"ciphertext": ciphertext,
	}, &resp)
	if err != nil {
		return "", err
	}
	return resp.Data.Ciphertext, nil
}

func (c *Client) post(ctx context.Context, path string, body any, out any) error {
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url(path), bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(req, out)
}

func (c *Client) get(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url(path), nil)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

func (c *Client) do(req *http.Request, out any) error {
	if c.Token != "" {
		req.Header.Set("X-Vault-Token", c.Token)
	}
	hc := c.HTTP
	if hc == nil {
		hc = http.DefaultClient
	}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		var er struct {
			Errors []string `json:"errors"`
		}
		_ = json.Unmarshal(payload, &er)
		msg := strings.Join(er.Errors, "; ")
		if msg == "" {
			msg = strings.TrimSpace(string(payload))
		}
		return fmt.Errorf("vaulttransit: %s: %s", resp.Status, msg)
	}
	if out == nil || len(payload) == 0 {
		return nil
	}
	return json.Unmarshal(payload, out)
}

func (c *Client) url(path string) string {
	return strings.TrimRight(c.BaseURL, "/") + path
}
