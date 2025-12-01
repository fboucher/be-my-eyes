package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fboucher/be-my-eyes/internal/api"
	"github.com/fboucher/be-my-eyes/internal/db"
	"github.com/fboucher/be-my-eyes/internal/models"
)

// Section represents which section is currently active in the left column
type Section int

const (
	StatusSection Section = iota
	LibrarySection
	HistorySection
)

// ViewMode represents the current view state
type ViewMode int

const (
	MainView ViewMode = iota
	QuestionDialogView
	MenuView
	HelpView
	AboutView
	UploadDialogView
)

// Model represents the TUI application state
type Model struct {
	// API and database
	apiClient *api.Client
	database  *db.DB

	// UI state
	width         int
	height        int
	activeSection Section
	viewMode      ViewMode

	// Components
	spinner       spinner.Model
	libraryList   list.Model
	historyList   list.Model
	detailsView   viewport.Model
	questionInput textarea.Model
	menuList      list.Model

	// Data
	videos        []models.Video
	history       []models.QueryHistory
	selectedVideo *models.Video
	selectedQuery *models.QueryHistory

	// Status
	statusMessage string
	isLoading     bool
	err           error

	// Upload dialog state
	uploadTitleInput textarea.Model
	uploadURLInput   textarea.Model
	uploadFocus      int // 0: title, 1: url
}

// videoItem implements list.Item for the library list
type videoItem struct {
	video models.Video
}

func (v videoItem) Title() string {
	title := v.video.Metadata.Title
	if title == "" {
		title = v.video.VideoID
	}
	return title
}

func (v videoItem) Description() string {
	status := strings.ToUpper(v.video.IndexingStatus)
	duration := fmt.Sprintf("%.1fs", v.video.Metadata.Duration)
	return fmt.Sprintf("%s • %s", status, duration)
}

func (v videoItem) FilterValue() string {
	return v.video.Metadata.Title
}

// historyItem implements list.Item for the history list
type historyItem struct {
	query models.QueryHistory
}

func (h historyItem) Title() string {
	// Truncate question to fit in the list
	question := h.query.Question
	if len(question) > 40 {
		question = question[:37] + "..."
	}
	return "Q: " + question
}

func (h historyItem) Description() string {
	if h.query.Error != nil && *h.query.Error != "" {
		return "❌ Error"
	}
	if h.query.Answer == "" {
		return "No content"
	}
	// Count sections or show truncated answer
	return "Response available"
}

func (h historyItem) FilterValue() string {
	return h.query.Question
}

// menuItem implements list.Item for the menu
type menuItem struct {
	title       string
	description string
	action      string
}

func (m menuItem) Title() string       { return m.title }
func (m menuItem) Description() string { return m.description }
func (m menuItem) FilterValue() string { return m.title }

// NewModel creates a new TUI model
func NewModel(apiClient *api.Client, database *db.DB) Model {
	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Initialize library list
	libraryDelegate := list.NewDefaultDelegate()
	libraryList := list.New([]list.Item{}, libraryDelegate, 0, 0)
	libraryList.Title = "Videos"
	libraryList.SetShowStatusBar(false)
	libraryList.SetFilteringEnabled(false)
	libraryList.Styles.Title = lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(lipgloss.Color("230")).
		Bold(true)

	// Initialize history list
	historyDelegate := list.NewDefaultDelegate()
	historyList := list.New([]list.Item{}, historyDelegate, 0, 0)
	historyList.Title = "History"
	historyList.SetShowStatusBar(false)
	historyList.SetFilteringEnabled(false)
	historyList.Styles.Title = lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(lipgloss.Color("230")).
		Bold(true)

	// Initialize details viewport
	detailsView := viewport.New(0, 0)

	// Initialize question input
	questionInput := textarea.New()
	questionInput.Placeholder = "Type your question here..."
	questionInput.Focus()

	// Initialize menu list
	menuDelegate := list.NewDefaultDelegate()
	menuList := list.New([]list.Item{}, menuDelegate, 0, 0)
	menuList.Title = "Menu"
	menuList.SetShowStatusBar(false)
	menuList.SetFilteringEnabled(false)
	menuList.Styles.Title = lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(lipgloss.Color("230")).
		Bold(true)

	// Initialize upload inputs
	uploadTitleInput := textarea.New()
	uploadTitleInput.Placeholder = "Enter video title..."
	uploadURLInput := textarea.New()
	uploadURLInput.Placeholder = "Enter video URL..."
	uploadTitleInput.Focus()
	uploadURLInput.Blur()

	return Model{
		apiClient:        apiClient,
		database:         database,
		activeSection:    LibrarySection,
		viewMode:         MainView,
		spinner:          s,
		libraryList:      libraryList,
		historyList:      historyList,
		detailsView:      detailsView,
		questionInput:    questionInput,
		menuList:         menuList,
		videos:           []models.Video{},
		history:          []models.QueryHistory{},
		statusMessage:    "Disconnected",
		isLoading:        false,
		uploadTitleInput: uploadTitleInput,
		uploadURLInput:   uploadURLInput,
		uploadFocus:      0,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.loadHistory(),
		m.testConnection(),
	)
}
