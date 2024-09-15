package postgres

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	model2 "zadanie-6105/internal/domain/model"
	sl "zadanie-6105/internal/lib/logger/slog"
)

func (s *Storage) CreateBid(name string, description string, tenderId string, authorType string, authorId string) (string, error) {
	const op = "Repo.CreateBid"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		insertQuery = `
        INSERT INTO bid (
                         name, description, status,
                         tenderId, authorType, authorId, version,
                         createdAt)
        VALUES ($1, $2, 'Created'::bid_status, $3, $4, $5, 1, CURRENT_TIMESTAMP)
        RETURNING id;
    `
		insertValues = []any{
			name, description, tenderId, authorType, authorId,
		}
		id string
	)

	log.Debug("beginning transaction")
	tx, err := s.db.Begin()
	if err != nil {
		log.Error("failed to begin transaction", sl.Err(err))
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer tx.Rollback()

	row := tx.QueryRow(insertQuery, insertValues...)
	err = row.Scan(&id)
	if err != nil {
		log.Error("failed to scan id", sl.Err(err))
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	log.Debug("trying to commit transaction")
	err = tx.Commit()
	if err != nil {
		log.Error("failed to commit transaction", sl.Err(err))
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return id, nil
}

func (s *Storage) GetBidsById(limit int32, offset int32, authorId string) ([]model2.BidDB, error) {
	const op = "Repo.GetBidsById"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		selectQuery = `
		SELECT 
		    b.id,
		    b.name,
		    b.description,
		    CAST(b.status AS text),
			b.tenderId,
			CAST(b.authorType AS text),
     		COALESCE(e.id, o.id) AS authorId,
     		b.version,
     		b.createdAt
		FROM bid b
		LEFT JOIN employee e ON  e.id = b.authorId  
		LEFT JOIN organization o ON o.id = b.authorId  
		WHERE COALESCE(e.id, o.id) = $3::uuid 
		ORDER BY name ASC
		LIMIT CASE WHEN $1 = 0 THEN NULL ELSE $1 END
		OFFSET COALESCE($2, 0);
`
		selectValues = []any{
			limit, offset, authorId,
		}
		bids []model2.BidDB
	)

	err := s.db.Select(&bids, selectQuery, selectValues...)
	if err != nil {
		log.Error("failed to select bids for user", sl.Err(err))
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return bids, nil
}

func (s *Storage) BidsForTender(tenderId string, limit int32, offset int32) ([]model2.BidDB, error) {
	const op = "Repo.BidsForTender"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		selectQuery = `
		SELECT id, name, description, CAST(status AS text),
		       tenderId, CAST(authorType AS text), authorId, version, createdAt
		FROM bid
		WHERE tenderId = $3
		ORDER BY name ASC
		LIMIT CASE WHEN $1 = 0 THEN NULL ELSE $1 END
		OFFSET COALESCE($2, 0)
		;
`
		selectValues = []any{
			limit, offset, tenderId,
		}
		bids []model2.BidDB
	)

	err := s.db.Select(&bids, selectQuery, selectValues...)
	if err != nil {
		log.Error("failed to select tenders for user", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return bids, nil
}

func (s *Storage) BidStatus(bidId string) (string, error) {
	const op = "Repo.BidStatus"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		selectQuery = `
		SELECT CAST(status AS TEXT)
		FROM bid
		WHERE id = $1::uuid;
`
		selectValues = []any{
			bidId,
		}
		status string
	)

	err := s.db.Get(&status, selectQuery, selectValues...)
	if err != nil {
		log.Error("failed to get status for bid", err)
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return status, nil
}

func (s *Storage) UpdateBidStatus(bidId string, status string) (string, error) {
	const op = "Repo.UpdateBidStatus"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		insertVersion = `
        INSERT INTO bid_version (bid_id, version, name, description, decision, status, tenderId, authorType, authorId, createdAt)
		SELECT id, version, name, description, decision, status, tenderId, authorType, authorId, createdAt
		FROM bid
		WHERE id = $1::uuid;
    `
		updateQuery = `
		UPDATE bid
		SET status = $2::bid_status, 
			version = version + 1
		WHERE id = $1::uuid
		RETURNING id;
`
		versionValues = []any{
			bidId,
		}
		updateValues = []any{
			bidId, status,
		}
		id string
	)

	log.Debug("beginning transaction")
	tx, err := s.db.Begin()
	if err != nil {
		log.Error("failed to begin transaction", sl.Err(err))
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(insertVersion, versionValues...)
	if err != nil {
		log.Error("failed to insert to history table", sl.Err(err))
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	row := tx.QueryRow(updateQuery, updateValues...)
	err = row.Scan(&id)
	if err != nil {
		log.Error("failed to scan id", sl.Err(err))
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	log.Debug("trying to commit transaction")
	err = tx.Commit()
	if err != nil {
		log.Error("failed to commit transaction", sl.Err(err))
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return id, nil
}

func (s *Storage) EditBid(bidId string, name string, description string) (string, error) {
	const op = "Repo.EditBid"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		insertVersion = `
        INSERT INTO bid_version (bid_id, version, name, description, decision, status, tenderId, authorType, authorId, createdAt)
		SELECT id, version, name, description, decision, status, tenderId, authorType, authorId, createdAt
		FROM bid
		WHERE id = $1::uuid;
    `
		updateQuery = `
		UPDATE bid
		SET name = COALESCE(NULLIF($2, ''), name),
		    description = COALESCE(NULLIF($3, ''), name),
			version = version + 1
		WHERE id = $1::uuid
		RETURNING id;
`
		versionValues = []any{
			bidId,
		}
		updateValues = []any{
			bidId, name, description,
		}
		id string
	)

	log.Debug("beginning transaction")
	tx, err := s.db.Begin()
	if err != nil {
		log.Error("failed to begin transaction", sl.Err(err))
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(insertVersion, versionValues...)
	if err != nil {
		log.Error("failed to insert to history table", sl.Err(err))
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	row := tx.QueryRow(updateQuery, updateValues...)
	err = row.Scan(&id)
	if err != nil {
		log.Error("failed to scan id", sl.Err(err))
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	log.Debug("trying to commit transaction")
	err = tx.Commit()
	if err != nil {
		log.Error("failed to commit transaction", sl.Err(err))
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return id, nil
}

func (s *Storage) SubmitDecision(bidId string, responsibleId string) error {
	const op = "Support.SubmitDecision"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		insertQuery = `
		INSERT INTO bid_approval (bid_id, responsible)
		VALUES ($1, $2);
`
		Values = []any{
			bidId, responsibleId,
		}
	)
	_, err := s.db.Exec(insertQuery, Values...)

	if err != nil {
		log.Info("failed to submit decision..", sl.Err(err))
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to submit decision.."))
	}
	return nil
}

func (s *Storage) ApplyDecision(bidId string, decision string) error {
	const op = "Support.ApplyDecision"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		applyQuery = `
		UPDATE bid
		SET decision = $2::bid_decision
		WHERE id = $1;
`
		Values = []any{
			bidId, decision,
		}
	)
	_, err := s.db.Exec(applyQuery, Values...)

	if err != nil {
		log.Info("failed to apply decision..", sl.Err(err))
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to apply decision.."))
	}
	return nil
}

func (s *Storage) Feedback(bidId string, feedback string) error {
	const op = "Support.Feedback"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		insertQuery = `
		INSERT INTO feedback (description, bidId)
		VALUES ($1, $2::uuid);
`
		insertValues = []any{
			feedback, bidId,
		}
	)
	_, err := s.db.Exec(insertQuery, insertValues...)

	if err != nil {
		log.Error("failed to leave feedback..", sl.Err(err))
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to leave feedback"))
	}
	return nil
}

func (s *Storage) RollbackBid(bidId string, version int32) (string, error) {
	const op = "Repo.EditTender"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		insertVersion = `
        INSERT INTO bid_version (bid_id, version, name, description, decision, status, tenderId, authorType, authorId, createdAt)
		SELECT id, version, name, description, decision, status, tenderId, authorType, authorId, createdAt
		FROM bid
		WHERE id = $1::uuid;
    `
		updateQuery = `
		UPDATE bid
		SET name = v.name,
			description = v.description,
			decision = v.decision,
			status = v.status::bid_status,
			tenderId = v.tenderId,
			authorType = v.authorType,
			authorId = v.authorId,
			version = bid.version + 1
		FROM bid_version v
		WHERE bid.id = v.bid_id
		  AND bid.id = $1
		  AND v.version = $2;
`
		versionValues = []any{
			bidId,
		}
		updateValues = []any{
			bidId, version,
		}
		id string
	)

	log.Debug("beginning transaction")
	tx, err := s.db.Begin()
	if err != nil {
		log.Error("failed to begin transaction", sl.Err(err))
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(insertVersion, versionValues...)
	if err != nil {
		log.Error("failed to insert to history table", sl.Err(err))
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	_, err = tx.Exec(updateQuery, updateValues...)
	if err != nil {
		log.Error("failed to scan id", sl.Err(err))
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	log.Debug("trying to commit transaction")
	err = tx.Commit()
	if err != nil {
		log.Error("failed to commit transaction", sl.Err(err))
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return id, nil
}

func (s *Storage) Reviews(authorUsername string, limit int32, offset int32) ([]model2.Feedback, error) {
	const op = "Repo.Reviews"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		selectQuery = `
		SELECT fb.id, fb.description, fb.createdAt
		FROM feedback fb
		JOIN bid b ON fb.bidId = b.id
		JOIN employee e ON b.authorId = e.id
		WHERE e.username = $1
		LIMIT CASE WHEN $2 = 0 THEN NULL ELSE $2 END
		OFFSET COALESCE($3, 0);
`
		selectValues = []any{
			authorUsername, limit, offset,
		}
		reviews []model2.Feedback
	)

	err := s.db.Select(&reviews, selectQuery, selectValues...)
	if err != nil {
		log.Error("failed to get reviews for author", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return reviews, nil
}
