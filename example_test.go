package menu_test

import (
	"fmt"

	menu "github.com/mateconpizza/go-fzf"
)

type Item struct {
	Name     string
	Quantity int
}

func (item Item) String() string {
	return fmt.Sprintf("%-12s (%d)", item.Name, item.Quantity)
}

var items = []Item{
	{Name: "Apples", Quantity: 12},
	{Name: "Bananas", Quantity: 8},
	{Name: "Oranges", Quantity: 5},
	{Name: "Pears", Quantity: 3},
	{Name: "Palandri", Quantity: 1},
}

func ExampleMenu() {
	m := menu.New[Item](
		// Uses `$FZF_DEFAULT_OPTS_FILE` and `$FZF_DEFAULT_OPTS`
		menu.WithDefaults(true),

		// Layout
		menu.WithHeight("40%"),
		menu.WithBorder(menu.BorderRounded),

		// Header
		menu.WithHeader("Inventory"),
		menu.WithHeaderKeymaps(),
		menu.WithHeaderBorder(menu.BorderBottom),

		// Interaction
		menu.WithMultiSelection(),
		menu.WithCycle(),
		menu.WithPrompt("Select> "),

		// Preview
		menu.WithPreviewCmd("echo {1}"),
		menu.WithPreviewBorder(menu.BorderRounded),

		// Appearance
		menu.WithColor("prompt", menu.ColorBrightCyan, menu.AttributeBold),
		menu.WithColor("footer", menu.ColorBrightBlue, menu.AttributeItalic),
	)

	selected, err := m.Select(items)
	if err != nil {
		return
	}

	for _, item := range selected {
		fmt.Println(item.Name)
		// Output: Apples
	}
}
