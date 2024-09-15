package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
	"time"
	"zadanie-6105/tests"
	"zadanie-6105/tests/client"
)

var (
	CreateTenderURL = client.BaseURL + "/tenders/new"
	CreateBidURL    = client.BaseURL + "/bids/new"
)

type TenderResponse struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	ServiceType string    `json:"serviceType"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
}
type BidResponse struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	AuthorType string    `json:"authorType"`
	AuthorId   string    `json:"authorId"`
	Version    int       `json:"version"`
	CreatedAt  time.Time `json:"createdAt"`
}

func CreateTenderJAMBO(c *client.Suite, t *testing.T) TenderResponse {
	body := tests.RandomTenderBodyJAMBO()
	bodyJSON, err := json.Marshal(body)
	require.NoError(t, err)
	bodyReq := bytes.NewReader(bodyJSON)
	req := client.FormRequest(http.MethodPost, CreateTenderURL, bodyReq)
	resp, err := c.Client.Do(req)
	require.NoError(t, err)

	var tender TenderResponse
	responseBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = json.Unmarshal(responseBytes, &tender)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	return tender
}

func Publish(c *client.Suite, t *testing.T, url string) {
	req := client.FormRequest(http.MethodPut, url, nil)
	resp, err := c.Client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func CreateBidEGER(c *client.Suite, t *testing.T, tenderId string) BidResponse {
	body := tests.RandomBidBodyEGER(tenderId)
	bodyJSON, err := json.Marshal(body)
	require.NoError(t, err)
	bodyReq := bytes.NewReader(bodyJSON)
	req := client.FormRequest(http.MethodPost, CreateBidURL, bodyReq)

	resp, err := c.Client.Do(req)
	require.NoError(t, err)

	var bid BidResponse
	responseBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = json.Unmarshal(responseBytes, &bid)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	return bid
}

func SubmitDecision(c *client.Suite, t *testing.T, url string) BidResponse {
	req := client.FormRequest(http.MethodPut, url, nil)
	resp, err := c.Client.Do(req)
	require.NoError(t, err)

	var bid BidResponse
	responseBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = json.Unmarshal(responseBytes, &bid)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	return bid
}

// GetTenders without service type because of time limits(now 19:30 13.09)
func GetTenders(c *client.Suite, t *testing.T, limit int32, offest int32) []TenderResponse {
	GetTendersURL := client.BaseURL + fmt.Sprintf("/tenders?")
	if limit != 0 {
		GetTendersURL += fmt.Sprintf("limit=%d", limit)
	}
	if offest != 0 {
		GetTendersURL += fmt.Sprintf("offset=%d", offest)
	}
	req := client.FormRequest(http.MethodGet, GetTendersURL, nil)
	resp, err := c.Client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var tenders []TenderResponse
	responseBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = json.Unmarshal(responseBytes, &tenders)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	return tenders
}

func BidsForTender(c *client.Suite, t *testing.T, tenderId string, username string) []BidResponse {
	GetBidsForTenderURL := client.BaseURL + fmt.Sprintf("/bids/%s/list?username=%s", tenderId, username)
	req := client.FormRequest(http.MethodGet, GetBidsForTenderURL, nil)
	resp, err := c.Client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var bids []BidResponse
	responseBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = json.Unmarshal(responseBytes, &bids)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	return bids
}
