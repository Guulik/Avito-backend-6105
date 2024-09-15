package e2e

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"zadanie-6105/tests/api"
	"zadanie-6105/tests/client"
)

var (
	responsibleUsername = "Jambo"
	PublishTenderURL    string
	bidId               string
	authorName          = "eger"
	PublishBidURL       string
	RejectBidURL        string
	ApproveBidURL       string
)

func TestReject(t *testing.T) {
	_, c := client.New(t)

	//CreateTenderJAMBO
	tender := api.CreateTenderJAMBO(c, t)
	//PublishTender
	PublishTenderURL = client.BaseURL +
		fmt.Sprintf("/tenders/%s/status?status=Published&username=%s", tender.Id, responsibleUsername)
	api.Publish(c, t, PublishTenderURL)
	//CreateBidEGER
	bid := api.CreateBidEGER(c, t, tender.Id)
	bidId = bid.Id
	//PublishBid
	PublishBidURL = client.BaseURL +
		fmt.Sprintf("/bids/%s/status?status=Published&username=%s", bidId, authorName)
	api.Publish(c, t, PublishBidURL)
	//RejectDecision
	RejectBidURL = client.BaseURL +
		fmt.Sprintf("/bids/%s/submit_decision?decision=Rejected&username=%s", bidId, responsibleUsername)
	bid = api.SubmitDecision(c, t, RejectBidURL)
	require.Equal(t, "Canceled", bid.Status)
	//3 version because we changed status 2 times: Created-Published-Canceled
	require.Equal(t, 3, bid.Version)

}

func TestApprove(t *testing.T) {
	_, c := client.New(t)

	//CreateTenderJAMBO
	tender := api.CreateTenderJAMBO(c, t)
	//PublishTender
	PublishTenderURL = client.BaseURL +
		fmt.Sprintf("/tenders/%s/status?status=Published&username=%s", tender.Id, responsibleUsername)
	api.Publish(c, t, PublishTenderURL)
	//CreateBids
	_ = api.CreateBidEGER(c, t, tender.Id)
	_ = api.CreateBidEGER(c, t, tender.Id)
	bid := api.CreateBidEGER(c, t, tender.Id)
	bidId = bid.Id
	//PublishBid
	PublishBidURL = client.BaseURL +
		fmt.Sprintf("/bids/%s/status?status=Published&username=%s", bidId, authorName)
	api.Publish(c, t, PublishBidURL)
	//Approve
	responsibles := []string{"Jambo", "ignat", "test_user"}
	for i := 0; i < 3; i++ {
		//hardcoded in DB
		responsibleUsername = responsibles[i]
		ApproveBidURL = client.BaseURL +
			fmt.Sprintf("/bids/%s/submit_decision?decision=Approved&username=%s", bidId, responsibleUsername)
		bid = api.SubmitDecision(c, t, ApproveBidURL)
	}
	require.Equal(t, "Published", bid.Status)
	//2 version because we changed status once: Created-Published
	require.Equal(t, 2, bid.Version)
	//check other bids is canceled
	otherBids := api.BidsForTender(c, t, tender.Id, "Jambo")
	for _, b := range otherBids {
		if b.Id != bid.Id {
			require.Equal(t, "Canceled", b.Status)
			//3 version because we changed status 2 times: Created-Published-Canceled
			require.Equal(t, 3, b.Version)
		}
	}
}
