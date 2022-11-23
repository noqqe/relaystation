package relaystation

import (
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

type Rules []string

// Loads rules from environment variables
// This is needed to create rules at twitter stream api
// export RULE_1="from:foo OR from:bar" -> [ from:foo OR from:bar ]
func loadRules() Rules {

	r := make(Rules, 5)
	envVars := [5]string{"RULE_1", "RULE_2", "RULE_3", "RULE_4", "RULE_5"}

	for i, envVar := range envVars {
		o := os.Getenv(envVar)

		// Rule length can be maximum of 150 chars
		if utf8.RuneCountInString(o) > 150 {
			log.Fatal("Rule contains more then 150 chars")
		}

		// Rule must contain at least 1 search parameter
		if utf8.RuneCountInString(o) > 5 {
			r[i] = o
		}
	}

	if len(r) > 5 {
		log.Fatal("Maximum number of rules (5) exceeded")
	}

	if len(r) < 1 {
		log.Fatal("Min number of rules required is 1, try setting RULE_1 env variable")
	}

	return r
}

// Loads ruleset from environment variables and converts them to a list of accountnames
// This is needed for caching account IDs at startup
// from:foo OR from:bar -> [ foo, bar ]
func loadAccounts() []string {

	var accounts string
	var accountList []string

	envVars := [5]string{"RULE_1", "RULE_2", "RULE_3", "RULE_4", "RULE_5"}

	for _, envVar := range envVars {
		o := os.Getenv(envVar)
		accounts = accounts + o
	}

	accounts = strings.Replace(accounts, "OR ", "", 1000)
	accounts = strings.Replace(accounts, "from:", "", 1000)
	accountList = strings.Split(accounts, " ")
	return accountList

}
