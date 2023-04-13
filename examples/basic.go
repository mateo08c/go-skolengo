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
		infos, err := service.GetInfos()
		if err != nil {
			return
		}

		golog.Info(infos)
	}
}
