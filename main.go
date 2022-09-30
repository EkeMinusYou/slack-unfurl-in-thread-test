package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var (
	SlackToken        string
	VerificationToken string
)

func init() {
	SlackToken = os.Getenv("SLACK_TOKEN")
	VerificationToken = os.Getenv("VERIFICATION_TOKEN")
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println(SlackToken)
	log.Println(VerificationToken)
	client := slack.New(SlackToken, slack.OptionDebug(true))
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read request body: %s", err.Error())
		return
	}
	event, err := slackevents.ParseEvent(
		json.RawMessage(body),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: VerificationToken}),
	)
	if err != nil {
		log.Printf("failed to parse event: %s", err.Error())
		return
	}

	linkSharedEvent, ok := event.InnerEvent.Data.(*slackevents.LinkSharedEvent)
	if !ok {
		log.Printf("event is not LinkSharedEvent")
		return
	}

	if event.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		json.Unmarshal([]byte(body), &r)
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}

	externalId := uuid.NewString()
	remoteFile, err := client.AddRemoteFile(slack.RemoteFileParameters{
		ExternalID:   "slack unfurl test: " + externalId,
		ExternalURL:  "http://example.com",
		Title:        "slack unfurl test",
		PreviewImage: "https://go.dev/images/gophers/ladder.svg",
	})
	if err != nil {
		log.Printf("failed to add remote file. err=%s", err)
		return
	}
	fmt.Println(remoteFile.ExternalID)

	blocks := make([]slack.Block, 0, 1)
	blocks = append(blocks, slack.NewFileBlock("", remoteFile.ExternalID, "remote"))
	link := linkSharedEvent.Links[0]
	_, _, _, err = client.UnfurlMessage(
		linkSharedEvent.Channel,
		linkSharedEvent.TimeStamp,
		map[string]slack.Attachment{
			link.URL: {
				Blocks: slack.Blocks{
					BlockSet: blocks,
				},
			},
		},
	)
	if err != nil {
		log.Printf("failed to unfurl remote file. URL=%s err=%s", link.URL, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
