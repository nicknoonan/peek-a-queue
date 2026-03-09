package main

import (
	"context"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/aws/aws-sdk-go-v2/config"
)

type model struct {
	awsClient     AWSClient
	styles        styles
	darkBG        bool
	width, height int
	list          *listModel
	keys          *listKeyMap
}

func initialModel(ctx context.Context) (*model, error) {
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	awsClient := NewAWSClient(awsConfig)

	// Initialize the model and list.
	m := model{
		awsClient: awsClient,
	}
	m.styles = newStyles(false) // default to dark background styles

	listKeys := newListKeyMap()

	// Setup list.
	queueList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	queueList.Title = "Queues"
	queueList.Styles.Title = m.styles.title
	queueList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.refreshItem,
			listKeys.refreshPage,
		}
	}
	queueList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
			listKeys.refreshItem,
			listKeys.refreshPage,
		}
	}

	m.list = &listModel{
		bubbleListModel: &queueList,
		awsClient:       &awsClient,
		styles:          &m.styles,
	}
	m.keys = listKeys

	return &m, nil
}


type refreshTickMsg time.Time

func refreshTick() tea.Cmd {
	return tea.Tick(10*time.Second, func(t time.Time) tea.Msg {
		return refreshTickMsg(t)
	})
}

type initialLoadMsg string

func initialLoad() tea.Cmd {
	return func() tea.Msg {
		return initialLoadMsg("")
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestBackgroundColor,
		refreshTick(),
		initialLoad(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case refreshTickMsg:
		cmds = append(cmds,
			m.list.refreshItemAttributes(context.TODO(), m.list.VisibleItems()...),
			refreshTick(),
		)
	case initialLoadMsg:
		cmds = append(cmds,
			m.list.StartSpinner(),
			m.awsClient.ListAllQueuesCmd(context.TODO()),
		)
	case queueListMsg:
		m.list.StopSpinner()
		cmds = append(cmds,
			m.list.SetItems(msg),
			m.awsClient.GetQueueAttributesCmd(context.TODO(), msg),
		)
	case tea.BackgroundColorMsg:
		m.darkBG = msg.IsDark()
		m.updateListProperties()
		return m, nil

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.updateListProperties()
		return m, nil
	case queueAttributesMsg:
		m.list.StopSpinner()

		if m.list.IsFiltered() {
			currentFilter := m.list.FilterValue()
			currentIndex := m.list.Index()
			m.list.ResetFilter()
			m.list.SetFilterText(currentFilter)
			m.list.Select(currentIndex)
		}

		if msg.err != nil {
			return m, m.list.NewStatusMessage(m.styles.statusMessage.Render("error: " + msg.err.Error()))
		}

		cmds = append(cmds, m.list.setItemsBatchCmd(msg.itemList))
	case tea.KeyPressMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.refreshItem):
			cmds = append(cmds, m.list.refreshItemAttributes(context.TODO(), m.list.SelectedItem()))
		case key.Matches(msg, m.keys.refreshPage):
			cmds = append(cmds, m.list.refreshItemAttributes(context.TODO(), m.list.VisibleItems()...))

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = &newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	v := tea.NewView(m.styles.app.Render(m.list.View()))
	v.AltScreen = true
	return v
}
