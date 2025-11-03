package service

import (
	"github.com/notLeoHirano/bartr/store"
)

type Service struct {
	repo *store.Store
}

func New(repo *store.Store) *Service {
	return &Service{repo: repo}
}
