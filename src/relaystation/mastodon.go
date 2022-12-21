package relaystation

import (
	"context"
	"log"
	"os"

	"github.com/mattn/go-mastodon"
)

type MastodonConfig struct {
	Email        string
	Password     string
	Server       string
	ClientID     string
	ClientSecret string
}

// holy shit this is crap.
// no validation whatsoever
//
//	curl -X POST \
//		-F 'client_name=Test Application' \
//		-F 'redirect_uris=urn:ietf:wg:oauth:2.0:oob' \
//		-F 'scopes=read write follow push' \
//		-F 'website=https://myapp.example' \
//		https://mastodon.example/api/v1/apps
func loadMastodonCredentials() *MastodonConfig {
	m := MastodonConfig{}
	m.Email = os.Getenv("MASTODON_EMAIL")
	m.Password = os.Getenv("MASTODON_PASSWORD")
	m.Server = os.Getenv("MASTODON_SERVER")
	m.ClientID = os.Getenv("MASTODON_CLIENTID")
	m.ClientSecret = os.Getenv("MASTODON_CLIENTSECRET")
	return &m
}

// This posts to mastodon
func postToMastodon(text string) (*mastodon.Status, error) {

	var m *MastodonConfig
	var c *mastodon.Client
	var t *mastodon.Toot

	m = loadMastodonCredentials()
	app, err := mastodon.RegisterApp(context.Background(), &mastodon.AppConfig{
		Server:       m.Server,
		ClientName:   "relaystation",
		Scopes:       "read write follow",
		Website:      "https://github.com/noqqe/relaystation",
		RedirectURIs: "urn:ietf:wg:oauth:2.0:oob",
	})
	if err != nil {
		log.Fatal(err)
	}
	c = mastodon.NewClient(&mastodon.Config{
		Server:       m.Server,
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
	})

	err = c.Authenticate(context.Background(), m.Email, m.Password)
	if err != nil {
		log.Fatal(err)
	}

	t = &mastodon.Toot{}
	t.Status = text

	if !dryrun {
		status, err := c.PostStatus(context.Background(), t)
		if err != nil {
			log.Fatal(err)
			return status, err
		}
	}

	return &mastodon.Status{}, nil
}
