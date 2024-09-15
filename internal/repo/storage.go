package repo

import (
	"zadanie-6105/internal/domain/model"
)

type Checkers interface {
	CheckResponsibleToTender(
		tenderId string,
		username string,
	) error
	CheckIdByName(
		username string,
	) (string, error)
	CheckCorporateById(
		userId string,
	) (string, error)
	CheckResponsibility(
		username string,
	) (string, error)
	CheckTender(
		tenderId string,
	) (model.TenderDB, error)
	CheckTenderVersion(
		tenderId string,
		version int32,
	) error
	CheckBid(
		bidId string,
	) (model.BidDB, error)
	CheckBidVersion(
		bidId string,
		version int32,
	) error
	CheckBidTenderOwner(
		bidId string,
		organizationId string,
	) (string, error)
	CheckAccessToBidByOrganizationId(
		bidId string,
		organizationId string,
	) error
	CheckBidAuthorByUsername(
		bidId string,
		username string,
	) error
	CheckStatusForbiddenForBid(
		bidId string,
		username string,
	) error
	CheckBidDecisionCount(
		bidId string,
	) (int, error)
	CheckResponsibleCount(
		organizationId string,
	) (int, error)
	CheckSameSubmitter(
		bidId string,
		username string,
	) error
	CheckBidAvailability(
		bidId string,
	) error
	CheckBidCanceled(
		bidId string,
	) error
}
