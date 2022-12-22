# relaystation

relaystation is a software to bring multiple twitter accounts into mastodon
using Twitter Stream API v2.

Our usecase is rather specific hence the name relaystation is borrowed from
pro cycling. I wanted to bring the 18 UCI Worldtour Teams to the Fediverse,
so I wrote this tool.

The goal of relaystation is to be a general purpose tool to forward any
Twitter Stream API expression into the Fediverse.

## Installation

You can fetch the Docker Image or checkout the repo

    task
    ./relaystation

## Configuration

```
export GOTWI_API_KEY=“xxx”
export GOTWI_API_KEY_SECRET=“xxx”
export MASTODON_EMAIL=mail@example.net
export MASTODON_PASSWORD=xxx
export MASTODON_SERVER=https://example.net/
export RULE_1="from:GroupamaFDJ OR from:AG2RCITROENTEAM OR from:AstanaQazTeam OR from:BHRVictorious OR from:BORAhansgrohe OR from:TeamCOFIDIS OR from:EFprocycling"
export RULE_2="from:INEOSGrenadiers OR from:IntermarcheWG OR from:IsraelPremTech OR from:JumboVismaRoad OR from:Lotto_Soudal OR from:Movistar_Team"
export RULE_3="from:qst_alphavinyl OR from:GreenEDGEteam OR from:TeamDSM OR from:TrekSegafredo OR from:TeamEmiratesUAE"
```

## Building

I have a release step built into my task file that relases a new version.

```
task release

```

## Development

Until I have a test mechanism, I use this to craft a handmade tweet
with a valid tweetid that fetches some images and posts a new status at
mastodon

``` golang
m := newMastodonClient()
t := newTwitterClient()

text := "test"
authorid := "12345"
tweetid := "1605518462985682946"

tweet := streamtypes.SearchStreamOutput{}
tweet.Data.Text = &text
tweet.Data.AuthorID = &authorid
tweet.Data.ID = &tweetid

toot := m.ComposeToot(&tweet, Accounts{}, t)
status, err := m.postToMastodon(toot)
log.Println(status, err)

os.Exit(1)
```



## License

Licensed under MIT license.

See [LICENSE.txt](https://raw.githubusercontent.com/noqqe/relaystation/master/LICENSE.txt) file for details.
