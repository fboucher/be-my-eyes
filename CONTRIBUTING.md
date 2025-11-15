# Contributing to Be My Eyes

Thank you for your interest in contributing to Be My Eyes! This document provides guidelines and information for contributors.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- A Reka AI API key for testing
- Basic understanding of Go and TUI applications

### Development Setup

1. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/be-my-eyes.git
   cd be-my-eyes
   ```

2. **Install Dependencies**
   ```bash
   make deps
   ```

3. **Set Up API Key**
   ```bash
   export REKA_API_KEY=your_api_key_here
   ```

4. **Build the Application**
   ```bash
   make build
   ```

5. **Run the Application**
   ```bash
   ./be-my-eyes
   ```

### Using Dev Container

If you use VS Code, you can use the included dev container:

1. Install Docker and the Remote-Containers extension
2. Open the project in VS Code
3. Click "Reopen in Container" when prompted
4. All dependencies will be installed automatically

## Development Workflow

### Making Changes

1. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Your Changes**
   - Write clear, commented code
   - Follow Go conventions (use `gofmt`)
   - Add comments for complex logic
   - Update documentation if needed

3. **Build and Test**
   ```bash
   make build
   ./be-my-eyes  # Manual testing
   ```

4. **Format Code**
   ```bash
   make fmt
   ```

5. **Commit Your Changes**
   ```bash
   git add .
   git commit -m "Brief description of changes"
   ```

   **Commit Message Guidelines:**
   - Use present tense ("Add feature" not "Added feature")
   - First line should be a brief summary (50 chars or less)
   - Add detailed description if needed after a blank line
   - Reference issues: "Fixes #123" or "Relates to #456"

6. **Push and Create PR**
   ```bash
   git push origin feature/your-feature-name
   ```
   Then create a Pull Request on GitHub.

## Code Style

### Go Conventions

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` to format code (run `make fmt`)
- Run `go vet` to catch common mistakes
- Use meaningful variable and function names

### Comments

Since this project is beginner-friendly, please add helpful comments:

```go
// Good: Explains what and why
// LoadHistory retrieves all query history from the database
// sorted by creation time (newest first) for display in the UI
func (m Model) loadHistory() tea.Cmd {
    return func() tea.Msg {
        history, err := m.database.GetAllHistory()
        return historyLoadedMsg{history: history, err: err}
    }
}

// Avoid: States the obvious
// Gets history
func (m Model) loadHistory() tea.Cmd {
    ...
}
```

### Package Organization

- `cmd/`: Application entry points only
- `internal/`: Private application code
- Keep packages focused and cohesive
- Avoid circular dependencies

## Adding Features

### New UI Components

1. Define component in `internal/ui/model.go`
2. Initialize in `NewModel()`
3. Update in `internal/ui/update.go`
4. Render in `internal/ui/view.go`

### New API Endpoints

1. Add models in `internal/models/models.go`
2. Add client method in `internal/api/client.go`
3. Create command function in `internal/ui/update.go`
4. Handle response in `Update()`

### New Database Operations

1. Add function to `internal/db/db.go`
2. Update schema if needed (in `initSchema()`)
3. Document the change

## Testing

Currently, the project doesn't have automated tests, but contributions adding tests are welcome!

### Manual Testing Checklist

When making changes, please test:

- [ ] Application starts without errors
- [ ] Configuration loading works
- [ ] Database operations work
- [ ] API calls succeed (with valid key)
- [ ] UI renders correctly
- [ ] Navigation works (keyboard and mouse)
- [ ] Dialogs open and close properly
- [ ] Application quits cleanly
- [ ] Screen clears on exit

### Future: Automated Tests

We plan to add:
- Unit tests for each package
- Integration tests for data flow
- UI state transition tests

Contributions in this area are especially welcome!

## Documentation

### Code Documentation

- Add package-level comments for each package
- Document all exported types and functions
- Include examples in comments where helpful
- Update ARCHITECTURE.md for significant changes

### User Documentation

- Update README.md for user-facing changes
- Add examples for new features
- Update key bindings in documentation
- Keep installation instructions current

## Pull Request Process

1. **Before Submitting**
   - Ensure code builds: `make build`
   - Format code: `make fmt`
   - Test manually
   - Update documentation
   - Write clear commit messages

2. **PR Description**
   Include:
   - What problem does this solve?
   - What changes were made?
   - How to test the changes?
   - Screenshots (if UI changes)
   - Related issues

3. **Review Process**
   - Maintainers will review your PR
   - Address any feedback
   - Make requested changes
   - PR will be merged when approved

4. **After Merge**
   - Delete your feature branch
   - Pull latest main
   - Celebrate! 🎉

## Project Structure

See [ARCHITECTURE.md](ARCHITECTURE.md) for detailed architecture information.

```
be-my-eyes/
├── cmd/be-my-eyes/          # Application entry point
├── internal/                # Private application code
│   ├── api/                # API client
│   ├── config/             # Configuration
│   ├── db/                 # Database
│   ├── models/             # Data models
│   └── ui/                 # TUI components
├── .devcontainer/          # Dev container config
├── ARCHITECTURE.md         # Architecture docs
├── CONTRIBUTING.md         # This file
├── Makefile               # Build automation
└── README.md              # User documentation
```

## Common Tasks

### Adding a New Key Binding

1. Update `updateMainView()` in `internal/ui/update.go`:
   ```go
   case "n":  // Your new key
       // Handle the action
       return m, someCommand()
   ```

2. Update footer in `renderFooter()` in `internal/ui/view.go`:
   ```go
   keys := []string{
       "n: new action",
       // ... other keys
   }
   ```

3. Update README.md with new binding
4. Update help screen in `viewHelp()`

### Adding a New View Mode

1. Add constant in `internal/ui/model.go`:
   ```go
   const (
       MainView ViewMode = iota
       YourNewView
   )
   ```

2. Add view function in `internal/ui/view.go`:
   ```go
   func (m Model) viewYourNew() string {
       // Render your view
   }
   ```

3. Update `View()` to include your mode:
   ```go
   case YourNewView:
       return m.viewYourNew()
   ```

4. Add navigation logic in `Update()`

### Adding Configuration Options

1. Update `Config` struct in `internal/config/config.go`:
   ```go
   type Config struct {
       APIKey  string `json:"api_key"`
       YourNew string `json:"your_new"`
   }
   ```

2. Update `Load()` and `Save()` if needed
3. Use in application code
4. Document in README.md

## Debugging

### Enable Logging

You can add debug logging using the standard library:

```go
import "log"

log.Printf("Debug: %+v", someVariable)
```

For TUI apps, log to a file:

```go
f, _ := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
log.SetOutput(f)
```

### Common Issues

**Build fails with "missing dependencies"**
```bash
make deps
```

**Database errors**
```bash
rm ~/.config/be-my-eyes/history.db
# Then restart the app
```

**API connection fails**
- Check API key is set correctly
- Verify internet connection
- Check Reka AI API status

## Communication

### Asking Questions

- Open a GitHub Issue for questions
- Use Discussions for general topics
- Be respectful and constructive

### Reporting Bugs

Include:
- Go version: `go version`
- OS and version
- Steps to reproduce
- Expected vs actual behavior
- Error messages or logs
- Screenshots if relevant

### Suggesting Features

Open an issue with:
- Clear description of the feature
- Use cases and benefits
- Possible implementation approach
- Mockups or examples if applicable

## Code of Conduct

### Our Pledge

We are committed to providing a welcoming and inclusive environment for all contributors.

### Expected Behavior

- Be respectful and professional
- Welcome newcomers and help them learn
- Accept constructive criticism gracefully
- Focus on what's best for the project

### Unacceptable Behavior

- Harassment or discrimination
- Trolling or insulting comments
- Personal attacks
- Publishing others' private information

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

## Recognition

Contributors will be recognized in:
- GitHub contributors page
- Release notes
- Special thanks in README (for significant contributions)

## Getting Help

If you need help contributing:

1. Check [ARCHITECTURE.md](ARCHITECTURE.md) for technical details
2. Look at existing code for examples
3. Open an issue with questions
4. Ask in GitHub Discussions

## Thank You!

Every contribution, no matter how small, is valuable and appreciated. Thank you for helping make Be My Eyes better!

---

Happy coding! 🚀
