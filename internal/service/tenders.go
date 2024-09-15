package service

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strings"
	"zadanie-6105/internal/domain/model"
)

type RepoTenderProvider interface {
	Tenders(
		limit int32,
		offset int32,
		serviceTypes []string,
	) ([]model.TenderDB, error)
	TendersByUser(
		limit int32,
		offset int32,
		username string,
	) ([]model.TenderDB, error)
	Status(
		id string,
	) (string, error)
}
type RepoTenderCreator interface {
	CreateTender(
		name string,
		description string,
		serviceType string,
		organizationId string,
		creatorUsername string,
	) (string, error)
}
type RepoTenderEditor interface {
	ChangeTenderStatus(
		id string,
		status string,
	) (string, error)
	EditTender(
		id string,
		name string,
		description string,
		serviceType string,
	) (string, error)
	RollbackTender(
		id string,
		version int32,
	) (string, error)
}

func (s *Service) Tenders(limit int32, offset int32, serviceType []string) ([]model.TenderResponse, error) {
	const op = "Service.Tenders"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		Tenders         []model.TenderDB
		TendersResponse []model.TenderResponse
		err             error
	)

	Tenders, err = s.repoTenderProvider.Tenders(limit, offset, serviceType)
	if err != nil {
		return nil, err
	}
	log.Info("Tenders from DB:", Tenders)

	var tendersToExport = make([]model.TenderDB, 0, len(Tenders))
	for _, tender := range Tenders {
		if !(strings.EqualFold(tender.Status, "Created") || strings.EqualFold(tender.Status, "Closed")) {
			tendersToExport = append(tendersToExport, tender)
		}
	}
	log.Info("Tenders to export:", tendersToExport)
	TendersResponse = model.ConvertTenders(tendersToExport)
	log.Info("Converted Tenders to response:", TendersResponse)

	return TendersResponse, nil
}

func (s *Service) CreateTender(name string, description string, serviceType string, organizationId string, creatorUsername string) (model.TenderResponse, error) {
	const op = "Service.CreateTender"
	log := s.log.With(
		slog.String("op", op),
	)
	var (
		TenderDB       model.TenderDB
		TenderResponse model.TenderResponse
		err            error
		orgId          string
	)

	// check status 401
	_, err = s.checkers.CheckIdByName(creatorUsername)
	if err != nil {
		return model.TenderResponse{}, err
	}
	// check status 401
	_, err = s.checkers.CheckCorporateById(organizationId)
	if err != nil {
		return model.TenderResponse{}, err
	}
	// check status 403
	orgId, err = s.checkers.CheckResponsibility(creatorUsername)
	if err != nil || orgId != organizationId {
		return model.TenderResponse{}, echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("user does not responsible for this organization"))
	}

	tenderId, err := s.repoTenderCreator.CreateTender(name, description, serviceType, organizationId, creatorUsername)
	if err != nil {
		return model.TenderResponse{}, err
	}

	tenderDB, err := s.checkers.CheckTender(tenderId)
	if err != nil {
		return model.TenderResponse{}, err
	}
	log.Info("Created tender from DB:", TenderDB)

	TenderResponse = model.ConvertTenderToResponse(tenderDB)
	log.Info("Converted Tender to response:", TenderResponse)

	return TenderResponse, nil
}

func (s *Service) GetTenderByUser(limit int32, offset int32, username string) ([]model.TenderResponse, error) {
	const op = "Service.GetTenderByUser"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		TendersDB       []model.TenderDB
		TendersResponse []model.TenderResponse
		err             error
	)
	// check status 401
	_, err = s.checkers.CheckIdByName(username)
	if err != nil {
		return nil, err
	}
	//check status 403
	_, err = s.checkers.CheckResponsibility(username)
	if err != nil {
		return nil, err
	}

	TendersDB, err = s.repoTenderProvider.TendersByUser(limit, offset, username)
	if err != nil {
		return nil, err
	}
	log.Info("Tenders from DB:", TendersDB)

	TendersResponse = model.ConvertTenders(TendersDB)
	log.Info("Converted Tenders to response:", TendersResponse)

	return TendersResponse, nil
}

func (s *Service) TenderStatus(tenderId string, username string) (string, error) {
	const op = "Service.TenderStatus"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		status string
		err    error
	)
	// check status 401
	_, err = s.checkers.CheckIdByName(username)
	if err != nil {
		return "", err
	}
	// check status 404
	_, err = s.checkers.CheckTender(tenderId)
	if err != nil {
		return "", err
	}

	status, err = s.repoTenderProvider.Status(tenderId)
	if err != nil {
		return "", err
	}
	log.Info("Status tender from DB:", status)

	// check status 403
	if strings.EqualFold(status, "Created") || strings.EqualFold(status, "Closed") {
		err = s.checkers.CheckResponsibleToTender(tenderId, username)
		if err != nil {
			return "", err
		}
	}

	return status, nil
}

func (s *Service) ChangeTenderStatus(tenderId string, status string, username string) (model.TenderResponse, error) {
	const op = "Service.ChangeTenderStatus"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		TenderDB       model.TenderDB
		TenderResponse model.TenderResponse
		err            error
		relatedBids    []model.BidDB
	)
	// check status 401
	_, err = s.checkers.CheckIdByName(username)
	if err != nil {
		return model.TenderResponse{}, err
	}
	// check status 404
	TenderDB, err = s.checkers.CheckTender(tenderId)
	if err != nil {
		return model.TenderResponse{}, err
	}
	// check status 403
	err = s.checkers.CheckResponsibleToTender(tenderId, username)
	if err != nil {
		return model.TenderResponse{}, err
	}
	// check status 403
	if strings.EqualFold(TenderDB.Status, "Closed") {
		return model.TenderResponse{}, echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("tender is closed"))
	}
	_, err = s.repoTenderEditor.ChangeTenderStatus(tenderId, status)
	if err != nil {
		return model.TenderResponse{}, err
	}

	//reject all related bids
	relatedBids, err = s.repoBidProvider.BidsForTender(tenderId, 0, 0)
	if err != nil {
		return model.TenderResponse{}, err
	}
	for _, bid := range relatedBids {
		_, err = s.repoBidEditor.UpdateBidStatus(bid.Id, "Canceled")
		if err != nil {
			return model.TenderResponse{}, err
		}
		err = s.repoBidDecisionMaker.ApplyDecision(bid.Id, "Rejected")
		if err != nil {
			return model.TenderResponse{}, err
		}
	}

	TenderDB, err = s.checkers.CheckTender(tenderId)
	if err != nil {
		return model.TenderResponse{}, err
	}
	log.Info("Tender with changed status from DB:", TenderDB)

	TenderResponse = model.ConvertTenderToResponse(TenderDB)
	log.Info("Converted Tender to response:", TenderResponse)

	return TenderResponse, nil
}

func (s *Service) EditTender(tenderId string, username string, name string, description string, serviceType string) (model.TenderResponse, error) {
	const op = "Service.EditTender"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		TenderDB       model.TenderDB
		TenderResponse model.TenderResponse
		err            error
	)
	// check status 401
	_, err = s.checkers.CheckIdByName(username)
	if err != nil {
		return model.TenderResponse{}, err
	}
	// check status 404
	TenderDB, err = s.checkers.CheckTender(tenderId)
	if err != nil {
		return model.TenderResponse{}, err
	}
	// check status 403
	err = s.checkers.CheckResponsibleToTender(tenderId, username)
	if err != nil {
		return model.TenderResponse{}, err
	}
	// check status 403
	if strings.EqualFold(TenderDB.Status, "closed") {
		return model.TenderResponse{}, echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("tender is closed"))
	}

	_, err = s.repoTenderEditor.EditTender(tenderId, name, description, serviceType)
	if err != nil {
		return model.TenderResponse{}, err
	}
	log.Info("Edited tender from DB:", TenderDB)

	TenderDB, err = s.checkers.CheckTender(tenderId)
	if err != nil {
		return model.TenderResponse{}, err
	}
	log.Info("Created tender from DB:", TenderDB)

	TenderResponse = model.ConvertTenderToResponse(TenderDB)
	log.Info("Converted Tender to response:", TenderResponse)

	return TenderResponse, nil
}

func (s *Service) RollbackTender(tenderId string, version int32, username string) (model.TenderResponse, error) {
	const op = "Service.RollbackTender"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		TenderDB       model.TenderDB
		TenderResponse model.TenderResponse
		err            error
	)

	// check status 401
	_, err = s.checkers.CheckIdByName(username)
	if err != nil {
		return model.TenderResponse{}, err
	}
	// check status 404
	TenderDB, err = s.checkers.CheckTender(tenderId)
	if err != nil {
		return model.TenderResponse{}, err
	}
	// check status 404
	err = s.checkers.CheckTenderVersion(tenderId, version)
	if err != nil {
		return model.TenderResponse{}, err
	}
	// check status 403
	err = s.checkers.CheckResponsibleToTender(tenderId, username)
	if err != nil {
		return model.TenderResponse{}, err
	}
	// check status 403
	if strings.EqualFold(TenderDB.Status, "Closed") {
		return model.TenderResponse{}, echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("tender is closed"))
	}

	_, err = s.repoTenderEditor.RollbackTender(tenderId, version)
	if err != nil {
		return model.TenderResponse{}, err
	}
	log.Info("Rollback tender from DB:", TenderDB)

	TenderDB, err = s.checkers.CheckTender(tenderId)
	if err != nil {
		return model.TenderResponse{}, err
	}
	log.Info("Created tender from DB:", TenderDB)

	TenderResponse = model.ConvertTenderToResponse(TenderDB)
	log.Info("Converted Tender to response:", TenderResponse)

	return TenderResponse, nil
}
