package service

import "sync"

type botService struct {
	BotToken      string
	BotPrefix     string
	BotID         string
	RssService    string
	ActiveStreams map[string]chan bool
	Mux           sync.Mutex
}

type FeedItem struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

var (
	BotID string
)

func NewBot(botPrefix string, botToken string, rssService string) botService {
	return botService{
		BotPrefix:  botPrefix,
		BotToken:   botToken,
		RssService: rssService,
	}
}
