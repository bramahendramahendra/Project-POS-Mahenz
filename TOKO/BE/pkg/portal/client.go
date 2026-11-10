// Package portal berisi klien HTTP untuk memanggil API Portal (server ke server).
package portal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	// ErrNotConfigured: alamat Portal belum diisi di config.
	ErrNotConfigured = errors.New("alamat Portal belum dikonfigurasi")
	// ErrRejected: Portal menolak permintaan (kode tidak valid, kedaluwarsa, kunci salah, dst).
	ErrRejected = errors.New("permintaan SSO ditolak Portal")
	// ErrUnavailable: Portal tidak dapat dihubungi atau membalas tidak wajar.
	ErrUnavailable = errors.New("Portal tidak dapat dihubungi")
)

// ExchangeUser: identitas user dari Portal.
type ExchangeUser struct {
	PortalUserID int    `json:"portal_user_id"`
	Username     string `json:"username"`
	FullName     string `json:"full_name"`
	IsActive     bool   `json:"is_active"`
}

// ExchangeResult: hasil penukaran kode SSO.
type ExchangeResult struct {
	Valid bool         `json:"valid"`
	AppID string       `json:"app_id"`
	User  ExchangeUser `json:"user"`
}

// Client memanggil API Portal.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type exchangeBody struct {
	AppID string `json:"app_id"`
	Code  string `json:"code"`
}

type envelope struct {
	Status bool           `json:"status"`
	Data   ExchangeResult `json:"data"`
}

// Exchange menukar kode SSO ke Portal. appKey dikirim lewat header X-App-Key.
func (c *Client) Exchange(ctx context.Context, appID, appKey, code string) (*ExchangeResult, error) {
	if c.baseURL == "" {
		return nil, ErrNotConfigured
	}

	payload, err := json.Marshal(exchangeBody{AppID: appID, Code: code})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/sso/exchange", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Key", appKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}

	switch {
	case resp.StatusCode == http.StatusOK:
		// lanjut parsing
	case resp.StatusCode >= 400 && resp.StatusCode < 500:
		return nil, ErrRejected
	default:
		return nil, fmt.Errorf("%w: status %d", ErrUnavailable, resp.StatusCode)
	}

	var env envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("%w: balasan tidak dapat dibaca", ErrUnavailable)
	}
	if !env.Status {
		return nil, ErrRejected
	}
	return &env.Data, nil
}
