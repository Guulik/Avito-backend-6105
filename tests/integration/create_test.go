package integration

import (
	"fmt"
	"testing"
	"zadanie-6105/tests/api"
	"zadanie-6105/tests/client"
)

var (
	tenderId string
)

func TestCreateTender_Happy(t *testing.T) {
	_, c := client.New(t)
	api.CreateTenderJAMBO(c, t)
}

func TestCreatePublishedTender(t *testing.T) {
	_, c := client.New(t)

	tender := api.CreateTenderJAMBO(c, t)
	tenderId = tender.Id
	PublishTenderURL := client.BaseURL +
		fmt.Sprintf("/tenders/%s/status?status=Published&username=%s", tender.Id, "Jambo")
	api.Publish(c, t, PublishTenderURL)
}
