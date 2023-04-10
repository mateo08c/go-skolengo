package calendar

import (
	"encoding/json"
	"strconv"
	"time"
)

type Period struct {
	ID        string    `json:"id"`
	Subject   string    `json:"subject"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	Teacher   string    `json:"teacher"`
	Room      string    `json:"room"`
	Cancelled bool      `json:"cancelled"`
}

func (p *Period) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p *Period) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Period) IsOngoing() bool {
	return p.Start.Before(time.Now().UTC()) && p.End.After(time.Now().UTC())
}

func (p *Period) IsUpcoming() bool {
	return p.Start.After(time.Now().UTC())
}

func (p *Period) IsFinished() bool {
	return p.End.Before(time.Now().UTC())
}

func (p *Period) TimePassed() time.Duration {
	return time.Now().UTC().Sub(p.Start)
}

func (p *Period) TimeLeft() time.Duration {
	return p.End.Sub(time.Now().UTC())
}

func (p *Period) InProgress() bool {
	return p.Start.Before(time.Now().UTC()) && p.End.After(time.Now().UTC())
}

func (p *Period) StartIn() bool {
	return p.Start.After(time.Now().UTC())
}

func TimestampToTime(unix string) time.Time {
	t, _ := strconv.ParseInt(unix, 10, 64)
	tm := time.Unix(t/1000, 0)
	return tm
}
