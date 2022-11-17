package relaystation

import (
	"log"
	"os"
	"unicode/utf8"
)

type Rules []string

func loadRules() Rules {

	r := make(Rules, 5)
	envvars := [5]string{"RULE_1", "RULE_2", "RULE_3", "RULE_4", "RULE_5"}

	for i, envvar := range envvars {
		o := os.Getenv(envvar)

		// Rule length can be maximum of 150 chars
		if utf8.RuneCountInString(o) > 150 {
			log.Fatal("Rule %d contains more then 150 chars", i+1)
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
