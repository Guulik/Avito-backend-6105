package model

import (
	"fmt"
	"time"
)

func ConvertBidToResponse(bidDB BidDB) BidResponse {
	bid := BidResponse{
		Id:         bidDB.Id,
		Name:       bidDB.Name,
		Status:     bidDB.Status,
		AuthorType: bidDB.AuthorType,
		AuthorId:   bidDB.AuthorId,
		Version:    bidDB.Version,
	}

	timestamp, err := time.Parse(time.RFC3339, bidDB.CreatedAt)
	if err != nil {
		fmt.Println("failed to parse time for tender")
		return BidResponse{}
	}
	bid.CreatedAt = time.Time.Format(timestamp, time.RFC3339)
	return bid
}

func ConvertBids(bidsDB []BidDB) []BidResponse {
	bids := make([]BidResponse, len(bidsDB))

	for i, bidDB := range bidsDB {
		b := ConvertBidToResponse(bidDB)

		bids[i] = b
	}
	return bids
}

func ConvertTenderToResponse(tenderDB TenderDB) TenderResponse {
	tender := TenderResponse{
		Id:          tenderDB.Id,
		Name:        tenderDB.Name,
		Description: tenderDB.Description,
		Status:      tenderDB.Status,
		ServiceType: tenderDB.ServiceType,
		Version:     tenderDB.Version,
	}

	timestamp, err := time.Parse(time.RFC3339, tenderDB.CreatedAt)
	if err != nil {
		fmt.Println("failed to parse time for tender")
		return TenderResponse{}
	}
	tender.CreatedAt = time.Time.Format(timestamp, time.RFC3339)
	return tender
}

func ConvertTenders(tendersDB []TenderDB) []TenderResponse {
	tenders := make([]TenderResponse, len(tendersDB))

	for i, tenderDB := range tendersDB {
		t := ConvertTenderToResponse(tenderDB)

		tenders[i] = t
	}
	return tenders
}
