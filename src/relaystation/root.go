package relaystation

import (
	"fmt"
	"log"
	"os"

	"github.com/michimani/gotwi"
	"github.com/spf13/cobra"
)

var Version = "unknown"
var dryrun, clean bool

var rootCmd = &cobra.Command{
	Version:               Version,
	Long:                  `relaystation - Twitter Mastodon Forwarder using Twitter stream`,
	Use:                   "relaystation",
	DisableFlagsInUseLine: true,
	SilenceErrors:         true,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetOutput(os.Stdout)
		log.Printf("Starting up relaystation %s", Version)

		var to_create Rules
		var to_delete Rules
		var accs []AccountMap

		if clean {

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

		log.Println("Starting stream...")
		for {
			execSearchStream(accs)
		}

	},
}

func init() {
	rootCmd.Flags().BoolVarP(&dryrun, "dryrun", "u", false, "Don't post anything from anywhere")
	rootCmd.Flags().BoolVarP(&clean, "clean", "c", false, "Redo all the rules from twitter stream api")
}

func Root() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newOAuth2Client() (*gotwi.Client, error) {

	in2 := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth2BearerToken,
	}

	return gotwi.NewClient(in2)
}
