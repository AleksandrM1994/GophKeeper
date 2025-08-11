package client

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/GophKeeper/internal/client/dto"
	"github.com/GophKeeper/internal/repository"
)

type fakeRoundTripper struct {
	resp *http.Response
	err  error
}

func (f *fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f.resp, f.err
}

func TestSavePrivateData_Success(t *testing.T) {
	// Prepare fake HTTP response
	body := ioutil.NopCloser(bytes.NewBufferString(`{"status":"ok"}`))
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       body,
	}
	// Inject fakeRoundTripper into http.Client
	fakeClient := &http.Client{
		Transport: &fakeRoundTripper{resp: resp, err: nil},
	}

	cfg := viper.New()
	cfg.Set("gophkeeper_server.save_host", "http://example.com/save")
	lg, _ := zap.NewDevelopment()
	sugar := lg.Sugar()

	c := NewClient(sugar, cfg, fakeClient)

	req := &dto.SavePrivateDataRequest{
		Data:  []byte("hello"),
		Type:  repository.PrivateDataTypeText,
		Nonce: []byte("nonce"),
		Login: "user1",
		JWT:   "token123",
	}

	err := c.SavePrivateData(context.Background(), req)
	require.NoError(t, err)

	// Verify that the request was built correctly
	// (could be inspected via fakeRoundTripper if it recorded req)
}

func TestSavePrivateData_DoError(t *testing.T) {
	fakeRT := &fakeRoundTripper{resp: nil, err: errors.New("network failure")}
	fakeClient := &http.Client{Transport: fakeRT}

	cfg := viper.New()
	cfg.Set("gophkeeper_server.save_host", "http://example.com/save")
	lg, _ := zap.NewDevelopment()
	sugar := lg.Sugar()

	c := NewClient(sugar, cfg, fakeClient)

	req := &dto.SavePrivateDataRequest{
		Data:  []byte("x"),
		Type:  repository.PrivateDataTypeText,
		Nonce: []byte("n"),
		Login: "u",
		JWT:   "jwt",
	}

	err := c.SavePrivateData(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client.Do")
}
