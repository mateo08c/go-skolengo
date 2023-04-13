# 🏫 go-skolengo

Go-skolengo est un package qui vous permet de récupérer facilement une multitude d'informations depuis la plateforme Skolengo (aussi connue sous le nom de Mon bureau numérique ou Kdecole). Grâce à une authentification basée sur l'utilisateur et le mot de passe, ce package rend disponible le code pour votre application.

[![wakatime](https://wakatime.com/badge/user/edc0f08e-3aca-4441-8b23-94a859fe119a/project/359c0ab2-2ba2-48c0-9044-5f27807f7e7c.svg)](https://wakatime.com/badge/user/edc0f08e-3aca-4441-8b23-94a859fe119a/project/359c0ab2-2ba2-48c0-9044-5f27807f7e7c)


## Fonctionnalités

- 💬 Récupération des messages ✅
- 🔐 Authentification sécurisée basée sur l'utilisateur et le mot de passe ✅
- 💻 Code facile à utiliser pour votre application ✅
- 📚 Récupération des notes ❌
- 📝 Récupération des devoirs ❌
- 📅 Récupération des emplois du temps ❌

## Utilisation

Pour utiliser go-skolengo, il vous suffit d'installer le package et d'ajouter les informations d'authentification de votre compte Skolengo. Ensuite, vous pouvez facilement récupérer les informations que vous souhaitez depuis Skolengo.

Voici un exemple de code pour récupérer les messages depuis Skolengo :

```go
import "github.com/mateo08c/go-skolengo/skolengo"

func main() {
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
}
```

# Contribuer
Si vous souhaitez contribuer à go-skolengo, n'hésitez pas à nous envoyer une pull request.
