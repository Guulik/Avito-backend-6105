package postgres

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strings"
	"zadanie-6105/internal/domain/model"
	sl "zadanie-6105/internal/lib/logger/slog"
)

func (s *Storage) CheckResponsibleToTender(tenderId string, username string) error {
	const op = "Support.CheckResponsibleToTender"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		selectQuery = `
		SELECT e.username
		FROM employee e
		JOIN organization_responsible resp ON resp.user_id = e.id
		JOIN organization o ON resp.organization_id = o.id
		JOIN tender t ON t.organization_id = o.id
		WHERE t.id = $1::uuid AND e.username = $2;
`
		selectValues = []any{
			tenderId, username,
		}
		id string
	)

	row := s.db.QueryRow(selectQuery, selectValues...)
	err := row.Scan(&id)
	if err != nil {
		log.Info("user have no access to tender ", err)
		return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("user have no access to tender "))
	}

	return nil
}

func (s *Storage) CheckIdByName(username string) (string, error) {
	const op = "Support.CheckIdByName"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		selectQuery = `
		SELECT id
		FROM employee
		WHERE username = $1
		
		UNION
		
		SELECT id
		FROM organization
		WHERE name = $1;
`
		selectValues = []any{
			username,
		}
		id string
	)

	row := s.db.QueryRow(selectQuery, selectValues...)
	err := row.Scan(&id)
	if err != nil {
		log.Error("user not found", sl.Err(err))
		return "", echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("user not found"))
	}

	return id, nil
}

func (s *Storage) CheckCorporateById(userId string) (string, error) {
	const op = "Support.CheckCorporateById"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		selectQuery = `
		SELECT username
		FROM employee
		WHERE id = $1
		
		UNION
		
		SELECT name
		FROM organization
		WHERE id = $1;
`
		selectValues = []any{
			userId,
		}
		username string
	)

	row := s.db.QueryRow(selectQuery, selectValues...)
	err := row.Scan(&username)
	if err != nil {
		log.Error("user or organization not found", sl.Err(err))
		return "", echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("user or organization not found"))
	}

	return username, nil
}

func (s *Storage) CheckResponsibility(username string) (string, error) {
	const op = "Support.CheckResponsibility"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		selectQuery = `
		SELECT resp.organization_id
		FROM employee e
		JOIN organization_responsible resp ON resp.user_id = e.id
		WHERE e.username = $1;
`
		selectValues = []any{
			username,
		}
		organisationId string
	)

	row := s.db.QueryRow(selectQuery, selectValues...)
	err := row.Scan(&organisationId)
	if err != nil {
		log.Info("user does not responsible for any organization", sl.Err(err))
		return "", echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("user does not responsible for any organization"))
	}

	return organisationId, nil
}

func (s *Storage) CheckResponsibleCount(organizationId string) (int, error) {
	const op = "Support.CheckResponsibleCount"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		selectQuery = `
		SELECT COUNT(*) AS responsible_count
		FROM organization_responsible
		WHERE organization_id = $1;
`
		selectValues = []any{
			organizationId,
		}
		responsibleCount int
	)

	row := s.db.QueryRow(selectQuery, selectValues...)
	err := row.Scan(&responsibleCount)
	if err != nil {
		log.Info("failed to get responsible count", sl.Err(err))
		return -1, echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("failed to get responsible count"))
	}

	return responsibleCount, nil
}

func (s *Storage) CheckTender(tenderId string) (model.TenderDB, error) {
	const op = "Support.CheckTender"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		selectQuery = `
		SELECT id, name, description, CAST(serviceType AS text) AS serviceType,
		       CAST(status AS text), organization_id,creator_username, version, created_at
		FROM tender
		WHERE id = $1::uuid;
`
		selectValues = []any{
			tenderId,
		}
		tender model.TenderDB
	)

	row := s.db.QueryRow(selectQuery, selectValues...)
	err := row.Scan(&tender.Id,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationId,
		&tender.CreatorUsername,
		&tender.Version,
		&tender.CreatedAt)
	if err != nil {
		log.Error("tender not found", sl.Err(err))
		return model.TenderDB{}, echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("tender not found"))
	}

	return tender, nil
}

func (s *Storage) CheckTenderVersion(tenderId string, version int32) error {
	const op = "Support.CheckTenderVersion"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		selectQuery = `
		SELECT version
		FROM tender_version
		WHERE tender_id = $1::uuid 
		AND version = $2;
`
		selectValues = []any{
			tenderId, version,
		}
		verstion int32 = -1
	)

	row := s.db.QueryRow(selectQuery, selectValues...)
	err := row.Scan(&verstion)
	if err != nil {
		log.Error("version not found", sl.Err(err))
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("no such tender version"))
	}

	return nil
}

func (s *Storage) CheckBid(bidId string) (model.BidDB, error) {
	const op = "Support.CheckBid"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		selectQuery = `
		SELECT id, name, description,  COALESCE(CAST(decision AS text),''),
		       CAST(status AS text), tenderId,authorType,authorId, version, createdAt
		FROM bid
		WHERE id = $1::uuid;
`
		selectValues = []any{
			bidId,
		}
		bid model.BidDB
	)

	row := s.db.QueryRow(selectQuery, selectValues...)
	err := row.Scan(&bid.Id,
		&bid.Name,
		&bid.Description,
		&bid.Decision,
		&bid.Status,
		&bid.TenderId,
		&bid.AuthorType,
		&bid.AuthorId,
		&bid.Version,
		&bid.CreatedAt)
	if err != nil {
		log.Error("bid not found", sl.Err(err))
		return model.BidDB{}, echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("bid not found"))
	}
	return bid, nil
}

func (s *Storage) CheckBidVersion(bidId string, version int32) error {
	const op = "Support.CheckBidVersion"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		selectQuery = `
		SELECT version
		FROM bid_version
		WHERE bid_id = $1::uuid 
		AND version = $2;
`
		selectValues = []any{
			bidId, version,
		}
		verstion int32 = -1
	)

	row := s.db.QueryRow(selectQuery, selectValues...)
	err := row.Scan(&verstion)
	if err != nil {
		log.Error("version not found", sl.Err(err))
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("no such bid version"))
	}

	return nil
}

func (s *Storage) CheckBidAuthorByUsername(bidId string, username string) error {
	const op = "Support.CheckBidAuthorByUsername"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		authorQuery = `
		SELECT e.username
		FROM bid b 
		JOIN employee e ON b.authorId = e.id
		WHERE b.id = $1::uuid AND e.username = $2;
`
		authorValues = []any{
			bidId, username,
		}
		name string
	)
	row := s.db.QueryRow(authorQuery, authorValues...)
	errAuthor := row.Scan(&name)

	if errAuthor != nil {
		log.Info("user have no access to bid ", sl.Err(errAuthor))
		return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("user have no access to bid "))
	}
	return nil
}

func (s *Storage) CheckBidTenderOwner(bidId string, organizationId string) (string, error) {
	const op = "Support.CheckBidTenderOwner"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		tenderOwnerQuery = `
		SELECT t.id
		FROM bid b 
		JOIN tender t ON t.id = b.tenderid
		JOIN organization o ON o.id = t.organization_id
		WHERE b.id = $1::uuid AND o.id = $2::uuid;
`
		tenderOwnerValues = []any{
			bidId, organizationId,
		}
		tenderId string
	)

	row := s.db.QueryRow(tenderOwnerQuery, tenderOwnerValues...)
	errTenderOwner := row.Scan(&tenderId)

	if errTenderOwner != nil {
		log.Info("organization have no access to bid ", sl.Err(errTenderOwner))
		return "", echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("organization have no access to bid "))
	}

	return tenderId, nil
}

func (s *Storage) CheckAccessToBidByOrganizationId(bidId string, organizationId string) error {
	const op = "Support.CheckAccessToBidByOrganizationId"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		authorQuery = `
		SELECT o.name
		FROM bid b 
		JOIN organization o ON b.authorId = o.id
		WHERE b.id = $1::uuid AND o.id = $2::uuid;
`
		authorValues = []any{
			bidId, organizationId,
		}
		name string
	)

	row := s.db.QueryRow(authorQuery, authorValues...)
	errAuthor := row.Scan(&name)

	_, errTenderOwner := s.CheckBidTenderOwner(bidId, organizationId)
	if errAuthor != nil || errTenderOwner != nil {
		log.Info("organization have no access to bid ", sl.Err(errAuthor))
		return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("organization have no access to bid "))
	}

	return nil
}

func (s *Storage) CheckStatusForbiddenForBid(bidId string, username string) error {
	const op = "Support.CheckStatusForbiddenForBid"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		errNotAuthor      error
		errNotResponsible error
	)
	//check status 403 (not author)
	errNotAuthor = s.CheckBidAuthorByUsername(bidId, username)
	//check status 403 (not responsible)
	organizationId, errNotResponsible := s.CheckResponsibility(username)
	if organizationId != "" {
		log.Debug("orgId", organizationId)
		errNotResponsible = s.CheckAccessToBidByOrganizationId(bidId, organizationId)
	}

	if errNotAuthor != nil {
		log.Debug("errNotAuthor", sl.Err(errNotAuthor))
	}
	if errNotResponsible != nil {
		log.Debug("errNotResponsible", sl.Err(errNotResponsible))
	}
	if errNotAuthor != nil && errNotResponsible != nil {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("user have no access to bid "))
	}
	return nil
}

func (s *Storage) CheckBidDecisionCount(bidId string) (int, error) {
	const op = "Support.CheckBidDecisionCount"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		countQuery = `
		SELECT COUNT(DISTINCT bd.responsible) AS decision_count
		FROM bid_approval bd
		JOIN organization_responsible resp ON bd.responsible = resp.user_id
		JOIN bid b ON bd.bid_id = b.id
		JOIN tender t ON b.tenderId = t.id
		WHERE b.id = $1::uuid
		AND resp.organization_id = t.organization_id 
		GROUP BY b.id;
`
		countValues = []any{
			bidId,
		}
		count int
	)
	row := s.db.QueryRow(countQuery, countValues...)
	err := row.Scan(&count)

	if err != nil {
		log.Error("failed to get decision count..", sl.Err(err))
		return -1, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get decision count.."))
	}
	return count, nil
}

func (s *Storage) CheckSameSubmitter(bidId string, username string) error {
	const op = "Support.CheckSameSubmitter"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		countQuery = `
		SELECT COUNT(*)
		FROM bid_approval ba
		JOIN employee e ON ba.responsible = e.id
		WHERE ba.bid_id = $1::uuid
		AND e.username = $2;
`
		countValues = []any{
			bidId, username,
		}
		count int
	)
	row := s.db.QueryRow(countQuery, countValues...)
	err := row.Scan(&count)
	if err != nil {
		log.Error("failed to get decision count..", sl.Err(err))
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get decision count.."))
	}
	if count > 0 {
		log.Info("user already sent decision")
		return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("user already sent decision"))
	}
	return nil
}

func (s *Storage) CheckBidAvailability(bidId string) error {
	const op = "Support.CheckBidAvailability"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		decisionQuery = `
		SELECT coalesce(decision::text, '') 
		FROM bid
		WHERE id = $1::uuid;
`
		canceledQuery = `
		SELECT status::text
		FROM bid
		WHERE id = $1::uuid;
`
		Values = []any{
			bidId,
		}
		decision         string
		status           string
		errDecisionTaken error
		errBidCanceled   error
	)
	row := s.db.QueryRow(decisionQuery, Values...)
	err := row.Scan(&decision)
	if err != nil {
		log.Error("failed to get decision", sl.Err(err))
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get decision count.."))
	}
	if decision != "" {
		errDecisionTaken = fmt.Errorf("decision on bid already taken")
	}

	row = s.db.QueryRow(canceledQuery, Values...)
	err = row.Scan(&status)
	if err != nil {
		log.Error("failed to get status", sl.Err(err))
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get decision count.."))
	}
	if strings.EqualFold(status, "Canceled") {
		errBidCanceled = fmt.Errorf("bid is canceled")
	}

	if errDecisionTaken != nil || errBidCanceled != nil {
		log.Info("bid is locked")
		return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("bid is locked"))
	}

	return nil
}

func (s *Storage) CheckBidCanceled(bidId string) error {
	const op = "Support.CheckBidCanceled"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		canceledQuery = `
		SELECT status::text
		FROM bid
		WHERE id = $1::uuid;
`
		Values = []any{
			bidId,
		}

		status string

		errBidCanceled error
	)

	row := s.db.QueryRow(canceledQuery, Values...)
	err := row.Scan(&status)
	if err != nil {
		log.Error("failed to get status", sl.Err(err))
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get decision count.."))
	}
	if strings.EqualFold(status, "Canceled") {
		errBidCanceled = fmt.Errorf("bid is canceled")
	}

	if errBidCanceled != nil {
		log.Info("bid is canceled")
		return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("bid is locked"))
	}

	return nil
}
