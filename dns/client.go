package dns

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/pkg/errors"
)

type Client struct {
	client   *http.Client
	username string
	password string
	url      string
}

func NewClient(username, password, url string) *Client {
	return &Client{
		client:   cleanhttp.DefaultClient(),
		username: username,
		password: password,
		url:      url,
	}
}

func (d Client) Set(r Request) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(r); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, d.url, &buf)
	if err != nil {
		return err
	}
	if d.username != "" && d.password != "" {
		req.SetBasicAuth(d.username, d.password)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.Errorf("Non 2xx http status: %d", resp.StatusCode)
	}

	return nil
}
