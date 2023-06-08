package inbox

import (
	"time"
)

type RedactorInfo struct {
	Code        string `json:"code"`
	Libelle     string `json:"libelle"`
	LibelleLong string `json:"libelleLong"`
}

type MessageRedactor struct {
	Id   string          `json:"id"`
	Name string          `json:"name"`
	Type string          `json:"type"`
	Info []*RedactorInfo `json:"info"`
}

type MessageRecipient struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type MessageContent struct {
	Subject         string              `json:"subject"`
	Body            string              `json:"body"`
	ParticipationID string              `json:"participation_id"`
	Redactor        *MessageRedactor    `json:"sender"`
	Attachments     []*Attachment       `json:"attachments"`
	Recipients      []*MessageRecipient `json:"recipient"`
	Groups          []string            `json:"groups"`
	Date            *time.Time          `json:"date"`
}

type Message struct {
	ID         string          `json:"id"`
	FolderID   string          `json:"folder_id"`
	ServiceURL string          `json:"service_url"`
	Content    *MessageContent `json:"content"`
	Type       string          `json:"type"`
}

func (r *MessageRecipient) SetFirstName(firstName string) {
	r.FirstName = firstName
}

func (r *MessageRecipient) SetLastName(lastName string) {
	r.LastName = lastName
}
