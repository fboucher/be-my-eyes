package ui

import (
	"github.com/charmbracelet/bubbles/list"
)

// updateSizes updates the sizes of UI components when window is resized
func (m *Model) updateSizes() {
	// Sizes will be calculated dynamically in the view
}

// updateLibraryList updates the library list with current videos
func (m *Model) updateLibraryList() {
	items := make([]list.Item, len(m.videos))
	for i, v := range m.videos {
		items[i] = videoItem{video: v}
	}
	m.libraryList.SetItems(items)

	// Select first item if none selected
	if len(items) > 0 && m.selectedVideo == nil {
		m.selectedVideo = &m.videos[0]
		m.updateDetailView()
	}
}

// updateHistoryList updates the history list with current history
func (m *Model) updateHistoryList() {
	items := make([]list.Item, len(m.history))
	for i, h := range m.history {
		items[i] = historyItem{query: h}
	}
	m.historyList.SetItems(items)
}

// updateSelectedVideo updates the currently selected video from the list
func (m *Model) updateSelectedVideo() {
	selectedItem := m.libraryList.SelectedItem()
	if item, ok := selectedItem.(videoItem); ok {
		m.selectedVideo = &item.video
	}
}

// updateSelectedQuery updates the currently selected query from the list
func (m *Model) updateSelectedQuery() {
	selectedItem := m.historyList.SelectedItem()
	if item, ok := selectedItem.(historyItem); ok {
		m.selectedQuery = &item.query
	}
}

// updateDetailView updates the detail view based on current selection
func (m *Model) updateDetailView() {
	content := m.renderDetails()
	m.detailsView.SetContent(content)
	m.detailsView.GotoTop()
}

// updateMenuList updates the menu list based on active section
func (m *Model) updateMenuList() {
	var items []list.Item

	// Always available actions
	items = append(items, menuItem{
		title:       "Help",
		description: "Show help screen",
		action:      "help",
	})

	items = append(items, menuItem{
		title:       "About",
		description: "About this application",
		action:      "about",
	})

	items = append(items, menuItem{
		title:       "Refresh Library",
		description: "Refresh the video library",
		action:      "refresh",
	})

	// Library-specific actions
	if m.activeSection == LibrarySection && m.selectedVideo != nil {
		items = append(items, menuItem{
			title:       "Ask a Question",
			description: "Ask a question about the selected video",
			action:      "ask",
		})
	}

	items = append(items, menuItem{
		title:       "Quit",
		description: "Exit the application",
		action:      "quit",
	})

	m.menuList.SetItems(items)
}
