package repo

import (
	"fmt"
	"strings"
)

// User represents a person in the social network
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Validate checks if a User is valid
func (u *User) Validate() error {
	if u.ID <= 0 {
		return fmt.Errorf("user id must be positive, got %d", u.ID)
	}

	if strings.TrimSpace(u.Name) == "" {
		return fmt.Errorf("user name cannot be empty")
	}

	if len(u.Name) > 100 {
		return fmt.Errorf("user name must be less than 100 characters, got %d", len(u.Name))
	}

	return nil
}

// Friendship represents a connection between two users
type Friendship struct {
	User1ID int `json:"user1_id"`
	User2ID int `json:"user2_id"`
	Strength int `json:"strength"` // 1-5 scale
}

// Validate checks if a Friendship is valid
func (f *Friendship) Validate() error {
	if f.User1ID <= 0 {
		return fmt.Errorf("user1_id must be positive, got %d", f.User1ID)
	}

	if f.User2ID <= 0 {
		return fmt.Errorf("user2_id must be positive, got %d", f.User2ID)
	}

	if f.User1ID == f.User2ID {
		return fmt.Errorf("cannot create friendship between same user (id: %d)", f.User1ID)
	}

	if f.Strength < 1 || f.Strength > 5 {
		return fmt.Errorf("friendship strength must be between 1-5, got %d", f.Strength)
	}

	return nil
}

// FriendRecommendation represents a recommendation result
type FriendRecommendation struct {
	UserID   int
	UserName string
	Degree   int // 1st degree, 2nd degree, etc.
	Strength int // Average connection strength
}