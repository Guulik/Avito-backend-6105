package service

import (
	"log/slog"
	"zadanie-6105/internal/repo"
)

type Service struct {
	log *slog.Logger

	repoTenderProvider RepoTenderProvider
	repoTenderCreator  RepoTenderCreator
	repoTenderEditor   RepoTenderEditor

	repoBidProvider      RepoBidProvider
	repoBidCreator       RepoBidCreator
	repoBidEditor        RepoBidEditor
	repoBidDecisionMaker RepoBidDecisionMaker
	repoBidFeedbacker    RepoBidFeedbacker

	checkers repo.Checkers
}

func New(
	log *slog.Logger,

	tenderProvider RepoTenderProvider,
	tenderCreator RepoTenderCreator,
	tenderEditor RepoTenderEditor,

	bidProvider RepoBidProvider,
	bidCreator RepoBidCreator,
	bidEditor RepoBidEditor,
	bidDecisionMaker RepoBidDecisionMaker,
	bidFeedbacker RepoBidFeedbacker,

	checkers repo.Checkers,
) *Service {
	return &Service{
		log:                  log,
		repoTenderProvider:   tenderProvider,
		repoTenderCreator:    tenderCreator,
		repoTenderEditor:     tenderEditor,
		repoBidProvider:      bidProvider,
		repoBidCreator:       bidCreator,
		repoBidEditor:        bidEditor,
		repoBidDecisionMaker: bidDecisionMaker,
		repoBidFeedbacker:    bidFeedbacker,
		checkers:             checkers,
	}
}
