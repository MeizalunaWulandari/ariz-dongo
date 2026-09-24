package fieldsa

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/MeizalunaWulandari/ariz-dongo/internal/logger"
)

type Config struct {
	BaseURL         string
	DispatchBaseURL string
	Username        string
	Password        string
	LoginBearer     string
}

type Client struct {
	baseURL         string
	dispatchBaseURL string
	username        string
	password        string
	loginBearer     string

	httpClient *http.Client

	mu          sync.RWMutex
	token       string
	tokenExpiry time.Time
}

func New(cfg Config) *Client {
	return &Client{
		baseURL: strings.TrimRight(
			cfg.BaseURL,
			"/",
		),

		dispatchBaseURL: strings.TrimRight(
			cfg.DispatchBaseURL,
			"/",
		),

		username:    cfg.Username,
		password:    cfg.Password,
		loginBearer: cfg.LoginBearer,

		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Get melakukan GET request ke Fieldsa.
// Token, Authorization, HTTP request, response,
// dan retry 401 ditangani secara internal.
func (c *Client) Get(
	ctx context.Context,
	baseURL string,
	path string,
) ([]byte, error) {

	token, err := c.Token(ctx)
	if err != nil {
		return nil, err
	}

	data, statusCode, err := c.doGet(
		ctx,
		baseURL,
		path,
		token,
	)

	if err != nil {
		return nil, err
	}

	// Token mungkin expired sebelum waktu
	// expiry yang kita simpan.
	if statusCode == http.StatusUnauthorized {
		logger.Warn(
			"Fieldsa returned HTTP 401, refreshing token",
		)

		c.clearToken()

		token, err = c.Login(ctx)
		if err != nil {
			return nil, err
		}

		data, statusCode, err = c.doGet(
			ctx,
			baseURL,
			path,
			token,
		)

		if err != nil {
			return nil, err
		}
	}

	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf(
			"Fieldsa request failed: HTTP %d: %s",
			statusCode,
			strings.TrimSpace(string(data)),
		)
	}

	return data, nil
}

// Post melakukan POST request ke Fieldsa.
// Disiapkan untuk endpoint Fieldsa berikutnya.
func (c *Client) Post(
	ctx context.Context,
	baseURL string,
	path string,
	payload any,
) ([]byte, error) {

	token, err := c.Token(ctx)
	if err != nil {
		return nil, err
	}

	data, statusCode, err := c.doPost(
		ctx,
		baseURL,
		path,
		token,
		payload,
	)

	if err != nil {
		return nil, err
	}

	if statusCode == http.StatusUnauthorized {
		logger.Warn(
			"Fieldsa returned HTTP 401, refreshing token",
		)

		c.clearToken()

		token, err = c.Login(ctx)
		if err != nil {
			return nil, err
		}

		data, statusCode, err = c.doPost(
			ctx,
			baseURL,
			path,
			token,
			payload,
		)

		if err != nil {
			return nil, err
		}
	}

	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf(
			"Fieldsa request failed: HTTP %d: %s",
			statusCode,
			strings.TrimSpace(string(data)),
		)
	}

	return data, nil
}

func (c *Client) doGet(
	ctx context.Context,
	baseURL string,
	path string,
	token string,
) ([]byte, int, error) {

	url := strings.TrimRight(
		baseURL,
		"/",
	) + "/" + strings.TrimLeft(
		path,
		"/",
	)

	logger.Info(
		"Fieldsa GET: %s",
		url,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		return nil, 0, fmt.Errorf(
			"create GET request: %w",
			err,
		)
	}

	c.setCommonHeaders(req, token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"GET request failed: %w",
			err,
		)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf(
			"read GET response: %w",
			err,
		)
	}

	logger.Info(
		"Fieldsa GET response: HTTP %d",
		resp.StatusCode,
	)

	return body, resp.StatusCode, nil
}

func (c *Client) doPost(
	ctx context.Context,
	baseURL string,
	path string,
	token string,
	payload any,
) ([]byte, int, error) {

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"marshal POST payload: %w",
			err,
		)
	}

	url := strings.TrimRight(
		baseURL,
		"/",
	) + "/" + strings.TrimLeft(
		path,
		"/",
	)

	logger.Info(
		"Fieldsa POST: %s",
		url,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		strings.NewReader(string(body)),
	)

	if err != nil {
		return nil, 0, fmt.Errorf(
			"create POST request: %w",
			err,
		)
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	c.setCommonHeaders(req, token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"POST request failed: %w",
			err,
		)
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf(
			"read POST response: %w",
			err,
		)
	}

	logger.Info(
		"Fieldsa POST response: HTTP %d",
		resp.StatusCode,
	)

	return responseBody, resp.StatusCode, nil
}

func (c *Client) setCommonHeaders(
	req *http.Request,
	token string,
) {
	req.Header.Set(
		"Accept",
		"application/json",
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)
}

func (c *Client) clearToken() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.token = ""
	c.tokenExpiry = time.Time{}
}
