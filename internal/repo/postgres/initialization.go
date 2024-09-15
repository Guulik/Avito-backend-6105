package postgres

import (
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"zadanie-6105/config"
	"zadanie-6105/internal/repo"
)

var _ repo.Checkers = (*Storage)(nil)

type Storage struct {
	log *slog.Logger
	db  *sqlx.DB
}

func New(log *slog.Logger, db *sqlx.DB) *Storage {
	return &Storage{
		db:  db,
		log: log,
	}
}

func ConnectPostgres(c *config.Config) (*sqlx.DB, error) {
	connectionUrl := c.POSTGRES_CONN

	db, err := sqlx.Connect("pgx", connectionUrl)
	if err != nil {
		fmt.Println("connection error: ", err)
	}

	fmt.Println("connection to db was successful!")

	return db, err
}

func CreateTables(db *sqlx.DB) error {
	var err error

	// add extension to generate uuid
	_, err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		fmt.Println("failed to create extension: ", err)
		return err
	}

	//считай миграция))
	//_ = DropTables(db)

	creationFunctions := []func(db *sqlx.DB) error{
		OrganizationTable, EmployeeTable, OrganizationResponsibleTable, TenderTable, BidTable, FeedbackTable,
		TenderVersionTable, BidVersionTable, BidDecisionTable,
	}

	for i, function := range creationFunctions {
		err = function(db)
		if err != nil {
			fmt.Println("failed to create table: ", i, err)
			return err
		}
	}

	return nil
}

func OrganizationTable(db *sqlx.DB) error {
	query :=
		`
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'organization_type') THEN
        CREATE TYPE organization_type AS ENUM ('IE', 'LLC', 'JSC');
    END IF;
END $$;


CREATE TABLE IF NOT EXISTS organization
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    type        organization_type,
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);
`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}

func EmployeeTable(db *sqlx.DB) error {
	query :=
		`
CREATE TABLE IF NOT EXISTS employee (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}

func OrganizationResponsibleTable(db *sqlx.DB) error {
	query :=
		`
CREATE TABLE IF NOT EXISTS organization_responsible (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);
`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}

func TenderTable(db *sqlx.DB) error {
	query :=
		`
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'service_type') THEN
        CREATE TYPE service_type AS ENUM ('Construction', 'Delivery', 'Manufacture');
    END IF;
END $$;


DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tender_status') THEN
        CREATE TYPE tender_status AS ENUM ('Created', 'Published', 'Closed');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS tender (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(100),
  description VARCHAR(500),
  serviceType service_type,
  status tender_status,
  organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
  creator_username VARCHAR(50) REFERENCES employee(username) ON DELETE CASCADE,
  version INT DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}

func BidTable(db *sqlx.DB) error {
	typesQuery :=
		`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'bid_status') THEN
				CREATE TYPE bid_status AS ENUM ('Created', 'Published', 'Canceled');
			END IF;
		END $$;
		
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'bid_decision') THEN
				CREATE TYPE bid_decision AS ENUM ('','Approved', 'Rejected');
			END IF;
		END $$;
		
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'author_type') THEN
				CREATE TYPE author_type AS ENUM ('Organization', 'User');
			END IF;
		END $$;
`

	query :=
		`
CREATE TABLE IF NOT EXISTS bid (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	name VARCHAR(250),
	description VARCHAR(500),
    decision bid_decision DEFAULT '',
	status bid_status,
	tenderId UUID REFERENCES tender(id) ON DELETE CASCADE,
	authorType author_type,
	authorId UUID,
	version INT DEFAULT 1,
	createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

	if _, err := db.Exec(typesQuery); err != nil {
		return err
	}

	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}

func FeedbackTable(db *sqlx.DB) error {
	query :=
		`
CREATE TABLE IF NOT EXISTS feedback (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  description VARCHAR(1000),
  createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  bidId UUID REFERENCES bid(id) ON DELETE CASCADE
);
`

	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}

func TenderVersionTable(db *sqlx.DB) error {
	query :=
		`
CREATE TABLE IF NOT EXISTS tender_version (
  tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
  version INT,
  name VARCHAR(250),
  description VARCHAR(500),
  serviceType service_type,
  status tender_status,
  organization_id UUID,
creator_username VARCHAR(50) REFERENCES employee(username) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (tender_id, version)
);
`

	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}

func BidVersionTable(db *sqlx.DB) error {
	query :=
		`
CREATE TABLE IF NOT EXISTS bid_version (
  bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
  version INT,
  name VARCHAR(250),
  description TEXT,
  decision bid_decision,
  status bid_status,
  tenderId UUID,
  authorType author_type,
  authorId UUID,
  createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (bid_id, version)
);
`

	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}

func BidDecisionTable(db *sqlx.DB) error {
	query :=
		`
CREATE TABLE IF NOT EXISTS bid_approval (
   bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
	responsible UUID,
    PRIMARY KEY (bid_id, responsible)
);
`

	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}

func UniqueIndex(db *sqlx.DB) error {
	// пока не уверен
	query :=
		`
CREATE UNIQUE INDEX unique_organization_user
    ON organization_responsible (organization_id, user_id);
);
`
	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}

func DropTables(db *sqlx.DB) error {
	//TODO: delete me
	var DEV_drop string
	DEV_drop = `
			DROP TYPE IF EXISTS tender_status CASCADE;
			DROP TYPE IF EXISTS service_type CASCADE;
			DROP TABLE IF EXISTS tender CASCADE;
			`
	if _, err := db.Exec(DEV_drop); err != nil {
		return err
	}

	DEV_drop = `
					DROP TYPE IF EXISTS bid_status CASCADE;
					DROP TYPE IF EXISTS bid_decision CASCADE;
					DROP TABLE IF EXISTS bid CASCADE;
				`
	if _, err := db.Exec(DEV_drop); err != nil {
		return err
	}

	DEV_drop = `DROP TABLE IF EXISTS feedback CASCADE;`
	if _, err := db.Exec(DEV_drop); err != nil {
		return err
	}

	//TODO: delete me
	DEV_drop = `DROP TABLE IF EXISTS tender_version CASCADE;`
	if _, err := db.Exec(DEV_drop); err != nil {
		return err
	}
	DEV_drop = `DROP TABLE IF EXISTS bid_version CASCADE;`
	if _, err := db.Exec(DEV_drop); err != nil {
		return err
	}

	DEV_drop = `DROP TABLE IF EXISTS bid_approval CASCADE;`
	if _, err := db.Exec(DEV_drop); err != nil {
		return err
	}

	return nil
}
