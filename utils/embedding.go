package utils

import (
	"context"
	"fmt"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

func Embedding(ctx context.Context, like string, client *arkruntime.Client) (model.EmbeddingResponse, error) {
	req := model.EmbeddingRequestStrings{
		Input: []string{
			like,
		},
		Model:          "doubao-embedding-large-text-240915",
		EncodingFormat: "float",
	}
	resp, err := client.CreateEmbeddings(ctx, req)
	if err != nil {
		fmt.Printf("embeddings error: %v\n", err)
		return model.EmbeddingResponse{}, err
	}
	return resp, nil
}
