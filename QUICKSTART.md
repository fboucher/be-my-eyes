# Quick Start Guide

This guide will help you get started with Be My Eyes quickly.

## Step 1: Installation

### Option A: Build from Source

```bash
# Clone the repository
git clone https://github.com/fboucher/be-my-eyes.git
cd be-my-eyes

# Build the application
make build

# The binary will be created in the current directory
```

### Option B: Download Pre-built Binary

Coming soon - releases will be available on GitHub.

## Step 2: Get a Reka API Key

1. Visit [reka.ai](https://www.reka.ai/)
2. Sign up for an account
3. Navigate to the API section
4. Create a new API key
5. Copy your API key (you'll need it in the next step)

## Step 3: Configure Your API Key

You have two options:

### Option A: Environment Variable (Temporary)

```bash
export REKA_API_KEY=your_api_key_here
```

This will only work for your current terminal session.

### Option B: Configuration File (Permanent)

```bash
# Create the config directory
mkdir -p ~/.config/be-my-eyes

# Create the config file
cat > ~/.config/be-my-eyes/config.json << EOF
{
  "api_key": "your_api_key_here"
}
EOF

# Secure the config file
chmod 600 ~/.config/be-my-eyes/config.json
```

This will persist across sessions.

## Step 4: Run the Application

```bash
./be-my-eyes
```

Or if you installed it to your PATH:

```bash
be-my-eyes
```

## First Time Use

When you first run the application:

1. The app will test the connection to the Reka API
2. Status will show "Connected" if successful
3. The Library will be empty initially
4. History will be empty initially

## Basic Workflow

### 1. Understanding the Interface

```
┌─────────── Status ──────────┐ ┌────────── Details ───────────┐
│ ⠋ Connected                 │ │ Select an item from the      │
└─────────────────────────────┘ │ left to see details          │
┌─────────── Videos ──────────┐ │                              │
│ (empty)                     │ │                              │
└─────────────────────────────┘ │                              │
┌─────────── History ─────────┐ │                              │
│ (empty)                     │ │                              │
└─────────────────────────────┘ └──────────────────────────────┘
 r: refresh, a: ask, x: menu, q: quit, tab: section, ↑↓: navigate
```

**Left Column (40%):**
- **Status**: Connection status and loading indicators
- **Videos**: Your video library
- **History**: Past questions and answers

**Right Column (60%):**
- **Details**: Information about selected item

**Footer:**
- Key bindings reference

### 2. Adding Videos

Currently, videos are discovered automatically when you ask questions about them. To add a video:

1. Get the video ID from Reka AI (after uploading a video through their API)
2. Ask a question about that video
3. The video will appear in your Library

Future versions will include direct upload functionality.

### 3. Asking Questions

Once you have a video in your library:

1. **Select a Video**
   - Press Tab until the Videos section is highlighted
   - Use ↑/↓ arrows to select a video
   - The Details panel will show video information

2. **Open Question Dialog**
   - Press `a` (or press `x` for menu, then select "Ask a Question")
   - A dialog will appear

3. **Type Your Question**
   - Type your question in the text area
   - Example: "What is happening in this video?"
   - Example: "Summarize the main points"
   - Example: "What objects are visible?"

4. **Submit**
   - Press `Ctrl+S` to submit
   - The app will show "Asking question..." in the Status
   - When complete, the answer will appear in History

5. **View Answer**
   - Press Tab to switch to the History section
   - Your question should be at the top
   - Select it to see the full answer in the Details panel

### 4. Navigation

**Keyboard:**
- `Tab`: Switch between sections (Status → Videos → History → Status)
- `↑`/`↓` or `j`/`k`: Navigate within a list
- `Enter`: Select an item
- `?`: Show help
- `x`: Open menu
- `r`: Refresh library
- `q`: Quit

**Mouse:**
- Click to select items
- Scroll to navigate lists

### 5. Using the Menu

Press `x` to open the menu. Available actions depend on context:

**Always Available:**
- Help: Show keyboard shortcuts and help
- About: Information about the app
- Refresh Library: Update video list from API
- Quit: Exit the application

**When Video Selected:**
- Ask a Question: Open question dialog for selected video

### 6. Viewing History

Your question history is automatically saved:

1. Press Tab to navigate to the History section
2. Use ↑/↓ to browse past questions
3. Press Enter to select a query
4. The full question and answer will appear in the Details panel

History is stored locally in `~/.config/be-my-eyes/history.db`

## Tips and Tricks

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `q` | Quit (from main view) |
| `r` | Refresh library |
| `a` | Ask question (video selected) |
| `x` | Open menu |
| `?` | Show help |
| `Tab` | Switch section |
| `↑`/`↓` | Navigate |
| `j`/`k` | Navigate (Vim-style) |
| `Enter` | Select item |
| `Esc` | Close dialog/menu |
| `Ctrl+C` | Force quit |

### In Question Dialog

| Key | Action |
|-----|--------|
| Type normally | Enter question text |
| `Ctrl+S` | Submit question |
| `Esc` | Cancel and return to main view |

### Best Practices

1. **Be Specific**: Ask clear, specific questions for better results
2. **Review History**: Check past answers before asking similar questions
3. **Organize**: Use the refresh feature to keep your library up to date
4. **Explore**: Try the help screen (`?`) to discover features

### Common Issues

**"No API key found"**
- Set the REKA_API_KEY environment variable
- Or create ~/.config/be-my-eyes/config.json with your key

**"Disconnected" status**
- Check your internet connection
- Verify your API key is valid
- Check Reka AI service status

**Empty library**
- Ask a question about a video (by video ID)
- The library populates from your question history

**Database errors**
- The database is at ~/.config/be-my-eyes/history.db
- You can delete it to start fresh (you'll lose history)

## Example Session

Here's what a typical session might look like:

```bash
# 1. Start the app
./be-my-eyes

# 2. The interface appears:
#    - Status shows "Connected"
#    - Videos section is active (highlighted)
#    - Library is empty
#    - History is empty

# 3. Ask a question about a video (press 'a')
#    Dialog opens: "Ask a Question"
#    Type: "What is shown in this video?"
#    Press Ctrl+S to submit

# 4. Status shows "Asking question..."
#    After a few seconds, Status shows "Question answered"

# 5. Switch to History (press Tab twice)
#    Your question appears at the top
#    Select it to see the answer

# 6. The Details panel shows:
#    Question: What is shown in this video?
#    Answer: [AI-generated response]
#    Status: success
#    Date: 2024-11-15 14:30:00

# 7. Continue asking questions or exploring
#    Press 'q' to quit when done
```

## Next Steps

- Read the full [README.md](README.md) for more details
- Check [ARCHITECTURE.md](ARCHITECTURE.md) for technical information
- See [CONTRIBUTING.md](CONTRIBUTING.md) to contribute
- Visit the [GitHub repository](https://github.com/fboucher/be-my-eyes) for issues and discussions

## Getting Help

If you run into issues:

1. Press `?` in the app for help
2. Check this guide
3. Review the [README.md](README.md)
4. Open an issue on [GitHub](https://github.com/fboucher/be-my-eyes/issues)

## Advanced Usage

### Multiple API Keys

You can switch between different API keys by:

1. Using different environment variables
2. Editing ~/.config/be-my-eyes/config.json
3. Restarting the application

### Database Location

History is stored at: `~/.config/be-my-eyes/history.db`

You can:
- Back it up for safekeeping
- Delete it to start fresh
- Copy it between machines

### Building from Source

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup.

---

Enjoy using Be My Eyes! 🎉
