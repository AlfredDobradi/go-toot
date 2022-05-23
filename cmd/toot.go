package main

import (
	"fmt"
	"os"

	"github.com/alfreddobradi/go-toot/config"
	"github.com/alfreddobradi/go-toot/mastodon"
)

var (
	instanceURL  string = ""
	clientID     string = ""
	clientSecret string = ""
	clientCode   string = ""
)

func main() {

	instanceURL = os.Getenv("TOOT_INSTANCE_URL")
	clientID = os.Getenv("TOOT_CLIENT_ID")
	clientSecret = os.Getenv("TOOT_CLIENT_SECRET")
	clientCode = os.Getenv("TOOT_CLIENT_CODE")

	visibility := "private"
	switch os.Getenv("TOOT_VISIBILITY") {
	case "public", "unlisted":
		visibility = os.Getenv("TOOT_VISIBILITY")
	}

	token := os.Getenv("TOOT_CLIENT_TOKEN")

	if instanceURL == "" {
		fmt.Println("TOOT_INSTANCE_URL variable missing, please set the variable to the URL of your instance without trailing backslashes.")
		os.Exit(1)
	}
	config.SetInstanceURL(instanceURL)

	if token == "" {
		if clientID == "" {
			fmt.Println("TOOT_CLIENT_ID variable missing, please set the variable in order to continue.")
			os.Exit(1)
		}
		config.SetClientID(clientID)

		if clientSecret == "" {
			fmt.Println("TOOT_CLIENT_SECRET variable missing, please set the variable in order to continue.")
			os.Exit(1)
		}
		config.SetClientSecret(clientSecret)

		if clientCode == "" {
			codeURI, err := mastodon.GetCode()
			if err != nil {
				panic(err)
			}

			fmt.Printf("TOOT_CLIENT_CODE variable missing. To get an authorization code, go to %s and authorize the application.", codeURI)
			os.Exit(1)
		} else {
			tokenEntity, err := mastodon.GetToken(clientCode)
			if err != nil {
				panic(err)
			}

			token = tokenEntity.Token
		}
	}

	if len(os.Args) < 2 {
		fmt.Println("Missing Toot body")
		os.Exit(1)
	}

	if err := mastodon.Post(token, visibility, os.Args[len(os.Args)-1]); err != nil {
		panic(err)
	}

	fmt.Println("The toot was posted. Yay!")
}
