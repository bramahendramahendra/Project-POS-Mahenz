// Package authclient berisi klien HTTP untuk memanggil API BE User (server ke server).
//
// Setiap permintaan membawa kunci aplikasi Portal (header X-App-Key). Permintaan atas nama
// user yang sedang login membawa X-Acting-User-Id (id akun di BE User); BE User sendiri yang
// memeriksa user itu masih aktif dan, untuk endpoint admin, ber-role admin.
package authclient

import (
	"bytes"
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"portal_api/errors"
)

// Pesan BE User bila kunci aplikasi salah, kosong, atau aplikasinya nonaktif (lihat
// AppKeyMiddleware di AUTH/BE). Kesalahan pemasangan, bukan kesalahan pengguna.
const keyRejectedMessage = "Kunci aplikasi tidak valid"

var (
	// ErrNotConfigured: alamat BE User atau kunci aplikasi belum diatur, atau kunci ditolak.
	ErrNotConfigured = stderrors.New("layanan akun belum diatur")
	// ErrUnavailable: BE User tidak dapat dihubungi atau membalas tidak wajar (5xx).
	ErrUnavailable = stderrors.New("layanan akun tidak dapat dihubungi")
)

// Credentials: identitas Portal di BE User (tabel auth_credentials).
type Credentials struct {
	AppID  string
	AppKey string
}

// CredentialsReader membaca kunci aplikasi. nil, nil = belum diatur.
type CredentialsReader interface {
	GetAuthCredentials() (*Credentials, error)
}

// Envelope: bentuk respons standar (sama di semua BE proyek ini).
type Envelope struct {
	Code    string          `json:"code"`
	Status  bool            `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Errors  json.RawMessage `json:"errors,omitempty"`
}

// Response: jawaban BE User yang berhasil (2xx). Body disimpan utuh agar bisa diteruskan.
type Response struct {
	StatusCode int
	Body       []byte
	Envelope   Envelope
}

// DecodeData membaca bagian data ke out.
func (r *Response) DecodeData(out any) error {
	if len(r.Envelope.Data) == 0 {
		return fmt.Errorf("%w: data kosong", ErrUnavailable)
	}
	if err := json.Unmarshal(r.Envelope.Data, out); err != nil {
		return fmt.Errorf("%w: data tidak dapat dibaca", ErrUnavailable)
	}
	return nil
}

// Client memanggil API BE User.
type Client struct {
	baseURL    string
	creds      CredentialsReader
	httpClient *http.Client
}

func NewClient(baseURL string, creds CredentialsReader) *Client {
	return &Client{
		baseURL:    baseURL,
		creds:      creds,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Configured: alamat BE User dan kunci aplikasi sudah diatur (belum menjamin kuncinya benar).
func (c *Client) Configured() (bool, error) {
	if c.baseURL == "" {
		return false, nil
	}
	creds, err := c.creds.GetAuthCredentials()
	if err != nil {
		return false, err
	}
	return creds != nil && creds.AppKey != "", nil
}

// Post mengirim body JSON ke path (mis. "/users/list"). actingUserID 0 = tanpa user
// (verifikasi login). Hasil:
//   - 2xx: *Response.
//   - 4xx: *errors.UpstreamError (penolakan dari BE User, diteruskan ke pengguna).
//   - kunci ditolak / belum diatur: ErrNotConfigured. 5xx / tidak terhubung: ErrUnavailable.
func (c *Client) Post(ctx context.Context, path string, actingUserID int, body []byte) (*Response, error) {
	if c.baseURL == "" {
		return nil, ErrNotConfigured
	}
	creds, err := c.creds.GetAuthCredentials()
	if err != nil {
		return nil, err
	}
	if creds == nil || creds.AppKey == "" {
		return nil, ErrNotConfigured
	}
	if len(body) == 0 {
		body = []byte("{}")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Key", creds.AppKey)
	if actingUserID > 0 {
		req.Header.Set("X-Acting-User-Id", strconv.Itoa(actingUserID))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("%w: status %d", ErrUnavailable, resp.StatusCode)
	}

	var env Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("%w: balasan tidak dapat dibaca (status %d)", ErrUnavailable, resp.StatusCode)
	}

	if resp.StatusCode >= 400 {
		if resp.StatusCode == http.StatusUnauthorized && env.Message == keyRejectedMessage {
			return nil, ErrNotConfigured
		}
		upstream := &errors.UpstreamError{HTTPCode: resp.StatusCode, Code: env.Code, Message: env.Message}
		if len(env.Errors) > 0 && string(env.Errors) != "null" {
			upstream.Errors = env.Errors
		}
		return nil, upstream
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%w: status %d", ErrUnavailable, resp.StatusCode)
	}
	return &Response{StatusCode: resp.StatusCode, Body: raw, Envelope: env}, nil
}

// ToServiceError mengubah kegagalan klien menjadi error yang aman ditampilkan ke pengguna.
// Penolakan dari BE User (UpstreamError) dan error lain dikembalikan apa adanya.
func ToServiceError(err error) error {
	switch {
	case err == nil:
		return nil
	case stderrors.Is(err, ErrNotConfigured):
		return &errors.ServiceUnavailableError{Message: "Layanan akun belum diatur di server ini"}
	case stderrors.Is(err, ErrUnavailable):
		return &errors.ServiceUnavailableError{Message: "Layanan akun tidak dapat dihubungi. Coba lagi beberapa saat"}
	default:
		return err
	}
}
