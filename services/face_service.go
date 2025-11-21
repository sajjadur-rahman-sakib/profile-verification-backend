package services

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"main.go/config"
)

type FaceService struct{}

func NewFaceService() *FaceService {
	return &FaceService{}
}

func (s *FaceService) CompareFaces(documentPath, selfiePath string) (bool, error) {
	configuration := config.GetConfig()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file1, err := os.Open(documentPath)
	if err != nil {
		return false, err
	}
	defer file1.Close()
	part1, _ := writer.CreateFormFile("image1", "document.jpg")
	io.Copy(part1, file1)

	file2, err := os.Open(selfiePath)
	if err != nil {
		return false, err
	}
	defer file2.Close()
	part2, _ := writer.CreateFormFile("image2", "selfie.jpg")
	io.Copy(part2, file2)

	writer.Close()

	req, err := http.NewRequest("POST", configuration.FaceService, body)
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result struct {
		IsMatch bool `json:"is_match"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result.IsMatch, nil
}
