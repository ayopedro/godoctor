# godoctor

`godoctor` is a project that provides a Model Context Protocol (MCP) server and a command-line interface (CLI) to interact with it.

## User Perspective

### godoctor MCP Server

`godoctor` is an MCP server that exposes a set of tools that can be used by an MCP client.

#### Building the server

To build the server, run the following command from the root of the project:

```bash
go build -o bin/godoctor ./cmd/godoctor
```

This will create the `godoctor` executable in the `bin` directory.

### godoctor-cli

`godoctor-cli` is a command-line interface that allows you to interact with the `godoctor` MCP server.

#### Building the CLI

To build the CLI, run the following command from the root of the project:

```bash
go build -o bin/godoctor-cli ./cmd/godoctor-cli
```

This will create the `godoctor-cli` executable in the `bin` directory.

#### Using the CLI

The CLI exposes the tools available on the `godoctor` server as subcommands.

##### hello_world

The `hello_world` tool returns a "Hello, World!" message. It has an optional `name` argument to customize the greeting.

**Usage:**

```bash
./bin/godoctor-cli hello_world [name=<your-name>]
```

**Example:**

```bash
./bin/godoctor-cli hello_world name=John
```

**Output:**

```
Hello, John!
```

##### godoc

The `godoc` tool retrieves the documentation for a Go package or symbol. It has a mandatory `package` argument and an optional `symbol` argument.

**Usage:**

```bash
./bin/godoctor-cli godoc package=<package-name> [symbol=<symbol-name>]
```

**Examples:**

*   Get documentation for a standard library package:

    ```bash
    ./bin/godoctor-cli godoc package=net/http
    ```

*   Get documentation for a symbol in a standard library package:

    ```bash
    ./bin/godoctor-cli godoc package=fmt symbol=Println
    ```

*   Get documentation for an external package:

    ```bash
    ./bin/godoctor-cli godoc package=github.com/mark3labs/mcp-go/server
    ```

*   Get documentation for a symbol in an external package:

    ```bash
    ./bin/godoctor-cli godoc package=github.com/mark3labs/mcp-go/server symbol=NewMCPServer
    ```

## Developer Perspective

### Project Structure

The project is organized as follows:

*   `cmd/godoctor`: Contains the source code for the `godoctor` MCP server.
*   `cmd/godoctor-cli`: Contains the source code for the `godoctor-cli`.
*   `bin`: Contains the compiled executables.

### Dependencies

The project uses the following main dependencies:

*   **`github.com/mark3labs/mcp-go`**: A Go implementation of the Model Context Protocol, used to create the `godoctor` server.
*   **`github.com/modelcontextprotocol/go-sdk`**: The official Go SDK for the Model Context Protocol, used to create the `godoctor-cli` client.

The `go.mod` file uses a `replace` directive to use a local copy of the `go-sdk`. This is because the version of the `go-sdk` that supports the client implementation is a pre-release version and is not yet available in the main branch.

### How to add a new tool

To add a new tool to the `godoctor` server, follow these steps:

1.  **Open `cmd/godoctor/main.go`:** This is the file that contains the server implementation.
2.  **Add a new `mcpServer.AddTool` call:** This function takes the tool definition and a handler function as arguments.

    The tool definition includes the tool's name, description, and arguments. The handler function is what gets executed when the tool is called.

    **Example:**

    ```go
    mcpServer.AddTool(mcp.NewTool("my-tool",
        mcp.WithDescription("This is my new tool."),
        mcp.WithString("my-arg",
            mcp.Description("This is an argument for my tool."),
            mcp.Required(),
        ),
    ), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        myArg := request.GetString("my-arg", "default-value")
        // Your tool logic here...
        return &mcp.CallToolResult{
            Content: []mcp.Content{
                mcp.TextContent{
                    Type: "text",
                    Text: "This is the result of my tool.",
                },
            },
        }, nil
    })
    ```
3.  **Rebuild the `godoctor` server:**

    ```bash
    go build -o bin/godoctor ./cmd/godoctor
    ```

The `godoctor-cli` is implemented to dynamically discover the available tools on the server, so there is no need to update the CLI to support the new tool. You can immediately use it:

```bash
./bin/godoctor-cli my-tool my-arg=some-value
```
