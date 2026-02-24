package main

import (
	"context"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)


func aiResponse(query string) string {
	ctx := context.Background()
	client := openai.NewClient(
		option.WithBaseURL("https://ai.hackclub.com/proxy/v1"),
		option.WithAPIKey(os.Getenv("HCAI_API_KEY")),
	)


	resp, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(query)},
		Model: openai.ChatModel("hi"),
	})

	if err != nil {
		panic(err)
	}

	return resp.OutputText()

	
}