package tests

import (
	"github.com/brianvoe/gofakeit/v7"
)

type TenderBody struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	ServiceType     string `json:"serviceType"`
	OrganizationId  string `json:"organizationId"`
	CreatorUsername string `json:"creatorUsername"`
}

type BidBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TenderId    string `json:"tenderId"`
	AuthorType  string `json:"authorType"`
	AuthorId    string `json:"authorId"`
}

func RandomTenderBodyJAMBO() TenderBody {
	return TenderBody{
		Name:            gofakeit.Name(),
		Description:     gofakeit.Bird(),
		ServiceType:     gofakeit.RandomString([]string{"Delivery", "Manufacture", "Construction"}),
		OrganizationId:  "550e8400-e29b-41d4-a716-446655440000",
		CreatorUsername: "Jambo",
	}
}

func RandomBidBodyEGER(tenderId string) BidBody {
	//authorId is hardcoded because i cannot get user id by username(it is available only in service layer)
	return BidBody{
		Name:        gofakeit.Name(),
		Description: gofakeit.Slogan(),
		TenderId:    tenderId,
		AuthorType:  "User",
		AuthorId:    "f3d40c06-d106-4ccb-bd5f-106e09eb8d48",
	}
}
