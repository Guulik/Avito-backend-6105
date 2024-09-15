package service

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strings"
	"zadanie-6105/internal/domain/model"
	sl "zadanie-6105/internal/lib/logger/slog"
)

type RepoBidProvider interface {
	GetBidsById(
		limit int32,
		offset int32,
		authorId string,
	) ([]model.BidDB, error)
	BidsForTender(
		tenderId string,
		limit int32,
		offset int32,
	) ([]model.BidDB, error)
	BidStatus(
		bidId string,
	) (string, error)
}
type RepoBidCreator interface {
	CreateBid(
		name string,
		description string,
		tenderId string,
		authorType string,
		authorId string,
	) (string, error)
}
type RepoBidEditor interface {
	UpdateBidStatus(
		bidId string,
		status string,
	) (string, error)
	EditBid(
		bidId string,
		name string,
		description string,
	) (string, error)
	RollbackBid(
		bidId string,
		version int32,
	) (string, error)
}
type RepoBidDecisionMaker interface {
	SubmitDecision(
		bidId string,
		responsibleId string,
	) error
	ApplyDecision(
		bidId string,
		decision string,
	) error
}
type RepoBidFeedbacker interface {
	Feedback(
		bidId string,
		feedback string,
	) error
	Reviews(
		authorUsername string,
		limit int32,
		offset int32,
	) ([]model.Feedback, error)
}

func (s *Service) CreateBid(name string, description string, tenderId string, authorType string, authorId string) (model.BidResponse, error) {
	const op = "Service.CreateBid"
	log := s.log.With(
		slog.String("op", op),
	)
	var (
		BidDB           model.BidDB
		BidResponse     model.BidResponse
		err             error
		bidId           string
		status          string
		orgranizationId string
	)
	// check status 401
	username, err := s.checkers.CheckCorporateById(authorId)
	if err != nil {
		return model.BidResponse{}, err
	}
	// check status 404
	_, err = s.checkers.CheckTender(tenderId)
	if err != nil {
		return model.BidResponse{}, err
	}
	// check status 403
	status, err = s.repoTenderProvider.Status(tenderId)
	if err != nil {
		return model.BidResponse{}, err
	}
	if strings.EqualFold(status, "Closed") || strings.EqualFold(status, "Created") {
		return model.BidResponse{}, echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("tender is not available"))
	}
	// если тип автора указан как организация,
	// то автором будет не ОТВТЕТСВЕННЫЙ за организацию,
	// а сама ОРГАНИЗАЦИЯ
	if strings.EqualFold(authorType, "Organization") {
		orgId, _ := s.checkers.CheckIdByName(username)
		relatedOrganizationId, _ := s.checkers.CheckResponsibility(username)
		if relatedOrganizationId != "" {
			orgranizationId = relatedOrganizationId
		} else {
			orgranizationId = orgId
		}
		bidId, err = s.repoBidCreator.CreateBid(name, description, tenderId, authorType, orgranizationId)
		if err != nil {
			return model.BidResponse{}, err
		}
	} else {
		bidId, err = s.repoBidCreator.CreateBid(name, description, tenderId, authorType, authorId)
		if err != nil {
			return model.BidResponse{}, err
		}
	}

	BidDB, err = s.checkers.CheckBid(bidId)
	if err != nil {
		return model.BidResponse{}, err
	}
	log.Info("Created Bid from DB:", BidDB)

	BidResponse = model.ConvertBidToResponse(BidDB)
	log.Info("Converted Bid to response:", BidResponse)

	return BidResponse, nil
}
func (s *Service) GetBidsByUser(limit int32, offset int32, username string) ([]model.BidResponse, error) {
	const op = "Service.GetBidsById"
	log := s.log.With(
		slog.String("op", op),
	)
	var (
		BidsDB         []model.BidDB
		BidsResponse   []model.BidResponse
		err            error
		organizationId string
	)
	//check status 401
	userId, err := s.checkers.CheckIdByName(username)
	if err != nil {
		return nil, err
	}
	BidsDB, err = s.repoBidProvider.GetBidsById(limit, offset, userId)
	if err != nil {
		return nil, err
	}
	// если пользователь -- ответственный за оргу, то добавятся ещё биды от организации
	organizationId, err = s.checkers.CheckResponsibility(username)
	log.Debug("org id", organizationId)
	if organizationId != "" {
		bidsByOrganization, err := s.repoBidProvider.GetBidsById(limit, offset, organizationId)
		BidsDB = append(BidsDB, bidsByOrganization...)
		if err != nil {
			return nil, err
		}
	}
	log.Info("Bids from DB:", BidsDB)

	BidsResponse = model.ConvertBids(BidsDB)
	log.Info("Converted Bids to response:", BidsResponse)

	return BidsResponse, nil
}

func (s *Service) BidsForTender(tenderId string, limit int32, offset int32, username string) ([]model.BidResponse, error) {
	const op = "Service.BidsForTender"
	log := s.log.With(
		slog.String("op", op),
	)
	var (
		BidsDB       []model.BidDB
		BidsResponse []model.BidResponse
		Tender       model.TenderDB
		err          error
	)
	//check status 401
	_, err = s.checkers.CheckIdByName(username)
	if err != nil {
		return nil, err
	}
	// check status 404
	Tender, err = s.checkers.CheckTender(tenderId)
	if err != nil {
		return nil, err
	}

	//check status 403
	if !strings.EqualFold(Tender.Status, "Published") {
		err = s.checkers.CheckResponsibleToTender(tenderId, username)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusForbidden, "you have no access to Tender -- it is not published")
		}
	}

	BidsDB, err = s.repoBidProvider.BidsForTender(tenderId, limit, offset)
	if err != nil {
		return nil, err
	}
	log.Info("Bids from DB:", BidsDB)

	//extract unavailable bid from response
	var exportBidsDB = make([]model.BidDB, 0, len(BidsDB))
	for _, bid := range BidsDB {
		status := bid.Status
		err = s.checkers.CheckBidAuthorByUsername(bid.Id, username)
		// автор может увидеть только своё, а ответственный за оргу только опубликованные и отмененные
		if !strings.EqualFold(status, "Created") || err == nil {
			exportBidsDB = append(exportBidsDB, bid)
		}
	}

	BidsResponse = model.ConvertBids(exportBidsDB)
	log.Info("Converted Bids to response:", BidsResponse)

	return BidsResponse, nil
}

func (s *Service) BidStatus(bidId string, username string) (string, error) {
	const op = "Service.BidStatus"
	log := s.log.With(
		slog.String("op", op),
	)
	var (
		status string
		err    error
	)
	//check status 404
	_, err = s.checkers.CheckBid(bidId)
	if err != nil {
		return "", err
	}
	// check status 401
	_, err = s.checkers.CheckIdByName(username)
	if err != nil {
		return "", err
	}

	//check status 403
	status, err = s.repoBidProvider.BidStatus(bidId)
	if err != nil {
		return "nil", err
	}
	log.Info("Status for bid from DB:", status)

	if !strings.EqualFold(status, "Published") {
		err = s.checkers.CheckBidAuthorByUsername(bidId, username)
		if err != nil {
			return "", err
		}
	}

	return status, nil
}

func (s *Service) UpdateBidStatus(bidId string, status string, username string) (model.BidResponse, error) {
	const op = "Service.UpdateBidStatus"
	log := s.log.With(
		slog.String("op", op),
	)
	var (
		BidDB       model.BidDB
		BidResponse model.BidResponse
		err         error
	)
	//check status 404
	_, err = s.checkers.CheckBid(bidId)
	if err != nil {
		return model.BidResponse{}, err
	}
	// check status 401
	_, err = s.checkers.CheckIdByName(username)
	if err != nil {
		return model.BidResponse{}, err
	}
	// check status 403
	// если бид отменен или по нему принято решение, статус нельзя изменить
	err = s.checkers.CheckBidAvailability(bidId)
	if err != nil {
		return model.BidResponse{}, err
	}
	//check status 403
	//author and responsilbe
	err = s.checkers.CheckStatusForbiddenForBid(bidId, username)
	if err != nil {
		return model.BidResponse{}, err
	}

	_, err = s.repoBidEditor.UpdateBidStatus(bidId, status)
	if err != nil {
		return model.BidResponse{}, err
	}

	BidDB, err = s.checkers.CheckBid(bidId)
	if err != nil {
		return model.BidResponse{}, err
	}
	log.Info("Created Bid from DB:", BidDB)

	BidResponse = model.ConvertBidToResponse(BidDB)
	log.Info("Converted Bid to response:", BidResponse)

	return BidResponse, nil
}

func (s *Service) EditBid(bidId string, username string, name string, description string) (model.BidResponse, error) {
	const op = "Service.EditBid"
	log := s.log.With(
		slog.String("op", op),
	)
	var (
		BidDB       model.BidDB
		BidResponse model.BidResponse
		err         error
	)

	//check status 404
	_, err = s.checkers.CheckBid(bidId)
	if err != nil {
		return model.BidResponse{}, err
	}
	// check status 401
	_, err = s.checkers.CheckIdByName(username)
	if err != nil {
		return model.BidResponse{}, err
	}
	// check status 403
	// decision taken or bid canceled
	err = s.checkers.CheckBidAvailability(bidId)
	if err != nil {
		return model.BidResponse{}, err
	}
	//check status 403
	//only author or responsible can edit
	err = s.checkers.CheckStatusForbiddenForBid(bidId, username)
	if err != nil {
		return model.BidResponse{}, err
	}

	_, err = s.repoBidEditor.EditBid(bidId, name, description)
	if err != nil {
		return model.BidResponse{}, err
	}

	BidDB, err = s.checkers.CheckBid(bidId)
	if err != nil {
		return model.BidResponse{}, err
	}
	log.Info("Created Bid from DB:", BidDB)

	BidResponse = model.ConvertBidToResponse(BidDB)
	log.Info("Converted Bid to response:", BidResponse)

	return BidResponse, nil
}

func (s *Service) SubmitDecision(bidId string, decision string, username string) (model.BidResponse, error) {
	const op = "Service.SubmitDecision"
	log := s.log.With(
		slog.String("op", op),
	)
	var (
		BidDB           model.BidDB
		BidResponse     model.BidResponse
		err             error
		decisionCount   int
		kworum          int
		organizationId  string
		relatedTenderId string
		responsibleId   string
		otherBids       []model.BidDB
	)
	//check status 404
	_, err = s.checkers.CheckBid(bidId)
	if err != nil {
		log.Error("failed to check bid", sl.Err(err))
		return model.BidResponse{}, err
	}

	// check status 401
	_, err = s.checkers.CheckIdByName(username)
	if err != nil {
		return model.BidResponse{}, err
	}
	// check status 403
	//if bid just created, organization cannot submit decision
	status, err := s.repoBidProvider.BidStatus(bidId)
	if err != nil {
		return model.BidResponse{}, err
	}
	if strings.EqualFold(status, "Created") {
		return model.BidResponse{}, echo.NewHTTPError(http.StatusForbidden, "organization have no access to just created bids")
	}
	// decision taken or bid canceled
	err = s.checkers.CheckBidAvailability(bidId)
	if err != nil {
		return model.BidResponse{}, err
	}
	//check status 403
	//check organization ownership to tender
	organizationId, err = s.checkers.CheckResponsibility(username)
	if err != nil {
		return model.BidResponse{}, err
	}
	relatedTenderId, err = s.checkers.CheckBidTenderOwner(bidId, organizationId)
	if err != nil {
		return model.BidResponse{}, err
	}
	err = s.checkers.CheckSameSubmitter(bidId, username)
	if err != nil {
		return model.BidResponse{}, err
	}

	// if decision = reject, apply decision without submitting
	if strings.EqualFold(decision, "Rejected") {
		err = s.repoBidDecisionMaker.ApplyDecision(bidId, decision)
		if err != nil {
			return model.BidResponse{}, err
		}
		_, err = s.repoBidEditor.UpdateBidStatus(bidId, "Canceled")
		if err != nil {
			return model.BidResponse{}, err
		}
	} else {
		responsibleId, err = s.checkers.CheckIdByName(username)
		if err != nil {
			return model.BidResponse{}, err
		}
		err = s.repoBidDecisionMaker.SubmitDecision(bidId, responsibleId)
		if err != nil {
			return model.BidResponse{}, err
		}
		log.Info("Bid from DB:", BidDB)

		decisionCount, err = s.checkers.CheckBidDecisionCount(bidId)
		if err != nil {
			return model.BidResponse{}, err
		}

		organizationId, err = s.checkers.CheckResponsibility(username)
		if err != nil {
			return model.BidResponse{}, err
		}
		kworum, err = s.checkers.CheckResponsibleCount(organizationId)
		if err != nil {
			return model.BidResponse{}, err
		}
		kworum = min(3, kworum)
		if decisionCount >= kworum {
			err = s.repoBidDecisionMaker.ApplyDecision(bidId, decision)
			if err != nil {
				return model.BidResponse{}, err
			}
			_, err = s.repoTenderEditor.ChangeTenderStatus(relatedTenderId, "Closed")
			if err != nil {
				return model.BidResponse{}, err
			}

			otherBids, err = s.repoBidProvider.BidsForTender(relatedTenderId, 0, 0)
			if err != nil {
				return model.BidResponse{}, err
			}
			for _, bid := range otherBids {
				if bid.Id != bidId {
					err = s.repoBidDecisionMaker.ApplyDecision(bid.Id, "Rejected")
					if err != nil {
						return model.BidResponse{}, err
					}
					_, err = s.repoBidEditor.UpdateBidStatus(bid.Id, "Canceled")
					if err != nil {
						return model.BidResponse{}, err
					}
				}
			}
		}
	}

	//Get bid
	BidDB, err = s.checkers.CheckBid(bidId)
	if err != nil {
		return model.BidResponse{}, err
	}

	BidResponse = model.ConvertBidToResponse(BidDB)
	log.Info("Converted Bid to response:", BidResponse)

	return BidResponse, nil
}

func (s *Service) Feedback(bidId string, feedback string, username string) (model.BidResponse, error) {
	const op = "Service.Feedback"
	log := s.log.With(
		slog.String("op", op),
	)
	var (
		BidDB          model.BidDB
		BidResponse    model.BidResponse
		err            error
		organizationId string
	)

	//check status 404
	_, err = s.checkers.CheckBid(bidId)
	if err != nil {
		return model.BidResponse{}, err
	}
	// check status 401
	_, err = s.checkers.CheckIdByName(username)
	if err != nil {
		return model.BidResponse{}, err
	}
	//check status 403
	//if bid just created, organization cannot submit decision
	status, err := s.repoBidProvider.BidStatus(bidId)
	if err != nil {
		return model.BidResponse{}, err
	}
	if strings.EqualFold(status, "Created") {
		return model.BidResponse{}, echo.NewHTTPError(http.StatusForbidden, "organization have no access to just created bids")
	}
	//author cannot leave feedback on his own bid
	errNotAuthor := s.checkers.CheckBidAuthorByUsername(bidId, username)
	if errNotAuthor == nil {
		return model.BidResponse{}, echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("author cannot leave feedback on his own bid"))
	}
	//check organization ownership to tender
	organizationId, err = s.checkers.CheckResponsibility(username)
	if err != nil {
		return model.BidResponse{}, err
	}
	_, err = s.checkers.CheckBidTenderOwner(bidId, organizationId)
	if err != nil {
		return model.BidResponse{}, err
	}

	err = s.repoBidFeedbacker.Feedback(bidId, feedback)
	if err != nil {
		return model.BidResponse{}, err
	}

	BidDB, err = s.checkers.CheckBid(bidId)
	if err != nil {
		return model.BidResponse{}, err
	}
	log.Info("Bid from DB:", BidDB)
	BidResponse = model.ConvertBidToResponse(BidDB)
	log.Info("Converted Bid to response:", BidResponse)

	return BidResponse, nil
}

func (s *Service) RollbackBid(bidId string, version int32, username string) (model.BidResponse, error) {
	const op = "Service.RollbackBid"
	log := s.log.With(
		slog.String("op", op),
	)
	var (
		BidDB       model.BidDB
		BidResponse model.BidResponse
		err         error
	)

	//check status 404
	_, err = s.checkers.CheckBid(bidId)
	if err != nil {
		log.Debug("bid not found", sl.Err(err))
		return model.BidResponse{}, err
	}
	// check status 401
	_, err = s.checkers.CheckIdByName(username)
	if err != nil {
		log.Debug("user not found", sl.Err(err))
		return model.BidResponse{}, err
	}
	// check status 403
	// decision taken or bid canceled
	err = s.checkers.CheckBidAvailability(bidId)
	if err != nil {
		return model.BidResponse{}, err
	}
	// check status 404
	err = s.checkers.CheckBidVersion(bidId, version)
	if err != nil {
		return model.BidResponse{}, err
	}
	//check status 403
	err = s.checkers.CheckStatusForbiddenForBid(bidId, username)
	if err != nil {
		return model.BidResponse{}, err
	}

	_, err = s.repoBidEditor.RollbackBid(bidId, version)
	if err != nil {
		return model.BidResponse{}, err
	}
	log.Info("Bid from DB:", BidDB)

	BidDB, err = s.checkers.CheckBid(bidId)
	if err != nil {
		return model.BidResponse{}, err
	}

	BidResponse = model.ConvertBidToResponse(BidDB)
	log.Info("Converted Bid to response:", BidResponse)

	return BidResponse, nil
}

func (s *Service) Reviews(tenderId string, authorUsername string, requesterUsername string, limit int32, offset int32) ([]model.Feedback, error) {
	const op = "Service.Reviews"
	log := s.log.With(
		slog.String("op", op),
	)
	var (
		Feedbacks  []model.Feedback
		err        error
		BidsByUser []model.BidDB
		authorId   string
	)

	//check status 404
	_, err = s.checkers.CheckTender(tenderId)
	if err != nil {
		log.Debug("bid not found", sl.Err(err))
		return nil, err
	}
	// check status 401
	authorId, err = s.checkers.CheckIdByName(authorUsername)
	if err != nil {
		log.Debug("author not found", sl.Err(err))
		return nil, err
	}
	log.Debug("authorId", authorId)
	// check status 401
	_, err = s.checkers.CheckIdByName(requesterUsername)
	if err != nil {
		log.Debug("requester not found", sl.Err(err))
		return nil, err
	}
	// check status 403
	// автор не может посмотреть отзывы на свои предложения.
	// ну да, странно, но в задании написано, что ток ответственный может...
	// а ещё теоретически автор может не знать айди тендера.
	err = s.checkers.CheckResponsibleToTender(tenderId, requesterUsername)
	if err != nil {
		return nil, err
	}
	// check 404
	BidsByUser, err = s.repoBidProvider.GetBidsById(0, 0, authorId)
	if err != nil {
		return nil, err
	}
	if len(BidsByUser) == 0 {
		return nil, echo.NewHTTPError(http.StatusNotFound, "no bids by this user")
	}
	log.Debug("user bids", BidsByUser)
	//это уже скорее костыль, но у меня нет времени...
	atLeastOneBid := false
	for _, bid := range BidsByUser {
		if strings.EqualFold(bid.TenderId, tenderId) {
			atLeastOneBid = true
			break
		}
	}
	if !atLeastOneBid {
		return nil, echo.NewHTTPError(http.StatusNotFound, "no bids for this tender by the specified author")
	}
	//extract just created Bids, because organization cannot leave feedback for it
	var bidsToCheckFeedback = make([]model.BidDB, 0, len(BidsByUser))
	for _, bid := range BidsByUser {
		status := bid.Status
		if !strings.EqualFold(status, "Created") {
			bidsToCheckFeedback = append(bidsToCheckFeedback, bid)
		}
	}

	Feedbacks, err = s.repoBidFeedbacker.Reviews(authorUsername, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(Feedbacks) == 0 {
		return nil, echo.NewHTTPError(http.StatusNotFound, "no feedbacks for bids by this author")
	}
	log.Info("Feedbacks from DB:", Feedbacks)

	return Feedbacks, nil
}
