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
	"path/filepath"
	"strings"
	"time"

	"food-backend/internal/config"
)

type StorageService struct {
	publicKey   string
	privateKey  string
	urlEndpoint string
	httpClient  *http.Client
}

type imageKitUploadResponse struct {
	FileID   string `json:"fileId"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	FilePath string `json:"filePath"`
}

func NewStorageService(cfg config.StorageConfig) *StorageService {
	return &StorageService{
		publicKey:   strings.TrimSpace(cfg.ImageKitPublicKey),
		privateKey:  strings.TrimSpace(cfg.ImageKitPrivateKey),
		urlEndpoint: strings.TrimSpace(cfg.ImageKitURLEndpoint),

		httpClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

func (s *StorageService) UploadFile(
	ctx context.Context,
	file []byte,
	fileName string,
) (*imageKitUploadResponse, error) {
	return s.UploadFileWithMime(
		ctx,
		file,
		fileName,
		"video/mp4",
	)
}

func (s *StorageService) UploadFileWithMime(
	ctx context.Context,
	file []byte,
	fileName string,
	mimeType string,
) (*imageKitUploadResponse, error) {

	if s.privateKey == "" {
		return nil, fmt.Errorf("imagekit private key is not configured")
	}

	if len(file) == 0 {
		return nil, fmt.Errorf("upload file is empty")
	}

	fileName = filepath.Base(strings.TrimSpace(fileName))

	if fileName == "" || fileName == "." {
		return nil, fmt.Errorf("invalid upload file name")
	}

	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// ImageKit expects the multipart field to be named "file".
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("create imagekit file part: %w", err)
	}

	if _, err := part.Write(file); err != nil {
		return nil, fmt.Errorf("write imagekit file: %w", err)
	}

	// Required by ImageKit.
	if err := writer.WriteField("fileName", fileName); err != nil {
		return nil, fmt.Errorf("write imagekit fileName: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close imagekit multipart writer: %w", err)
	}

	const uploadURL = "https://upload.imagekit.io/api/v1/files/upload"

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		uploadURL,
		&body,
	)
	if err != nil {
		return nil, fmt.Errorf("create imagekit request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString(
		[]byte(s.privateKey + ":"),
	)

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"imagekit upload request failed: %w",
			err,
		)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"read imagekit response failed: %w",
			err,
		)
	}

	// IMPORTANT: ImageKit request ID.
	requestID := resp.Header.Get("x-ik-requestId")

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"imagekit upload failed: status=%d request_id=%s response=%s",
			resp.StatusCode,
			requestID,
			strings.TrimSpace(string(respBody)),
		)
	}

	var result imageKitUploadResponse

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf(
			"decode imagekit response failed: %w; response=%s",
			err,
			string(respBody),
		)
	}

	if result.URL == "" {
		return nil, fmt.Errorf(
			"imagekit upload succeeded but URL missing: %s",
			string(respBody),
		)
	}

	return &result, nil
}
