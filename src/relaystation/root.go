package relaystation

import (
	"context"
	"fmt"
	"log"
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

	r := loadRules()
	// log.Println(r)
	os.Exit(1)
	log.Println("Cleaning rules")
	_, to_clear := listSearchStreamRules()

	deleteSearchStreamRules(to_clear)
	//listSearchStreamRules
	createSearchStreamRules(os.Getenv("RULE_1"))
	//listSearchStreamRules
	execSearchStream()

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
func listSearchStreamRules() (error, Rules) {

	rules := make(Rules, 5)
	c, err := newOAuth2Client()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Could not get OAuth2Client. Check credentials"), rules
	}

	p := &types.ListRulesInput{}
	res, err := filteredstream.ListRules(context.Background(), c, p)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("Could not get rules"), rules
	}

	for i, r := range res.Data {
		rules[i] = gotwi.StringValue(r.ID)
		fmt.Printf("ID: %s, Value: %s, Tag: %s\n", gotwi.StringValue(r.ID), gotwi.StringValue(r.Value), gotwi.StringValue(r.Tag))
	}
	return nil, rules
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
