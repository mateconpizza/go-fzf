// Package menu provides a flexible wrapper for the fzf interactive filter,
// enabling customizable selection menus.
package menu

import (
	"context"
	"errors"
	"log/slog"
)

var (
	// ErrFzf reports a generic fzf error (return code: 2).
	ErrFzf = errors.New("error: code 2")

	// ErrNoMatching reports that no items matched the query (return code: 1).
	ErrNoMatching = errors.New("no matching record: code 1")

	// ErrInvalidShellCommand reports an invalid shell command for a become
	// action (return code: 126).
	ErrInvalidShellCommand = errors.New("invalid shell command for become action: code 126")

	// ErrActionAborted reports that the user aborted the fzf session (return code:
	// 130).
	ErrActionAborted = errors.New("action aborted: code 130")

	// ErrPermissionDenied reports a permission error from a become action
	// (return code: 127).
	ErrPermissionDenied = errors.New("permission denied from become action: code 127")

	// ErrNoItems reports that no items were provided.
	ErrNoItems = errors.New("no items found")
)

type Option func(*Options)

type Options struct {
	// useDefaults
	useDefaults bool

	// header contains header lines displayed in the FZF interface.
	header            []string
	headerSeparator   string
	noHeader          bool
	showKeybindHeader bool

	// arg holds the command-line arguments passed to FZF.
	// These are built from various options and configurations.
	argsBuilder *ArgsBuilder

	// interruptFn handles FZF cancellation signals (Ctrl-C, ESC, etc.).
	interruptFn func(error)

	// runner executes the FZF command and handles I/O.
	// Can be customized for testing or different execution environments.
	runner MenuRunner

	// keymaps manages the keyboard shortcuts and their actions.
	// Provides methods to register and manage keybindings.
	keymaps *KeymapManager

	// headerKeymapFmt formats a key binding and its description for display in
	// the menu header.
	headerKeymapFmt KeymapFormatter

	// headerSeparatorFmt formats the separator used between
	// keymap entries in the menu header.
	headerSeparatorFmt SeparatorFormatter
}

type FmtFunc[T any] func(item T) string

// Items holds the data and transformation logic for menu items.
type Items[T any] struct {
	// Formatter converts items to display strings for FZF.
	// If nil, a default Formatter will be used that calls String() method.
	// The function should return ANSI-formatted strings for rich display.
	Formatter FmtFunc[T]
}

type Menu[T any] struct {
	Options
	Items[T]
}

// Select executes Fzf with the set elements and returns the selected item/s.
func (m *Menu[T]) Select(ctx context.Context, items []T) ([]T, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := m.Validate(); err != nil {
		return nil, err
	}

	if err := m.buildArgs(); err != nil {
		return nil, err
	}

	selected, err := selectFromItems(ctx, m, items)
	if err != nil {
		return nil, err
	}

	if len(selected) == 0 {
		return nil, ErrNoItems
	}

	return selected, nil
}

// callInterruptFn safely executes the interrupt callback if set.
func (m *Menu[T]) callInterruptFn(err error) {
	if m.interruptFn != nil {
		slog.Debug("calling interruptFn with err", "err", err)
		m.interruptFn(err)
	}

	slog.Debug("interruptFn is nil")
}

// SetInterruptFn sets the interrupt function for the menu.
func (m *Menu[T]) SetInterruptFn(fn func(error)) { m.interruptFn = fn }

// SetFormatter sets a function to format items for display in fzf.
func (m *Menu[T]) SetFormatter(preprocessor FmtFunc[T]) { m.Formatter = preprocessor }

// New returns a new Menu.
func New[T any](opts ...Option) *Menu[T] {
	o := Options{
		header:             make([]string, 0),
		runner:             &defaultRunner{},
		keymaps:            NewKeymapManager(),
		argsBuilder:        NewArgsBuilder(),
		headerKeymapFmt:    defaultKeymapFormatter,
		headerSeparatorFmt: defaultSeparatorFormatter,
	}

	for _, fn := range opts {
		fn(&o)
	}

	return &Menu[T]{
		Options: o,
	}
}

func (m Menu[T]) Validate() error {
	return m.argsBuilder.Validate()
}

// Select displays the given items in an interactive menu and returns the
// selected items.
func Select[T any](ctx context.Context, items []T, opts ...Option) ([]T, error) {
	m := New[T](opts...)

	items, err := m.Select(ctx, items)
	if err != nil {
		return nil, err
	}

	return items, err
}

// SelectOne displays the given items in an interactive menu and returns the
// single selected item.
//
//nolint:ireturn
func SelectOne[T any](ctx context.Context, items []T, opts ...Option) (T, error) {
	var zero T

	selected, err := Select(ctx, items, opts...)
	if err != nil {
		return zero, err
	}

	if len(selected) == 0 {
		return zero, ErrNoItems
	}

	return selected[0], nil
}

// SelectMany displays the given items in an interactive menu and returns the
// selected items.
func SelectMany[T any](ctx context.Context, items []T, opts ...Option) ([]T, error) {
	selected, err := Select(ctx, items, append(opts, WithMultiSelection())...)
	if err != nil {
		return nil, err
	}

	if len(selected) == 0 {
		return nil, ErrNoItems
	}

	return selected, nil
}
