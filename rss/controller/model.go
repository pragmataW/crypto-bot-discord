package controller

import "rrs/model"

type RSSController struct {
	Src IService
}

type IService interface {
	ParseNewFeeds(ret chan model.FeedItem, errs chan error)
}

func NewController(src IService) RSSController {
	return RSSController{
		Src: src,
	}
}
