package xtream_codes_go

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"log/slog"
)

type contextKey string

const (
	defaultPlayerApi string = "player_api.php"
)

type credentials struct {
	host     *url.URL
	username string
	password string
}

type options struct {
	logger  *slog.Logger
	client  *http.Client
	dumper  io.Writer
	apiPath string
}

type Option func(*options) error

func WithLogger(logger *slog.Logger) Option {
	return func(o *options) error {
		o.logger = logger
		return nil
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(o *options) error {
		if client == nil {
			return fmt.Errorf("http client cannot be nil")
		}
		o.client = client
		return nil
	}
}

func WithDumper(dumper io.Writer) Option {
	return func(o *options) error {
		o.dumper = dumper
		return nil
	}
}

func WithAPIPath(apiPath string) Option {
	return func(o *options) error {
		if apiPath == "" {
			return fmt.Errorf("apiPath cannot be empty")
		}
		o.apiPath = apiPath
		return nil
	}
}

func NewApiClient(host, username, password string, opts ...Option) (*ApiClient, error) {
	uri, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	o := options{
		logger:  slog.Default(),
		client:  &http.Client{},
		apiPath: defaultPlayerApi,
	}

	for _, opt := range opts {
		if err := opt(&o); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
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

	var api = &ApiClient{
		client:  o.client,
		config:  creds,
		apiPath: o.apiPath,
		logger:  o.logger,
	}

	o.client.Transport = &ApiTransport{
		inner:  transport,
		logger: o.logger,
		dumper: o.dumper,
		api:    api,
	}

	return api, nil
}

type ApiClient struct {
	client    *http.Client
	mu        sync.RWMutex
	loginInfo *LoginInfo
	config    *credentials
	apiPath   string
	logger    *slog.Logger
}

func (a *ApiClient) getLoginInfo() *LoginInfo {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.loginInfo
}

func (a *ApiClient) setLoginInfo(info *LoginInfo) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.loginInfo = info
}

func (a *ApiClient) fetch(ctx context.Context, action string, params map[string]string, path string, data any) error {

	request, err := http.NewRequestWithContext(ctx, "GET", path, nil)

	if err != nil {
		return err
	}

	query := request.URL.Query()
	if action != "" {
		query.Set("action", action)
	}
	for k, v := range params {
		query.Set(k, v)
	}

	loginInfo := a.getLoginInfo()
	if loginInfo == nil {
		query.Set("username", a.config.username)
		query.Set("password", a.config.password)
		request.URL.Scheme = a.config.host.Scheme
		request.URL.Host = a.config.host.Host
		request.URL.Path = a.config.host.Path + "/" + path
	} else {
		query.Set("username", loginInfo.UserInfo.Username)
		query.Set("password", loginInfo.UserInfo.Password)
		request.URL.Scheme = loginInfo.ServerInfo.ServerProtocol
		request.URL.Host = loginInfo.ServerInfo.Url
		request.URL.Path = "/" + path
	}

	request.URL.RawQuery = query.Encode()

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
	info := a.getLoginInfo()
	if info == nil {
		return ""
	}

	ext := ""
	if extension != "" {
		ext = "." + extension
	}

	return fmt.Sprintf(
		"%s://%s/%s/%s/%s/%d%s",
		info.ServerInfo.ServerProtocol,
		info.ServerInfo.Url,
		stream,
		info.UserInfo.Username,
		info.UserInfo.Password,
		id,
		ext,
	)
}
