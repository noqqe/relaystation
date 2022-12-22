package relaystation

import (
	"context"
	"log"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/fields"
	"github.com/michimani/gotwi/tweet/tweetlookup"
	tweettypes "github.com/michimani/gotwi/tweet/tweetlookup/types"
)

func fetchTweet(id string) string {

	c, err := newOAuth2Client()
	if err != nil {
		log.Println(err)
		return ""
	}

	// construct input and output for tweet fetching
	input := &tweettypes.GetInput{}
	output := &tweettypes.GetOutput{}
	input.SetAccessToken(gotwi.AuthenMethodOAuth2BearerToken)
	input.Expansions = append(input.Expansions, fields.ExpansionAttachmentsMediaKeys)
	input.ID = id
	input.MediaFields = append(input.MediaFields, fields.MediaFieldUrl)

	output, err = tweetlookup.Get(context.TODO(), c, input)
	if err != nil {
		log.Println(err)
		return ""
	}

	for _, v := range output.Includes.Media {
		log.Printf(gotwi.StringValue(v.URL))
	}

	return ""
}
