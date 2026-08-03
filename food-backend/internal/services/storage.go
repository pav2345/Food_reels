package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"food-backend/internal/config"
)

type StorageService struct {
	publicKey   string
	privateKey  string
	urlEndpoint string
	httpClient  *http.Client
}

type imageKitUploadResponse struct {
	URL string `json:"url"`
}

func NewStorageService(cfg config.StorageConfig) *StorageService {
	return &StorageService{
		publicKey:   cfg.ImageKitPublicKey,
		privateKey:  cfg.ImageKitPrivateKey,
		urlEndpoint: cfg.ImageKitURLEndpoint,
		httpClient:  &http.Client{},
	}
}

func (s *StorageService) UploadFile(ctx context.Context, file []byte, fileName string) (*imageKitUploadResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("create form file: %w", err)
	}

	if _, err := part.Write(file); err != nil {
		return nil, fmt.Errorf("write file: %w", err)
	}

	if err := writer.WriteField("fileName", fileName); err != nil {
		return nil, fmt.Errorf("write fileName field: %w", err)
	}

	if err := writer.WriteField("mimeType", "video/mp4"); err != nil {
		return nil, fmt.Errorf("write mimeType field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://upload.imagekit.io/api/v1/files/upload", body)
	if err != nil {
		return nil, fmt.Errorf("create upload request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(s.privateKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upload request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read upload response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("imagekit upload failed: status %d body %s", resp.StatusCode, string(respBody))
	}

	var result imageKitUploadResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decode upload response: %w", err)
	}

	return &result, nil
}
