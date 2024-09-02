package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type ApiClient struct {
	Http    *http.Client
	Jwt     string
	BaseUrl string
	OrgId   string
}

type tokenResponse struct {
	AccessToken string
}

type loggingRoundTripper struct{
	Transport   http.RoundTripper
}

func (t *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	slog.Debug("-->", "method", req.Method, "url", req.URL)

	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	slog.Debug("<--", "status", resp.StatusCode, "url", resp.Request.URL, "t", time.Since(start))

	return resp, err
}

func NewApiClient(ctx context.Context, host string, apikey string, orgid string) (*ApiClient, error) {

	h := &http.Client{
		Transport: &loggingRoundTripper{
			Transport: http.DefaultTransport,
		},
	}

	baseUrl := fmt.Sprintf("%s/api/msc/v1/organizations/%s", host, orgid)

	url := fmt.Sprintf("%s/credentials/v2/token", baseUrl)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-API-KEY", apikey)
	resp, err := h.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("bad token request")
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var token tokenResponse
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	return &ApiClient{Http: h, Jwt: token.AccessToken, BaseUrl: baseUrl, OrgId: orgid}, nil
}
