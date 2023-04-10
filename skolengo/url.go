package skolengo

import (
	"net/url"
)

type UrlBuilder struct {
	URL *url.URL
}

func (s *Service) MessageContentURL(folderID string, messageID string) (*url.URL, error) {
	builder := NewURLBuilder(s.URL)
	builder.SetPath("sg.do")
	builder.AddParam("PROC", "MESSAGERIE")
	builder.AddParam("ACTION", "CONSULTER_COMMUNICATION")
	builder.AddParam("ID_DOSSIER", folderID)
	builder.AddParam("ID_COMMUNICATION", messageID)

	return builder.Build()
}

func (s *Service) MessageRecipientURL(messageID string) (*url.URL, error) {
	builder := NewURLBuilder(s.URL)
	builder.SetPath("sg.do")
	builder.AddParam("PROC", "MESSAGERIE")
	builder.AddParam("ACTION", "LISTER_DESTINATAIRES_GROUPE")
	builder.AddParam("ID_COMMUNICATION", messageID)
	return builder.Build()
}

func (s *Service) PeriodUrl() (*url.URL, error) {
	builder := NewURLBuilder(s.URL)
	builder.SetPath("sg.do")
	builder.AddParam("PROC", "CDT_AFFICHAGE")
	builder.AddParam("VUE", "E")
	return builder.Build()
}

func (s *Service) MessagerieURL() (*url.URL, error) {
	builder := NewURLBuilder(s.URL)
	builder.SetPath("sg.do")
	builder.AddParam("PROC", "MESSAGERIE")
	return builder.Build()
}

func NewURLBuilder(u *url.URL) *UrlBuilder {
	clone := *u
	return &UrlBuilder{
		URL: &clone,
	}
}

func (s *UrlBuilder) Build() (*url.URL, error) {
	s.URL.RawQuery = s.URL.Query().Encode()
	return s.URL, nil
}

func (s *UrlBuilder) SetURL(u *url.URL) *UrlBuilder {
	s.URL = u
	return s
}

func (s *UrlBuilder) SetPath(path string) *UrlBuilder {
	s.URL.Path = path
	return s
}

func (s *UrlBuilder) SetParams(params *url.Values) *UrlBuilder {
	s.URL.RawQuery = params.Encode()
	return s
}

func (s *UrlBuilder) AddParam(key string, value string) *UrlBuilder {
	q := s.URL.Query()
	q.Add(key, value)
	s.URL.RawQuery = q.Encode()
	return s
}

func (s *UrlBuilder) AddParams(query url.Values) {
	for key, values := range query {
		for _, value := range values {
			s.AddParam(key, value)
		}
	}
}
