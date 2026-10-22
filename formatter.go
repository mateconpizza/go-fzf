package menu

import (
	"context"
	"fmt"
)

type (
	// KeymapFormatter formats a single keymap entry.
	//
	// The arguments are the key binding, the separator between the
	// binding and description (typically ":"), and the description.
	KeymapFormatter func(bind, sep, desc string) string

	// SeparatorFormatter formats the separator used between
	// keymap entries in the menu header.
	SeparatorFormatter func(sep string) string
)

var (
	// DefaultKeymapFormatter returns an unformatted keymap entry.
	defaultKeymapFormatter = func(bind, sep, desc string) string {
		return bind + sep + desc
	}

	// DefaultSeparatorFormatter returns the separator unchanged.
	defaultSeparatorFormatter = func(sep string) string {
		return sep
	}
)

// defaultPreprocessor provides fallback formatting item for display in fzf.
func defaultPreprocessor[T any](item T) string { return fmt.Sprintf("%+v", item) }

// formatItemsPreprocessed formats each item in the slice using the preprocessor function
// and returns a channel of formatted strings.
func formatItemsPreprocessed(ctx context.Context, formattedItems []string) chan string {
	inputChan := make(chan string)

	go func() {
		defer close(inputChan)
		for _, formatted := range formattedItems {
			if ctx.Err() != nil {
				return
			}

			select {
			case inputChan <- formatted:
			case <-ctx.Done():
				return
			}
		}
	}()

	return inputChan
}

// processOutputPreprocessed formats items, maps them to their original values, and sends
// the filtered results to resultChan.
func processOutputPreprocessed[T any](itemMap map[string]T, outputChan <-chan string, resultChan chan<- []T) {
	var result []T

	for s := range outputChan {
		if item, exists := itemMap[s]; exists {
			result = append(result, item)
		}
	}

	resultChan <- result
}
