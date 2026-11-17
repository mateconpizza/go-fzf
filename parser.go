package menu

import (
	"context"
	"fmt"
	"regexp"
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
	log := m.log().With("fn", "select-from-items")

	if len(items) == 0 {
		log.Debug("no items to select")
		return nil, ErrNoItems
	}

	log.Debug("preparing menu items", "count", len(items))

	if m.Formatter == nil {
		log.Debug("formatter is nil, using default preprocessor")
		m.Formatter = defaultPreprocessor
	}

	for i, arg := range m.argsBuilder.list {
		log.Debug("menu argument", "index", i, "value", arg)
	}

	// Pre-process all items once for better performance.
	formattedItems := make([]string, len(items))
	itemMap := make(map[string]T, len(items))

	for i, item := range items {
		formatted := m.Formatter(item)
		formattedItems[i] = formatted
		itemMap[ansiCodeRemover(formatted)] = item
	}

	log.Debug("menu items formatted", "count", len(formattedItems))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Channels.
	inputChan := formatItemsPreprocessed(ctx, formattedItems)
	outputChan := make(chan string)
	resultChan := make(chan []T, 1)

	go processOutputPreprocessed(itemMap, outputChan, resultChan)

	closeOutput := sync.OnceFunc(func() {
		close(outputChan)
	})
	defer closeOutput()

	// Build Fzf.Options.
	options, err := m.runner.Parse(m.withDefaults(), m.argsBuilder.Build())
	if err != nil {
		log.Error("failed to parse fzf options", "error", err)
		return nil, fmt.Errorf("fzf: %w", err)
	}

	options.Input = inputChan
	options.Output = outputChan

	log.Debug("running fzf", "items", len(items))

	retcode, err := m.runner.Run(options)
	if retcode != ExitSuccess {
		if err != nil {
			log.Debug("fzf runner returned error",
				"exit_code", retcode,
				"error", err,
			)
		}

		err = handleFzfErr(retcode)

		log.Debug("fzf exited without successful selection",
			"exit_code", retcode,
			"error", err,
		)

		m.callInterruptFn(err)

		return nil, err
	}

	closeOutput()

	result := <-resultChan

	if err != nil {
		log.Error("fzf execution failed", "error", err)
		return nil, err
	}

	log.Debug("fzf selection completed", "selected", len(result))

	return result, nil
}

// ansiCodeRemover removes ANSI codes from a given string.
func ansiCodeRemover(s string) string {
	return ansiEscapeRe.ReplaceAllString(s, "")
}
