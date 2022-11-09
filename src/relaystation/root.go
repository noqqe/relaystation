package relaystation

import (
	"context"
	"fmt"
	"os"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/filteredstream"
	"github.com/michimani/gotwi/tweet/filteredstream/types"
)

func newOAuth2Client() (*gotwi.Client, error) {

	in2 := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth2BearerToken,
	}

	return gotwi.NewClient(in2)
}

func Root() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("The 1st parameter for command is required. (create|stream)")
		os.Exit(1)
	}

	command := args[1]

	switch command {
	case "list":
		// list search stream rules
		listSearchStreamRules()
	case "delete":
		// delete a specified rule
		if len(args) < 3 {
			fmt.Println("The 2nd parameter for rule ID to delete is required.")
			os.Exit(1)
		}

		ruleID := args[2]
		deleteSearchStreamRules(ruleID)
	case "create":
		// create a search stream rule
		if len(args) < 3 {
			fmt.Println("The 2nd parameter for keyword of search stream rule is required.")
			os.Exit(1)
		}

		keyword := args[2]
		createSearchStreamRules(keyword)
	case "stream":
		// exec filtered stream API
		for {
			execSearchStream()
		}
	default:
		fmt.Println("Undefined command. Command should be 'create' or 'stream'.")
		os.Exit(1)
	}
	//createSearchStreamRules("-is:retweet (from:zeitonline OR from:elhotzo")

}

func execSearchStream() {
	c, err := newOAuth2Client()
	if err != nil {
		fmt.Println(err)
		return
	}

	p := &types.SearchStreamInput{}

	s, err := filteredstream.SearchStream(context.Background(), c, p)
	if err != nil {
		fmt.Println(err)
		return
	}

	for s.Receive() {
		t, err := s.Read()
		if err != nil {
			fmt.Println(err)
		} else {
			if t != nil {
				fmt.Println(gotwi.StringValue(t.Data.ID), gotwi.StringValue(t.Data.Text))
			}
		}
	}
}

func createSearchStreamRules(keyword string) {
	c, err := newOAuth2Client()
	if err != nil {
		fmt.Println(err)
		return
	}

	p := &types.CreateRulesInput{
		Add: []types.AddingRule{
			{Value: gotwi.String(keyword), Tag: gotwi.String(keyword)},
		},
	}

	res, err := filteredstream.CreateRules(context.TODO(), c, p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, r := range res.Data {
		fmt.Printf("ID: %s, Value: %s, Tag: %s\n", gotwi.StringValue(r.ID), gotwi.StringValue(r.Value), gotwi.StringValue(r.Tag))
	}
}

// createSearchStreamRules lists search stream rules.
func listSearchStreamRules() {
	c, err := newOAuth2Client()
	if err != nil {
		fmt.Println(err)
		return
	}

	p := &types.ListRulesInput{}
	res, err := filteredstream.ListRules(context.Background(), c, p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, r := range res.Data {
		fmt.Printf("ID: %s, Value: %s, Tag: %s\n", gotwi.StringValue(r.ID), gotwi.StringValue(r.Value), gotwi.StringValue(r.Tag))
	}
}

func deleteSearchStreamRules(ruleID string) {
	c, err := newOAuth2Client()
	if err != nil {
		fmt.Println(err)
		return
	}

	p := &types.DeleteRulesInput{
		Delete: &types.DeletingRules{
			IDs: []string{
				ruleID,
			},
		},
	}

	res, err := filteredstream.DeleteRules(context.TODO(), c, p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, r := range res.Data {
		fmt.Printf("ID: %s, Value: %s, Tag: %s\n", gotwi.StringValue(r.ID), gotwi.StringValue(r.Value), gotwi.StringValue(r.Tag))
	}
}
