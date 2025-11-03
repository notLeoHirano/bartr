package service

import (
	"fmt"
	"log"

	"github.com/notLeoHirano/bartr/models"
)

func (s *Service) CreateSwipe(swipe *models.Swipe) error {
	if swipe.Direction != "left" && swipe.Direction != "right" {
		return fmt.Errorf("direction must be 'left' or 'right'")
	}

	if err := s.repo.CreateSwipe(swipe); err != nil {
		return err
	}

	log.Printf("User %d swiped %s on item %d", swipe.UserID, swipe.Direction, swipe.ItemID)

	// Check for matches only if swipe was right
	if swipe.Direction == "right" {
		if err := s.checkAndCreateMatches(swipe.UserID, swipe.ItemID); err != nil {
			log.Printf("Error checking for matches: %v", err)
		}
	}

	return nil
}

func (s *Service) checkAndCreateMatches(swipingUserID, swipedItemID int) error {
	print("checking for matches")
	itemOwnerID, err := s.repo.GetItemOwnerID(swipedItemID)
	if err != nil {
		return fmt.Errorf("error finding item owner: %w", err)
	}

	userItems, err := s.repo.GetItems(swipingUserID, false)
	if err != nil {
		return fmt.Errorf("error fetching user's items: %w", err)
	}

	for _, userItem := range userItems {
		if userItem.UserID != swipingUserID {
			continue
		}

		// Check if the item owner already swiped right on this user's item
		ownerSwiped, err := s.repo.UserSwipedRight(itemOwnerID, userItem.ID)
		if err != nil {
			log.Printf("Error checking swipe: %v", err)
			continue
		}
		if !ownerSwiped {
			continue
		}

		// Create match using transaction
		if err := s.repo.CreateMatchIfNeeded(swipingUserID, itemOwnerID, userItem.ID, swipedItemID); err != nil {
			log.Printf("Error creating match: %v", err)
			continue
		}

		log.Printf("Match created! User %d item %d <-> User %d item %d",
			swipingUserID, userItem.ID, itemOwnerID, swipedItemID)
	}

	return nil
}


func (s *Service) GetMatches(userID int) ([]models.MatchResponse, error) {
	return s.repo.GetMatches(userID)
}