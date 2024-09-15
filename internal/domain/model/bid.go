package model

type BidDB struct {
	Id          string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Decision    string `db:"decision"`
	Status      string `db:"status"`
	TenderId    string `db:"tenderid"`

	AuthorType string `db:"authortype"`
	AuthorId   string `db:"authorid"`
	Version    int    `db:"version"`
	CreatedAt  string `db:"createdat"`

	Feedback []Feedback
}

type BidResponse struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	AuthorType string `json:"authorType"`
	AuthorId   string `json:"authorId"`
	Version    int    `json:"version"`
	CreatedAt  string `json:"createdAt"`
}
