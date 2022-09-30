module slack-unfurl-in-thread-test

go 1.19

require github.com/slack-go/slack v0.11.3

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/pkg/errors v0.8.0 // indirect
)

replace github.com/slack-go/slack => github.com/ledmonster/slack v0.7.3-0.20210101120547-3eb3ab316abf
