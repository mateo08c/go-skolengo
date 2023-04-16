# 🏫 go-skolengo

Go-skolengo est un package qui vous permet de récupérer facilement une multitude d'informations depuis la plateforme Skolengo (aussi connue sous le nom de Mon bureau numérique ou Kdecole).

[![wakatime](https://wakatime.com/badge/user/edc0f08e-3aca-4441-8b23-94a859fe119a/project/359c0ab2-2ba2-48c0-9044-5f27807f7e7c.svg)](https://wakatime.com/badge/user/edc0f08e-3aca-4441-8b23-94a859fe119a/project/359c0ab2-2ba2-48c0-9044-5f27807f7e7c)


## Fonctionnalités

- 💬 Récupération des messages ✅
- 🔐 Authentification sécurisée basée sur l'utilisateur et le mot de passe ✅
- 💻 Code facile à utiliser pour votre application ✅
- 📚 Récupération des notes ❌
- 📝 Récupération des devoirs ❌
- 📅 Récupération des emplois du temps ❌

## Utilisation

Pour utiliser go-skolengo, il vous suffit d'installer le package et d'ajouter les informations d'authentification de votre compte. Ensuite, vous pouvez facilement récupérer les informations que vous souhaitez.

Voici un exemple de code pour récupérer les messages et enregistrer les pièces jointes dans un dossier :
```GO
package main

import (
	"github.com/kataras/golog"
	"github.com/mateo08c/go-skolengo/skolengo"
	"os"
)

func main() {
	golog.SetLevel("debug")

	_ = os.Mkdir("attachments", os.ModePerm)

	client, err := skolengo.NewClient(os.Getenv("SKOLENGO_USERNAME"), os.Getenv("SKOLENGO_PASSWORD"))
	if err != nil {
		panic(err)
	}

	services, err := client.GetServices()
	if err != nil {
		panic(err)
	}

	for _, service := range services {
		messages, err := service.GetMessages(-1, true)
		if err != nil {
			panic(err)
		}

		for _, message := range messages {
			for _, attachment := range message.Content.Attachments {
				err := attachment.SaveToFile("attachments/" + attachment.Name + "." + attachment.Extension)
				if err != nil {
					golog.Error(err)
				}
			}
		}
	}
}
```

# Explications
1. [ ] TODO : Ajouter des explications sur le code 🥸

# Contribuer
Si vous souhaitez contribuer à **go-skolengo**, n'hésitez pas à nous envoyer une pull request.

## Crédits
- [Mateo](https://github.com/mateo08c) - Développeur principal