package telegram

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
)

type Client struct {
	token   string
	baseURL string
}

func NewClient(token string) *Client {
	baseURL := fmt.Sprintf("https://api.telegram.org/bot%s", token)
	log.Println("Constructed Telegram API base URL:", baseURL)
	return &Client{
		token:   token,
		baseURL: baseURL,
	}
}

func (c *Client) SendMessage(chatID, message string) error {
	apiURL := fmt.Sprintf("%s/sendMessage", c.baseURL)
	log.Println("Sending to Telegram API:", apiURL, "with chat_id:", chatID)
	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id": {chatID},
		"text":    {message},
	})
	if err != nil {
		log.Println("HTTP request failed:", err)
		return err
	}
	defer resp.Body.Close()
	log.Println("Telegram API response status:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api returned status %s", resp.Status)
	}
	return nil
}

func (c *Client) SendPhoto(chatID string, photoData []byte, filename, caption string) error {
	apiURL := fmt.Sprintf("%s/sendPhoto", c.baseURL)
	return c.sendFile(apiURL, chatID, "photo", photoData, filename, caption)
}

func (c *Client) SendDocument(chatID string, docData []byte, filename, caption string) error {
	apiURL := fmt.Sprintf("%s/sendDocument", c.baseURL)
	return c.sendFile(apiURL, chatID, "document", docData, filename, caption)
}

func (c *Client) SendAudio(chatID string, audioData []byte, filename, caption string) error {
	apiURL := fmt.Sprintf("%s/sendAudio", c.baseURL)
	return c.sendFile(apiURL, chatID, "audio", audioData, filename, caption)
}

func (c *Client) sendFile(apiURL, chatID, fieldName string, fileData []byte, filename, caption string) error {
	log.Println("Sending file to Telegram API:", apiURL, "with chat_id:", chatID, "filename:", filename, "caption:", caption)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, filepath.Base(filename))
	if err != nil {
		log.Println("Failed to create form file:", err)
		return fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = part.Write(fileData)
	if err != nil {
		log.Println("Failed to write file data:", err)
		return fmt.Errorf("failed to write file data: %v", err)
	}

	err = writer.WriteField("chat_id", chatID)
	if err != nil {
		log.Println("Failed to write chat_id:", err)
		return fmt.Errorf("failed to write chat_id: %v", err)
	}

	if caption != "" {
		err = writer.WriteField("caption", caption)
		if err != nil {
			log.Println("Failed to write caption:", err)
			return fmt.Errorf("failed to write caption: %v", err)
		}
	}

	err = writer.Close()
	if err != nil {
		log.Println("Failed to close writer:", err)
		return fmt.Errorf("failed to close writer: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		log.Println("Failed to create request:", err)
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("HTTP request failed:", err)
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	log.Println("Telegram API response status:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api returned status %s", resp.Status)
	}
	return nil
}
