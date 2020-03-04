package rss

type CreateForm struct {
	DomainID uint   `json:"domain_id" form:"domain_id"`
	RssLink  string `json:"rss_link" form:"rss_link"`
}

type UpdateForm struct {
	DomainID uint   `json:"domain_id" form:"domain_id"`
	RssLink  string `json:"rss_link" form:"rss_link"`
	RssID    uint   `json:"rss_id" form:"rss_id"`
}
