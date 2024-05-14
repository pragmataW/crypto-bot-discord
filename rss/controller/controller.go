package controller

import (
	"log"
	"rrs/model"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetUpRoutes(app *fiber.App, ctrl RSSController) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(ctrl.HandleWebsocket))
}

func (ctrl RSSController) HandleWebsocket(c *websocket.Conn) {
	ret := make(chan model.FeedItem)
	errs := make(chan error)

	go ctrl.Src.ParseNewFeeds(ret, errs)

	for {
		select {
		case feedItem := <-ret:
			if err := c.WriteJSON(feedItem); err != nil {
				log.Println("Write error:", err)
				return
			}
		case err := <-errs:
			log.Println("Error received:", err)
			if err := c.WriteMessage(websocket.TextMessage, []byte(err.Error())); err != nil {
				log.Println("Error sending error message:", err)
				return
			}
		}
	}
}
