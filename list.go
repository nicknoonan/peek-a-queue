package main

import (
	"context"

	// "charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type listModel struct {
	bubbleListModel *list.Model
	awsClient *AWSClient
	styles *styles
}

func (lm *listModel) Update(msg tea.Msg) (listModel, tea.Cmd) {
	bubbleListModel, cmd := lm.bubbleListModel.Update(msg)
	return listModel{
		bubbleListModel: &bubbleListModel,
		awsClient: lm.awsClient,
		styles: lm.styles,
	}, cmd
}

func (lm *listModel) View() string {
	return lm.bubbleListModel.View()
}

func (lm *listModel) SetSize(width, height int) {
	lm.bubbleListModel.SetSize(width, height)
}

func (lm *listModel) SetTitleStyle(title lipgloss.Style) {
	lm.bubbleListModel.Styles.Title = title
}

func (lm *listModel) VisibleItems() []list.Item {
	return lm.bubbleListModel.VisibleItems()
}

func (lm *listModel) StopSpinner() {
	lm.bubbleListModel.StopSpinner()
}

func (lm *listModel) IsFiltered() bool {
	return lm.bubbleListModel.IsFiltered()
}

func (lm *listModel) FilterValue() string {
	return lm.bubbleListModel.FilterValue()
}

func (lm *listModel) ResetFilter() {
	lm.bubbleListModel.ResetFilter()
}

func (lm *listModel) SetFilterText(filter string) {
	lm.bubbleListModel.SetFilterText(filter)
}

func (lm *listModel) Select(index int) {
	lm.bubbleListModel.Select(index)
}

func (lm *listModel) Index() int {
	return lm.bubbleListModel.Index()
}

func (lm *listModel) StartSpinner() tea.Cmd {
	return lm.bubbleListModel.StartSpinner()
}

func (lm *listModel) NewStatusMessage(message string) tea.Cmd {
	return lm.bubbleListModel.NewStatusMessage(message)
}

func (lm *listModel) Items() []list.Item {
	return lm.bubbleListModel.Items()
}

func (lm *listModel) SetItem(index int, listItem list.Item) tea.Cmd {
	return lm.bubbleListModel.SetItem(index, listItem)
}

func (lm *listModel) FilterState() list.FilterState {
	return lm.bubbleListModel.FilterState()
}

func (lm *listModel) ToggleSpinner() tea.Cmd {
	return lm.bubbleListModel.ToggleSpinner()
}

func (lm *listModel) ShowTitle() bool {
	return lm.bubbleListModel.ShowTitle()
}

func (lm *listModel) SetShowTitle(value bool) {
	lm.bubbleListModel.SetShowTitle(value)
}

func (lm *listModel) SetShowFilter(value bool) {
	lm.bubbleListModel.SetShowFilter(value)
}

func (lm *listModel) SetFilteringEnabled(value bool) {
	lm.bubbleListModel.SetFilteringEnabled(value)
}

func (lm *listModel) SelectedItem() list.Item {
	return lm.bubbleListModel.SelectedItem()
}

func (lm *listModel) SetShowStatusBar(value bool) {
	lm.bubbleListModel.SetShowStatusBar(value)
}

func (lm *listModel) ShowStatusBar() bool {
	return lm.bubbleListModel.ShowStatusBar()
}

func (lm *listModel) SetShowPagination(value bool) {
	lm.bubbleListModel.SetShowPagination(value)
}

func (lm *listModel) ShowPagination() bool {
	return lm.bubbleListModel.ShowPagination()
}

func (lm *listModel) SetShowHelp(value bool) {
	lm.bubbleListModel.SetShowHelp(value)
}

func (lm *listModel) ShowHelp() bool {
	return lm.bubbleListModel.ShowHelp()
}

type batchItem struct {
	index   int
	setItem item
}

func (lm *listModel) setItemsBatchCmd(setItemsList []batchItem) tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			var cmds []tea.Cmd

			for _, setItem := range setItemsList {
				cmds = append(cmds, lm.SetItem(setItem.index, setItem.setItem))
			}

			return tea.Batch(cmds...)
		},
	)
}

func (lm *listModel) loadPageAttributes(ctx context.Context, listItems ...list.Item) tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, 
		lm.StartSpinner(),
		lm.awsClient.GetQueueAttributesCmd(ctx, lm.Items(), listItems),
	)

	if lm.styles != nil {
		cmds = append(cmds, lm.NewStatusMessage(lm.styles.statusMessage.Render("refreshing...")))
	}

	return tea.Batch(cmds...)
}

// func newItemDelegate(ctx context.Context, awsClient *AWSClient, keys *delegateKeyMap, styles *styles) list.DefaultDelegate {
// 	d := list.NewDefaultDelegate()

// 	// d.UpdateFunc = func(msg tea.Msg, listModel *list.Model) tea.Cmd {
// 	// 	// switch msg := msg.(type) {
// 	// 	// case tea.KeyPressMsg:
// 	// 	// 	switch {
// 	// 	// 	case key.Matches(msg, keys.choose):
// 	// 	// 		return loadPageAttributes(ctx, listModel, styles, awsClient, listModel.SelectedItem())
// 	// 	// 	case key.Matches(msg, keys.refresh):
// 	// 	// 		return loadPageAttributes(ctx, listModel, styles, awsClient, listModel.VisibleItems()...)
// 	// 	// 	}
// 	// 	// }

// 	// 	return nil
// 	// }

// 	help := []key.Binding{keys.choose, keys.refresh}

// 	d.ShortHelpFunc = func() []key.Binding {
// 		return help
// 	}

// 	d.FullHelpFunc = func() [][]key.Binding {
// 		return [][]key.Binding{help}
// 	}

// 	return d
// }

// type delegateKeyMap struct {
// 	choose  key.Binding
// 	refresh key.Binding
// }

// // Additional short help entries. This satisfies the help.KeyMap interface and
// // is entirely optional.
// func (d delegateKeyMap) ShortHelp() []key.Binding {
// 	return []key.Binding{
// 		d.choose,
// 		d.refresh,
// 	}
// }

// // Additional full help entries. This satisfies the help.KeyMap interface and
// // is entirely optional.
// func (d delegateKeyMap) FullHelp() [][]key.Binding {
// 	return [][]key.Binding{
// 		{
// 			d.choose,
// 			d.refresh,
// 		},
// 	}
// }

// func newDelegateKeyMap() *delegateKeyMap {
// 	return &delegateKeyMap{
// 		choose: key.NewBinding(
// 			key.WithKeys("enter"),
// 			key.WithHelp("enter", "refreshes selected item attributes"),
// 		),
// 		refresh: key.NewBinding(
// 			key.WithKeys("r"),
// 			key.WithHelp("r", "refreshes current page attributes"),
// 		),
// 	}
// }
