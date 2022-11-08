package xrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	Client *http.Client
	Auth   *AuthInfo
	Host   string
}

func (c *Client) getClient() *http.Client {
	if c.Client == nil {
		return http.DefaultClient
	}
	return c.Client
}

type XRPCRequestType int

type AuthInfo struct {
	AccessJwt  string `json:"accessJwt"`
	RefreshJwt string `json:"refreshJwt"`
	Handle     string `json:"handle"`
	Did        string `json:"did"`
}

const (
	Query = XRPCRequestType(iota)
	Procedure
)

func makeParams(p map[string]interface{}) string {
	var parts []string
	for k, v := range p {
		parts = append(parts, fmt.Sprintf("%s=%s", k, url.QueryEscape(fmt.Sprint(v))))
	}

	return strings.Join(parts, "&")
}

func (c *Client) Do(ctx context.Context, kind XRPCRequestType, method string, params map[string]interface{}, bodyobj interface{}, out interface{}) error {
	var body io.Reader
	if bodyobj != nil {
		b, err := json.Marshal(bodyobj)
		if err != nil {
			return err
		}

		body = bytes.NewReader(b)
	}

	var m string
	switch kind {
	case Query:
		m = "GET"
	case Procedure:
		m = "POST"
	default:
		return fmt.Errorf("unsupported request kind: %d", kind)
	}

	var paramStr string
	if len(params) > 0 {
		paramStr = "?" + makeParams(params)
	}

	req, err := http.NewRequest(m, c.Host+"/xrpc/"+method+paramStr, body)
	if err != nil {
		return err
	}

	if bodyobj != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.Auth != nil {
		req.Header.Set("Authorization", "Bearer "+c.Auth.AccessJwt)
	}

	resp, err := c.getClient().Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		var i interface{}
		_ = json.NewDecoder(resp.Body).Decode(&i)
		fmt.Println("debug body response: ", i)
		return fmt.Errorf("XRPC ERROR %d: %s", resp.StatusCode, resp.Status)
	}

	if out != nil {
		if buf, ok := out.(*bytes.Buffer); ok {
			io.Copy(buf, resp.Body)
		} else {
			if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
				return fmt.Errorf("decoding xrpc response: %w", err)
			}
		}
	}

	return nil
}
