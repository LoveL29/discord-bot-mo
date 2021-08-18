package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
  "net/http"
	"syscall"
  "regexp"
  "io/ioutil"
  "strings"
  "encoding/json"
	"github.com/bwmarrin/discordgo"
)

const TOKYO_LATITUDE string = "35.652832"
const TOKYO_LONGITUDE string = "139.839478"

var (
	Token string
  TenorKey string
  YelpKey string
)

type YelpStruct struct {
  YelpBusinesses []YelpBusiness `json:"businesses"`
}

type YelpBusiness struct {
  Id            string        `json:"id"`
  alias         string        `json:"alias"`
  Name          string        `json:"name"`
  ImageUrl      string        `json:"image_url"`
  IsClosed      bool          `json:"is_closed"`
  Url           string        `json:"url"`
  ReviewCount   int           `json:"review_count"`
  Categories    []interface{} `json:"id"`
  Rating        float64       `json:"rating"`
  Coordinates   interface{}   `json:"coordinates"`
  Transactions  []string      `json:"transactions"`
  Price         string        `json:"price"`
  Location      interface{}   `json:"location"`
  Phone         string        `json:"phone"`
  DisplayPhone  string        `json:"display_phone"`
  Distance      float64       `json:"distance"`
}

func init() {
	flag.StringVar(&Token, "t", os.Getenv("bot_token"), "Bot Token")
  flag.StringVar(&TenorKey, "g", os.Getenv("tenor_key"), "Tenor key")
  flag.StringVar(&YelpKey, "y", os.Getenv("yelp_key"), "Yelp key")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.AddHandler(messageCreate)
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
  }

  matched_serach_gif, _ := regexp.MatchString("!gif .*" , m.Content)
  random_serach_meal, _ := regexp.MatchString("!menu .*" , m.Content)

	if matched_serach_gif {
    gif_search_keyword := strings.ReplaceAll(m.Content, "!gif ", "")
    url := "https://g.tenor.com/v1/search?q=" + gif_search_keyword +
      "&key=" + TenorKey +
      "&limit=1"

    resp, _ := http.Get(url)
    defer resp.Body.Close()
    byteArray, _ := ioutil.ReadAll(resp.Body)
    reuslt_map := make(map[string]interface{})
    json.Unmarshal([]byte(byteArray), &reuslt_map)
    gif_url := reuslt_map["results"].([]interface{})[0].(map[string]interface{})["media"].([]interface{})[0].(map[string]interface{})["gif"].(map[string]interface{})["url"]
    gif_url_str, err := json.Marshal(gif_url)
    if err != nil {
      fmt.Println(err)
    }
    s.ChannelMessageSend(m.ChannelID, strings.Trim(string(gif_url_str), "\""))
  }

  if random_serach_meal {
    meal_search_keyword := strings.ReplaceAll(m.Content, "!menu ", "")
    url := "https://api.yelp.com/v3/businesses/search?term=" + meal_search_keyword +
      "&latitude=" + TOKYO_LATITUDE +
      "&longitude=" + TOKYO_LONGITUDE
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
      fmt.Println(err)
    }
    req.Header.Add("Authorization", "Bearer " + YelpKey)
    client := &http.Client{}
    resp, err := client.Do(req)
    byteArray, _ := ioutil.ReadAll(resp.Body)
    result_yelp_struct := new(YelpStruct)
    json.Unmarshal([]byte(byteArray), &result_yelp_struct)

    if result_yelp_struct.YelpBusinesses[0] {
      s.ChannelMessageSend(m.ChannelID, strings.Trim(string(result_yelp_struct.YelpBusinesses[0].Name), "\""))
      s.ChannelMessageSend(m.ChannelID, strings.Trim(string(result_yelp_struct.YelpBusinesses[0].Url), "\""))
      s.ChannelMessageSend(m.ChannelID, strings.Trim(string(result_yelp_struct.YelpBusinesses[0].ImageUrl), "\""))
    }
  }
}

