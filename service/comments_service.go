package service

import (
	"fmt"

	"github.com/notLeoHirano/bartr/models"
)

func (s *Service) CreateComment(comment *models.Comment) error {
	if comment.Content == "" {
		return fmt.Errorf("comment content is required")
	}

	// Verify user is part of the match
	inMatch, err := s.repo.UserInMatch(comment.MatchID, comment.UserID)
	if err != nil {
		return err
	}
	if !inMatch {
		return fmt.Errorf("you are not part of this match")
	}

	return s.repo.CreateComment(comment)
}

func (s *Service) GetComments(matchID int) ([]models.Comment, error) {
	return s.repo.GetComments(matchID)
}