package model

type WebResponse[T any] struct {
	Data    T             `json:"data"`
	Paging  *PageMetadata `json:"paging,omitempty"`
	Message string        `json:"messages,omitempty"`
}

type PageResponse[T any] struct {
	Data         []T          `json:"data,omitempty"`
	PageMetadata PageMetadata `json:"paging,omitempty"`
}

type PageMetadata struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}

type Auth struct {
	ID       string
	Username string
}
