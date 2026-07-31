# go-fzf

Provides a wrapper around fzf for interactive menus.

Extracted from [gomarks](https://github.com/mateconpizza/gm)

## example

```go
package main

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
		menu.WithDefaults(true), // uses `$FZF_DEFAULT_OPTS_FILE` and `$FZF_DEFAULT_OPTS`
		menu.WithPrompt("Select> "),
		menu.WithHeader("Inventory"),
		menu.WithBorder(menu.BorderRounded),
		menu.WithBorderLabel("= demo ="),
		menu.WithPreviewCmd("echo {1}"),
		menu.WithHeight("40%"),
		menu.WithFooter("~~ Give me apples ~~"),
	)

	selected, err := m.Select(items)
	if err != nil {
		return
	}

	for _, item := range selected {
		fmt.Println(item.Name)
	}
}
```
