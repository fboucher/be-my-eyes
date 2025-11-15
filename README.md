# Be My Eyes

A Terminal User Interface (TUI) for interacting with the [Reka Vision AI API](https://www.reka.ai/). This application allows you to manage video libraries, ask questions about video content, and view query history through an intuitive terminal interface.

## Features

- 🎬 **Video Library Management**: Browse and manage your indexed videos
- ❓ **Interactive Q&A**: Ask questions about video content using AI
- 📜 **Query History**: Review past questions and answers stored locally
- 🖱️ **Mouse & Keyboard**: Navigate using keyboard shortcuts or mouse
- 💾 **Local Storage**: SQLite database for persistent query history
- 🎨 **Beautiful TUI**: Clean interface inspired by [Lazydocker](https://github.com/jesseduffield/lazydocker)

## Installation

### Prerequisites

- Go 1.21 or later
- A Reka AI API key (get one at [reka.ai](https://www.reka.ai/))

### From Source

```bash
# Clone the repository
git clone https://github.com/fboucher/be-my-eyes.git
cd be-my-eyes

# Build the application
make build

# Install to $GOPATH/bin
make install
```

### Using Package Managers

#### Homebrew (macOS/Linux)

```bash
# Coming soon
brew install fboucher/tap/be-my-eyes
```

#### APT (Debian/Ubuntu)

```bash
# Coming soon
```

#### Yay (Arch Linux)

```bash
# Coming soon
yay -S be-my-eyes
```

## Configuration

Before running the application, you need to configure your Reka API key. You have two options:

### Option 1: Environment Variable

```bash
export REKA_API_KEY=your_api_key_here
```

### Option 2: Configuration File

Create a configuration file at `~/.config/be-my-eyes/config.json`:

```json
{
  "api_key": "your_api_key_here"
}
```

The application will automatically create the configuration directory on first run if it doesn't exist.

## Usage

Run the application:

```bash
be-my-eyes
```

### Navigation

- **Arrow Keys** / **j/k**: Navigate up and down in lists
- **Tab**: Switch between sections (Status, Library, History)
- **Enter**: Select an item
- **Mouse**: Click to select items and navigate

### Key Bindings

- **r**: Refresh the video library
- **a**: Ask a question about the selected video
- **x**: Open the menu
- **?**: Show help screen
- **q**: Quit the application

### Asking Questions

1. Select a video from the Library section
2. Press **a** or use the menu (**x**)
3. Type your question in the dialog
4. Press **Ctrl+S** to submit
5. The answer will appear in the History section

### Menu Actions

The menu (**x** key) provides access to:

- **Ask a Question**: When a video is selected in Library
- **Refresh Library**: Update the video list from the API
- **Help**: Display the help screen
- **About**: Show information about the application
- **Quit**: Exit the application

## Interface Layout

```text
┌─────────── Status ──────────┐ ┌────────── Details ───────────┐
│ ⠋ Connected                 │ │ Title: Example Video         │
└─────────────────────────────┘ │ ID: vid_123456               │
┌─────────── Videos ──────────┐ │ Status: INDEXED              │
│ ▸ Video A    INDEXED • 123s │ │ Duration: 123.4s             │
│   Video B    PROCESSING     │ │                              │
└─────────────────────────────┘ │                              │
┌─────────── History ─────────┐ │                              │
│ Q: "What's covered?"        │ │                              │
│ Q: "Summarize the intro"    │ │                              │
└─────────────────────────────┘ └──────────────────────────────┘
 r: refresh, a: ask, x: menu, q: quit, tab: section, ↑↓: navigate
```

### Left Column (40% width)

- **Status**: Shows connection status and loading indicators
- **Library**: List of videos with their indexing status and duration
- **History**: List of previous questions and answers

### Right Column (60% width)

- **Details**: Shows detailed information about the selected item
  - For videos: metadata, description, resolution, etc.
  - For history: full question and answer text

## Development

### Project Structure

```
be-my-eyes/
├── cmd/
│   └── be-my-eyes/       # Main application entry point
│       └── main.go
├── internal/
│   ├── api/              # Reka API client
│   │   └── client.go
│   ├── config/           # Configuration management
│   │   └── config.go
│   ├── db/               # SQLite database layer
│   │   └── db.go
│   ├── models/           # Data models
│   │   └── models.go
│   └── ui/               # TUI components
│       ├── model.go      # Main model
│       ├── update.go     # Update logic
│       ├── view.go       # View rendering
│       └── helpers.go    # Helper functions
├── .devcontainer/        # VS Code dev container
│   └── devcontainer.json
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### Building from Source

```bash
# Download dependencies
make deps

# Build the application
make build

# Run tests (when available)
make test

# Format code
make fmt

# Build for all platforms
make build-all
```

### Dev Container

This project includes a VS Code dev container configuration for easy development:

1. Install Docker and VS Code with the Remote-Containers extension
2. Open the project in VS Code
3. Click "Reopen in Container" when prompted
4. All dependencies will be installed automatically

### Code Style

The codebase includes extensive comments for developers who are new to Go:

- Each package has clear documentation
- Complex functions include explanatory comments
- All public types and functions are documented

## Database Schema

The application uses SQLite for local storage:

### query_history table

```sql
CREATE TABLE query_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    video_id TEXT NOT NULL,
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    error TEXT,
    status TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

Database location: `~/.config/be-my-eyes/history.db`

## API Integration

This application integrates with the Reka Vision AI API:

- **Videos API**: Retrieve video metadata and indexing status
- **QA Chat API**: Ask questions about video content

### API Endpoints

- `POST /videos/get`: Get video information by IDs
- `POST /qa/chat`: Ask questions about videos

See the [Reka AI documentation](https://docs.reka.ai/) for more details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Credits

- **Creator**: [fboucher](https://github.com/fboucher)
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **UI Components**: [Bubbles](https://github.com/charmbracelet/bubbles)
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Inspiration**: [Lazydocker](https://github.com/jesseduffield/lazydocker)

## Support

If you encounter any issues or have questions, please open an issue on GitHub:

https://github.com/fboucher/be-my-eyes/issues

## Roadmap

- [ ] Video upload functionality
- [ ] Advanced filtering and search
- [ ] Export query history
- [ ] Custom themes
- [ ] Configurable key bindings
- [ ] Multiple API key profiles

---

Built with ❤️ using Go and Bubble Tea
