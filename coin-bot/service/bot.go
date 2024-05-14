package service

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
)

func (s *botService) Start() error {
	goBot, err := discordgo.New("Bot " + s.BotToken)
	if err != nil {
		return err
	}

	user, err := goBot.User("@me")
	if err != nil {
		return err
	}

	s.BotID = user.ID

	goBot.AddHandler(s.handleMessages)

	err = goBot.Open()
	if err != nil {
		return err
	}

	s.ActiveStreams = make(map[string]chan bool)
	fmt.Println("Go bot is running...")
	return nil
}

func (s *botService) handleMessages(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == s.BotID {
		return
	}

	switch message.Content {
	case s.BotPrefix + "news":
		stopChan := make(chan bool)
		s.Mux.Lock()
		s.ActiveStreams[message.ChannelID] = stopChan
		s.Mux.Unlock()
		go s.startNewsStream(session, message.ChannelID, stopChan)

	case s.BotPrefix + "stop":
		s.Mux.Lock()
		if stopChan, ok := s.ActiveStreams[message.ChannelID]; ok {
			stopChan <- true
			delete(s.ActiveStreams, message.ChannelID)
		}
		s.Mux.Unlock()
	}
}

func (s *botService) startNewsStream(session *discordgo.Session, channelID string, stopChan chan bool) {
	c, _, err := websocket.DefaultDialer.Dial(s.RssService, nil)
	if err != nil {
		log.Printf("Error connecting to WebSocket: %s", err)
		return
	}
	defer c.Close()

	var item FeedItem
	errorCount := 0
	for {
		select {
		case <-stopChan:
			fmt.Println("Stopping news stream for channel", channelID)
			return
		default:
			if errorCount > 5 {
				fmt.Println("Too many errors, stopping the news stream")
				return
			}

			err := c.ReadJSON(&item)
			if err != nil {
				log.Printf("Error reading JSON: %s\n", err)
				errorCount++
				continue
			}

			message := fmt.Sprintf("Title: %s\nLink: %s", item.Title, item.Link)
			session.ChannelMessageSend(channelID, message)
		}
	}
}
