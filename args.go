package menu

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	shellwords "github.com/junegunn/go-shellwords"
)

var ErrArgEmpty = errors.New("argument cannot be empty")

// Args holds the FZF arguments.
type Args []string

func (a Args) Validate() error {
	_, err := shellwords.Parse(strings.Join(a, " "))
	return err
}

// ArgsBuilder constructs command-line arguments for FZF.
type ArgsBuilder struct {
	list          Args
	ansi          string // Enable processing of ANSI color codes
	bind          string // Comma-separated list of custom key/event bindings
	border        string // Border around the window
	borderLabel   string // Label to print on the horizontal border line
	color         string // Color configuration
	footer        string // The given string will be printed as the sticky footer
	footerBorder  string // Border around the footer window
	header        string // The given string will be printed as the sticky header
	height        string // Set the height of the menu
	highlightLine string // Highlight the whole current line (bold)
	info          string // Determines the display style of the finder info.
	layout        string // Choose the layout (default: default)
	multi         string // Enable multi-select with tab/shift-tab
	noColor       string // Disable output color
	noScrollbar   string // Remove scrollbar
	pointer       string // Pointer to the current line
	preview       string // Execute the given command for the current line
	previewWindow string // Determines the layout of the preview window.
	previewBorder string // Draws a single separator line
	prompt        string // Input prompt
	read0         string // Read input delimited by ASCII NUL characters instead of newline characters
	sync          string // Synchronous search for multi-staged filtering
	tac           string // Reverse the order of the input
	cycle         string // Enable cyclic scroll
	headerFirst   string // Print header before the prompt line
	headerLabel   string // Label to print on the header border
	headerBorder  string // Draw border around the header section.
	withNth       string // Transform the presentation of each line using the field index expressions
}

func NewArgsBuilder() *ArgsBuilder {
	return &ArgsBuilder{
		ansi:          "--ansi",
		bind:          "--bind",
		border:        "--border",
		borderLabel:   "--border-label",
		color:         "--color",
		footer:        "--footer",
		footerBorder:  "--footer-border",
		header:        "--header",
		height:        "--height",
		highlightLine: "--highlight-line",
		info:          "--info",
		layout:        "--layout",
		multi:         "--multi",
		noColor:       "--no-color",
		noScrollbar:   "--no-scrollbar",
		pointer:       "--pointer",
		preview:       "--preview",
		previewWindow: "--preview-window",
		previewBorder: "--preview-border",
		prompt:        "--prompt",
		read0:         "--read0",
		sync:          "--sync",
		tac:           "--tac",
		cycle:         "--cycle",
		headerFirst:   "--header-first",
		headerLabel:   "--header-label",
		headerBorder:  "--header-border",
		withNth:       "--with-nth",
	}
}

// Validate checks that all string fields in ArgsBuilder are not empty.
func (a *ArgsBuilder) Validate() error {
	if err := a.list.Validate(); err != nil {
		return err
	}

	v := reflect.ValueOf(*a)
	t := reflect.TypeFor[ArgsBuilder]()

	for i := range v.NumField() {
		field := v.Field(i)
		fieldName := t.Field(i).Name

		if field.Kind() == reflect.String {
			if field.String() == "" {
				return fmt.Errorf("%w: %q", ErrArgEmpty, fieldName)
			}
		}
	}

	return nil
}

func (a *ArgsBuilder) Parse() error {
	_, err := shellwords.Parse(strings.Join(a.list, " "))
	return err
}

func (a *ArgsBuilder) Add(s ...string) *ArgsBuilder {
	a.list = append(a.list, s...)
	return a
}

func (a *ArgsBuilder) Build() Args                           { return a.list }
func (a *ArgsBuilder) String() string                        { return strings.Join(a.list, " ") }
func (a *ArgsBuilder) Custom(s ...string) *ArgsBuilder       { return a.Add(s...) }
func (a *ArgsBuilder) WithAnsi() *ArgsBuilder                { return a.Add(a.ansi) }
func (a *ArgsBuilder) WithHeight(s string) *ArgsBuilder      { return a.Add(a.height + "=" + s) }
func (a *ArgsBuilder) WithInfo(is InfoStyle) *ArgsBuilder    { return a.Add(a.info + "=" + string(is)) }
func (a *ArgsBuilder) WithLayout(l Layout) *ArgsBuilder      { return a.Add(a.layout + "=" + string(l)) }
func (a *ArgsBuilder) WithNoColor() *ArgsBuilder             { return a.Add(a.noColor) }
func (a *ArgsBuilder) WithNoScrollbar() *ArgsBuilder         { return a.Add(a.noScrollbar) }
func (a *ArgsBuilder) WithPointer(s string) *ArgsBuilder     { return a.Add(a.pointer + "=" + s) }
func (a *ArgsBuilder) WithMultiSelection() *ArgsBuilder      { return a.Add(a.multi) }
func (a *ArgsBuilder) WithHighlightLine() *ArgsBuilder       { return a.Add(a.highlightLine) }
func (a *ArgsBuilder) WithPrompt(s string) *ArgsBuilder      { return a.Add(a.prompt + "=" + s) }
func (a *ArgsBuilder) WithSync() *ArgsBuilder                { return a.Add(a.sync) }
func (a *ArgsBuilder) WithTac() *ArgsBuilder                 { return a.Add(a.tac) }
func (a *ArgsBuilder) WithCycle() *ArgsBuilder               { return a.Add(a.cycle) }
func (a *ArgsBuilder) WithBorder(b Border) *ArgsBuilder      { return a.Add(a.border + "=" + string(b)) }
func (a *ArgsBuilder) WithFooter(footer string) *ArgsBuilder { return a.Add(a.footer + "=" + footer) }
func (a *ArgsBuilder) WithPreview(s string) *ArgsBuilder     { return a.Add(a.preview + "=" + s) }
func (a *ArgsBuilder) WithHeader(s string) *ArgsBuilder      { return a.Add(a.header + "=" + s) }
func (a *ArgsBuilder) WithHeaderBorder(b Border) *ArgsBuilder {
	return a.Add(a.headerBorder + "=" + string(b))
}

func (a *ArgsBuilder) WithNth(idx ...string) *ArgsBuilder {
	return a.Add(a.withNth + "=" + strings.Join(idx, ","))
}

func (a *ArgsBuilder) WithFooterBorder(b Border) *ArgsBuilder {
	return a.Add(a.footerBorder + "=" + string(b))
}

func (a *ArgsBuilder) WithPreviewWindow(s string) *ArgsBuilder {
	return a.Add(a.previewWindow + "=" + s)
}

func (a *ArgsBuilder) WithPreviewBorder(b Border) *ArgsBuilder {
	return a.Add(a.previewBorder + "=" + string(b))
}

func (a *ArgsBuilder) WithBorderLabel(s string) *ArgsBuilder {
	return a.Add(a.border, a.borderLabel+"="+s)
}

func (a *ArgsBuilder) WithColor(target string, styles ...string) *ArgsBuilder {
	color := a.color + "=" + target
	if len(styles) > 0 {
		color += ":" + strings.Join(styles, ":")
	}

	return a.Add(color)
}
