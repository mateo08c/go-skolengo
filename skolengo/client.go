package skolengo

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/kataras/golog"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	Username string
	Password string

	cookies []*http.Cookie
}

// GetServices return the list of services.
// This method is used to log in
func (c *Client) GetServices() ([]*Service, error) {
	start := time.Now()
	var services []*Service

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	//---------------- 1st request ----------------
	req, err := http.NewRequest("GET", "https://cas.monbureaunumerique.fr/delegate/redirect/EDU", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	f := doc.Find("input[name=SAMLRequest]")
	samlRequest := f.AttrOr("value", "")

	//---------------- 2nd request ----------------
	req, err = http.NewRequest("POST", "https://educonnect.education.gouv.fr/idp/profile/SAML2/POST/SSO", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("RelayState", "https://cas.monbureaunumerique.fr/saml/SAMLAssertionConsumer")
	q.Add("SAMLRequest", samlRequest)
	req.URL.RawQuery = q.Encode()

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	//---------------- 3rd request ----------------
	req, err = http.NewRequest("POST", resp.Header.Get("Location"), nil)
	if err != nil {
		return nil, err
	}
	for _, cookie := range resp.Cookies() {
		req.AddCookie(cookie)
	}

	q = req.URL.Query()
	q.Add("j_username", c.Username)
	q.Add("j_password", c.Password)
	q.Add("_eventId_proceed", "")
	req.URL.RawQuery = q.Encode()

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	f = doc.Find("input[name=SAMLResponse]")
	samlResponse := f.AttrOr("value", "")

	if samlResponse == "" {
		return nil, errors.New("invalid credentials")
	}

	//---------------- 4th request ----------------
	//This query return the TGC cookie is used to log in to the services
	resp, err = http.PostForm("https://cas.monbureaunumerique.fr/saml/SAMLAssertionConsumer", url.Values{
		"RelayState": {
			"https://cas.monbureaunumerique.fr/saml/SAMLAssertionConsumer",
		},
		"SAMLResponse": {
			samlResponse,
		},
	})
	if err != nil {
		return nil, err
	}
	c.SetCookies(resp.Cookies()) //Set the TGC cookie

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	//Parse all the services
	doc.Find(".msg .p-like a").Each(func(i int, s *goquery.Selection) {
		href, exist := s.Attr("href")
		u, err := url.Parse(href)
		if err != nil {
			return
		}

		name := s.Text()

		if !exist {
			return
		}

		req, err := http.NewRequest("GET", href, nil)
		if err != nil {
			return
		}

		for _, cookie := range resp.Cookies() {
			req.AddCookie(cookie)
		}

		resp, err = client.Do(req)
		if err != nil {
			return
		}

		ticketURL := resp.Header.Get("Location")
		if ticketURL == "" {
			return
		}

		req, err = http.NewRequest("GET", ticketURL, nil)
		if err != nil {
			return
		}

		resp, err = client.Do(req)
		if err != nil {
			return
		}

		//get service in u query
		q := u.Query()
		service := q.Get("service")
		if service == "" {
			return
		}

		b, err := url.Parse(service)
		if err != nil {
			return
		}

		services = append(services, NewService(name, resp.Cookies(), b))
	})

	golog.Debugf("GetServices took %s", time.Since(start))

	return services, nil
}

func (c *Client) Post(u *url.URL) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, err
	}

	if len(c.cookies) != 0 {
		req.Header.Set("cookie", c.GetCookiesString())
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	golog.Debugf("POST %s", u.String())

	return resp, nil
}

func (c *Client) SetUsername(u string) {
	c.Username = u
}

func (c *Client) SetPassword(s string) {
	c.Password = s
}

func (c *Client) SetCookies(cs []*http.Cookie) {
	c.cookies = cs
}

func (c *Client) GetCookies() []*http.Cookie {
	return c.cookies
}

func (c *Client) GetCookiesString() string {
	var s string
	for _, c := range c.cookies {
		s += fmt.Sprintf("%s=%s;", c.Name, c.Value)
	}
	return s
}
