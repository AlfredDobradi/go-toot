package mastodon

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alfreddobradi/go-toot/config"
	"github.com/google/uuid"
)

type Token struct {
	Token     string
	Type      string
	Scope     string
	Timestamp time.Time
}

func (t *Token) UnmarshalJSON(b []byte) error {
	raw := map[string]interface{}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	t.Token = raw["access_token"].(string)
	t.Type = raw["token_type"].(string)
	t.Scope = raw["scope"].(string)
	ts := raw["created_at"].(float64)
	t.Timestamp = time.Unix(int64(ts), 0)

	return nil
}

func GetCode() (string, error) {
	uri := fmt.Sprintf("%s/oauth/authorize", config.InstanceURL())
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Add("client_id", config.ClientID())
	q.Add("redirect_uri", config.RedirectURI())
	q.Add("response_type", "code")
	q.Add("scope", config.Scope())
	uri = fmt.Sprintf("%s?%s", u.String(), q.Encode())

	return uri, nil
}

func GetToken(code string) (*Token, error) {
	uri := fmt.Sprintf("%s/oauth/token", config.InstanceURL())
	resp, err := http.PostForm(uri, url.Values{
		"client_id":     {config.ClientID()},
		"client_secret": {config.ClientSecret()},
		"redirect_uri":  {config.RedirectURI()},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"scope":         {config.Scope()},
	})

	if err != nil {
		return nil, err
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Non-OK error code (%d) received: %s", resp.StatusCode, string(raw))
	}

	var t Token
	if err := json.Unmarshal(raw, &t); err != nil {
		return nil, err
	}

	return &t, nil
}

func Post(token string, visibility string, msg string) error {
	uri := fmt.Sprintf("%s/api/v1/statuses", config.InstanceURL())

	values := url.Values{
		"status":     {msg},
		"visibility": {visibility},
	}

	body := strings.NewReader(values.Encode())

	req, err := http.NewRequest(http.MethodPost, uri, body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Idempotency-Key", uuid.New().String())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Go-Toot v0.1.0")

	tr := &http.Transport{}
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	cl := http.DefaultClient
	cl.Transport = tr

	res, err := cl.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		raw, _ := io.ReadAll(res.Body)
		return fmt.Errorf("Non-OK status code (%d) received: %s", res.StatusCode, string(raw))
	}

	return nil
}
