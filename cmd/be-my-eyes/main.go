package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fboucher/be-my-eyes/internal/api"
	"github.com/fboucher/be-my-eyes/internal/config"
	"github.com/fboucher/be-my-eyes/internal/db"
	"github.com/fboucher/be-my-eyes/internal/ui"
	"github.com/fboucher/be-my-eyes/internal/version"
)

func main() {
	// Lightweight flag handling for version/help before doing any setup
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--version", "-v":
			fmt.Printf("be-my-eyes %s\n", version.Version)
			return
		case "--help", "-h":
			printHelp()
			return
		}
		if strings.HasPrefix(arg, "version") { // allow 'version' subcommand style
			fmt.Printf("be-my-eyes %s\n", version.Version)
			return
		}
	}

	// Load configuration and ensure API key is available
	apiKey, err := config.EnsureAPIKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nPlease set your Reka API key:\n")
		fmt.Fprintf(os.Stderr, "  export REKA_API_KEY=your_api_key_here\n")
		fmt.Fprintf(os.Stderr, "\nOr add it to ~/.config/be-my-eyes/config.json:\n")
		fmt.Fprintf(os.Stderr, "  {\"api_key\": \"your_api_key_here\"}\n")
		os.Exit(1)
	}

	// Initialize API client
	apiClient := api.NewClient(apiKey)

	// Open database
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Create TUI model
	model := ui.NewModel(apiClient, database)

	// Create program with alternate screen buffer (clears on exit)
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer (clears on exit)
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("be-my-eyes - TUI for interacting with the Reka Vision AI API")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  be-my-eyes [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -h, --help       Show this help message")
	fmt.Println("  -v, --version    Show version information")
	fmt.Println()
	fmt.Println("Environment:")
	fmt.Println("  REKA_API_KEY     Your Reka API key (or use config file)")
	fmt.Println()
	fmt.Println("Config file:")
	fmt.Println("  ~/.config/be-my-eyes/config.json containing {\"api_key\": \"...\"}")
}
