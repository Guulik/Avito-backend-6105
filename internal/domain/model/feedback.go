package model

type Feedback struct {
	Id          string `db:"id"`
	Description string `db:"description"`
	CreatedAt   string `db:"createdat"`
}
