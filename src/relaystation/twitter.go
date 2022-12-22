package relaystation

import (
	"context"
	"fmt"
	"html"
	"log"
	"strings"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/fields"
	"github.com/michimani/gotwi/tweet/filteredstream"
	streamtypes "github.com/michimani/gotwi/tweet/filteredstream/types"
	"github.com/michimani/gotwi/tweet/tweetlookup"
	tweettypes "github.com/michimani/gotwi/tweet/tweetlookup/types"
	"github.com/michimani/gotwi/user/userlookup"
	usertypes "github.com/michimani/gotwi/user/userlookup/types"
)

type Accounts []AccountMap
type AccountMap struct {
	ID       string
	Username string
}

type Twitter struct {
	Client *gotwi.Client
}

func newTwitterClient() (Twitter, error) {

	in2 := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth2BearerToken,
	}

	var t Twitter
	var err error
	t.Client, err = gotwi.NewClient(in2)
	if err != nil {
		log.Fatalln(err)
	}

	return t, nil
}

// Stream API

func (t Twitter) execSearchStream(accs Accounts) {
	m := newMastodonClient()

	p := &streamtypes.SearchStreamInput{
		Expansions: fields.ExpansionList{fields.ExpansionAuthorID, fields.ExpansionAttachmentsMediaKeys},
	}

	s, err := filteredstream.SearchStream(context.Background(), t.Client, p)
	if err != nil {
		log.Println(err)
		return
	}

	for s.Receive() {
		t, err := s.Read()
		if err != nil {
			log.Println(err)
		} else {
			if t != nil {

				username := accs.translateIDtoUsername(gotwi.StringValue(t.Data.AuthorID))
				log.Printf("Found Tweet from %s (%s): %s", username, gotwi.StringValue(t.Data.AuthorID), gotwi.StringValue(t.Data.ID))
				toottext := html.UnescapeString(gotwi.StringValue(t.Data.Text))

				if !dryrun {
					status, err := m.postToMastodon(username + ": " + toottext)
					if err != nil {
						log.Printf("Error posting tweet to mastodon. Error: %s\n", err)
					} else {
						log.Printf("Posted tweet from %s to mastodon: %s\n", username, status.URL)
					}
				} else {
					log.Printf("Not posting tweet to mastodon. Because --dryrun is active.\n")
				}
			}
		}
	}
}

func (t Twitter) createSearchStreamRules(keyword string) {

	p := &streamtypes.CreateRulesInput{
		Add: []streamtypes.AddingRule{

			{Value: gotwi.String(keyword), Tag: gotwi.String(keyword)},
		},
	}

	res, err := filteredstream.CreateRules(context.TODO(), t.Client, p)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for _, r := range res.Data {
		log.Printf("Rule: %s\n", gotwi.StringValue(r.Value))
	}
}

func (t Twitter) listSearchStreamRules() (error, Rules) {

	rules := make(Rules, 5)

	p := &streamtypes.ListRulesInput{}
	res, err := filteredstream.ListRules(context.Background(), t.Client, p)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("Could not get rules"), rules
	}

	for i, r := range res.Data {
		rules[i] = gotwi.StringValue(r.ID)
		// log.Printf("ID: %s, Value: %s\n", gotwi.StringValue(r.ID), gotwi.StringValue(r.Value))
	}

	return nil, rules
}

func (t Twitter) deleteSearchStreamRules(ruleID string) {

	if ruleID == "" {
		return
	}

	log.Println("Deleting", ruleID)
	p := &streamtypes.DeleteRulesInput{
		Delete: &streamtypes.DeletingRules{
			IDs: []string{
				ruleID,
			},
		},
	}

	res, err := filteredstream.DeleteRules(context.TODO(), t.Client, p)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for _, r := range res.Data {
		log.Println(gotwi.StringValue(r.Value))
	}
}

// Tweet API

func (t Twitter) fetchTweet(id string) []string {

	// construct input and output for tweet fetching
	input := &tweettypes.GetInput{}
	output := &tweettypes.GetOutput{}
	var urls []string

	// configure input
	input.SetAccessToken(gotwi.AuthenMethodOAuth2BearerToken)
	input.Expansions = append(input.Expansions, fields.ExpansionAttachmentsMediaKeys)
	input.ID = id
	input.MediaFields = append(input.MediaFields, fields.MediaFieldUrl)

	output, err := tweetlookup.Get(context.TODO(), t.Client, input)
	if err != nil {
		log.Println(err)
		return urls
	}

	for _, v := range output.Includes.Media {
		log.Printf(gotwi.StringValue(v.URL))
		urls = append(urls, gotwi.StringValue(v.URL))
	}

	return urls
}

// User API

func (t Twitter) fetchUsernames(usernames []string) Accounts {

	accs := make(Accounts, len(usernames))

	input := &usertypes.GetByUsernameInput{}
	input.SetAccessToken(gotwi.AuthenMethodOAuth2BearerToken)

	for i, v := range usernames {
		input.Username = v
		output, err := userlookup.GetByUsername(context.TODO(), t.Client, input)
		if err != nil {
			log.Println(err)
		}
		accs[i].ID = *output.Data.ID
		accs[i].Username = *output.Data.Name
		log.Printf("Tracking: " + *output.Data.Name)
	}

	return accs
}

// Find Username by ID in Slice
func (accs Accounts) translateIDtoUsername(id string) string {
	for _, v := range accs {
		if strings.Contains(v.ID, id) {
			return v.Username
		}
	}
	return ""
}
