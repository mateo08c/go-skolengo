package inbox

import (
	"encoding/json"
	"time"
)

type MessageType string

const (
	MessageTypeInternal      MessageType = "internal"
	MessageTypeExternal      MessageType = "external"
	MessageTypeInstitutional MessageType = "institutional"
)

type Sender struct {
	Name string `json:"name"`
}

type Message struct {
	ID          string        `json:"id"`
	FolderID    string        `json:"folder_id"`
	Subject     string        `json:"subject"`
	Sender      *Sender       `json:"sender"`
	Recipient   []string      `json:"recipient"`
	Groups      []string      `json:"groups"`
	Date        *time.Time    `json:"date"`
	Content     string        `json:"content"`
	Type        string        `json:"type"`
	ServiceURL  string        `json:"service_url"`
	Attachments []*Attachment `json:"-"`
}

func (m *Message) SetSubject(subject string) {
	m.Subject = subject
}

func (m *Message) SetSender(sender *Sender) {
	m.Sender = sender
}

func (m *Message) SetGroups(groups []string) {
	m.Groups = groups
}

func (m *Message) AddGroup(group string) {
	m.Groups = append(m.Groups, group)
}

func (m *Message) SetRecipient(recipient []string) {
	m.Recipient = recipient
}

func (m *Message) SetDate(date *time.Time) {
	m.Date = date
}

func (m *Message) SetContent(content string) {
	m.Content = content
}

func (m *Message) AddContent(content string) {
	m.Content += content
}

func (m *Message) SetType(messageType string) {
	m.Type = messageType
}

func (m *Message) AddAttachment(attachment *Attachment) {
	m.Attachments = append(m.Attachments, attachment)
}

func (m *Message) SetServiceURL(service string) {
	m.ServiceURL = service
}

func (m *Message) SetID(get string) {
	m.ID = get
}

func (m *Message) SetFolderID(get string) {
	m.FolderID = get
}

func (m *Message) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}
