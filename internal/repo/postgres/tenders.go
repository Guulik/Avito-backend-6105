package postgres

import (
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"log/slog"
	"net/http"
	"zadanie-6105/internal/domain/model"
	sl "zadanie-6105/internal/lib/logger/slog"
)

func (s *Storage) Tenders(limit int32, offset int32, serviceTypes []string) ([]model.TenderDB, error) {
	const op = "Repo.Tenders"
	log := s.log.With(
		slog.String("op", op),
	)

	log.Debug("types: ", pq.Array(serviceTypes))
	var (
		selectQuery = `
		SELECT id, name, description, CAST(serviceType AS text),
		       CAST(status AS text), organization_id, creator_username, version, created_at
		FROM tender
		WHERE ($3::TEXT[] IS NULL OR serviceType::TEXT = ANY($3))
		ORDER BY name ASC
		LIMIT CASE WHEN $1 = 0 THEN NULL ELSE $1 END
		OFFSET COALESCE($2, 0);
`
		selectValues = []any{
			limit, offset, pq.Array(serviceTypes),
		}
		tenders []model.TenderDB
	)

	err := s.db.Select(&tenders, selectQuery, selectValues...)
	if err != nil {
		log.Error("failed to select tenders", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return tenders, nil
}

func (s *Storage) CreateTender(name string, description string, serviceType string, organizationId string, creatorUsername string) (string, error) {
	const op = "Repo.CreateTender"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		insertQuery = `
        INSERT INTO tender (name, description, serviceType, status, organization_id, creator_username, version, created_at)
        VALUES ($1, $2, $3, 'Created'::tender_status, $4,$5, 1, CURRENT_TIMESTAMP)
        RETURNING id;
    `
		insertValues = []any{
			name, description, serviceType, organizationId, creatorUsername,
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

func (s *Storage) TendersByUser(limit int32, offset int32, username string) ([]model.TenderDB, error) {
	const op = "Repo.TendersByUser"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		selectQuery = `
		SELECT t.id, t.name, t.description, CAST(t.serviceType AS text),
		       CAST(t.status AS text), t.organization_id, t.creator_username, t.version, t.created_at
		FROM tender t
		JOIN organization o ON t.organization_id = o.id
		JOIN organization_responsible resp ON resp.organization_id = o.id
		JOIN employee e ON e.id = resp.user_id
		WHERE e.username = $3
		ORDER BY t.name ASC
		LIMIT CASE WHEN $1 = 0 THEN NULL ELSE $1 END
		OFFSET COALESCE($2, 0);
`
		selectValues = []any{
			limit, offset, username,
		}
		tenders []model.TenderDB
	)

	err := s.db.Select(&tenders, selectQuery, selectValues...)
	if err != nil {
		log.Error("failed to select tenders for user", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return tenders, nil
}

func (s *Storage) Status(tenderID string) (string, error) {
	const op = "Repo.Status"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		selectQuery = `
		SELECT CAST(status AS TEXT)
		FROM tender
		WHERE id = $1::uuid;
`
		selectValues = []any{
			tenderID,
		}
		status string
	)

	err := s.db.Get(&status, selectQuery, selectValues...)
	if err != nil {
		log.Error("failed to get status for tender", err)
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return status, nil
}

func (s *Storage) ChangeTenderStatus(tenderId string, status string) (string, error) {
	const op = "Repo.ChangeTenderStatus"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		insertVersion = `
        INSERT INTO tender_version (tender_id, version, name, description, serviceType, status, organization_id, creator_username, created_at)
		SELECT id, version, name, description, serviceType, status, organization_id, creator_username, created_at
		FROM tender
		WHERE id = $1::uuid;
    `
		updateQuery = `
		UPDATE tender
		SET status = $2::tender_status, 
			version = version + 1
		WHERE id = $1::uuid
		RETURNING id;
`
		versionValues = []any{
			tenderId,
		}
		updateValues = []any{
			tenderId, status,
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

func (s *Storage) EditTender(tenderId string, name string, description string, serviceType string) (string, error) {
	const op = "Repo.EditTender"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		insertVersion = `
        INSERT INTO tender_version (tender_id, version, name, description, serviceType, status, organization_id, creator_username,  created_at)
		SELECT id, version, name, description, serviceType, status, organization_id, creator_username, created_at
		FROM tender
		WHERE id = $1::uuid;
    `
		updateQuery = `
		UPDATE tender
		SET name = COALESCE(NULLIF($2, ''), name),
		    description = COALESCE(NULLIF($3, ''), description),
		    serviceType = COALESCE(NULLIF($4, '')::service_type, serviceType), 
			version = version + 1
		WHERE id = $1::uuid
		RETURNING id;
`
		versionValues = []any{
			tenderId,
		}
		updateValues = []any{
			tenderId, name, description, serviceType,
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

func (s *Storage) RollbackTender(tenderId string, version int32) (string, error) {
	const op = "Repo.RollbackTender"
	log := s.log.With(
		slog.String("op", op),
	)

	var (
		insertVersion = `
        INSERT INTO tender_version (tender_id, version, name, description, serviceType, status, organization_id, creator_username, created_at)
		SELECT id, version, name, description, serviceType, status, organization_id, creator_username, created_at
		FROM tender
		WHERE id = $1::uuid;
    `
		updateQuery = `
		UPDATE tender
		SET name = v.name,
			description = v.description,
			serviceType = v.serviceType,
			status = v.status::tender_status,
			organization_id = v.organization_id,
			creator_username = v.creator_username,
			version = tender.version + 1
		FROM tender_version v
		WHERE tender.id = v.tender_id
		  AND tender.id = $1
		  AND v.version = $2;
`
		versionValues = []any{
			tenderId,
		}
		updateValues = []any{
			tenderId, version,
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
		log.Error("failed to update", sl.Err(err))
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
