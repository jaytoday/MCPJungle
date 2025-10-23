package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mcpjungle/mcpjungle/pkg/types"
	"github.com/spf13/cobra"
)

var (
	registerCmdServerName  string
	registerCmdServerURL   string
	registerCmdServerDesc  string
	registerCmdBearerToken string

	registerCmdServerConfigFilePath string
)

var registerMCPServerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register an MCP Server",
	Long: "Register an MCP Server in mcpjungle.\n" +
		"The recommended way is to specify the json configuration file for your mcp server.\n" +
		"Flags are provided for convenience if you want to register a streamable http based server.\n" +
		"But a config file is *required* if you want to register a server using stdio or sse transport.\n" +
		"\nNOTE: A server's name is unique across mcpjungle and must not contain\nany whitespaces, special characters or multiple consecutive underscores '__'.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip flag validation if config file is provided
		if registerCmdServerConfigFilePath != "" {
			return nil
		}
		// Otherwise, validate required flags
		if registerCmdServerName == "" {
			return fmt.Errorf("either supply a configuration file or set the required flag \"name\"")
		}
		if registerCmdServerURL == "" {
			return fmt.Errorf("required flag \"url\" not set")
		}
		return nil
	},
	RunE: runRegisterMCPServer,
	Annotations: map[string]string{
		"group": string(subCommandGroupBasic),
		"order": "2",
	},
}

func init() {
	registerMCPServerCmd.Flags().StringVar(
		&registerCmdServerName,
		"name",
		"",
		"MCP server name",
	)
	registerMCPServerCmd.Flags().StringVar(
		&registerCmdServerURL,
		"url",
		"",
		"URL of the streamable http MCP server (eg- http://localhost:8000/mcp)",
	)
	registerMCPServerCmd.Flags().StringVar(
		&registerCmdServerDesc,
		"description",
		"",
		"Server description",
	)
	registerMCPServerCmd.Flags().StringVar(
		&registerCmdBearerToken,
		"bearer-token",
		"",
		"If provided, MCPJungle will use this token to authenticate with the http MCP server for all requests."+
			" This is useful if the MCP server requires static tokens (eg- your API token) for authentication.",
	)
	registerMCPServerCmd.Flags().StringVarP(
		&registerCmdServerConfigFilePath,
		"conf",
		"c",
		"",
		"Path to a JSON configuration file for the MCP server.\n"+
			"If provided, the mcp server will be registered using the configuration in the file.\n"+
			"All other flags will be ignored.",
	)

	rootCmd.AddCommand(registerMCPServerCmd)
}

func readMcpServerConfig(filePath string) (types.RegisterServerInput, error) {
	var input types.RegisterServerInput

	data, err := os.ReadFile(filePath)
	if err != nil {
		return input, fmt.Errorf("failed to read config file %s: %w", filePath, err)
	}
	// Parse JSON config
	if err := json.Unmarshal(data, &input); err != nil {
		return input, fmt.Errorf("failed to parse config file: %w", err)
	}

	return input, nil
}

func runRegisterMCPServer(cmd *cobra.Command, args []string) error {
	var input types.RegisterServerInput

	if registerCmdServerConfigFilePath == "" {
		// If no config file is provided, use the flags to create the input for server registration
		input = types.RegisterServerInput{
			Name:        registerCmdServerName,
			Transport:   string(types.TransportStreamableHTTP),
			URL:         registerCmdServerURL,
			Description: registerCmdServerDesc,
			BearerToken: registerCmdBearerToken,
		}
	} else {
		// If a config file is provided, read the configuration from the file
		var err error
		input, err = readMcpServerConfig(registerCmdServerConfigFilePath)
		if err != nil {
			return err
		}
	}

	s, err := apiClient.RegisterServer(&input)
	if err != nil {
		return fmt.Errorf("failed to register server: %w", err)
	}
	fmt.Printf("Server %s registered successfully!\n", s.Name)

	if types.McpServerTransport(s.Transport) == types.TransportSSE {
		cmd.Println()
		cmd.Println("This MCP server uses the SSE (Server-sent events) transport.")
		cmd.Println("So its tools will be accessible at the '/sse' endpoint")
		cmd.Println("WARNING: SSE is deprecated, consider migrating this MCP server to streamable http transport.")
	}

	tools, err := apiClient.ListTools(s.Name)
	if err != nil {
		// if we fail to fetch tool list, fail silently because this is not a must-have output
		return nil
	}

	cmd.Println()
	if len(tools) == 0 {
		cmd.Println("This server does not provide any tools.")
		return nil
	}
	cmd.Println("The following tools are now available from this server:")
	for i, tool := range tools {
		cmd.Printf("%d. %s: %s\n\n", i+1, tool.Name, tool.Description)
	}

	prompts, err := apiClient.ListPrompts(s.Name)
	if err != nil {
		return nil
	}
	if len(prompts) > 0 {
		cmd.Println()
		cmd.Println("The following prompts are now available from this server:")
		for i, prompt := range prompts {
			cmd.Printf("%d. %s\n", i+1, prompt.Name)
			if prompt.Description != "" {
				cmd.Printf("   %s\n", prompt.Description)
			}
			cmd.Println()
		}
	}

	return nil
}
