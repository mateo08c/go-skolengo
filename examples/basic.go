package examples

import (
	"github.com/kataras/golog"
	"github.com/mateo08c/go-skolengo/skolengo"
)

func Start() {
	golog.SetLevel("debug")

	client, err := skolengo.NewClient("username", "password")
	if err != nil {
		golog.Error(err)
		return
	}
	services, err := client.GetServices()
	if err != nil {
		golog.Error(err)
		return
	}

	for _, service := range services {
		messages, err := service.GetMessages(10)
		if err != nil {
			golog.Error(err)
			return
		}

		for _, message := range messages {
			golog.Infof("Message: %s - %s", message.ID, message.Subject)
		}
	}

	golog.Info("Done")
}
