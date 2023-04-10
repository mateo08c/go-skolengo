package skolengo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/kataras/golog"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	MonBureauNumeriqueLoginURL         = "https://cas.monbureaunumerique.fr/login"
	MonBureauNumeriqueHomeURL          = "https://www.monbureaunumerique.fr"
	MonBureauNumeriqueLoginURlRedirect = "https://cas.monbureaunumerique.fr/login?service=https%3A%2F%2Flyc-monge.monbureaunumerique.fr%2Fsg.do%3FPROC%3DMESSAGERIE"
	SeleniumTypeLocal                  = iota
	SeleniumTypeRemote                 = iota
)

type SeleniumType int

type Client struct {
	SeleniumURL  string
	SeleniumType SeleniumType

	Username string
	Password string

	AutoLogin bool

	cookies []http.Cookie
}

func (c *Client) Login() error {
	golog.Debugf("login with %s", c.Username)

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{"--no-sandbox", "--disable-dev-shm-usage", "--remote-allow-origins=*", "-start-maximized"}})

	u := c.SeleniumURL

	if c.SeleniumType == SeleniumTypeLocal {
		d, err := selenium.NewChromeDriverService("./driver.exe", 4444)
		if err != nil {
			return err
		}
		defer d.Stop()
		u = ""

		golog.Debug("using local selenium")
	}

	driver, err := selenium.NewRemote(caps, u)
	if err != nil {
		return err
	}
	defer driver.Quit()

	//get the login page
	err = driver.Get(MonBureauNumeriqueLoginURlRedirect)
	if err != nil {
		return err
	}

	//select the user type
	elem1, err := driver.FindElement(selenium.ByXPATH, "//label[@for=\"idp-EDU\"]")
	if err != nil {
		return err
	}
	if err = elem1.Click(); err != nil {
		return err
	}

	//get button and click
	elem, err := driver.FindElement(selenium.ByID, "button-submit")
	if err != nil {
		return err
	}
	if err = elem.Click(); err != nil {
		return err
	}

	//click on bouton_eleve
	elem, err = driver.FindElement(selenium.ByID, "bouton_eleve")
	if err != nil {
		return err
	}
	if err = elem.Click(); err != nil {
		return err
	}

	//set username input
	elem, err = driver.FindElement(selenium.ByID, "username")
	if err != nil {
		return err
	}
	if err = elem.SendKeys(c.Username); err != nil {
		return err
	}

	//set password input
	elem, err = driver.FindElement(selenium.ByID, "password")
	if err != nil {
		return err
	}
	if err = elem.SendKeys(c.Password); err != nil {
		return err
	}

	//click on bouton_valider
	elem, err = driver.FindElement(selenium.ByID, "bouton_valider")
	if err != nil {
		return err
	}

	if err = elem.Click(); err != nil {
		return err
	}

	_, err = driver.FindElement(selenium.ByCSSSelector, ".user")
	if err != nil {
		return errors.New("invalid credentials")
	}

	//get cookie
	cs, err := driver.GetCookies()
	if err != nil {
		return err
	}

	//convert cookie to skolengo cookie
	var skolengoCookies []http.Cookie
	for _, c := range cs {
		a := http.Cookie{
			Name:   c.Name,
			Value:  c.Value,
			Path:   c.Path,
			Domain: c.Domain,
			Secure: c.Secure,
		}
		a.Expires = time.Unix(int64(c.Expiry), 0)
		skolengoCookies = append(skolengoCookies, a)
	}

	//set cookie
	c.SetCookies(skolengoCookies)

	println(c.GetCookiesString())

	return nil
}

func (c *Client) CheckCookieValidity() bool {
	if len(c.cookies) == 0 {
		return false
	}

	d, err := url.Parse(MonBureauNumeriqueHomeURL)
	if err != nil {
		return false
	}

	builder := NewURLBuilder(d)
	builder.AddParam("PROC", "PAGE_ACCUEIL")
	builder.AddParam("ACTION", "VALIDER")

	u, err := builder.Build()
	if err != nil {
		return false
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return false
	}

	req.Header.Set("cookie", c.GetCookiesString())

	resp, err := client.Do(req)
	if err != nil {
		return false
	}

	if resp.StatusCode != 200 {
		return false
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return false
	}

	if doc.Find(".user").Length() == 0 {
		return false
	}

	return true
}

func (c *Client) Get(u *url.URL) (*http.Response, error) {
	start := time.Now()
	if c.AutoLogin && !c.CheckCookieValidity() {
		golog.Debug("cookie invalid, logging in")
		err := c.Login()
		if err != nil {
			return nil, errors.New("error while logging in: " + err.Error())
		}
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
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

	var sinceString string
	since := time.Since(start)
	switch {
	case since < 500*time.Millisecond:
		sinceString = color.HiGreenString(time.Since(start).String())
	case since < 1*time.Second:
		sinceString = color.HiYellowString(time.Since(start).String())
	case since < 2*time.Second:
		sinceString = color.HiRedString(time.Since(start).String())
	default:
		sinceString = color.HiMagentaString(time.Since(start).String())
	}

	g := color.New(color.BgCyan)
	g.Add(color.FgWhite)
	g.Add(color.Bold)

	golog.Debugf("%s %s %s", g.Sprintf("GET"), u.String(), sinceString)

	return resp, nil
}

func (c *Client) Post(u *url.URL) (*http.Response, error) {
	start := time.Now()
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

	var sinceString string
	since := time.Since(start)
	switch {
	case since < 500*time.Millisecond:
		sinceString = color.HiGreenString(time.Since(start).String())
	case since < 1*time.Second:
		sinceString = color.HiYellowString(time.Since(start).String())
	case since < 2*time.Second:
		sinceString = color.HiRedString(time.Since(start).String())
	default:
		sinceString = color.HiMagentaString(time.Since(start).String())
	}

	g := color.New(color.BgMagenta)
	g.Add(color.FgWhite)
	g.Add(color.Bold)

	golog.Debugf("%s %s %s", g.Sprintf("POST"), u.String(), sinceString)

	return resp, nil
}

func (c *Client) GetName() (string, error) {
	d, err := url.Parse(MonBureauNumeriqueHomeURL)
	if err != nil {
		return "", err
	}

	builder := NewURLBuilder(d)
	builder.AddParam("PROC", "PAGE_ACCUEIL")
	builder.AddParam("ACTION", "VALIDER")

	u, err := builder.Build()
	if err != nil {
		return "", err
	}

	resp, err := c.Get(u)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New("invalid status code")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	name := doc.Find("body > div.header > nav > ul.user > li:nth-child(1)").Text()
	return name, nil
}

func GetServiceByID(services []*Service, id string) (*Service, error) {
	for _, service := range services {
		if service.Id() == id {
			return service, nil
		}
	}
	return nil, errors.New("service not found")
}

func (c *Client) GetServices() ([]*Service, error) {
	d, err := url.Parse(MonBureauNumeriqueHomeURL)
	if err != nil {
		return nil, err
	}

	builder := NewURLBuilder(d)
	builder.SetPath("sg.do")
	builder.AddParam("PROC", "PAGE_ACCUEIL")
	builder.AddParam("ACTION", "VALIDER")
	homeURL, err := builder.Build()
	if err != nil {
		return nil, err
	}

	get, err := c.Get(homeURL)
	if err != nil {
		return nil, err
	}

	if get.StatusCode != 200 {
		return nil, errors.New("invalid cookie")
	}
	defer get.Body.Close()

	doc, err := goquery.NewDocumentFromReader(get.Body)
	if err != nil {
		return nil, err
	}

	var services []*Service
	doc.Find("body > div.header > div.header__set > div.header__set2 > nav > div > div > ul > li > a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		name := s.Text()

		u, err := url.Parse(href)
		if err != nil {
			return
		}

		h := u.Query().Get("service")
		if h == "" {
			return
		}

		parse, err := url.Parse(h)
		if err != nil {
			return
		}

		f := NewService(name, parse)

		services = append(services, f)
	})

	return services, nil
}

// SetAutoLogin enable function to check if the user is logged in on ALL requests
// This methode is not recommended because it will make a request on every request and add a lot of time (20ms per request)
func (c *Client) SetAutoLogin(autoLogin bool) {
	c.AutoLogin = autoLogin
}

func (c *Client) SetSeleniumType(seleniumType SeleniumType) {
	c.SeleniumType = seleniumType
}

func (c *Client) SetUsername(u string) {
	c.Username = u
}

func (c *Client) SetPassword(s string) {
	c.Password = s
}

func (c *Client) GetCookiesJSON() (string, error) {
	d, err := json.Marshal(c.cookies)
	return string(d), err
}

func (c *Client) SetCookies(cs []http.Cookie) {
	c.cookies = cs
}

func (c *Client) GetCookies() []http.Cookie {
	return c.cookies
}

func (c *Client) GetCookiesString() string {
	var s string
	for _, c := range c.cookies {
		s += fmt.Sprintf("%s=%s;", c.Name, c.Value)
	}
	return s
}

func Screenshot(driver selenium.WebDriver) error {
	//create screenshot
	screenshot, err := driver.Screenshot()
	if err != nil {
		return err
	}

	f, err := os.Create("screenshot.png")
	if err != nil {
		return err
	}
	_, err = f.Write(screenshot)
	return err
}
