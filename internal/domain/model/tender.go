package model

type TenderDB struct {
	Id              string `db:"id"`
	Name            string `db:"name"`
	Description     string `db:"description"`
	ServiceType     string `db:"servicetype"`
	Status          string `db:"status"`
	OrganizationId  string `db:"organization_id"`
	CreatorUsername string `db:"creator_username"`
	Version         int    `db:"version"`
	CreatedAt       string `db:"created_at"`
}

type TenderResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ServiceType string `json:"serviceType"`
	Version     int    `json:"version"`
	CreatedAt   string `json:"createdAt"`
}
