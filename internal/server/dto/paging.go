package dto

type Paging struct {
	From      int    `json:"from"`
	Limit     int    `json:"limit"`
	Page      int    `json:"page"`
	To        int    `json:"to"`
	Total     int    `json:"total"`
	TotalPage int    `json:"totalPage"`
	Next      string `json:"next"`
	Prev      string `json:"prev"`
}
