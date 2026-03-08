package main

import "fmt"

type item struct {
	name      string
	available string
	inFlight  string
	url       string
}

func (i item) Title() string { return i.name }
func (i item) Description() string {
	if i.available == "" || i.inFlight == "" {
		return "loading..."
	}
	return fmt.Sprintf("Available: %s | In Flight: %s", i.available, i.inFlight)
}

func (i item) FilterValue() string { return i.name }
