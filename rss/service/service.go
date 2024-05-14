package service

import (
	"log"
	"rrs/model"
	"time"

	"github.com/mmcdole/gofeed"
)

func (s rssService) ParseNewFeeds(ret chan model.FeedItem, errs chan error) {
	fp := gofeed.NewParser()

	items, err := fp.ParseURL(s.Url)
	if err != nil {
		errs <- err
		return
	}

	currentItems := make(map[string]model.FeedItem)
	for _, item := range items.Items {
		log.Println("items filling...")
		i := model.FeedItem{
			Title:       item.Title,
			Link:        item.Link,
		}
		currentItems[item.GUID] = i
		select {
		case ret <- i:
		default:
			log.Println("channel is closed or full")
		}
		time.Sleep(time.Second * 1)
	}

	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()
	log.Println("ticker created...")
	for range ticker.C {
		log.Println("RSS is reading...")
		items, err := fp.ParseURL(s.Url)
		if err != nil {
			errs <- err
			return
		}
		for _, item := range items.Items {
			if _, exists := currentItems[item.GUID]; !exists {
				newItem := model.FeedItem{
					Title:       item.Title,
					Link:        item.Link,
				}
				currentItems[item.GUID] = newItem
				select {
				case ret <- newItem:
				default:
					log.Println("channel is closed or full")
				}
				time.Sleep(time.Second * 1)
			}
		}
	}
}
