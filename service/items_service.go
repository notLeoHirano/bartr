package service

import (
	"fmt"

	"github.com/notLeoHirano/bartr/models"
)

func (s *Service) GetItems(userID int, excludeOwn bool) ([]models.ItemWithOwner, error) {
	return s.repo.GetItems(userID, excludeOwn)
}

func (s *Service) CreateItem(item *models.Item) error {
	if item.Title == "" {
		return fmt.Errorf("title is required")
	}
	return s.repo.CreateItem(item)
}

func (s *Service) DeleteItem(id int, userID int) error {
	deleted, err := s.repo.DeleteItem(id, userID)
	if err != nil {
		return err
	}
	if !deleted {
		return fmt.Errorf("item not found or you don't have permission to delete it")
	}
	return nil
}