# How It Works

This document provides technical details about how Be My Eyes works internally.

## Architecture Overview

Be My Eyes follows a clean architecture pattern with clear separation of concerns:

```text
┌─────────────────────────────────────────────────────────┐
│                     cmd/be-my-eyes                      │
│                   (Application Entry)                   │
└────────────────────────┬────────────────────────────────┘
                         │
         ┌───────────────┼───────────────┐
         │               │               │
         ▼               ▼               ▼
    ┌────────┐      ┌────────┐     ┌────────┐
    │  API   │      │   DB   │     │   UI   │
    │ Client │      │ Layer  │     │  (TUI) │
    └────────┘      └────────┘     └────────┘
         │               │               │
         └───────────────┼───────────────┘
                         ▼
                   ┌──────────┐
                   │  Models  │
                   └──────────┘
```

## Components

### 1. Entry Point (`cmd/be-my-eyes/main.go`)

The main function orchestrates startup:

1. Checks for `--version` or `--help` flags
2. Loads API key from environment or config file
3. Initializes the Reka API client
4. Opens SQLite database connection
5. Creates the TUI model
6. Starts the Bubble Tea program with alternate screen buffer

### 2. Configuration (`internal/config/`)

Manages API key storage:

- **Location**: `~/.config/be-my-eyes/config.json`
- **Permissions**: `0600` (owner read/write only)
- **Format**: JSON with `api_key` field
- **Fallback**: Reads from `REKA_API_KEY` environment variable
- **Auto-save**: Saves env var to config on first run

### 3. Database Layer (`internal/db/`)

Uses SQLite for persistent storage:

#### Database Location

`~/.config/be-my-eyes/history.db`

#### Schema

**`query_history` table:**

```sql
CREATE TABLE IF NOT EXISTS query_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    video_id TEXT NOT NULL,
    video_title TEXT NOT NULL,
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    error TEXT,
    status TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

**`video_clips` table:**

```sql
CREATE TABLE IF NOT EXISTS video_clips (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    query_id INTEGER NOT NULL,
    clip_id TEXT NOT NULL,
    start_time REAL NOT NULL,
    end_time REAL NOT NULL,
    info TEXT NOT NULL,
    FOREIGN KEY (query_id) REFERENCES query_history(id) ON DELETE CASCADE
);
```

**Indexes:**

- `idx_query_history_video_id` on `query_history(video_id)`
- `idx_query_history_created_at` on `query_history(created_at)`
- `idx_video_clips_query_id` on `video_clips(query_id)`

#### Key Operations

- **`SaveQuery()`**: Stores question, answer, video clips in a transaction
- **`GetAllHistory()`**: Retrieves all queries ordered by date (newest first)
- **`GetHistoryByVideoID()`**: Retrieves queries for a specific video
- **`getVideoClips()`**: Internal method to load clips for a query

### 4. API Client (`internal/api/`)

Communicates with the Reka Vision AI API:

#### Base URL

`https://vision-agent.api.reka.ai`

#### Authentication

Uses `X-Api-Key` header for all requests.

#### Endpoints

**`POST /videos/get`**

Retrieves video metadata:

- Request: `{"video_ids": ["vid_123", ...]}` (empty array returns all videos)
- Response: `{"results": [{"video_id": "...", "metadata": {...}, ...}]}`

**`POST /qa/chat`**

Asks questions about videos:

- Request: `{"video_id": "vid_123", "messages": [{"role": "user", "content": "..."}]}`
- Response: `{"chat_response": "...", "status": "success", ...}`

#### Response Processing

The `chat_response` field contains escaped JSON with:

- **Global answer**: Markdown text (section_id: "1", section_type: "markdown")
- **Video clips**: Array of timestamped clips (section_type: "video-clips-info")

The client parses this JSON and extracts:

1. The main answer text
2. Video clips with timestamps and descriptions

### 5. UI Layer (`internal/ui/`)

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) (Elm architecture for Go):

#### Model (`model.go`)

The `Model` struct contains all application state:

- **Dependencies**: API client, database connection
- **UI state**: Window dimensions, active section, view mode
- **Components**: Spinner, lists, viewport, textarea, menu
- **Data**: Videos, history, selected items
- **Status**: Messages, loading state, errors

#### Update (`update.go`)

Handles all events and business logic:

**Message Types:**

- `historyLoadedMsg`: History fetched from database
- `videosLoadedMsg`: Videos fetched from API
- `questionAskedMsg`: Q&A response received
- `connectionTestedMsg`: Connection status checked

**Key Functions:**

- `updateMainView()`: Handles keyboard input in main view
- `updateQuestionDialog()`: Handles question input dialog
- `updateMenuView()`: Handles menu navigation
- `loadHistory()`: Async command to fetch history
- `refreshLibrary()`: Async command to fetch videos
- `askQuestion()`: Async command to submit Q&A

#### View (`view.go`)

Renders the UI using [Lipgloss](https://github.com/charmbracelet/lipgloss):

**View Modes:**

- `MainView`: Primary interface (40/60 split)
- `QuestionDialogView`: Question input dialog
- `MenuView`: Action menu
- `HelpView`: Help screen
- `AboutView`: About screen

**Layout:**

- Left column (40%): Status, Videos, History
- Right column (60%): Details panel
- Footer: Key bindings

#### Helpers (`helpers.go`)

Utility functions for UI operations:

- `updateSizes()`: Adjusts component sizes on window resize
- `updateLibraryList()`: Syncs library list with video data
- `updateHistoryList()`: Syncs history list with query data
- `updateSelectedVideo()`: Updates selected video from list
- `updateSelectedQuery()`: Updates selected query from list
- `updateDetailView()`: Updates detail panel content
- `updateMenuList()`: Populates context-aware menu items

### 6. Models (`internal/models/`)

Shared data structures:

- **`Video`**: Video metadata from API
- **`VideoMetadata`**: Detailed video information
- **`QueryHistory`**: Stored Q&A pair
- **`VideoClip`**: Timestamped video segment
- **`QARequest/QAResponse`**: API request/response types
- **`VideosGetRequest/VideosGetResponse`**: API types

## Data Flow

### Application Startup

```text
1. main.go
   ↓
2. Load config → Get API key
   ↓
3. Create API client
   ↓
4. Open database
   ↓
5. Create UI model → Initialize components
   ↓
6. Start Bubble Tea program
   ↓
7. Model.Init() → Test connection, Load history
   ↓
8. Auto-refresh library on history load
```

### Refreshing Library

```text
User presses 'r'
   ↓
updateMainView() sets loading state
   ↓
refreshLibrary() command spawned
   ↓
Goroutine calls API client
   ↓
API: POST /videos/get (empty body)
   ↓
videosLoadedMsg sent to Update()
   ↓
Update() stores videos, updates list
   ↓
View() re-renders with new data
```

### Asking a Question

```text
User presses 'a'
   ↓
Switch to QuestionDialogView
   ↓
User types question, presses ctrl+s
   ↓
askQuestion() command spawned
   ↓
Goroutine calls API client
   ↓
API: POST /qa/chat
   ↓
Parse chat_response JSON
   ↓
Extract answer and video clips
   ↓
Save to database (transaction)
   ↓
questionAskedMsg sent to Update()
   ↓
Update() reloads history
   ↓
historyLoadedMsg sent
   ↓
View() shows new query in History
```

### Section Navigation

```text
User presses Tab
   ↓
updateMainView() cycles activeSection
   ↓
Videos → History → Videos
   ↓
updateDetailView() called
   ↓
Based on activeSection:
  - Videos: renderVideoDetails()
  - History: renderQueryDetails()
   ↓
Details panel updates
```

## Async Operations

Bubble Tea uses commands for async work:

```go
// Command function returns a message
func (m Model) refreshLibrary() tea.Cmd {
    return func() tea.Msg {
        // This runs in a goroutine
        videos, err := m.apiClient.GetAllVideos()
        return videosLoadedMsg{videos: videos, err: err}
    }
}

// Update receives the message
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case videosLoadedMsg:
        m.videos = msg.videos
        m.updateLibraryList()
        // Continue processing...
    }
}
```

## Error Handling

Errors are handled gracefully at each layer:

1. **API Layer**: Returns errors from HTTP calls
2. **DB Layer**: Returns errors from SQL operations
3. **UI Layer**: 
   - Stores error in `model.err`
   - Displays in status message
   - Continues operation (doesn't crash)

Example error flow:

```text
API call fails
   ↓
Command returns msg with err != nil
   ↓
Update() checks msg.err
   ↓
If error: Set statusMessage, display error
   ↓
Application continues running
```

## Styling System

Uses Lipgloss for declarative styling:

```go
boxStyle = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("240")).
    Padding(0, 1)
```

Active sections get highlighted borders:

```go
if m.activeSection == LibrarySection {
    libraryStyle = activeBoxStyle  // Color 205 (pink)
}
```

## Mouse Support

Enabled with `tea.WithMouseCellMotion()` option:

- Click to select items (basic implementation)
- Scroll in viewports
- Full mouse support planned for future

## Terminal Handling

Uses alternate screen buffer:

- `tea.WithAltScreen()` option
- Screen clears on exit
- No terminal pollution
- Restore previous content

## Performance Considerations

- **Lazy rendering**: Only renders visible items in lists
- **Viewport optimization**: Only renders visible lines in details
- **Indexed queries**: Database uses indexes for fast lookups
- **Async operations**: API calls don't block UI

## Security

- **API key**: Stored with `0600` permissions (owner only)
- **Config directory**: Created with `0755` permissions
- **No logging**: API key never logged or printed
- **Database**: Contains query text but no sensitive data beyond that

## Dependencies

### Direct Dependencies

```go
require (
    github.com/charmbracelet/bubbles
    github.com/charmbracelet/bubbletea
    github.com/charmbracelet/lipgloss
    github.com/mattn/go-sqlite3
)
```

### Why These Dependencies?

- **Bubble Tea**: Industry-standard TUI framework, Elm architecture
- **Bubbles**: Pre-built components (list, viewport, textarea, spinner)
- **Lipgloss**: Declarative terminal styling, layout helpers
- **go-sqlite3**: CGo-based SQLite driver, mature and performant

## Build Process

### Development Build

```bash
make build
# Produces: ./be-my-eyes
```

### Release Build

Uses GoReleaser for cross-platform builds:

```bash
# Triggered by pushing a tag
git tag v1.0.0
git push origin v1.0.0

# GitHub Actions workflow:
# 1. Builds for Linux (amd64, arm64)
# 2. Builds for macOS (amd64, arm64)  
# 3. Creates .deb packages
# 4. Generates checksums
# 5. Creates GitHub release
# 6. Updates Homebrew tap
```

Version is injected at build time:

```bash
go build -ldflags "-X github.com/fboucher/be-my-eyes/internal/version.Version=v1.0.0"
```

## Configuration Files

### `~/.config/be-my-eyes/config.json`

```json
{
  "api_key": "reka_abc123..."
}
```

### `~/.config/be-my-eyes/history.db`

SQLite database with two tables (see schema above).

## Extending the Application

### Adding a New Key Binding

1. Add case in `updateMainView()` or appropriate update function
2. Implement the action (may need new command function)
3. Update `renderFooter()` to show the key
4. Update help screen in `viewHelp()`

### Adding a New API Endpoint

1. Add request/response types in `internal/models/`
2. Add method in `internal/api/client.go`
3. Create command function in `internal/ui/update.go`
4. Add message type for the response
5. Handle message in `Update()`

### Adding a New Database Table

1. Update `initSchema()` in `internal/db/db.go`
2. Add functions to query/insert data
3. Update models in `internal/models/`
4. Use in UI layer as needed

### Adding a New View Mode

1. Add constant in `internal/ui/model.go`
2. Create view function in `internal/ui/view.go`
3. Add case in `View()` switch
4. Add update handler in `Update()`
5. Add navigation logic

## Testing

Currently no automated tests, but here's the recommended approach:

### Unit Tests

- Mock HTTP responses for API client
- Use in-memory SQLite for database tests
- Test parsing logic separately

### Integration Tests

- Test full data flow (API → DB → UI state)
- Test async command handling
- Test error propagation

### Manual Testing Checklist

- [ ] Start without API key → Shows error
- [ ] Start with API key → Connects
- [ ] Press 'r' → Loads videos
- [ ] Select video → Shows details
- [ ] Press 'a' → Opens dialog
- [ ] Submit question → Saves to history
- [ ] Tab navigation → Switches sections
- [ ] Press 'x' → Opens menu
- [ ] Press '?' → Shows help
- [ ] Press 'q' → Quits cleanly

## Troubleshooting

### "No API key found"

- Set `REKA_API_KEY` environment variable
- Or create `~/.config/be-my-eyes/config.json` manually

### "Disconnected" status

- Check internet connection
- Verify API key is valid
- Check Reka API status

### Database errors

- Delete `~/.config/be-my-eyes/history.db` to reset
- Check file permissions
- Ensure SQLite driver is compiled (CGo required)

### Build errors

- Run `make deps` to fetch dependencies
- Ensure Go 1.21+ is installed
- Check that CGo is enabled (for SQLite)

## Future Enhancements

Potential improvements:

- [ ] Video upload functionality
- [ ] Search/filter in lists
- [ ] Export history to JSON/CSV
- [ ] Custom color themes
- [ ] Configurable key bindings
- [ ] Multi-video selection
- [ ] Video playback integration
- [ ] Automated tests
- [ ] Performance profiling

## References

- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
- [Lipgloss Styling](https://github.com/charmbracelet/lipgloss)
- [Reka AI API Docs](https://docs.reka.ai/)
- [Go SQLite3 Driver](https://github.com/mattn/go-sqlite3)
- [Elm Architecture](https://guide.elm-lang.org/architecture/)
