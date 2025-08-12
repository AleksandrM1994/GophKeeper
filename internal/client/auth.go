package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	api "github.com/GophKeeper/pkg/api"
)

func (c *ClientImpl) AuthUser(ctx context.Context, req *api.AuthUserRequest) (*api.AuthUserResponse, error) {
	url := c.cfg.GetString("gophkeeper_server.auth_host")

	bytesData, errMarshal := json.Marshal(req)
	if errMarshal != nil {
		return nil, fmt.Errorf("json.Marshal: %w", errMarshal)
	}

	reader := bytes.NewReader(bytesData)

	requestWithContext, err := http.NewRequestWithContext(ctx, "POST", url, reader)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	res, err := c.httpClient.Do(requestWithContext)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %w", err)
	}
	defer res.Body.Close()

	resBytes, errReadAll := io.ReadAll(res.Body)
	if errReadAll != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", errReadAll)
	}

	var authUserResponse *api.AuthUserResponse

	errUnmarshal := json.Unmarshal(resBytes, &authUserResponse)
	if errUnmarshal != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", errUnmarshal)
	}

	fmt.Println("response :", authUserResponse)

	return authUserResponse, nil
}
