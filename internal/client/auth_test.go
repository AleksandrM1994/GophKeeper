package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	api "github.com/GophKeeper/pkg/api"
)

func TestAuthUser_Success(t *testing.T) {
	// Ожидаемый ответ
	expected := &api.AuthUserResponse{Jwt: "token_xyz"}
	respBytes, err := json.Marshal(expected)
	require.NoError(t, err)

	fakeResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBuffer(respBytes)),
	}
	fakeClient := &http.Client{
		Transport: &fakeRoundTripper{resp: fakeResp, err: nil},
	}

	cfg := viper.New()
	cfg.Set("gophkeeper_server.auth_host", "http://example.com/auth")
	lg, _ := zap.NewDevelopment()
	sugar := lg.Sugar()

	c := NewClient(sugar, cfg, fakeClient)

	req := &api.AuthUserRequest{Login: "user1", Password: "pass1"}
	got, err := c.AuthUser(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, expected.Jwt, got.Jwt)
}

func TestAuthUser_DoError(t *testing.T) {
	fakeClient := &http.Client{
		Transport: &fakeRoundTripper{resp: nil, err: errors.New("network fail")},
	}

	cfg := viper.New()
	cfg.Set("gophkeeper_server.auth_host", "http://example.com/auth")
	lg, _ := zap.NewDevelopment()
	sugar := lg.Sugar()

	c := NewClient(sugar, cfg, fakeClient)

	_, err := c.AuthUser(context.Background(), &api.AuthUserRequest{Login: "", Password: ""})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client.Do")
}

func TestAuthUser_UnmarshalError(t *testing.T) {
	fakeResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{not: "json"]`)),
	}
	fakeClient := &http.Client{
		Transport: &fakeRoundTripper{resp: fakeResp, err: nil},
	}

	cfg := viper.New()
	cfg.Set("gophkeeper_server.auth_host", "http://example.com/auth")
	lg, _ := zap.NewDevelopment()
	sugar := lg.Sugar()

	c := NewClient(sugar, cfg, fakeClient)

	_, err := c.AuthUser(context.Background(), &api.AuthUserRequest{Login: "u", Password: "p"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "json.Unmarshal")
}
