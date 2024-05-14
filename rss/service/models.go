package service

type rssService struct {
	Url string
}

func NewRSSService(url string) rssService {
	return rssService{
		Url: url,
	}
}