package menu

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"sync"
)

var ansiEscapeRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

const ExitSuccess = 0

// handleFzfErr returns an error based on the exit code of fzf.
//
//	0      Normal exit
//	1      No match
//	2      Error
//	126    Permission denied error from become action
//	127    Invalid shell command for become action.
//	130    Interrupted with CTRL-C or ESC.
func handleFzfErr(retcode int) error {
	switch retcode {
	case 1:
		return ErrNoMatching
	case 2:
		return ErrFzf
	case 126:
		return ErrInvalidShellCommand
	case 127:
		return ErrPermissionDenied
	case 130:
		return ErrActionAborted
	}

	return nil
}

// selectFromItems runs Fzf with the given items and returns the selected item/s.
func selectFromItems[T any](ctx context.Context, m *Menu[T], items []T) ([]T, error) {
	if len(items) == 0 {
		return nil, ErrNoItems
	}

	if m.Formatter == nil {
		slog.Debug("preprocessor is nil, defaulting to 'defaultPreprocessor'")
		m.Formatter = defaultPreprocessor
	}

	for i, arg := range m.argsBuilder.list {
		slog.Debug("menu args", strconv.Itoa(i), arg)
	}

	// Pre-process all items once for better performance
	formattedItems := make([]string, len(items))
	itemMap := make(map[string]T, len(items))
	for i, item := range items {
		ti := item
		formatted := m.Formatter(ti)
		formattedItems[i] = formatted
		itemMap[ansiCodeRemover(formatted)] = item
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// channels
	inputChan := formatItemsPreprocessed(ctx, formattedItems)
	outputChan := make(chan string)
	resultChan := make(chan []T, 1)
	go processOutputPreprocessed(itemMap, outputChan, resultChan)

	closeOutput := sync.OnceFunc(func() { close(outputChan) })
	defer closeOutput() // covers Parse errors and the non-success retcode

	// Build Fzf.Options
	options, err := m.runner.Parse(m.withDefaults(), m.argsBuilder.Build())
	if err != nil {
		return nil, fmt.Errorf("fzf: %w", err)
	}

	// Set up input and output channels
	options.Input = inputChan
	options.Output = outputChan

	// Run Fzf
	retcode, err := m.runner.Run(options)
	if retcode != ExitSuccess {
		// regardless of what kind of error, always call `callInterruptFn`
		err = handleFzfErr(retcode)
		m.callInterruptFn(err)

		return nil, err
	}

	closeOutput()
	result := <-resultChan

	if err != nil {
		return nil, err
	}

	return result, nil
}

// ansiCodeRemover removes ANSI codes from a given string.
func ansiCodeRemover(s string) string {
	return ansiEscapeRe.ReplaceAllString(s, "")
}
