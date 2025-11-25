package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"google.golang.org/genai"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type HelloWorldParams struct {
	Name string `json:"name"`
}

type GodocParams struct {
	Package string `json:"package"`
	Symbol  string `json:"symbol"`
}

type CodeReviewParams struct {
	Code string `json:"code"`
	Hint string `json:"hint"`
}

func main() {
	mcpServer := mcp.NewServer(&mcp.Implementation{Name: "godoctor", Version: "0.0.1"}, nil)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "hello_world",
		Description: "A simple hello world tool.",
	}, func(ctx context.Context, request *mcp.CallToolRequest, params HelloWorldParams) (*mcp.CallToolResult, any, error) {
		name := params.Name
		if name == "" {
			name = "World"
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Hello, %s!", name)},
			},
		}, nil, nil
	})

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "godoc",
		Description: "Gets the documentation for a Go package or symbol.",
	}, func(ctx context.Context, request *mcp.CallToolRequest, params GodocParams) (*mcp.CallToolResult, any, error) {
		pkg := params.Package
		symbol := params.Symbol

		args := []string{"doc"}
		if pkg != "" {
			args = append(args, pkg)
		}
		if symbol != "" {
			args = append(args, symbol)
		}

		cmd := exec.CommandContext(ctx, "go", args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to execute go doc: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(output)},
			},
		}, nil, nil
	})

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "code_review",
		Description: "Analyzes Go code and provides a list of improvements in JSON format.",
	}, func(ctx context.Context, request *mcp.CallToolRequest, params CodeReviewParams) (*mcp.CallToolResult, any, error) {
		code := params.Code
		hint := params.Hint

		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey: os.Getenv("GEMINI_API_KEY"),
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create genai client: %w", err)
		}

		prompt := fmt.Sprintf("Please review the following Go code and provide a list of improvements in JSON format. Each improvement should have a 'line', 'description', and 'suggestion'.\n\nCode:\n```go\n%s\n```", code)
		if hint != "" {
			prompt += fmt.Sprintf("\n\nHint: %s", hint)
		}

		resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), nil)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate content: %w", err)
		}

		var responseText string
		for _, cand := range resp.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					responseText += part.Text
				}
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText},
			},
		}, nil, nil
	})

	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return mcpServer
	}, nil)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}