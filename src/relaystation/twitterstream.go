package relaystation

import (
	"context"
	"fmt"
	"log"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/filteredstream"
	"github.com/michimani/gotwi/tweet/filteredstream/types"
)

func execSearchStream(accs Accounts) {
	c, err := newOAuth2Client()
	if err != nil {
		log.Println(err)
		return
	}

	p := &types.SearchStreamInput{}

	s, err := filteredstream.SearchStream(context.Background(), c, p)
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
				log.Println("Found Tweet: " + gotwi.StringValue(t.Data.AuthorID) + " " + gotwi.StringValue(t.Data.ID))
				username := accs.translateIDtoUsername(gotwi.StringValue(t.Data.AuthorID))
				err := postToMastodon(username + ": " + gotwi.StringValue(t.Data.Text))
				if err != nil {
					log.Printf("Error posting tweet to mastodon. Error: %s\n", err)
				} else {
					log.Printf("Posted tweet from %s to mastodon\n", username)
				}
			}
		}
	}
}

func createSearchStreamRules(keyword string) {
	c, err := newOAuth2Client()
	if err != nil {
		log.Println(err)
		return
	}

	p := &types.CreateRulesInput{
		Add: []types.AddingRule{
			{Value: gotwi.String(keyword), Tag: gotwi.String(keyword)},
		},
	}

	res, err := filteredstream.CreateRules(context.TODO(), c, p)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for _, r := range res.Data {
		log.Printf("Rule: %s\n", gotwi.StringValue(r.Value))
	}
}

func listSearchStreamRules() (error, Rules) {

	rules := make(Rules, 5)

	c, err := newOAuth2Client()
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not get OAuth2Client. Check credentials"), rules
	}

	p := &types.ListRulesInput{}
	res, err := filteredstream.ListRules(context.Background(), c, p)
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

func deleteSearchStreamRules(ruleID string) {
	c, err := newOAuth2Client()
	if err != nil {
		log.Println(err)
		return
	}

	if ruleID == "" {
		return
	}

	log.Println("Deleting", ruleID)
	p := &types.DeleteRulesInput{
		Delete: &types.DeletingRules{
			IDs: []string{
				ruleID,
			},
		},
	}

	res, err := filteredstream.DeleteRules(context.TODO(), c, p)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for _, r := range res.Data {
		log.Println(gotwi.StringValue(r.Value))
	}
}
