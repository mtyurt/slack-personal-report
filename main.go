package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mtyurt/slack"
)

func main() {
	user := os.Args[1]
	token := os.Args[2]
	cli := slack.New(token)
	fmt.Println("searching slack")

	searchParams := slack.SearchParameters{Sort: "timestamp", Count: 100, Highlight: true}

	resp, err := cli.SearchMessages("from:me", searchParams)
	if err != nil {
		panic(err)
	}
	msgMap := make(map[string][]slack.SearchMessage)
	for _, m := range resp.Matches {
		channelName := normalizeChannelName(cli, m.Channel)
		addMessageToMap(&m, channelName, msgMap)
	}
	for channel, msgList := range msgMap {
		fmt.Println("In channel " + channel + ":")
		for _, m := range msgList {
			fmt.Printf("At %s: %s\n", normalizeTimestamp(m.Timestamp), m.Text)
		}
		fmt.Println("-----------\n")
	}

}
func addMessageToMap(msg *slack.SearchMessage, channelName string, msgMap map[string][]slack.SearchMessage) {
	var msgList []slack.SearchMessage
	if msgList, ok := msgMap[channelName]; !ok {
		msgList = make([]slack.SearchMessage, 0)
		msgMap[channelName] = msgList
	}
	msgList = append(msgList, *msg)
	msgMap[channelName] = msgList
}
func normalizeTimestamp(timestamp string) string {
	ts := strings.Split(timestamp, ".")
	var tsSec, tsNsec int
	var err error
	if tsSec, err = strconv.Atoi(ts[0]); err != nil {
		panic(err)
	}
	if tsNsec, err = strconv.Atoi(ts[1]); err != nil {
		panic(err)
	}
	return time.Unix(int64(tsSec), int64(tsNsec)).Format("Jan 2 15:04")
}
func normalizeChannelName(cli *slack.Client, ctxChannel slack.CtxChannel) string {
	name := ctxChannel.Name
	if strings.HasPrefix(name, "U") {
		user, err := cli.GetUserInfo(name)
		if err != nil {
			fmt.Printf("User is not found! channel: %v", ctxChannel)
			panic(err)
		}
		return user.Name
	} else if strings.HasPrefix(name, "mpdm") {
		name = name[5 : len(name)-2]
		return strings.Replace(name, "--", ",", -1)
	}
	return name
}
