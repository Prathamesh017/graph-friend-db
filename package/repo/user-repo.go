package repo
 
import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	
)
 
// UserRepository handles all user-related database operations
type UserRepository struct {
	session neo4j.SessionWithContext
}
 
// NewUserRepository creates a new user repository
func NewUserRepository(session neo4j.SessionWithContext) *UserRepository {
	return &UserRepository{session: session}
}
 
// CreateUser creates a new user node in the database
func (r *UserRepository) CreateUser(ctx context.Context, user *User) error {
	// Validate before creating
	if err := user.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
 
	// Create node in database
	result, err := r.session.Run(
		ctx,
		"CREATE (u:User {id: $id, name: $name}) RETURN u",
		map[string]interface{}{
			"id":   user.ID,
			"name": user.Name,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
 
	// Check if creation was successful
	if result.Next(ctx) {
		return nil
	}
 
	return fmt.Errorf("failed to create user: no result returned")
}
 
// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*User, error) {
	result, err := r.session.Run(
		ctx,
		"MATCH (u:User {id: $id}) RETURN u.id, u.name",
		map[string]interface{}{"id": id},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
 
	if result.Next(ctx) {
		record := result.Record()
		idVal, _ := record.Get("u.id")
		nameVal, _ := record.Get("u.name")
 
		return &User{
			ID:   int(idVal.(int64)), // Convert int64 to int
			Name: nameVal.(string),
		}, nil
	}
 
	return nil, fmt.Errorf("user with id %d not found", id)
}