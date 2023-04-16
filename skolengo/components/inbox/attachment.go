package inbox

import "os"

type Attachment struct {
	ID        int    `json:"id"`
	MessageID string `json:"message_id"`
	Name      string `json:"name"`
	Size      int    `json:"size"`
	Extension string `json:"extension"`
	Data      []byte `json:"data"`
}

func (a Attachment) SaveToFile(s string) error {
	create, err := os.Create(s)
	if err != nil {
		return nil
	}

	_, err = create.Write(a.Data)
	if err != nil {
		return nil
	}

	return nil
}
