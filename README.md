# Be My Eyes

A Terminal User Interface (TUI) for interacting with the [Reka Vision AI API](https://www.reka.ai/). Ask questions about your videos through an intuitive terminal interface.

## Features

- 🎬 **Video Library Management**: Browse your indexed videos from the Reka API
- ❓ **Interactive Q&A**: Ask questions about video content using AI
- 📜 **Query History**: Review past questions and answers stored locally in SQLite
- 🖱️ **Mouse & Keyboard**: Navigate using keyboard shortcuts or mouse
- 💾 **Local Storage**: SQLite database for persistent query history with video clips
- 🎨 **Beautiful TUI**: Clean interface built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)

## Installation

### Prerequisites

- Go 1.21 or later
- A Reka AI API key (get one at [reka.ai](https://www.reka.ai/))

### From Source

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

### Using Package Managers

#### Homebrew (macOS/Linux)

```bash
brew tap fboucher/tap
brew install be-my-eyes
```

#### APT (Debian/Ubuntu)

Download the `.deb` file from the [latest release](https://github.com/fboucher/be-my-eyes/releases/latest) and install:

```bash
sudo dpkg -i be-my-eyes_*_amd64.deb
```

## Configuration

Before running the application, you need to configure your Reka API key.

### Option 1: Environment Variable

```bash
export REKA_API_KEY=your_api_key_here
```

### Option 2: Configuration File

The application will automatically save your API key from the environment variable to `~/.config/be-my-eyes/config.json` on first run. Alternatively, create it manually:

```json
{
  "api_key": "your_api_key_here"
}
```

## Usage

Run the application:

```bash
be-my-eyes
```

### Key Bindings

#### Main View

| Key | Action |
|-----|--------|
| `q` | Quit the application |
| `r` | Refresh video library from API |
| `a` | Ask a question about the selected video |
| `x` | Open the menu |
| `?` | Show help screen |
| `tab` | Switch between sections (Videos → History → Videos) |
| `↑` / `↓` | Navigate up/down in lists |
| `j` / `k` | Navigate up/down (Vim-style) |
| `enter` | Select an item |
| `ctrl+c` | Force quit |

#### Question Dialog

| Key | Action |
|-----|--------|
| `ctrl+s` | Submit the question |
| `esc` | Cancel and return to main view |

#### Menu, Help, About Screens

| Key | Action |
|-----|--------|
| `esc` | Return to main view |
| `enter` | Execute selected menu action |
| `↑` / `↓` | Navigate menu items |

## Interface Layout

```text
┌─────────── Status ──────────┐ ┌────────── Details ───────────┐
│ Connected                   │ │ Title: Example Video         │
└─────────────────────────────┘ │ ID: vid_123456               │
┌─────────── Videos ──────────┐ │ Status: INDEXED              │
│ ▸ Video A    INDEXED • 123s │ │ Duration: 123.4s             │
│   Video B    PROCESSING     │ │ Resolution: 1920x1080        │
└─────────────────────────────┘ │ FPS: 30.0                    │
┌─────────── History ─────────┐ │                              │
│ Q: "What's in this video?"  │ │                              │
│ Q: "Summarize the content"  │ │                              │
└─────────────────────────────┘ └──────────────────────────────┘
 r: refresh, a: ask question, x: menu, q: quit, tab: section, ↑↓: navigate
```

### Asking Questions

1. Press `r` to refresh and load your video library from the Reka API
2. Use `↑` / `↓` to select a video (Videos section should be active by default)
3. Press `a` to open the question dialog
4. Type your question
5. Press `ctrl+s` to submit
6. The answer will appear in the History section
7. Press `tab` to switch to History and view the full response

### Menu Actions

Press `x` to open the menu. Available actions:

- **Ask a Question**: When a video is selected in the Videos section
- **Refresh Library**: Update the video list from the Reka API
- **Help**: Display the help screen with keyboard shortcuts
- **About**: Show information about the application
- **Quit**: Exit the application

## Development

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

See [HOW-IT-WORKS.md](HOW-IT-WORKS.md) for detailed technical documentation.

### Building from Source

```bash
# Download dependencies
make deps

# Build the application
make build

# Format code
make fmt

# Build for all platforms
make build-all
```

### Dev Container

This project includes a VS Code dev container configuration:

1. Install Docker and VS Code with the Remote-Containers extension
2. Open the project in VS Code
3. Click "Reopen in Container" when prompted

## API Integration

This application uses the Reka Vision AI API:

- `POST /videos/get`: Retrieve video metadata and indexing status
- `POST /qa/chat`: Ask questions about videos with timestamped clips

See [HOW-IT-WORKS.md](HOW-IT-WORKS.md) for details on the database schema and API integration.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Credits

- **Creator**: [fboucher](https://github.com/fboucher)
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **UI Components**: [Bubbles](https://github.com/charmbracelet/bubbles)
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss)

## Support

If you encounter any issues or have questions, please [open an issue on GitHub](https://github.com/fboucher/be-my-eyes/issues).

## License

MIT License - see [LICENSE](LICENSE) for details

---

Built with ❤️ using Go and Bubble Tea
