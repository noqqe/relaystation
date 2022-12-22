package relaystation

import (
	"context"
	"html"
	"log"
	"net/http"
	"os"

	"github.com/mattn/go-mastodon"
	"github.com/michimani/gotwi"
	streamtypes "github.com/michimani/gotwi/tweet/filteredstream/types"
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

// this is a bit of a monster because we need both mastodon and twitter clients
// here (to fetch tw media and upload mastodon media) not sure how to solve
// this best, so far but will work it out in the future
func (m Mastodon) ComposeToot(t *streamtypes.SearchStreamOutput, accs Accounts, tw Twitter) *mastodon.Toot {

	var toot mastodon.Toot
	var attachments Attachments

	// Convert and add text to toot
	username := accs.translateIDtoUsername(gotwi.StringValue(t.Data.AuthorID))
	log.Printf("extracted username: %s", username)
	text := html.UnescapeString(gotwi.StringValue(t.Data.Text))

	toot.Status = text
	// toot.Status = username + ": " + text
	log.Printf("Composed text: %s", toot.Status)

	image_urls := tw.fetchTweet(gotwi.StringValue(t.Data.ID))
	log.Printf("image_urls: %s", image_urls)
	attachments = m.uploadMedia(image_urls)
	log.Printf("attachments: %s", attachments)
	for _, v := range attachments {
		toot.MediaIDs = append(toot.MediaIDs, v.ID)
	}

	return &toot
}

// This posts to mastodon
func (m Mastodon) postToMastodon(toot *mastodon.Toot) (*mastodon.Status, error) {

	status, err := m.Client.PostStatus(context.Background(), toot)
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
