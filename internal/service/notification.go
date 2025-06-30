package service

import "notify_hub/internal/telegram"

type NotificationService struct {
	tgClient *telegram.Client
}

func NewNotificationService(tgClient *telegram.Client) *NotificationService {
	return &NotificationService{tgClient: tgClient}
}

func (s *NotificationService) SendToTelegram(userID, message string) error {
	return s.tgClient.SendMessage(userID, message)
}

func (s *NotificationService) SendPhotoToTelegram(userID, content string, photoData []byte, filename string) error {
	return s.tgClient.SendPhoto(userID, photoData, filename, content)
}

func (s *NotificationService) SendDocumentToTelegram(userID, content string, docData []byte, filename string) error {
	return s.tgClient.SendDocument(userID, docData, filename, content)
}

func (s *NotificationService) SendAudioToTelegram(userID, content string, audioData []byte, filename string) error {
	return s.tgClient.SendAudio(userID, audioData, filename, content)
}
