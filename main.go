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


	//-------------- storing data -------------
	// Storing each Q&A pair in Redis
	for i, qa := range qaPairs {
		err := database.StoreQAPair(ctx, rdb, i+1, qa)
		if err != nil {
			fmt.Println("Error storing Q&A pair:", err)
			os.Exit(1)
		}
		fmt.Printf("Stored Q&A pair: qa:%d\n", i+1)
	}

	// Creating index for vector search to have similar answers
	err = database.CreateIndex(ctx, rdb)
	if err != nil {
		fmt.Println("Error creating index:", err)
		os.Exit(1)
	}
	fmt.Println("Index created")

	//-------searching data via vector search ------

	//queryQuestion := "What is CCIE?"
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a question to search for: ")
	queryQuestion, _ := reader.ReadString('\n')
	queryQuestion = strings.TrimSpace(queryQuestion)

	answer, err := service.PerformVectorSearch(ctx, rdb, queryQuestion)
	if err != nil {
		fmt.Println("Error performing vector search:", err)
		os.Exit(1)
	}
	fmt.Println("Answer:", answer)
}
