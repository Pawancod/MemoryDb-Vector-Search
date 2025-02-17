package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Pawancod/MemoryDb-Vector-Search/internal/database"
	"github.com/Pawancod/MemoryDb-Vector-Search/internal/model"
	"github.com/Pawancod/MemoryDb-Vector-Search/internal/service"
)

func main() {
	ctx := context.Background()

	// Initializing Redis client
	rdb := database.NewRedisClient(ctx)

	// Check connection
	err := database.CheckRedisConnection(ctx, rdb)
	if err != nil {
		fmt.Println("Error connecting to MemoryDB:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to MemoryDB")

	// Defining Q&A pairs
	qaPairs := model.GetQAPairs(ctx)


	//--------Storing data and Indexing-----//
	for i, qa := range qaPairs {
		err := database.StoreQAPair(ctx, rdb, i+1, qa)
		if err != nil {
			fmt.Println("Error storing Q&A pair:", err)
			os.Exit(1)
		}
		fmt.Printf("Stored Q&A pair: qa:%d\n", i+1)
	}

	// Creating index for vector search
	err = database.CreateIndex(ctx, rdb)
	if err != nil {
		fmt.Println("Error creating index:", err)
		os.Exit(1)
	}
	fmt.Println("Index created")


	//--------Searching: Vector search-----//
	reader := bufio.NewReader(os.Stdin)
	defer func() {
		// Cleanup before exit
		fmt.Println("Flushing MemoryDB...")
		err := rdb.FlushAll(ctx).Err()
		if err != nil {
			fmt.Println("Error flushing MemoryDB:", err)
		}
		fmt.Println("MemoryDB flushed. Exiting...")
	}()

	for {
		
		fmt.Print("Enter a question to search for (type 'exit' to quit): ")
		queryQuestion, _ := reader.ReadString('\n')
		queryQuestion = strings.TrimSpace(queryQuestion)

		if queryQuestion == "exit" {
			break
		}

		answer, err := service.PerformVectorSearch(ctx, rdb, queryQuestion)
		if err != nil {
			fmt.Println("Error performing vector search:", err)
			continue // Continue to prompt for next question
		}

		fmt.Println("Answer:", answer)
	}

	fmt.Println("Program exited.")
}
