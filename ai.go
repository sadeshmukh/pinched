package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"

	"github.com/sadeshmukh/pinched/tools"
)

func aiResponseWithTools(query string, toolList []tools.Tool) string {
	ctx := context.Background()
	client := openai.NewClient(
		option.WithBaseURL("https://ai.hackclub.com/proxy/v1"),
		option.WithAPIKey(os.Getenv("HCAI_API_KEY")),
	)

	apiTools := []responses.ToolUnionParam{}
	for _, tool := range toolList {
		apiTools = append(apiTools,
			responses.ToolUnionParam{
				OfFunction: &responses.FunctionToolParam{
					Name:        tool.Name,
					Description: openai.String(tool.Description),
					Parameters: map[string]any{
						"type":       "object",
						"properties": tool.Parameters,
						"required":   tool.Required,
					},
				},
			})
	}

	resp, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(query)},
		Model: openai.ChatModel("qwen3-32b"),
		Tools: apiTools,
	})

	if err != nil {
		panic(err)
	}

	for _, item := range resp.Output {
		if item.Type == "function_call" {
			toolCall := item.AsFunctionCall()

			var whichTool *tools.Tool
			for _, t := range toolList {
				if t.Name == toolCall.Name {
					whichTool = &t
					break
				}
			}

			var args map[string]interface{}
			json.Unmarshal([]byte(toolCall.Arguments), &args)
			result, err := whichTool.Exec(args)

			if err != nil {
				panic(err)
			}

			resp, _ = client.Responses.New(ctx, responses.ResponseNewParams{
				PreviousResponseID: openai.String(resp.ID),
				Input: responses.ResponseNewParamsInputUnion{
					OfInputItemList: []responses.ResponseInputItemUnionParam{{
						OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
							CallID: toolCall.CallID,
							Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{
								OfString: openai.String(result),
							},
						},
					}},
				},
			})

		}
	}

	return resp.OutputText()
}
