package relaystation

import (
	"fmt"
	"log"
	"os"

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
		t, _ := newTwitterClient()

		if clean {

			log.Println("Loading rules from environment")
			to_create = loadRules()

			log.Println("Fetching current rules to delete")
			_, to_delete = t.listSearchStreamRules()

			log.Println("Clearing current rules")
			for _, rule := range to_delete {
				t.deleteSearchStreamRules(rule)
			}

			log.Println("Create new rules")
			for _, rule := range to_create {
				t.createSearchStreamRules(rule)
			}
		}

		urls := t.fetchTweet("1605583309475643393")
		m := newMastodonClient()
		m.uploadMedia(urls)
		os.Exit(1)
		accountids := loadAccounts()
		accs = t.fetchUsernames(accountids)

		log.Println("Starting stream...")
		for {
			t.execSearchStream(accs)
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
