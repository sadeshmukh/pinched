package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

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

	messages := []responses.ResponseInputItemUnionParam{
		{
			OfMessage: &responses.EasyInputMessageParam{
				Role: "user",
				Content: responses.EasyInputMessageContentUnionParam{
					OfString: openai.String(query),
				},
			},
		},
	}

	for {

		if messages[len(messages)-1].GetContent().AsAny() != nil {
			fmt.Println("ai: " + *messages[len(messages)-1].GetContent().AsAny().(*string))
		}

		resp, err := client.Responses.New(ctx, responses.ResponseNewParams{
			Input: responses.ResponseNewParamsInputUnion{OfInputItemList: messages},
			Model: openai.ChatModel("qwen3-32b"),
			Tools: apiTools,
		})

		if err != nil {
			panic(err)
		}
		hasToolCall := false
		for _, item := range resp.Output {
			if item.Type == "function_call" {
				hasToolCall = true
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

				if strings.HasPrefix(result, "|END|") {
					return strings.TrimPrefix(result, "|END|")
				}

				messages = append(messages,
					responses.ResponseInputItemParamOfFunctionCall(toolCall.Arguments, toolCall.CallID, toolCall.Name),
					responses.ResponseInputItemParamOfFunctionCallOutput(toolCall.CallID, result))

			}
		}

		if !hasToolCall {
			return resp.OutputText()
		}
	}

}
