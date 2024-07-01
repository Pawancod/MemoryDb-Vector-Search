package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
)

func PerformVectorSearch(ctx context.Context, rdb *redis.Client, queryQuestion string) (string, error) {
	queryVector, err := GenerateVector(ctx, queryQuestion)
	if err != nil {
		return "", err
	}
	queryVectorString := floatsToString(queryVector)

	// Perform a vector search
	queryResult, err := rdb.Do(ctx, "FT.SEARCH", "idx:qa", "*=>[KNN 1 @vector $vec AS score]", "PARAMS", "2", "vec", queryVectorString, "SORTBY", "score", "ASC").Result()
	if err != nil {
		return "", err
	}

	// Extract the answer from the search result
	answerKey := queryResult.([]interface{})[1].(string)
	answer, err := rdb.HGet(ctx, answerKey, "answer").Result()
	if err != nil {
		return "", err
	}

	return answer, nil
}

// GenerateVector calls an external service to generate a vector for a given question
func GenerateVector(ctx context.Context, question string) ([]float64, error) {
	payload := map[string]string{"text": question}
	jsonPayload, _ := json.Marshal(payload)
	resp, err := http.Post("http://localhost:5001/vectorize", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error generating vector: %v", err)
	}
	defer resp.Body.Close()

	var result map[string][]float64
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}
	return result["vector"], nil
}

// Helper function to convert a slice of floats to a space-separated string
func floatsToString(floats []float64) string {
	strs := make([]string, len(floats))
	for i, f := range floats {
		strs[i] = fmt.Sprintf("%f", f)
	}
	return strings.Join(strs, " ")
}
