# Architecture Documentation

## Overview

Be My Eyes is a Terminal User Interface (TUI) application built with Go that provides an interactive way to work with the Reka Vision AI API. The application follows a clean architecture pattern with separation of concerns.

## Project Structure

```
be-my-eyes/
├── cmd/be-my-eyes/          # Application entry point
│   └── main.go             # Main function, initialization
├── internal/               # Private application code
│   ├── api/               # External API integration
│   │   └── client.go      # Reka API client implementation
│   ├── config/            # Configuration management
│   │   └── config.go      # Config loading/saving
│   ├── db/                # Database layer
│   │   └── db.go          # SQLite operations
│   ├── models/            # Data models
│   │   └── models.go      # Shared data structures
│   └── ui/                # User interface
│       ├── model.go       # TUI state management
│       ├── update.go      # Event handling & business logic
│       ├── view.go        # Rendering logic
│       └── helpers.go     # UI helper functions
├── .devcontainer/         # Development container config
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
├── Makefile              # Build automation
└── README.md             # User documentation
```

## Architecture Layers

### 1. Entry Point (`cmd/be-my-eyes/main.go`)

The main function is responsible for:
- Loading configuration (API key)
- Initializing the API client
- Opening the database connection
- Creating and running the TUI application
- Handling cleanup on exit

### 2. Configuration Layer (`internal/config/`)

Manages application configuration:
- Loads API key from config file or environment variable
- Stores configuration in `~/.config/be-my-eyes/config.json`
- Provides API key validation

**Key Functions:**
- `Load()`: Read configuration from disk
- `Save()`: Write configuration to disk
- `EnsureAPIKey()`: Validate API key is available

### 3. Database Layer (`internal/db/`)

Handles persistent storage using SQLite:
- Stores query history (questions and answers)
- Provides CRUD operations for history
- Manages database schema initialization

**Key Functions:**
- `Open()`: Initialize database connection
- `SaveQuery()`: Store a Q&A pair
- `GetAllHistory()`: Retrieve all history
- `GetHistoryByVideoID()`: Get history for specific video

**Schema:**
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

### 4. API Client Layer (`internal/api/`)

Communicates with the Reka Vision AI API:
- Makes HTTP requests with API key authentication
- Handles JSON serialization/deserialization
- Provides type-safe API methods

**Key Functions:**
- `GetVideos(videoIDs)`: Fetch video metadata
- `AskQuestion(videoID, question)`: Submit Q&A request

**API Endpoints:**
- `POST /videos/get`: Get video information
- `POST /qa/chat`: Ask questions about videos

### 5. Models Layer (`internal/models/`)

Defines shared data structures:
- `Video`: Video metadata and indexing status
- `QueryHistory`: Stored Q&A pairs
- `QARequest`/`QAResponse`: API request/response types

### 6. UI Layer (`internal/ui/`)

Implements the TUI using Bubble Tea framework:

#### `model.go` - State Management
- Defines the `Model` struct (application state)
- Manages UI components (lists, viewports, inputs)
- Stores videos, history, and selections

#### `update.go` - Event Handling
- Processes keyboard input
- Handles mouse events
- Manages async operations (API calls, DB queries)
- Implements view mode switching

**Key Message Types:**
- `historyLoadedMsg`: History fetched from DB
- `videosLoadedMsg`: Videos fetched from API
- `questionAskedMsg`: Q&A completed
- `connectionTestedMsg`: Connection status checked

#### `view.go` - Rendering
- Renders different view modes (main, dialog, menu, help, about)
- Implements layout (40-60 column split)
- Styles components using Lipgloss

**View Modes:**
- `MainView`: Primary interface
- `QuestionDialogView`: Question input dialog
- `MenuView`: Action menu
- `HelpView`: Help screen
- `AboutView`: About screen

#### `helpers.go` - UI Utilities
- Updates component sizes on window resize
- Synchronizes lists with data
- Updates detail view content

## Data Flow

### Application Startup
```
main.go
  → Load config (API key)
  → Create API client
  → Open database
  → Initialize UI model
  → Start Bubble Tea program
```

### Refreshing Library
```
User presses 'r'
  → update.go: Set loading state
  → Start async API call
  → API client fetches videos
  → videosLoadedMsg received
  → Update model.videos
  → Refresh library list
  → Update status to "Connected"
```

### Asking a Question
```
User presses 'a'
  → Switch to QuestionDialogView
  → User types question
  → User presses Ctrl+S
  → Start async API call
  → API client submits question
  → Parse response
  → Save to database
  → questionAskedMsg received
  → Reload history from DB
  → Switch back to MainView
```

### Navigation Flow
```
User presses Tab
  → Cycle active section (Status → Library → History)
  → Update detail view for new selection
  
User presses ↑/↓
  → Navigate within active section list
  → Update selected item
  → Update detail view
```

## UI Component Tree

```
Model (Root)
├── Spinner (Status loading indicator)
├── Library List
│   └── Video Items
├── History List
│   └── Query Items
├── Details Viewport
├── Question Input (Textarea)
└── Menu List
    └── Menu Items
```

## State Management

The `Model` struct maintains all application state:

```go
type Model struct {
    // External dependencies
    apiClient *api.Client
    database  *db.DB

    // UI state
    width, height    int
    activeSection    Section
    viewMode         ViewMode

    // Components
    spinner          spinner.Model
    libraryList      list.Model
    historyList      list.Model
    detailsView      viewport.Model
    questionInput    textarea.Model
    menuList         list.Model

    // Data
    videos           []models.Video
    history          []models.QueryHistory
    selectedVideo    *models.Video
    selectedQuery    *models.QueryHistory

    // Status
    statusMessage    string
    isLoading        bool
    err              error
}
```

## Async Operations

The application uses Bubble Tea's command pattern for async operations:

1. **Tea Command**: Function that returns `tea.Msg`
2. **Goroutine**: Runs in background
3. **Message**: Sent back to Update() when complete
4. **Update**: Processes message and updates state

Example:
```go
func (m Model) refreshLibrary() tea.Cmd {
    return func() tea.Msg {
        // This runs in a goroutine
        videos, err := m.apiClient.GetVideos(ids)
        return videosLoadedMsg{videos: videos, err: err}
    }
}
```

## Error Handling

Errors are handled at each layer:

1. **API Layer**: Returns errors from HTTP calls
2. **DB Layer**: Returns errors from SQL operations
3. **UI Layer**: 
   - Stores errors in model.err
   - Displays in status message
   - Continues operation (doesn't crash)

## Configuration Files

### `~/.config/be-my-eyes/config.json`
```json
{
  "api_key": "your_api_key_here"
}
```

### `~/.config/be-my-eyes/history.db`
SQLite database storing query history

## Key Bindings

| Key | Action | Context |
|-----|--------|---------|
| q | Quit | Main view |
| r | Refresh library | Main view |
| a | Ask question | Main view (video selected) |
| x | Open menu | Main view |
| ? | Help | Main view |
| Tab | Switch section | Main view |
| ↑/↓ | Navigate | Main view |
| Enter | Select | Main view |
| Ctrl+S | Submit | Question dialog |
| Esc | Cancel | Dialogs |

## Design Patterns

### Model-View-Update (MVU)
The Bubble Tea framework uses the Elm architecture:
- **Model**: Application state
- **Update**: Event handlers that modify state
- **View**: Pure rendering functions

### Repository Pattern
Database operations are abstracted in the `db` package, providing a clean interface for data access.

### Client Pattern
API operations are encapsulated in the `api.Client`, hiding HTTP implementation details.

## Dependencies

### Core Dependencies
- `github.com/charmbracelet/bubbletea`: TUI framework
- `github.com/charmbracelet/bubbles`: UI components
- `github.com/charmbracelet/lipgloss`: Styling
- `github.com/mattn/go-sqlite3`: SQLite driver

### Why These Libraries?

**Bubble Tea**: Industry-standard TUI framework for Go
- Clean MVU architecture
- Great terminal handling
- Active community

**Bubbles**: Pre-built components
- List, Viewport, Textarea, Spinner
- Consistent styling
- Keyboard/mouse support

**Lipgloss**: Terminal styling
- Declarative style definitions
- Layout helpers
- Color support

**go-sqlite3**: Mature SQLite driver
- CGo-based for performance
- Standard database/sql interface

## Testing Strategy

While tests aren't currently implemented, here's the recommended approach:

### Unit Tests
- `api/`: Mock HTTP responses
- `db/`: Use in-memory SQLite
- `config/`: Use temp directories
- `models/`: Test JSON marshaling

### Integration Tests
- Test full API → DB flow
- Test UI state transitions
- Test async message handling

### Manual Testing
1. Start app without API key → Should show error
2. Set API key → Should connect
3. Refresh library → Should load videos
4. Ask question → Should save to history
5. Navigate sections → Should update details
6. Open menu → Should show actions
7. Press help → Should show help screen
8. Quit → Should clear screen

## Extension Points

### Adding New Actions
1. Add key binding in `updateMainView()`
2. Implement command function
3. Add message type
4. Handle message in `Update()`
5. Update footer with new binding

### Adding New Views
1. Add view mode constant
2. Implement view function in `view.go`
3. Add switch case in `View()`
4. Add navigation logic in `Update()`

### Adding New API Endpoints
1. Add method to `api.Client`
2. Define request/response models
3. Create command function in UI
4. Add message type
5. Handle in Update()

## Performance Considerations

- **Lazy Loading**: Load videos on-demand
- **Pagination**: History uses indexes for fast queries
- **Viewport**: Only renders visible content
- **Debouncing**: Could be added for search/filter

## Security Considerations

- API key stored with 0600 permissions
- Config directory has 0755 permissions
- No API key in logs or error messages
- Database has no sensitive data beyond query text

## Future Improvements

1. **Video Upload**: Add upload functionality
2. **Search/Filter**: Add filtering to lists
3. **Themes**: Customizable color schemes
4. **Export**: Export history to JSON/CSV
5. **Pagination**: Handle large video libraries
6. **Notifications**: Toast messages for actions
7. **Video Playback**: Integrate with media players
8. **Multi-Select**: Batch operations on videos
