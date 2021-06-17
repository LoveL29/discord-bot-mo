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

var (
	Token string
  TenorKey string
)

func init() {
	flag.StringVar(&Token, "t", os.Getenv("bot_token"), "Bot Token")
  flag.StringVar(&TenorKey, "g", os.Getenv("tenor_key"), "Tenor key")
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

	if matched_serach_gif {
    gif_search_keyword := strings.ReplaceAll(m.Content, "!gif ", "")
    url := "https://g.tenor.com/v1/search?q=" + gif_search_keyword + "&key=" + TenorKey + "&limit=1"

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
}

