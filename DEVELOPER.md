# Information for developers

This page provides more information for developers who would like to look into the codebase, contribute, or build from source.

> #### Dev Container
> 
> If you don't have Go installed locally, this project includes a VS Code dev container configuration:
> 
> 1. Install Docker and VS Code with the Remote-Containers extension
> 2. Open the project in VS Code
> 3. Click "Reopen in Container" when prompted

## Build From Source

```bash
# Clone the repository
git clone https://github.com/fboucher/be-my-eyes.git
cd be-my-eyes

# Download dependencies
make deps

# Option 1: Build and run directly
make run

# Option 2: Build the binary first, then run it
make build
./be-my-eyes

# Option 3: Run with Go directly (no build step)
go run ./cmd/be-my-eyes

# Option 4: Install to $GOPATH/bin (makes it available system-wide)
make install
be-my-eyes  # Now available anywhere in your terminal
```

### Project Structure

```text
be-my-eyes/
├── cmd/be-my-eyes/       # Main application entry point
├── internal/
│   ├── api/              # Reka API client
│   ├── config/           # Configuration management
│   ├── db/               # SQLite database operations
│   ├── models/           # Data models
│   ├── ui/               # TUI components (Bubble Tea)
│   └── version/          # Version information
├── Makefile              # Build automation
└── go.mod                # Go module definition
```

## References

- [Reka AI API Docs](https://link.reka.ai/doc-vision)
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
- [Lipgloss Styling](https://github.com/charmbracelet/lipgloss)
- [Go SQLite3 Driver](https://github.com/mattn/go-sqlite3)
