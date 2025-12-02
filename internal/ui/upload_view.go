package ui

import "github.com/charmbracelet/lipgloss"

var focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true)

// viewUploadDialog renders the upload input dialog
func (m Model) viewUploadDialog() string {
	// Highlight the focused field
	titleLabel := "Title:"
	urlLabel := "URL:"
	titleInput := m.uploadTitleInput.View()
	urlInput := m.uploadURLInput.View()
	if m.uploadFocus == 0 {
		titleLabel = focusedStyle.Render(titleLabel)
		titleInput = focusedStyle.Render(titleInput)
	} else {
		urlLabel = focusedStyle.Render(urlLabel)
		urlInput = focusedStyle.Render(urlInput)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Upload a Video"),
		"",
		lipgloss.JoinHorizontal(lipgloss.Top, titleLabel, " ", titleInput),
		lipgloss.JoinHorizontal(lipgloss.Top, urlLabel, " ", urlInput),
		"",
		footerStyle.Render("tab/shift+tab: switch, enter: upload, esc: cancel"),
	)
	dialog := dialogStyle.Render(content)
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)
}
