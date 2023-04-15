package skolengo

import (
	"net/url"
)

type UrlBuilder struct {
	URL *url.URL
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
