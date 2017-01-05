package main

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/koding/multiconfig"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

type DefaultConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	Token          string
	TokenSecret    string
}

type Tweets struct {
	id    int64
	tweet string
}

var conf = loadConfig()

func loadConfig() *DefaultConfig {
	m := multiconfig.NewWithPath("config.json")
	conf := new(DefaultConfig)
	m.MustLoad(conf)

	return conf
}

var client = initClient()

func initClient() *twitter.Client {
	config := oauth1.NewConfig(conf.ConsumerKey, conf.ConsumerSecret)
	token := oauth1.NewToken(conf.Token, conf.TokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient)
}

func main() {

	id := getSinceID()

	// Always fetch favorite 200 tweets from id
	params := &twitter.FavoriteListParams{
		SinceID: id,
		Count:   200,
	}

	tweets, _, _ := client.Favorites.List(params)
	tweetsMap := map[int64]string{}

	for key, tweet := range tweets {
		if key == 0 {
			saveID(tweet.ID)
		}

		regex, _ := regexp.Compile("\n")
		tweet.Text = regex.ReplaceAllString(tweet.Text, "")
		tweetsMap[tweet.ID] = tweet.Text
	}

	saveTweets(tweetsMap)
}

func saveTweets(tweetsMap map[int64]string) {
	f, err := os.OpenFile("./data/tweets.md", os.O_APPEND|os.O_WRONLY, 0600)
	defer f.Close()
	check(err)

	for k, v := range tweetsMap {
		s := "- " + strconv.FormatInt(k, 10) + "," + v + "\n"
		f.WriteString(s)
	}
}

func saveID(id int64) {
	println("Found new tweets")
	f, err := os.Create("./data/since.txt")
	defer f.Close()
	check(err)

	n3, err := f.WriteString(strconv.FormatInt(id, 10))
	fmt.Printf("wrote %d bytes\n", n3)
	f.Sync()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getSinceID() int64 {
	f, err := ioutil.ReadFile("./data/since.txt")
	check(err)
	return stringToInt64(string(f))
}

func stringToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}
