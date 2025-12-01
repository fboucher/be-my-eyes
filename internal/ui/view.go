// ...existing code...
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	// Box styles
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	activeBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")).
			Padding(0, 1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	helpStyle = lipgloss.NewStyle().
			Padding(1, 2)

	dialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")).
			Padding(1, 2)
)

// View renders the TUI
func (m Model) View() string {
	switch m.viewMode {
	case QuestionDialogView:
		return m.viewQuestionDialog()
	case MenuView:
		return m.viewMenu()
	case HelpView:
		return m.viewHelp()
	case AboutView:
		return m.viewAbout()
	case UploadDialogView:
		return m.viewUploadDialog()
	default:
		return m.viewMain()
	}
}

// viewMain renders the main view
func (m Model) viewMain() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Calculate dimensions (40-60 split)
	leftWidth := int(float64(m.width) * 0.4)
	rightWidth := m.width - leftWidth - 4 // Account for borders and padding

	// Build left column
	leftCol := m.renderLeftColumn(leftWidth)

	// Build right column
	rightCol := m.renderRightColumn(rightWidth)

	// Combine columns
	main := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)

	// Add footer
	footer := m.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Left, main, footer)
}

// renderLeftColumn renders the left column with status, library, and history
func (m Model) renderLeftColumn(width int) string {
	// Status section
	statusHeight := 3
	statusContent := m.renderStatus()
	statusBox := boxStyle.Width(width - 4).Height(statusHeight).Render(statusContent)

	// Calculate remaining height for library and history
	remainingHeight := m.height - statusHeight - 8 // Account for borders, footer, etc.
	libraryHeight := remainingHeight / 2
	historyHeight := remainingHeight - libraryHeight

	// Library section
	m.libraryList.SetSize(width-6, libraryHeight-2)
	libraryContent := m.libraryList.View()
	libraryStyle := boxStyle
	if m.activeSection == LibrarySection {
		libraryStyle = activeBoxStyle
	}
	libraryBox := libraryStyle.Width(width - 4).Height(libraryHeight).Render(libraryContent)

	// History section
	m.historyList.SetSize(width-6, historyHeight-2)
	historyContent := m.historyList.View()
	historyStyle := boxStyle
	if m.activeSection == HistorySection {
		historyStyle = activeBoxStyle
	}
	historyBox := historyStyle.Width(width - 4).Height(historyHeight).Render(historyContent)

	return lipgloss.JoinVertical(lipgloss.Left, statusBox, libraryBox, historyBox)
}

// renderRightColumn renders the right column with details
func (m Model) renderRightColumn(width int) string {
	detailHeight := m.height - 5 // Account for footer

	content := m.renderDetails()
	// Wrap the content to fit the width
	wrappedContent := lipgloss.NewStyle().Width(width - 6).Render(content)
	m.detailsView.SetContent(wrappedContent)
	m.detailsView.Width = width - 4
	m.detailsView.Height = detailHeight - 2

	detailBox := boxStyle.Width(width - 4).Height(detailHeight).Render(
		titleStyle.Render("Details") + "\n" + m.detailsView.View(),
	)

	return detailBox
}

// renderStatus renders the status section
func (m Model) renderStatus() string {
	status := m.statusMessage
	if m.isLoading {
		status = m.spinner.View() + " " + status
	}
	return titleStyle.Render("Status") + "\n" + statusStyle.Render(status)
}

// renderDetails renders the details panel based on what's selected
func (m Model) renderDetails() string {
	switch m.activeSection {
	case LibrarySection:
		if m.selectedVideo != nil {
			return m.renderVideoDetails()
		}
		return "Select a video to see details"

	case HistorySection:
		if m.selectedQuery != nil {
			return m.renderQueryDetails()
		}
		return "Select a query to see details"

	default:
		return "No details available"
	}
}

// renderVideoDetails renders details for the selected video
func (m Model) renderVideoDetails() string {
	v := m.selectedVideo
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Title: %s\n", v.Metadata.Title))
	b.WriteString(fmt.Sprintf("ID: %s\n", v.VideoID))
	b.WriteString(fmt.Sprintf("Status: %s\n", strings.ToUpper(v.IndexingStatus)))
	b.WriteString(fmt.Sprintf("Duration: %.1fs\n\n", v.Metadata.Duration))

	if v.Metadata.Description != "" {
		b.WriteString(fmt.Sprintf("Description:\n%s\n\n", v.Metadata.Description))
	}

	b.WriteString(fmt.Sprintf("Resolution: %dx%d\n", v.Metadata.Width, v.Metadata.Height))
	b.WriteString(fmt.Sprintf("FPS: %.1f\n", v.Metadata.AvgFPS))
	b.WriteString(fmt.Sprintf("Source: %s\n", v.Metadata.Source))

	return b.String()
}

// renderQueryDetails renders details for the selected query
func (m Model) renderQueryDetails() string {
	q := m.selectedQuery
	var b strings.Builder

	b.WriteString("Question:\n")
	b.WriteString(q.Question)
	b.WriteString("\n\n")

	if q.Error != nil && *q.Error != "" {
		b.WriteString("Error:\n")
		b.WriteString(*q.Error)
		b.WriteString("\n")
	} else {
		b.WriteString("Answer:\n\n")
		b.WriteString(q.Answer)
	}

	return b.String()
}

// renderFooter renders the footer with key bindings
func (m Model) renderFooter() string {
	keys := []string{
		"u: upload",
		"r: refresh",
		"a: ask question",
		"x: menu",
		"q: quit",
		"tab: change section",
		"↑↓: navigate/scroll",
	}
	return footerStyle.Render(strings.Join(keys, ", "))
}

// viewQuestionDialog renders the question input dialog
func (m Model) viewQuestionDialog() string {
	title := "Ask a Question"
	if m.selectedVideo != nil {
		title += fmt.Sprintf(" (Video: %s)", m.selectedVideo.Metadata.Title)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(title),
		"",
		m.questionInput.View(),
		"",
		footerStyle.Render("ctrl+s: submit, esc: cancel"),
	)

	dialog := dialogStyle.Render(content)

	// Center the dialog
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)
}

// viewMenu renders the menu
func (m Model) viewMenu() string {
	m.menuList.SetSize(40, 15)
	content := m.menuList.View()

	dialog := dialogStyle.Render(content)

	// Center the dialog
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)
}

// viewHelp renders the help screen
func (m Model) viewHelp() string {
	help := `Be My Eyes - Help

Navigation:
  ↑/↓, j/k    - Navigate lists
  tab         - Switch between sections
  enter       - Select item

Actions:
  r           - Refresh library
  a           - Ask a question about selected video
  u           - Upload video (not yet implemented)
  x           - Open menu
  ?           - Show this help
  q           - Quit

Question Dialog:
  ctrl+s      - Submit question
  esc         - Cancel

Press any key to return to main view.`

	content := helpStyle.Render(help)
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialogStyle.Render(content),
	)
}

// viewAbout renders the about screen
func (m Model) viewAbout() string {
	about := `Be My Eyes - About

A TUI for interacting with Reka Vision AI API

GitHub: https://github.com/fboucher/be-my-eyes
Creator: fboucher

Built with:
  - Bubble Tea (TUI framework)
  - Reka Vision AI API
  - SQLite (local storage)

Press any key to return to main view.`

	content := helpStyle.Render(about)
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialogStyle.Render(content),
	)
}
