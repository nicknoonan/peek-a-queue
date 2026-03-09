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
	awsClient       *AWSClient
	styles          *styles
}

func (lm *listModel) Update(msg tea.Msg) (listModel, tea.Cmd) {
	bubbleListModel, cmd := lm.bubbleListModel.Update(msg)
	return listModel{
		bubbleListModel: &bubbleListModel,
		awsClient:       lm.awsClient,
		styles:          lm.styles,
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

func (lm *listModel) SetItems(listItems []list.Item) tea.Cmd {
	return lm.bubbleListModel.SetItems(listItems)
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

func (lm *listModel) setItemsBatchCmd(listItems []list.Item) tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			indexMap := make(map[string]int)

			for i, cur := range lm.Items() {
				cur := cur.(item)
				indexMap[cur.url] = i
			}

			var cmds []tea.Cmd

			for _, cur := range listItems {
				cur := cur.(item)
				cmds = append(cmds, lm.SetItem(indexMap[cur.url], cur))
			}

			return tea.Batch(cmds...)
		},
	)
}

func (lm *listModel) loadPageAttributes(ctx context.Context, listItems ...list.Item) tea.Cmd {
	if len(listItems) == 0 || (len(listItems) == 1 && listItems[0] == nil) {
		return nil
	}

	var cmds []tea.Cmd

	cmds = append(cmds,
		lm.StartSpinner(),
		lm.awsClient.GetQueueAttributesCmd(ctx, listItems),
	)

	statusMessage := "refreshing page..."

	if len(listItems) == 1 {
		curItem := listItems[0].(item)
		statusMessage = "refreshing " + curItem.name + "..."
	}

	cmds = append(cmds, lm.NewStatusMessage(lm.styles.statusMessage.Render(statusMessage)))

	return tea.Batch(cmds...)
}
