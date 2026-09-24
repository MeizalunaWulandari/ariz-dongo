package fieldsa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/MeizalunaWulandari/ariz-dongo/internal/logger"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// Login melakukan login ke Fieldsa.
//
// INTERNAL ONLY.
// Tidak pernah menjadi endpoint ARIZ DONGO.
func (c *Client) Login(
	ctx context.Context,
) (string, error) {

	logger.Info("Fieldsa login started")

	if c.username == "" || c.password == "" {
		logger.Error(
			"Fieldsa login failed: credential belum dikonfigurasi",
		)

		return "", fmt.Errorf(
			"Fieldsa credential belum dikonfigurasi",
		)
	}

	payload := loginRequest{
		Username: c.username,
		Password: c.password,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf(
			"marshal login request: %w",
			err,
		)
	}

	url := c.baseURL + "/api/login"

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)

	if err != nil {
		return "", fmt.Errorf(
			"create login request: %w",
			err,
		)
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	req.Header.Set(
		"Accept",
		"application/json",
	)

	if c.loginBearer != "" {
		req.Header.Set(
			"Authorization",
			"Bearer "+c.loginBearer,
		)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf(
			"login request failed: %w",
			err,
		)
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(
		io.LimitReader(resp.Body, 4096),
	)

	if err != nil {
		return "", fmt.Errorf(
			"read login response: %w",
			err,
		)
	}

	if resp.StatusCode < 200 ||
		resp.StatusCode >= 300 {

		logger.Error(
			"Fieldsa login failed: HTTP %d",
			resp.StatusCode,
		)

		return "", fmt.Errorf(
			"login failed: HTTP %d: %s",
			resp.StatusCode,
			strings.TrimSpace(
				string(responseBody),
			),
		)
	}

	var result loginResponse

	if err := json.Unmarshal(
		responseBody,
		&result,
	); err != nil {
		return "", fmt.Errorf(
			"decode login response: %w",
			err,
		)
	}

	if result.Token == "" {
		return "", fmt.Errorf(
			"login response tidak memiliki token",
		)
	}

	c.mu.Lock()

	c.token = result.Token

	// Sementara.
	// Nanti bisa diganti dengan parsing JWT exp.
	c.tokenExpiry = time.Now().Add(
		50 * time.Minute,
	)

	c.mu.Unlock()

	logger.Info(
		"Fieldsa login successful",
	)

	return result.Token, nil
}

// Token mengembalikan token yang tersimpan.
//
// Jika token belum ada atau expired,
// Login() dipanggil otomatis.
func (c *Client) Token(
	ctx context.Context,
) (string, error) {

	c.mu.RLock()

	token := c.token
	expiry := c.tokenExpiry

	c.mu.RUnlock()

	if token != "" &&
		time.Now().Before(expiry) {

		return token, nil
	}

	logger.Info(
		"Fieldsa token unavailable or expired, login required",
	)

	return c.Login(ctx)
}
