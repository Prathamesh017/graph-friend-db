package db

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"fmt"
	
)

func InitDB() neo4j.DriverWithContext {
	ctx := context.Background()
	driver, err := neo4j.NewDriverWithContext(
		"neo4j://localhost:7688",
		neo4j.BasicAuth("neo4j", "password123", ""),
	)
	if err != nil {
		log.Fatalf("Failed to create Neo4j driver: %v", err)
	}

	if err := driver.VerifyConnectivity(ctx); err != nil {
		_ = driver.Close(ctx)
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}

	// Test connection
	err = driver.VerifyConnectivity(context.Background())
	if err != nil {
		log.Fatalf("Failed to verify connectivity: %v", err)
	}
 
	fmt.Println("✓ Connected to Neo4j successfully!")

	return driver
}
