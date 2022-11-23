package relaystation

import (
	"log"
	"os"

	"github.com/michimani/gotwi"
)

func newOAuth2Client() (*gotwi.Client, error) {

	in2 := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth2BearerToken,
	}

	return gotwi.NewClient(in2)
}

func Root() {

	var to_create Rules
	var to_delete Rules
	var accs []AccountMap

	// Silly hack for not writing an option parser
	// TODO: make it clean
	if len(os.Args) > 1 {

		log.Println("Loading rules from environment")
		to_create = loadRules()

		log.Println("Fetching current rules to delete")
		_, to_delete = listSearchStreamRules()

		log.Println("Clearing current rules")
		for _, rule := range to_delete {
			deleteSearchStreamRules(rule)
		}

		log.Println("Create new rules")
		for _, rule := range to_create {
			createSearchStreamRules(rule)
		}
	}

	accountids := loadAccounts()
	accs = fetchUsernames(accountids)
	log.Println(accs)

	os.Exit(1)
	log.Println("Current rules configuration:")
	listSearchStreamRules()

	log.Println("Starting stream...")
	for {
		execSearchStream()
	}

}
