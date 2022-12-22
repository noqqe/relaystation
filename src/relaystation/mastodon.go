package relaystation

import (
	"context"
	"log"
	"net/http"
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

type Mastodon struct {
	Client *mastodon.Client
}

type Attachments []*mastodon.Attachment

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

func newMastodonClient() *Mastodon {
	var m *MastodonConfig
	var c *mastodon.Client
	var mc Mastodon

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
	mc.Client = c
	return &mc
}

// This posts to mastodon
func (m Mastodon) postToMastodon(text string, attachments Attachments) (*mastodon.Status, error) {

	var t *mastodon.Toot
	t.Status = text
	t.MediaIDs = attachments

	status, err := m.Client.PostStatus(context.Background(), t)
	if err != nil {
		log.Fatal(err)
		return &mastodon.Status{}, err
	}

	return status, nil
}

func (m Mastodon) uploadMedia(URLs []string) Attachments {

	var attachments Attachments
	for _, URL := range URLs {

		//Get the response bytes from the url
		response, err := http.Get(URL)
		if err != nil {
			log.Printf("Could not download %s from Twitter\n", URL)
			continue
		}
		defer response.Body.Close()

		if response.StatusCode != 200 {
			log.Printf("Received non 200 response code while trying to download %s\n", URL)
		}

		a, err := m.Client.UploadMediaFromReader(context.Background(), response.Body)
		if err != nil {
			log.Printf("Could not upload %s to Mastodon\n", URL)
			continue
		}

		attachments = append(attachments, a)
		log.Printf("Uploaded image %s to %s\n", URL, a.URL)
	}

	return attachments
}
