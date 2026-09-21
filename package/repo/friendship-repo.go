package repo


import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
)



type FriendshipRepository struct {
	session neo4j.SessionWithContext
}
 
// NewFriendshipRepository creates a new friendship repository
func NewFriendshipRepository(session neo4j.SessionWithContext) *FriendshipRepository {
	return &FriendshipRepository{session: session}
}
 
// CreateFriendship creates a bidirectional friendship between two users
func (r *FriendshipRepository) CreateFriendship(ctx context.Context, friendship *Friendship) error {
	// Validate before creating
	if err := friendship.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
 
	// Check if users exist
	user1, err := NewUserRepository(r.session).GetUserByID(ctx, friendship.User1ID)
	if err != nil {
		return fmt.Errorf("user1 does not exist: %w", err)
	}
 
	user2, err := NewUserRepository(r.session).GetUserByID(ctx, friendship.User2ID)
	if err != nil {
		return fmt.Errorf("user2 does not exist: %w", err)
	}
 
	// Create bidirectional edges
	_, err = r.session.Run(
		ctx,
		`MATCH (u1:User {id: $user1_id}), (u2:User {id: $user2_id})
		 CREATE (u1)-[:FRIENDS_WITH {strength: $strength}]->(u2),
		        (u2)-[:FRIENDS_WITH {strength: $strength}]->(u1)`,
		map[string]interface{}{
			"user1_id": friendship.User1ID,
			"user2_id": friendship.User2ID,
			"strength": friendship.Strength,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create friendship: %w", err)
	}
 
	log.Printf("Created friendship: %s (id: %d) ↔ %s (id: %d) with strength %d",
		user1.Name, user1.ID, user2.Name, user2.ID, friendship.Strength)
 
	return nil
}
 
// GetDirectFriends gets all direct friends (1st degree) of a user
func (r *FriendshipRepository) GetDirectFriends(ctx context.Context, userID int) ([]User, error) {
	result, err := r.session.Run(
		ctx,
		`MATCH (u:User {id: $user_id})-[:FRIENDS_WITH]-(friend:User)
		 RETURN friend.id, friend.name`,
		map[string]interface{}{"user_id": userID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get direct friends: %w", err)
	}
 
	var friends []User
	for result.Next(ctx) {
		record := result.Record()
		id, _ := record.Get("friend.id")
		name, _ := record.Get("friend.name")
 
		friends = append(friends, User{
			ID:   int(id.(int64)),
			Name: name.(string),
		})
	}
 
	return friends, nil
}
 
// GetFriendsOfFriends gets 2nd degree friends
func (r *FriendshipRepository) GetFriendsOfFriends(ctx context.Context, userID int) ([]User, error) {
	result, err := r.session.Run(
		ctx,
		`MATCH (u:User {id: $user_id})-[:FRIENDS_WITH]-(friend)-[:FRIENDS_WITH]-(friend_of_friend)
		 WHERE friend_of_friend.id <> $user_id
		 RETURN DISTINCT friend_of_friend.id, friend_of_friend.name`,
		map[string]interface{}{"user_id": userID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get friends of friends: %w", err)
	}
 
	var friends []User
	seen := make(map[int]bool) // Track seen IDs to avoid duplicates
 
	for result.Next(ctx) {
		record := result.Record()
		id, _ := record.Get("friend_of_friend.id")
		name, _ := record.Get("friend_of_friend.name")
 
		idInt := int(id.(int64))
		if !seen[idInt] {
			friends = append(friends, User{
				ID:   idInt,
				Name: name.(string),
			})
			seen[idInt] = true
		}
	}
 
	return friends, nil
}
 