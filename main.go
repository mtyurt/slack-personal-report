package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mtyurt/slack"
)

type TimeSortableMessages []slack.SearchMessage

func (a TimeSortableMessages) Len() int      { return len(a) }
func (a TimeSortableMessages) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a TimeSortableMessages) Less(i, j int) bool {
	return convertStrTimestampToTime(a[i].Timestamp).Before(convertStrTimestampToTime(a[j].Timestamp))
}

var users map[string]*slack.User

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Printf("SLACK_TOKEN=your-token %s [OPTIONAL ARGUMENTS]\n\n   Optional Arguments:\n", os.Args[0])
		flag.PrintDefaults()
	}

	daily := flag.Bool("daily", true, "To print only previous day's messages")
	days := flag.Int("days", 1, "Number of days to search for in daily mode. Because day search starts from midnight by Slack.")
	weekly := flag.Bool("weekly", false, "To print only previous week's messages")
	short := flag.Bool("short", false, "Print only short output")
	extraSearch := flag.String("extra-search", " ", "Default search mode is 'from:me', use this flag if you want extra conditions on top of it, e.g.: '-extra-search=in:#channel'; in the end the search filter will be: 'from:me in:#channel'")
	flag.Parse()
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		fmt.Println("You should specify SLACK_TOKEN as environment variable!")
		os.Exit(1)
	}
	cli := slack.New(token)

	users = make(map[string]*slack.User)

	msgs := getMessages(cli, *daily, *days, *weekly, *extraSearch)
	msgMap := generateChannelMap(cli, msgs)

	generalMsgSuffix := ""
	if *weekly {
		generalMsgSuffix = " since last week"
	} else if *daily {
		generalMsgSuffix = fmt.Sprintf(" since %d day/s ago", *days)
	}
	fmt.Printf("You have posted %d messages in %d channels%s:\n", len(msgs), len(msgMap), generalMsgSuffix)

	printOutMessages(cli, msgMap, *short)
}

func getMessages(cli *slack.Client, daily bool, days int, weekly bool, extraSearch string) []slack.SearchMessage {
	searchParams := slack.SearchParameters{Sort: "timestamp", Count: 20, Highlight: true}
	if weekly {
		lastWeek := time.Now().Add(time.Hour * 24 * -7).Format("2006-01-02")
		extraSearch = "after:" + lastWeek + " " + extraSearch
	} else if daily {
		yesterday := time.Now().Add(time.Hour * -24 * time.Duration(days)).Format("2006-01-02")
		extraSearch = "after:" + yesterday + " " + extraSearch
	}
	resp, err := cli.SearchMessages("from:me "+extraSearch, searchParams)
	if err != nil {
		panic(err)
	}
	messages := resp.Matches
	for resp.Paging.Page < resp.Paging.Pages {
		searchParams.Page = resp.Paging.Page + 1
		resp, err = cli.SearchMessages("from:me "+extraSearch, searchParams)
		if err != nil {
			panic(err)
		}
		messages = append(messages, resp.Matches...)
	}
	return messages

}
func generateChannelMap(cli *slack.Client, messageList []slack.SearchMessage) map[string][]slack.SearchMessage {
	msgMap := make(map[string][]slack.SearchMessage)
	for _, m := range messageList {
		channelName := normalizeChannelName(cli, m.Channel)
		addMessageToMap(&m, channelName, msgMap)
	}
	return msgMap

}
func printOutMessages(cli *slack.Client, msgMap map[string][]slack.SearchMessage, short bool) {
	fmt.Println("-----------")
	for channel, channelMessages := range msgMap {
		if short {
			fmt.Printf("%d messages in channel %s\n", len(channelMessages), channel)
			continue
		}
		fmt.Println("In channel " + channel + ":")
		sort.Sort(TimeSortableMessages(channelMessages))
		for _, m := range channelMessages {
			fmt.Printf("At %s: %s\n", normalizeTimestamp(m.Timestamp), m.Text)
		}
		fmt.Println("-----------")
	}
}

func addMessageToMap(msg *slack.SearchMessage, channelName string, msgMap map[string][]slack.SearchMessage) {
	var msgList []slack.SearchMessage
	ok := false
	if msgList, ok = msgMap[channelName]; !ok {
		msgList = make([]slack.SearchMessage, 0)
		msgMap[channelName] = msgList
	}
	msgList = append(msgList, *msg)
	msgMap[channelName] = msgList
}

func convertStrTimestampToTime(timestamp string) time.Time {
	ts := strings.Split(timestamp, ".")
	var tsSec, tsNsec int
	var err error
	if tsSec, err = strconv.Atoi(ts[0]); err != nil {
		panic(err)
	}
	if tsNsec, err = strconv.Atoi(ts[1]); err != nil {
		panic(err)
	}
	return time.Unix(int64(tsSec), int64(tsNsec))

}

func normalizeTimestamp(timestamp string) string {
	return convertStrTimestampToTime(timestamp).Format("Jan 2 15:04")
}

func normalizeChannelName(cli *slack.Client, ctxChannel slack.CtxChannel) string {
	name := ctxChannel.Name
	if strings.HasPrefix(name, "U") {
		if user, ok := users[name]; ok {
			return user.Name
		}
		user, err := cli.GetUserInfo(name)
		if err != nil {
			fmt.Printf("User is not found! channel: %v", ctxChannel)
			panic(err)
		}
		users[name] = user
		return user.Name
	} else if strings.HasPrefix(name, "mpdm") {
		name = name[5 : len(name)-2]
		return strings.Replace(name, "--", ",", -1)
	}
	return name
}
