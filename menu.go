// Package menu provides a flexible wrapper for the fzf interactive filter,
// enabling customizable selection menus.
package menu

import (
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
func (m *Menu[T]) Select(items []T) ([]T, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}

	if err := m.buildArgs(); err != nil {
		return nil, err
	}

	selected, err := selectFromItems(m, items)
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
func (m *Menu[T]) SetInterruptFn(fn func(error)) {
	m.interruptFn = fn
}

// SetFormatter sets a function to format items for display in fzf.
func (m *Menu[T]) SetFormatter(preprocessor FmtFunc[T]) {
	m.Formatter = preprocessor
}

func (m *Menu[T]) withDefaults() bool {
	return m.useDefaults
}

// WithDefaults uses $FZF_DEFAULT_OPTS_FILE and $FZF_DEFAULT_OPTS.
func WithDefaults(b bool) Option {
	return func(o *Options) {
		o.useDefaults = b
	}
}

// WithAnsi enable processing of ANSI color codes.
func WithAnsi() Option {
	return func(o *Options) {
		o.argsBuilder.WithAnsi()
	}
}

// WithHeight sets the height of the menu.
func WithHeight(s string) Option {
	return func(o *Options) {
		o.argsBuilder.WithHeight(s)
	}
}

func WithInfo(is InfoStyle) Option {
	return func(o *Options) {
		o.argsBuilder.WithInfo(is)
	}
}

func WithLayout(l Layout) Option {
	return func(o *Options) {
		o.argsBuilder.WithLayout(l)
	}
}

// WithBorder sets border around the window.
func WithBorder(b Border) Option {
	return func(o *Options) {
		o.argsBuilder.WithBorder(b)
	}
}

func WithFooterBorder(b Border) Option {
	return func(o *Options) {
		o.argsBuilder.WithFooterBorder(b)
	}
}

// WithSync synchronous search for multi-staged filtering.
func WithSync() Option {
	return func(o *Options) {
		o.argsBuilder.WithSync()
	}
}

// WithTac reverse the order of the input.
func WithTac() Option {
	return func(o *Options) {
		o.argsBuilder.WithTac()
	}
}

// WithNoScrollbar do not display scrollbar.
func WithNoScrollbar() Option {
	return func(o *Options) {
		o.argsBuilder.WithNoScrollbar()
	}
}

// WithCycle enable cyclic scroll.
func WithCycle() Option {
	return func(o *Options) {
		o.argsBuilder.WithCycle()
	}
}

// WithInterruptFn sets a callback that executes on fzf interruption.
// Use for cleanup or custom error handling when user cancels selection.
func WithInterruptFn(fn func(error)) Option {
	return func(o *Options) {
		o.interruptFn = fn
	}
}

// WithArgsCustom adds new args to Fzf.
func WithArgsCustom(args ...string) Option {
	return func(o *Options) {
		o.argsBuilder.Add(args...)
	}
}

func WithArgs(fn func(b *ArgsBuilder) *ArgsBuilder) Option {
	return func(o *Options) {
		o.argsBuilder.Add(fn(NewArgsBuilder()).Build()...)
	}
}

// WithColor configures the color and text attributes for an fzf UI element.
//
// The target specifies the UI element to style (for example, "prompt",
// "header", or "border"). The values specify one or more ANSI colors or text
// attributes supported by fzf.
func WithColor(target string, values ...ColorValue) Option {
	return func(o *Options) {
		args := make([]string, len(values))
		for i, v := range values {
			args[i] = string(v)
		}
		o.argsBuilder.WithColor(target, args...)
	}
}

// WithKeybinds adds a keybind to Fzf.
func WithKeybinds(keys ...*Keymap) Option {
	return func(o *Options) {
		o.keymaps.Register(keys...)
	}
}

// WithMultiSelection adds a keybind to select multiple records.
func WithMultiSelection() Option {
	return func(o *Options) {
		o.argsBuilder.WithMultiSelection()
	}
}

// WithRunner Add new OptionFn for test configuration.
func WithRunner(r MenuRunner) Option {
	return func(o *Options) {
		o.runner = r
	}
}

// func WithPreview(b bool) Option {
// 	return func(o *Options) {
// 		o.preview = b
// 	}
// }

// WithPreviewCmd adds preview with a custom command.
func WithPreviewCmd(cmd string) Option {
	return func(o *Options) {
		o.argsBuilder.WithPreview(cmd)
	}
}

// WithPreviewWindow determines the layout of the preview window.
func WithPreviewWindow(args string) Option {
	return func(o *Options) {
		o.argsBuilder.WithPreviewWindow(args)
	}
}

// WithMultilineView adds multiline view and highlights the entire current line
// in Fzf.
func WithMultilineView() Option {
	return func(o *Options) {
		o.argsBuilder.Add(o.argsBuilder.highlightLine, o.argsBuilder.read0)
	}
}

// WithHeader adds a header to FZF, appending to existing headers.
func WithHeader(header string) Option {
	return func(o *Options) {
		o.header = append([]string{header}, o.header...)
	}
}

func WithoutHeader(b bool) Option {
	return func(o *Options) {
		o.noHeader = b
	}
}

// WithHeaderFirst print header before the prompt line.
func WithHeaderFirst() Option {
	return func(o *Options) {
		o.argsBuilder.Add(o.argsBuilder.headerFirst)
	}
}

// WithHeaderBorder draw border around the header section.
func WithHeaderBorder(b Border) Option {
	return func(o *Options) {
		o.argsBuilder.WithHeaderBorder(b)
	}
}

// WithHeaderLabel label to print on the header border.
func WithHeaderLabel(s string) Option {
	return func(o *Options) {
		o.argsBuilder.Add(o.argsBuilder.headerLabel + "=" + s)
	}
}

// WithPreviewBorder draws a single separator line.
func WithPreviewBorder(b Border) Option {
	return func(o *Options) {
		o.argsBuilder.WithPreviewBorder(b)
	}
}

// WithNth transform the presentation of each line using the field index
// expressions.
func WithNth(idx ...string) Option {
	return func(o *Options) {
		o.argsBuilder.WithNth(idx...)
	}
}

func WithFooter(footer string) Option {
	return func(o *Options) {
		o.argsBuilder.WithFooter(footer)
	}
}

// WithPrompt adds a prompt to Fzf.
func WithPrompt(s string) Option {
	return func(o *Options) {
		o.argsBuilder.WithPrompt(s)
	}
}

func WithOutputColor(b bool) Option {
	return func(o *Options) {
		if !b {
			o.argsBuilder.WithNoColor()
		}
	}
}

func WithNoOutputColor() Option {
	return func(o *Options) {
		o.argsBuilder.WithNoColor()
	}
}

func WithBorderLabel(s string) Option {
	return func(o *Options) {
		o.argsBuilder.WithBorderLabel(s)
	}
}

func WithPointer(s string) Option {
	return func(o *Options) {
		o.argsBuilder.WithPointer(s)
	}
}

func WithHeaderKeymaps() Option {
	return func(o *Options) {
		o.showKeybindHeader = true
	}
}

func WithHeaderKeymapFmt(fn KeymapFormatter) Option {
	return func(o *Options) {
		o.headerKeymapFmt = fn
	}
}

func WithHeaderSeparator(sep string) Option {
	return func(o *Options) {
		o.headerSeparator = sep
	}
}

func WithHeaderSeparatorFmt(fn SeparatorFormatter) Option {
	return func(o *Options) {
		o.headerSeparatorFmt = fn
	}
}

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
func Select[T any](items []T, opts ...Option) ([]T, error) {
	m := New[T](opts...)

	items, err := m.Select(items)
	if err != nil {
		return nil, err
	}

	return items, err
}
