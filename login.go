package xtream_codes_go

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type UserInfo struct {
	Username             string   `json:"username"`
	Password             string   `json:"password"`
	Message              string   `json:"message"`
	Auth                 boolean  `json:"auth"`
	Status               string   `json:"status"`
	ExpDate              int      `json:"exp_date,string"`
	IsTrial              boolean  `json:"is_trial"`
	ActiveCons           int      `json:"active_cons,string"`
	CreatedAt            int      `json:"created_at,string"`
	MaxConnections       int      `json:"max_connections,string"`
	AllowedOutputFormats []string `json:"allowed_output_formats"`
}

type ServerInfo struct {
	Url            string  `json:"url"`
	Port           int     `json:"port,string"`
	HttpsPort      int     `json:"https_port,string"`
	ServerProtocol string  `json:"server_protocol"`
	RtmpPort       int     `json:"rtmp_port,string"`
	Timezone       string  `json:"timezone"`
	TimestampNow   int     `json:"timestamp_now"`
	TimeNow        string  `json:"time_now"`
	Process        boolean `json:"process"`
}

type LoginInfo struct {
	UserInfo   *UserInfo   `json:"user_info"`
	ServerInfo *ServerInfo `json:"server_info"`
}

func (a *ApiClient) Login(ctx context.Context) (*LoginInfo, error) {

	var info *LoginInfo

	if err := a.fetch(ctx, playerApi, &info); err != nil {
		return nil, err
	}

	return info, nil
}

func authenticate(ctx context.Context, client *ApiClient, logger *slog.Logger) error {
	info, err := client.Login(ctx)

	if err != nil {
		return &Error{message: "failed to login", prev: err}
	}

	if nil != logger {

		expires := time.Unix(int64(info.UserInfo.ExpDate), 0)

		if loc, err := time.LoadLocation(info.ServerInfo.Timezone); err == nil {
			expires.In(loc)
		}

		logger.Debug(info.UserInfo.Message)
		logger.Debug(fmt.Sprintf("account (status: %s expires: %s)", strings.ToLower(info.UserInfo.Status), expires))
	}

	if info.UserInfo.Status != "Active" {
		return &Error{message: "user account not active"}
	}

	client.loginInfo = info

	return nil

}
