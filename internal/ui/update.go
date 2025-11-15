package ui

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fboucher/be-my-eyes/internal/models"
)

// Messages for async operations

// historyLoadedMsg is sent when history is loaded from database
type historyLoadedMsg struct {
	history []models.QueryHistory
	err     error
}

// videosLoadedMsg is sent when videos are loaded from API
type videosLoadedMsg struct {
	videos []models.Video
	err    error
}

// questionAskedMsg is sent when a question is asked
type questionAskedMsg struct {
	response *models.QAResponse
	err      error
}

// connectionTestedMsg is sent when connection test completes
type connectionTestedMsg struct {
	success bool
	err     error
}

// loadHistory loads history from the database
func (m Model) loadHistory() tea.Cmd {
	return func() tea.Msg {
		history, err := m.database.GetAllHistory()
		return historyLoadedMsg{history: history, err: err}
	}
}

// testConnection tests the API connection
func (m Model) testConnection() tea.Cmd {
	return func() tea.Msg {
		// Try to get an empty list of videos to test connection
		_, err := m.apiClient.GetVideos([]string{})
		return connectionTestedMsg{success: err == nil, err: err}
	}
}

// refreshLibrary refreshes the video library from API
func (m Model) refreshLibrary() tea.Cmd {
	return func() tea.Msg {
		// For now, we'll get videos from the history
		// In a real app, you'd maintain a list of video IDs
		videoIDs := make(map[string]bool)
		for _, h := range m.history {
			videoIDs[h.VideoID] = true
		}

		ids := make([]string, 0, len(videoIDs))
		for id := range videoIDs {
			ids = append(ids, id)
		}

		if len(ids) == 0 {
			return videosLoadedMsg{videos: []models.Video{}, err: nil}
		}

		response, err := m.apiClient.GetVideos(ids)
		if err != nil {
			return videosLoadedMsg{videos: nil, err: err}
		}

		return videosLoadedMsg{videos: response.Results, err: nil}
	}
}

// askQuestion asks a question about the current video
func (m Model) askQuestion(question string) tea.Cmd {
	if m.selectedVideo == nil {
		return nil
	}

	videoID := m.selectedVideo.VideoID

	return func() tea.Msg {
		response, err := m.apiClient.AskQuestion(videoID, question)
		if err != nil {
			return questionAskedMsg{response: nil, err: err}
		}

		// Parse the chat response to extract the answer
		var chatData struct {
			Sections []struct {
				SectionContent string `json:"section_content"`
			} `json:"sections"`
		}
		
		answer := response.ChatResponse
		if err := json.Unmarshal([]byte(response.ChatResponse), &chatData); err == nil {
			if len(chatData.Sections) > 0 {
				answer = chatData.Sections[0].SectionContent
			}
		}

		// Save to database
		if err := m.database.SaveQuery(videoID, question, answer, response.Error, response.Status); err != nil {
			return questionAskedMsg{response: response, err: err}
		}

		return questionAskedMsg{response: response, err: nil}
	}
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSizes()

	case tea.KeyMsg:
		// Handle global keys
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Handle view-specific keys
		switch m.viewMode {
		case MainView:
			return m.updateMainView(msg)
		case QuestionDialogView:
			return m.updateQuestionDialog(msg)
		case MenuView:
			return m.updateMenuView(msg)
		case HelpView:
			return m.updateHelpView(msg)
		case AboutView:
			return m.updateAboutView(msg)
		}

	case tea.MouseMsg:
		if m.viewMode == MainView {
			return m.handleMouse(msg)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case historyLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.statusMessage = fmt.Sprintf("Error loading history: %v", msg.err)
		} else {
			m.history = msg.history
			m.updateHistoryList()
			
			// Automatically load library if there are video IDs in history
			if len(m.history) > 0 {
				m.isLoading = true
				m.statusMessage = "Loading library..."
				cmds = append(cmds, m.refreshLibrary())
			}
		}

	case videosLoadedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.err = msg.err
			m.statusMessage = fmt.Sprintf("Error loading videos: %v", msg.err)
		} else {
			m.videos = msg.videos
			m.updateLibraryList()
			m.statusMessage = "Connected"
		}

	case questionAskedMsg:
		m.isLoading = false
		m.viewMode = MainView
		if msg.err != nil {
			m.err = msg.err
			m.statusMessage = fmt.Sprintf("Error asking question: %v", msg.err)
		} else {
			m.statusMessage = "Question answered"
			// Reload history
			cmds = append(cmds, m.loadHistory())
		}

	case connectionTestedMsg:
		if msg.success {
			m.statusMessage = "Connected"
			// If we have history with video IDs, the library will be loaded automatically
			// when historyLoadedMsg is processed
		} else {
			m.statusMessage = "Disconnected"
			if msg.err != nil {
				m.err = msg.err
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// updateMainView handles key input in the main view
func (m Model) updateMainView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "tab":
		// Switch active section
		m.activeSection = (m.activeSection + 1) % 3
		m.updateDetailView()

	case "r":
		// Refresh library
		m.isLoading = true
		m.statusMessage = "Refreshing..."
		cmds = append(cmds, m.refreshLibrary())

	case "a":
		// Ask question
		if m.selectedVideo != nil {
			m.viewMode = QuestionDialogView
			m.questionInput.Reset()
			m.questionInput.Focus()
		}

	case "x":
		// Open menu
		m.viewMode = MenuView
		m.updateMenuList()

	case "?":
		// Show help
		m.viewMode = HelpView

	case "up", "k":
		// Navigate up in active section
		switch m.activeSection {
		case LibrarySection:
			m.libraryList, _ = m.libraryList.Update(msg)
			m.updateSelectedVideo()
		case HistorySection:
			m.historyList, _ = m.historyList.Update(msg)
			m.updateSelectedQuery()
		}
		m.updateDetailView()

	case "down", "j":
		// Navigate down in active section
		switch m.activeSection {
		case LibrarySection:
			m.libraryList, _ = m.libraryList.Update(msg)
			m.updateSelectedVideo()
		case HistorySection:
			m.historyList, _ = m.historyList.Update(msg)
			m.updateSelectedQuery()
		}
		m.updateDetailView()

	case "enter":
		// Select item
		switch m.activeSection {
		case LibrarySection:
			m.updateSelectedVideo()
		case HistorySection:
			m.updateSelectedQuery()
		}
		m.updateDetailView()
	}

	return m, tea.Batch(cmds...)
}

// updateQuestionDialog handles input in the question dialog
func (m Model) updateQuestionDialog(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg.String() {
	case "esc":
		m.viewMode = MainView
		return m, nil

	case "ctrl+s":
		// Submit question
		question := m.questionInput.Value()
		if question != "" {
			m.isLoading = true
			m.statusMessage = "Asking question..."
			cmds = append(cmds, m.askQuestion(question))
		}
		return m, tea.Batch(cmds...)
	}

	// Update textarea
	var cmd tea.Cmd
	m.questionInput, cmd = m.questionInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// updateMenuView handles input in the menu view
func (m Model) updateMenuView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg.String() {
	case "esc", "x":
		m.viewMode = MainView
		return m, nil

	case "enter":
		// Execute menu action
		selected := m.menuList.SelectedItem()
		if item, ok := selected.(menuItem); ok {
			switch item.action {
			case "quit":
				return m, tea.Quit
			case "help":
				m.viewMode = HelpView
			case "about":
				m.viewMode = AboutView
			case "ask":
				if m.selectedVideo != nil {
					m.viewMode = QuestionDialogView
					m.questionInput.Reset()
					m.questionInput.Focus()
				}
			case "refresh":
				m.isLoading = true
				m.statusMessage = "Refreshing..."
				cmds = append(cmds, m.refreshLibrary())
				m.viewMode = MainView
			}
		}
		return m, tea.Batch(cmds...)
	}

	// Update list
	var cmd tea.Cmd
	m.menuList, cmd = m.menuList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// updateHelpView handles input in the help view
func (m Model) updateHelpView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, key.NewBinding(key.WithKeys("esc", "q", "?"))) {
		m.viewMode = MainView
	}
	return m, nil
}

// updateAboutView handles input in the about view
func (m Model) updateAboutView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, key.NewBinding(key.WithKeys("esc", "q"))) {
		m.viewMode = MainView
	}
	return m, nil
}

// handleMouse handles mouse events
func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Basic mouse support - clicking in different sections
	// This is a simplified implementation
	return m, nil
}
