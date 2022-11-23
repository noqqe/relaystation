package relaystation

import (
	"context"
	"log"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/user/userlookup"
	"github.com/michimani/gotwi/user/userlookup/types"
)

type AccountMap struct {
	ID       string
	Username string
}

func fetchUsernames(usernames []string) []AccountMap {

	accs := make([]AccountMap, len(usernames))

	c, err := newOAuth2Client()
	if err != nil {
		log.Println(err)
		return accs
	}

	input := &types.GetByUsernameInput{}
	output := &types.GetByUsernameOutput{}
	input.SetAccessToken(gotwi.AuthenMethodOAuth2BearerToken)

	for i, v := range usernames {
		input.Username = v
		output, err = userlookup.GetByUsername(context.TODO(), c, input)
		if err != nil {
			log.Println(err)
		}
		accs[i].ID = *output.Data.ID
		accs[i].Username = *output.Data.Username
	}

	return accs
}
