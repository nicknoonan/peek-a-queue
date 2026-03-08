package main

import (
	"context"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type setItemInput struct {
	index int
	setItem item
}

func setItemsCmd(listModel *list.Model, setItemsList []setItemInput) tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			var cmds []tea.Cmd
	
			for _, setItem := range setItemsList {
				cmds = append(cmds, listModel.SetItem(setItem.index, setItem.setItem))
			}
	
			return tea.Batch(cmds...)
		},
	)
}

func loadPageAttributes(ctx context.Context, model *model, listModel *list.Model, styles *styles, awsClient *AWSClient, listItems ...list.Item) tea.Cmd {
	return tea.Batch(
		func() tea.Msg{
			model.isLoading.Set(true)
			return nil
		},
		listModel.StartSpinner(),
		listModel.NewStatusMessage(styles.statusMessage.Render("refreshing...")),
		awsClient.GetQueueAttributesCmd(ctx, listModel.Items(), listItems),
		func() tea.Msg{
			model.isLoading.Set(false)
			return nil
		},
	)
}

func newItemDelegate(ctx context.Context, model *model, awsClient *AWSClient, keys *delegateKeyMap, styles *styles) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, listModel *list.Model) tea.Cmd {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch {
			case key.Matches(msg, keys.choose):
				return loadPageAttributes(ctx, model, listModel, styles, awsClient, listModel.SelectedItem())
			case key.Matches(msg, keys.refresh):
				return loadPageAttributes(ctx, model, listModel, styles, awsClient, listModel.VisibleItems()...)
			}
		}

		if !model.isLoading.Get() {
			for _, curItem := range listModel.VisibleItems() {
				curItem := curItem.(item)
				if curItem.lengthString == "" {
					return loadPageAttributes(ctx, model, listModel, styles, awsClient, listModel.VisibleItems()...)
				}
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose, keys.refresh}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type delegateKeyMap struct {
	choose key.Binding
	refresh key.Binding
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.refresh,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.refresh,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "refreshes selected item attributes"),
		),
		refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refreshes current page attributes"),
		),
	}
}
