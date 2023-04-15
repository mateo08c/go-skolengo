package skolengo

import (
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/kataras/golog"
	"github.com/mateo08c/go-skolengo/skolengo/components/inbox"
	"github.com/mateo08c/go-skolengo/skolengo/components/user"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	Name    string
	URL     *url.URL
	Cookies []*http.Cookie
}

func (s *Service) GetCookies() []*http.Cookie {
	return s.Cookies
}

func (s *Service) Get(u *url.URL) (*http.Response, error) {
	golog.Debug("GET: ", u.String())
	client := http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	for _, cookie := range s.Cookies {
		req.AddCookie(cookie)
	}
	return client.Do(req)
}

func (s *Service) Post(u *url.URL) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, err
	}
	for _, cookie := range s.Cookies {
		req.AddCookie(cookie)
	}
	return client.Do(req)
}

func NewService(name string, cookies []*http.Cookie, url *url.URL) *Service {
	url.Path = ""
	url.RawQuery = ""

	return &Service{
		Name:    name,
		URL:     url,
		Cookies: cookies,
	}
}

func (s *Service) GetUID() (*string, error) {
	builderFiche := NewURLBuilder(s.URL)
	builderFiche.SetPath("/kdecole/activation_service.jsp") //This url redirect to the right url
	builderFiche.AddParam("service", "FICHE_ELEVE")
	u, err := builderFiche.Build()
	if err != nil {
		return nil, err
	}

	resp, err := s.Get(u)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	uid := doc.Find("#UID_ELEVE")
	uidVal, _ := uid.Attr("value")

	return &uidVal, nil
}

func (s *Service) GetInbox() (*inbox.Inbox, error) {
	//before get information, we need to get first time the inbox
	err := s.InitInbox()
	if err != nil {
		return nil, errors.New("error while init inbox")
	}

	builder := NewURLBuilder(s.URL)
	builder.SetPath("sg.do")
	builder.AddParam("ACTION", "UPDATE_PAGINATION")
	builder.AddParam("FROM_AJAX", "true")
	builder.AddParam("PROC", "MESSAGERIE")

	u, err := builder.Build()
	if err != nil {
		return nil, err
	}

	resp, err := s.Get(u)
	if err != nil {
		return nil, err
	}

	var infos inbox.Inbox
	err = json.NewDecoder(resp.Body).Decode(&infos)
	if err != nil {
		return nil, err
	}

	return &infos, nil
}

func (s *Service) InitInbox() error {
	builder := NewURLBuilder(s.URL)
	builder.SetPath("sg.do")
	builder.AddParam("PROC", "MESSAGERIE")

	u, err := builder.Build()
	if err != nil {
		return err
	}

	_, err = s.Get(u)
	if err != nil {
		return err
	}

	return nil
}

// GetMessages return the messages of the inbox
// if max is -1, it will return all messages
func (s *Service) GetMessages(max int) ([]*inbox.Message, error) {
	start := time.Now()
	i, err := s.GetInbox()
	if err != nil {
		return nil, err
	}

	builder := NewURLBuilder(s.URL)
	builder.SetPath("sg.do")
	builder.AddParam("PROC", "MESSAGERIE")
	builder.AddParam("ACTION", "REFRESH_FILTER")
	builder.AddParam("NB_ELEMENTS", strconv.Itoa(i.Total))
	builder.AddParam("FROM_AJAX", "true")
	builder.AddParam("TYPE_TRI", "DATE_DESC")

	u, err := builder.Build()
	if err != nil {
		return nil, err
	}

	golog.Debugf(u.String())

	resp, err := s.Post(u)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var messages []*inbox.Message
	doc.Find("#js_boite_reception li").Each(func(i int, se *goquery.Selection) {
		if max != -1 && i >= max {
			return
		}

		href, _ := se.Find("a").Attr("href")

		//parse href
		u, err := url.Parse(href)
		if err != nil {
			golog.Error(err)
		}

		builder := NewURLBuilder(s.URL)
		builder.SetPath("sg.do")
		builder.AddParams(u.Query())

		u, err = builder.Build()
		if err != nil {
			return
		}

		get, err := s.Get(u)
		if err != nil {
			return
		}

		m := new(inbox.Message)
		content, err := goquery.NewDocumentFromReader(get.Body)
		if err != nil {
			return
		}

		m.SetID(u.Query().Get("ID_COMMUNICATION"))
		m.SetFolderID(u.Query().Get("ID_DOSSIER"))
		m.SetSubject(content.Find("#titreCommunication").Text())

		messages = append(messages, m)
	})

	golog.Infof("Get %d messages in %s", len(messages), time.Since(start))

	return messages, nil
}

func (s *Service) GetFolderID() (string, error) {
	builder := NewURLBuilder(s.URL)
	builder.SetPath("sg.do")
	builder.AddParam("PROC", "MESSAGERIE")

	u, err := builder.Build()
	if err != nil {
		return "", err
	}

	resp, err := s.Get(u)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	folderID := doc.Find("#HIDDEN_ID_DOSSIER_COURANT")
	folderIDVal, _ := folderID.Attr("value")

	return folderIDVal, nil
}

func (s *Service) GetInfos() (*user.Info, error) {
	start := time.Now()

	builderCord := NewURLBuilder(s.URL)
	builderCord.SetPath("sg.do")
	builderCord.AddParam("PROC", "COORDONNEES_UTILISATEUR")
	u, err := builderCord.Build()
	if err != nil {
		return nil, err
	}

	resp, err := s.Get(u)
	if err != nil {
		return nil, err
	}

	infos := new(user.Info)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return infos, err
	}

	//parse name and gender
	username := doc.Find("#js-prefsForm > div.container > div.grid--md.grid--template-columns-xl-12 > div > div > div > div:nth-child(1) > div.flex\\@xsb.gap.gap--column-xs.slug.slug--lg > div:nth-child(2) > div > div:nth-child(1) > h3").Text()

	d := strings.Split(username, ".")
	if len(d) != 0 {
		gender := d[0]
		switch gender {
		case "M":
			gender = "Monsieur"
		case "Mme":
			gender = "Madame"
		default:
			gender = ""
		}

		infos.SetGender(gender)
	}

	//get name and last name
	name := strings.Split(username, " ")
	if len(name) > 1 {
		infos.SetFirstName(name[1])
	}
	if len(name) > 2 {
		infos.SetLastName(name[2])
	}

	//get address
	address, _ := doc.Find(".panel address.h6-like").Html()
	add := strings.Split(address, "<br/>")
	var a string
	for _, s := range add {
		a += strings.TrimSpace(s) + " "
	}
	homePhone := doc.Find("#js-prefsForm > div.container > div.grid--md.grid--template-columns-xl-12 > div > div > div > div:nth-child(1) > div.grid--sm.grid--template-columns-sm-2.gap--md > div:nth-child(2)").Text()
	homePhone = strings.TrimSpace(homePhone)

	professionalPhone := doc.Find("#js-prefsForm > div.container > div.grid--md.grid--template-columns-xl-12 > div > div > div > div:nth-child(1) > div.grid--sm.grid--template-columns-sm-2.gap--md > div:nth-child(5)").Text()
	professionalPhone = strings.TrimSpace(professionalPhone)

	mobilePhone := doc.Find("#telephoneMobile")
	mobilePhoneValue, _ := mobilePhone.Attr("value")

	acceptToReceiveSMS := doc.Find("#accepteSMS")
	acceptToReceiveSMSValue, _ := acceptToReceiveSMS.Attr("value")

	redList := doc.Find("#listerouge")
	redListValue, _ := redList.Attr("value")

	infos.SetHomePhone(homePhone)
	infos.SetProfessionalPhone(professionalPhone)
	infos.SetMobilePhone(mobilePhoneValue)

	infos.SetAddress(a)

	infos.SetAgreeToReceiveSMS(acceptToReceiveSMSValue == "0")
	infos.SetShowInformationToPublic(redListValue == "0")

	//get more info
	builderFiche := NewURLBuilder(s.URL)
	builderFiche.SetPath("/kdecole/activation_service.jsp")
	builderFiche.AddParam("service", "FICHE_ELEVE")
	u, err = builderFiche.Build()
	if err != nil {
		return nil, err
	}

	resp, err = s.Get(u)
	if err != nil {
		return nil, err
	}

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return infos, err
	}

	uid := doc.Find("#UID_ELEVE")
	uidVal, _ := uid.Attr("value")

	birthDate := doc.Find("#sele--ficheeleve_resume > p.h6-like.h6-like--lead-xs > time")
	birthDateValue, _ := birthDate.Attr("datetime")

	t, err := time.Parse("2006-01-02T15:04:05.000Z07:00", birthDateValue)
	if err != nil {
		return infos, err
	}

	class := doc.Find("#sele--ficheeleve_resume > p.h6-like.h6-like--lead-xs > strong:nth-child(3)")
	classValue := strings.TrimSpace(class.Text())

	groups := doc.Find("#sele--ficheeleve_resume > p:nth-child(3)")
	groups.Find("span").Remove()
	groupsValue := strings.Split(strings.TrimSpace(groups.Text()), ";")

	regime := doc.Find("#sele--ficheeleve_resume > p:nth-child(4) > span.cartouche.cartouche--bg-primary.text--uppercase")
	regimeValue := strings.TrimSpace(regime.Text())

	responsable := doc.Find("#sele--ficheeleve_contacts > div > div.col.col--gutter-sm.col--xs-12.col--md-6.col--lg-12.col--xl-7 > ul > li")
	responsable.Each(func(i int, s *goquery.Selection) {
		var firstname, lastname, gender string

		name := s.Find("a.js-tuteur__lien").Text()
		split := strings.Split(name, " ")

		if len(split) != 0 {
			switch split[0] {
			case "M.":
				gender = "Monsieur"
			case "Mme":
				gender = "Madame"
			default:
				gender = "Autre"
			}
		}

		if len(split) > 1 {
			firstname = split[1]
		}

		if len(split) > 2 {
			lastname = split[2]
		}

		var places []string
		place := s.Find("span.cartouche")
		place.Each(func(i int, s *goquery.Selection) {
			places = append(places, s.Text())
		})

		infos.AddLegalResponsible(lastname, places, firstname, gender)
	})

	infos.SetRegime(regimeValue)
	infos.SetGroups(groupsValue)
	infos.SetUID(uidVal)
	infos.SetBirthDate(t)
	infos.SetClass(classValue)

	builderInbox := NewURLBuilder(s.URL)
	builderInbox.SetPath("sg.do")
	builderInbox.AddParam("PROC", "PARAMETRAGE_GENERAL")
	u, err = builderInbox.Build()
	if err != nil {
		return nil, err
	}

	resp, err = s.Get(u)
	if err != nil {
		return nil, err
	}

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return infos, err
	}

	email := doc.Find("#screenreader-contenu > form > div.container > div:nth-child(2) > div > div > p")
	email.Find("label").Remove()
	emailValue := strings.TrimSpace(email.Text())

	messageCancellation := doc.Find("#ANNULATION_ENVOI_MESSAGE_COCHE")
	messageCancellationValue, _ := messageCancellation.Attr("value")

	infos.SetEmail(emailValue)
	infos.SetMessageCancellation(messageCancellationValue == "0")

	builderPref := NewURLBuilder(s.URL)
	builderPref.SetPath("sg.do")
	builderPref.AddParam("PROC", "PREFERENCES_UTILISATEUR")
	builderPref.AddParam("ACTION", "PREFS")
	u, err = builderPref.Build()
	if err != nil {
		return nil, err
	}

	resp, err = s.Get(u)
	if err != nil {
		return nil, err
	}

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return infos, err
	}

	//TODO : parse css to get the value... but no found for the moment
	notificationEmail := doc.Find("#INFOS_UTILISATEUR\\.email")
	notificationEmailValue, _ := notificationEmail.Attr("value")

	infos.SetNotificationEmail(notificationEmailValue)

	elapsed := time.Since(start)
	golog.Infof("Scraping took %s", elapsed.String())

	return infos, nil
}

func (s *Service) Id() string {
	return strings.ToLower(strings.ReplaceAll(s.Name, " ", "-"))
}
