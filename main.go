package main



import (
	"context"
	db "friend-recommendation/package/db"
	repo "friend-recommendation/package/repo"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"fmt"
	"log"
)



func main() {
	driver := db.InitDB()
	defer driver.Close(context.Background())


	session := driver.NewSession(context.Background(), neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(context.Background())
	setupGraphData(session)
}



func setupGraphData(session neo4j.SessionWithContext) {
	ctx := context.Background()
 
	// Create user repository
	userRepo := repo.NewUserRepository(session)
	friendshipRepo := repo.NewFriendshipRepository(session)
 
	// Create users with validation
	users := []*repo.User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Carol"},
		{ID: 4, Name: "Dave"},
		{ID: 5, Name: "Eve"},
		{ID: 6, Name: "Frank"},
		{ID: 7, Name: "Grace"},
	}
 
	fmt.Println("Creating users...")
	for _, user := range users {
		err := userRepo.CreateUser(ctx, user)
		if err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
	}
	fmt.Printf("✓ Created %d users\n", len(users))
 
	// Create friendships with validation
	friendships := []*repo.Friendship{
		{User1ID: 1, User2ID: 2, Strength: 5},  // Alice - Bob
		{User1ID: 1, User2ID: 3, Strength: 2},  // Alice - Carol
		{User1ID: 2, User2ID: 5, Strength: 4},  // Bob - Eve
		{User1ID: 2, User2ID: 6, Strength: 3},  // Bob - Frank
		{User1ID: 3, User2ID: 6, Strength: 4},  // Carol - Frank
		{User1ID: 3, User2ID: 7, Strength: 3},  // Carol - Grace
		{User1ID: 5, User2ID: 6, Strength: 2},  // Eve - Frank
	}
 
	fmt.Println("\nCreating friendships...")
	for _, friendship := range friendships {
		err := friendshipRepo.CreateFriendship(ctx, friendship)
		if err != nil {
			log.Fatalf("Failed to create friendship: %v", err)
		}
	}
	fmt.Printf("✓ Created %d friendships (bidirectional)\n", len(friendships))
 
	// Query recommendations

 
	// 1st degree friends
	fmt.Println("\n1st Degree Friends (Direct):")
	friends1, err := friendshipRepo.GetDirectFriends(ctx, 1)
	if err != nil {
		log.Fatalf("Failed to get direct friends: %v", err)
	}
	for _, friend := range friends1 {
		fmt.Printf("  - %s (ID: %d)\n", friend.Name, friend.ID)
	}
 
	// 2nd degree friends
	fmt.Println("\n2nd Degree Friends (Friends of Friends):")
	friends2, err := friendshipRepo.GetFriendsOfFriends(ctx, 1)
	if err != nil {
		log.Fatalf("Failed to get friends of friends: %v", err)
	}
	for _, friend := range friends2 {
		fmt.Printf("  - %s (ID: %d)\n", friend.Name, friend.ID)
	}
 
	// fmt.Println("\n" + "="*60)
}