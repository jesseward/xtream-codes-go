package xtream_codes_go

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"log/slog"
)

type contextKey string

const valuesKey contextKey = "values"
const loginInfoKey contextKey = "loginInfo"

const (
	playerApi string = "player_api.php"
)

type credentials struct {
	host     *url.URL
	username string
	password string
}

type options struct {
	logger *slog.Logger
	client *http.Client
	dumper io.Writer
}

type Option func(*options)

func WithLogger(logger *slog.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(o *options) {
		o.client = client
	}
}

func WithDumper(dumper io.Writer) Option {
	return func(o *options) {
		o.dumper = dumper
	}
}

func NewApiClient(host, username, password string, opts ...Option) (*ApiClient, error) {
	uri, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	o := options{
		logger: slog.Default(),
		client: &http.Client{},
	}

	for _, opt := range opts {
		opt(&o)
	}

	var transport = o.client.Transport

	if transport == nil {
		transport = http.DefaultTransport
	}

	creds := &credentials{
		host:     uri,
		username: username,
		password: password,
	}

	o.client.Transport = &ApiTransport{
		inner:  transport,
		logger: o.logger,
		dumper: o.dumper,
		config: creds,
	}

	var api = &ApiClient{client: o.client}

	if err := authenticate(context.Background(), api, o.logger); err != nil {
		return nil, err
	}

	return api, nil
}

type ApiClient struct {
	client    *http.Client
	loginInfo *LoginInfo
}

func (a *ApiClient) context(ctx context.Context, action string, params map[string]string) context.Context {
	var values = make(url.Values)

	values.Set("action", action)

	for k, v := range params {
		values.Set(k, v)
	}

	return context.WithValue(ctx, valuesKey, values)
}

func (a *ApiClient) fetch(ctx context.Context, path string, data any) error {

	if a.loginInfo != nil {
		ctx = context.WithValue(ctx, loginInfoKey, a.loginInfo)
	}

	request, err := http.NewRequestWithContext(ctx, "GET", path, nil)

	if err != nil {
		return err
	}

	response, err := a.client.Do(request)

	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode >= 300 || response.StatusCode < 200 {
		return fmt.Errorf("unexpected status code (%d) returned for '%s'", response.StatusCode, response.Request.URL.RequestURI())
	}

	if err := json.NewDecoder(response.Body).Decode(data); err != nil {
		return err
	}

	return nil
}

func (a *ApiClient) streamUrl(stream string, id int, extension string) string {
	return fmt.Sprintf(
		"%s://%s/%s/%s/%s/%d.%s",
		a.loginInfo.ServerInfo.ServerProtocol,
		a.loginInfo.ServerInfo.Url,
		stream,
		a.loginInfo.UserInfo.Username,
		a.loginInfo.UserInfo.Password,
		id,
		extension,
	)
}
