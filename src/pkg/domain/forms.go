package domain

type CreateForm struct {
	Name             string `json:"name" form:"name"`
	RssLink          string `json:"rss_link" form:"rss_link"`
	TelegramUsername string `json:"telegram_username" form:"telegram_username"`
}

type UpdateForm struct {
	Name             string `json:"name" form:"name"`
	RssLink          string `json:"rss_link" form:"rss_link"`
	TelegramUsername string `json:"telegram_username" form:"telegram_username"`
}
