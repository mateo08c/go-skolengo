package user

import (
	"encoding/json"
	"time"
)

type Info struct {
	UID               string    `json:"uid"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Gender            string    `json:"gender"`
	BirthDate         time.Time `json:"birth_date"`
	Email             string    `json:"email"`
	Address           string    `json:"address"`
	HomePhone         string    `json:"home_phone"`
	ProfessionalPhone string    `json:"professional_phone"`
	MobilePhone       string    `json:"mobile_phone"`
	Regime            string    `json:"regime"`
	Class             string    `json:"class"`
	Groups            []string  `json:"groups"`
	LegalResponsible  []struct {
		Gender    string   `json:"gender"`
		FirstName string   `json:"first_name"`
		LastName  string   `json:"last_name"`
		Places    []string `json:"places"`
	}
	Settings `json:"settings"`
}

type Settings struct {
	AgreeToReceiveSMS       bool `json:"agree_to_receive_sms"`
	MessageCancellation     bool `json:"enable_message_cancellation"`
	ShowInformationToPublic bool `json:"show_information_to_public"`
	SmsEnabled              bool `json:"sms_enabled"`
	Notification            `json:"notification"`
}

type Notification struct {
	Email               string `json:"notification_email"`
	NewAutomaticMessage bool   `json:"new_automatic_message"`
	NewENTMessage       bool   `json:"new_ent_message"`
	NewExternalMessage  bool   `json:"new_external_message"`
	AcceptRegionalEmail bool   `json:"accept_regional_email"`
}

func (i *Info) SetBirthDate(birthDate time.Time) {
	i.BirthDate = birthDate
}

func (i *Info) AddGroup(group string) {
	i.Groups = append(i.Groups, group)
}

func (i *Info) SetRegime(regime string) {
	i.Regime = regime
}

func (i *Info) SetClass(class string) {
	i.Class = class
}

func (i *Info) AddLegalResponsible(lastname string, places []string, firstname string, gender string) {
	i.LegalResponsible = append(i.LegalResponsible, struct {
		Gender    string   `json:"gender"`
		FirstName string   `json:"first_name"`
		LastName  string   `json:"last_name"`
		Places    []string `json:"places"`
	}{
		Gender:    gender,
		FirstName: firstname,
		LastName:  lastname,
		Places:    places,
	})
}

func (i *Settings) SetMessageCancellation(enabled bool) {
	i.MessageCancellation = enabled
}

func (i *Info) SetFirstName(firstName string) {
	i.FirstName = firstName
}

func (i *Info) SetLastName(lastName string) {
	i.LastName = lastName
}

func (i *Info) SetEmail(email string) {
	i.Email = email
}

func (i *Info) SetAddress(address string) {
	i.Address = address
}

func (i *Info) SetHomePhone(homePhone string) {
	i.HomePhone = homePhone
}

func (i *Info) SetProfessionalPhone(professionalPhone string) {
	i.ProfessionalPhone = professionalPhone
}

func (i *Info) SetMobilePhone(mobilePhone string) {
	i.MobilePhone = mobilePhone
}

func (i *Info) SetAgreeToReceiveSMS(agreeToReceiveSMS bool) {
	i.Settings.AgreeToReceiveSMS = agreeToReceiveSMS
}

func (i *Info) SetShowInformationToPublic(showInformationToPublic bool) {
	i.Settings.ShowInformationToPublic = showInformationToPublic
}

func (i *Info) SetSmsEnabled(smsEnabled bool) {
	i.Settings.SmsEnabled = smsEnabled
}

func (i *Info) SetNotificationEmail(notificationEmail string) {
	i.Settings.Notification.Email = notificationEmail
}

func (i *Info) SetNewAutomaticMessage(newAutomaticMessage bool) {
	i.Settings.Notification.NewAutomaticMessage = newAutomaticMessage
}

func (i *Info) SetNewMessageFromEnt(newENTMessage bool) {
	i.Settings.Notification.NewENTMessage = newENTMessage
}

func (i *Info) SetNewMessageFromExternal(newExternalMessage bool) {
	i.Settings.Notification.NewExternalMessage = newExternalMessage
}

func (i *Info) SetAcceptRegionalEmail(acceptRegionalEmail bool) {
	i.Settings.Notification.AcceptRegionalEmail = acceptRegionalEmail
}

func (i *Info) SetGender(s string) {
	i.Gender = s
}

func (i *Info) SetUID(val string) {
	i.UID = val
}

func (i *Info) SetGroups(groups []string) {
	i.Groups = groups
}

func (i *Info) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, i)
}

func (i *Info) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}
