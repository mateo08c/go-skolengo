package examples

import (
	"encoding/json"
	"github.com/kataras/golog"
	"github.com/mateo08c/go-skolengo/skolengo"
	"net/http"
	"os"
	"strings"
)

func main() {
	golog.SetLevel("debug")

	client, err := skolengo.NewClient("http://selenoid.autotests.cloud:4444/wd/hub")
	if err != nil {
		golog.Error(err)
		return
	}

	client.SetSeleniumType(skolengo.SeleniumTypeLocal)

	client.SetUsername("username")
	client.SetPassword("password")
	client.SetAutoLogin(true)

	golog.Info("created client")
	cookies := ReadCookieFromFile("cookies.json")
	if len(cookies) != 0 {
		golog.Infof("Loaded %d cookies from file", len(cookies))
		client.SetCookies(cookies)
	}

	if len(client.GetCookies()) == 0 {
		golog.Error("no cookies found")
	}

	services, err := client.GetServices()
	if err != nil {
		golog.Error(err)
		return
	}

	service, err := skolengo.GetServiceByID(services, "service-id")
	if err != nil {
		golog.Error(err)
		return
	}

	uid, err := service.GetUID(client)
	if err != nil {
		return
	}
	golog.Info(strings.Repeat("-", 50))
	golog.Infof("UID %s", uid)

	infos, err := service.GetInfos(client)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info(strings.Repeat("-", 50))
	golog.Infof("UID: %s", infos.UID)
	golog.Infof("First name: %s", infos.FirstName)
	golog.Infof("Last name: %s", infos.LastName)
	golog.Infof("Birth date: %s", infos.BirthDate.String())
	golog.Infof("Class: %s", infos.Class)
	golog.Infof("Email: %s", infos.Email)
	golog.Infof("Address: %s", infos.Address)
	golog.Infof("Home phone: %s", infos.HomePhone)
	golog.Infof("Professional phone: %s", infos.ProfessionalPhone)
	golog.Infof("Mobile phone: %s", infos.MobilePhone)
	golog.Infof("Gender: %s", infos.Gender)
	golog.Infof("Regime: %s", infos.Regime)
	golog.Infof("Groups: %s", strings.Join(infos.Groups, ", "))
	golog.Info(strings.Repeat("-", 50))

	for _, legal := range infos.LegalResponsible {
		golog.Infof("First name: %s", legal.FirstName)
		golog.Infof("Last name: %s", legal.LastName)
		golog.Infof("Gender: %s", legal.Gender)
		golog.Infof("Place: %s", strings.Join(legal.Places, ", "))
	}

	golog.Info(strings.Repeat("-", 50))

	golog.Infof("Message cancellation: %t", infos.Settings.MessageCancellation)
	golog.Infof("Email notifications: %s", infos.Settings.Notification.Email)
	golog.Infof("SMS notifications: %t", infos.Settings.AgreeToReceiveSMS)
	golog.Infof("Push notifications: %t", infos.Settings.SmsEnabled)
	golog.Infof("Show information to public: %t", infos.Settings.ShowInformationToPublic)
	golog.Infof("New automatic message: %t", infos.Settings.Notification.NewAutomaticMessage)
	golog.Infof("New ENT message: %t", infos.Settings.Notification.NewENTMessage)
	golog.Infof("New external message: %t", infos.Settings.Notification.NewExternalMessage)
	golog.Infof("Accept regional email: %t", infos.Settings.Notification.AcceptRegionalEmail)
	golog.Info(strings.Repeat("-", 50))

	folderID, err := service.GetFolderID(client)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Infof("Folder ID: %s", folderID)
	golog.Info(strings.Repeat("-", 50))

	inbox, err := service.GetInbox(client)
	if err != nil {
		golog.Error(err)
		return
	}
	golog.Infof("Total messages: %d", inbox.Total)
	golog.Infof("Number of elements: %d", inbox.NbElements)
	golog.Infof("First message: %s", inbox.Premier)

	golog.Info(strings.Repeat("-", 50))

	// Get all messages
	_, err = service.GetMessages(client)
	if err != nil {
		golog.Error(err)
		return
	}

	if len(client.GetCookies()) != 0 {
		SaveCookieJSON(client.GetCookies(), "cookies.json")
	}
}

func SaveCookieJSON(cookies []http.Cookie, filename string) {
	cookieFile, _ := os.Create(filename)
	defer cookieFile.Close()
	json.NewEncoder(cookieFile).Encode(cookies)
}

func ReadCookieFromFile(filename string) []http.Cookie {
	cookieFile, _ := os.Open(filename)
	defer cookieFile.Close()
	var cookies []http.Cookie
	json.NewDecoder(cookieFile).Decode(&cookies)
	return cookies
}
