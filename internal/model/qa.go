package model

import (
	"context"
	"fmt"

	"github.com/Pawancod/MemoryDb-Vector-Search/internal/service"
)

type QAPair struct {
	Question string
	Answer   string
	Vector   []float64
}

func GetQAPairs(ctx context.Context) []QAPair {
	qaPairs := []QAPair{
		{"What is Pawan?", "Pawan is air", nil},
		{"What is CCIE?", "It is a certification", nil},
		{"What is AWS?", "AWS is a cloud service provider", nil},
	}

	// Generate vectors for each Q&A pair
	for i := range qaPairs {
		vector, err := service.GenerateVector(ctx, qaPairs[i].Question)
		if err != nil {
			fmt.Println("vector not generated for :", i, err)
		}
		qaPairs[i].Vector = vector
	}

	return qaPairs
}
