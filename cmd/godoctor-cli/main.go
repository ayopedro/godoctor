package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	reviewCmd := flag.NewFlagSet("review", flag.ExitOnError)
	reviewHint := reviewCmd.String("hint", "", "An optional hint to provide additional guidance to the AI reviewer.")

	docCmd := flag.NewFlagSet("doc", flag.ExitOnError)
	docSymbol := docCmd.String("symbol", "", "The symbol to get the documentation for.")

	if len(os.Args) < 2 {
		log.Fatal("Usage: godoctor-cli <command> [args...]")
	}

	switch os.Args[1] {
	case "review":
		reviewCmd.Parse(os.Args[2:])
		if reviewCmd.NArg() != 1 {
			log.Fatal("Usage: godoctor-cli review <file-path> [-hint <hint>]")
		}

		filePath := reviewCmd.Arg(0)
		code, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}

		arguments := map[string]any{
			"code": string(code),
		}
		if *reviewHint != "" {
			arguments["hint"] = *reviewHint
		}

		callTool("code_review", arguments)
	case "doc":
		docCmd.Parse(os.Args[2:])
		if docCmd.NArg() != 1 {
			log.Fatal("Usage: godoctor-cli doc <package> [-symbol <symbol>]")
		}

		arguments := map[string]any{
			"package": docCmd.Arg(0),
		}
		if *docSymbol != "" {
			arguments["symbol"] = *docSymbol
		}

		callTool("godoc", arguments)
	default:
		log.Fatalf("Unknown command: %s", os.Args[1])
	}
}

func callTool(toolName string, arguments map[string]any) {
	ctx := context.Background()

	transport := &mcp.StreamableClientTransport{
		Endpoint: "http://localhost:8080",
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "godoctor-cli", Version: "v1.0.0"}, nil)

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatalf("Failed to connect to MCP server: %v", err)
	}
	defer session.Close()

	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      toolName,
		Arguments: arguments,
	})
	if err != nil {
		log.Fatalf("Failed to call tool '%s': %v", toolName, err)
	}

	for _, content := range res.Content {
		if textContent, ok := content.(*mcp.TextContent); ok {
			fmt.Println(textContent.Text)
		}
	}
}

