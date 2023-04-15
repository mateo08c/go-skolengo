package inbox

type Attachment struct {
	ID        string `json:"id"`
	MessageID string `json:"message_id"`
	Name      string `json:"name"`
	Size      int    `json:"size"`
	Extension string `json:"extension"`
	Data      []byte `json:"data"`
}

func NewAttachment(messageID string, attachmentID string) *Attachment {
	return &Attachment{
		ID:        attachmentID,
		MessageID: messageID,
	}
}

func (a *Attachment) SetName(name string) {
	a.Name = name
}

func (a *Attachment) SetSize(size int) {
	a.Size = size
}

func (a *Attachment) SetExtension(extension string) {
	a.Extension = extension
}

func (a *Attachment) SetData(data []byte) {
	a.Data = data
}
