package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GophKeeper/internal/client/dto"
	"github.com/GophKeeper/internal/handlers/private_data"
)

func (c *ClientImpl) SavePrivateData(ctx context.Context, req *dto.SavePrivateDataRequest) error {
	c.lg.Infow("client save private data request", "req", req)
	url := c.cfg.GetString("gophkeeper_server.save_host")

	saveData := &private_data.SavePrivateDataRequest{
		Data:  req.Data,
		Type:  req.Type,
		Nonce: req.Nonce,
	}

	bytesData, errMarshal := json.Marshal(saveData)
	if errMarshal != nil {
		return fmt.Errorf("json.Marshal: %w", errMarshal)
	}

	reader := bytes.NewReader(bytesData)

	requestWithContext, err := http.NewRequestWithContext(ctx, "POST", url, reader)
	if err != nil {
		return fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	// Добавляем заголовок Authorization
	requestWithContext.Header.Set("Authorization", "Bearer "+req.JWT)

	res, err := c.httpClient.Do(requestWithContext)
	if err != nil {
		return fmt.Errorf("client.Do: %w", err)
	}
	defer res.Body.Close()

	fmt.Println("response :", res.Body)

	return nil
}
