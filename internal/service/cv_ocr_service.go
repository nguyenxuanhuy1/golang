package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

type CVOCRService struct {
	httpClient *http.Client
}

func NewCVOCRService() *CVOCRService {
	return &CVOCRService{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *CVOCRService) ReadCVImage(image []byte) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", errors.New("GEMINI_API_KEY not set")
	}

	imgBase64 := base64.StdEncoding.EncodeToString(image)

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": "Read the CV image carefully and return ALL readable text in plain text. Do not summarize. Keep line breaks.",
					},
					{
						"inline_data": map[string]interface{}{
							"mime_type": "image/jpeg",
							"data":      imgBase64,
						},
					},
				},
			},
		},
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(
		"POST",
		"https://generativelanguage.googleapis.com/v1/models/gemini-1.5-flash:generateContent?key="+apiKey,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	if len(result.Candidates) == 0 ||
		len(result.Candidates[0].Content.Parts) == 0 {
		return "", nil
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
