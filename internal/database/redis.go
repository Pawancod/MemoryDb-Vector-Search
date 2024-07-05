package database

import (
	"context"
	"fmt"
	"strings"
	"github.com/Pawancod/MemoryDb-Vector-Search/internal/model"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(ctx context.Context) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "clustercfg.test.u9fd0i.memorydb.ap-south-1.amazonaws.com:6379", // Replace with your MemoryDB endpoint
		Password: "",      // Replace with your MemoryDB password
		DB:       0,
	})
}

func CheckRedisConnection(ctx context.Context, rdb *redis.Client) error {
	_, err := rdb.Ping(ctx).Result()
	return err
}

func StoreQAPair(ctx context.Context, rdb *redis.Client, id int, qa model.QAPair) error {
	vectorString := floatsToString(qa.Vector)
	qaKey := fmt.Sprintf("qa:%d", id)

	return rdb.HSet(ctx, qaKey, "question", qa.Question, "answer", qa.Answer, "vector", vectorString).Err()
}

func CreateIndex(ctx context.Context, rdb *redis.Client) error {
	_, err := rdb.Do(ctx, "FT.CREATE", "idx:qa", "ON", "HASH", "PREFIX", "1", "qa:", "SCHEMA", "vector", "VECTOR", "HNSW", "6", "TYPE", "FLOAT32", "DIM", "384", "DISTANCE_METRIC", "COSINE").Result()
	return err
}

// Function for converting a slice of floats to a space-separated string
func floatsToString(floats []float64) string {
	strs := make([]string, len(floats))
	for i, f := range floats {
	// Convert float64 to float32
	f32 := float32(f)
	// Format as string with %f for float32
	strs[i] = fmt.Sprintf("%f", f32)
}
	return strings.Join(strs, " ")
}
