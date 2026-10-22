package menu

func (m *Menu[T]) withDefaults() bool { return m.useDefaults }

// WithDefaults uses $FZF_DEFAULT_OPTS_FILE and $FZF_DEFAULT_OPTS.
func WithDefaults(b bool) Option { return func(o *Options) { o.useDefaults = b } }

// WithAnsi enable processing of ANSI color codes.
func WithAnsi() Option { return func(o *Options) { o.argsBuilder.WithAnsi() } }

// WithHeight sets the height of the menu.
func WithHeight(s string) Option { return func(o *Options) { o.argsBuilder.WithHeight(s) } }

// WithInfo determines the display style of the finder info.
func WithInfo(is InfoStyle) Option { return func(o *Options) { o.argsBuilder.WithInfo(is) } }

func WithLayout(l Layout) Option { return func(o *Options) { o.argsBuilder.WithLayout(l) } }

// WithBorder sets border around the window.
func WithBorder(b Border) Option { return func(o *Options) { o.argsBuilder.WithBorder(b) } }

func WithFooterBorder(b Border) Option { return func(o *Options) { o.argsBuilder.WithFooterBorder(b) } }

// WithSync synchronous search for multi-staged filtering.
func WithSync() Option { return func(o *Options) { o.argsBuilder.WithSync() } }

// WithTac reverse the order of the input.
func WithTac() Option { return func(o *Options) { o.argsBuilder.WithTac() } }

// WithNoScrollbar do not display scrollbar.
func WithNoScrollbar() Option { return func(o *Options) { o.argsBuilder.WithNoScrollbar() } }

// WithCycle enable cyclic scroll.
func WithCycle() Option { return func(o *Options) { o.argsBuilder.WithCycle() } }

// WithInterruptFn sets a callback that executes on fzf interruption.
// Use for cleanup or custom error handling when user cancels selection.
func WithInterruptFn(fn func(error)) Option { return func(o *Options) { o.interruptFn = fn } }

// WithArgsCustom adds new args to Fzf.
func WithArgsCustom(args ...string) Option { return func(o *Options) { o.argsBuilder.Add(args...) } }

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
func WithKeybinds(keys ...*Keymap) Option { return func(o *Options) { o.keymaps.Register(keys...) } }

// WithMultiSelection adds a keybind to select multiple records.
func WithMultiSelection() Option { return func(o *Options) { o.argsBuilder.WithMultiSelection() } }

// WithRunner Add new OptionFn for test configuration.
func WithRunner(r MenuRunner) Option { return func(o *Options) { o.runner = r } }

// WithPreviewCmd adds preview with a custom command.
func WithPreviewCmd(cmd string) Option { return func(o *Options) { o.argsBuilder.WithPreview(cmd) } }

// WithPreviewWindow determines the layout of the preview window.
func WithPreviewWindow(args string) Option {
	return func(o *Options) { o.argsBuilder.WithPreviewWindow(args) }
}

// WithMultilineView adds multiline view and highlights the entire current line in Fzf.
func WithMultilineView() Option {
	return func(o *Options) { o.argsBuilder.Add(o.argsBuilder.highlightLine, o.argsBuilder.read0) }
}

// WithHeader adds a header to FZF, appending to existing headers.
func WithHeader(header string) Option {
	return func(o *Options) { o.header = append([]string{header}, o.header...) }
}
func WithoutHeader(b bool) Option { return func(o *Options) { o.noHeader = b } }

// WithHeaderFirst print header before the prompt line.
func WithHeaderFirst() Option {
	return func(o *Options) { o.argsBuilder.Add(o.argsBuilder.headerFirst) }
}

// WithHeaderBorder draw border around the header section.
func WithHeaderBorder(b Border) Option { return func(o *Options) { o.argsBuilder.WithHeaderBorder(b) } }

// WithHeaderLabel label to print on the header border.
func WithHeaderLabel(s string) Option {
	return func(o *Options) { o.argsBuilder.Add(o.argsBuilder.headerLabel + "=" + s) }
}

// WithPreviewBorder draws a single separator line.
func WithPreviewBorder(b Border) Option {
	return func(o *Options) { o.argsBuilder.WithPreviewBorder(b) }
}

// WithNth transform the presentation of each line using the field index expressions.
func WithNth(idx ...string) Option { return func(o *Options) { o.argsBuilder.WithNth(idx...) } }

func WithFooter(footer string) Option { return func(o *Options) { o.argsBuilder.WithFooter(footer) } }

// WithPrompt adds a prompt to Fzf.
func WithPrompt(s string) Option { return func(o *Options) { o.argsBuilder.WithPrompt(s) } }

func WithOutputColor(b bool) Option {
	return func(o *Options) {
		if !b {
			o.argsBuilder.WithNoColor()
		}
	}
}
func WithNoOutputColor() Option       { return func(o *Options) { o.argsBuilder.WithNoColor() } }
func WithBorderLabel(s string) Option { return func(o *Options) { o.argsBuilder.WithBorderLabel(s) } }
func WithPointer(s string) Option     { return func(o *Options) { o.argsBuilder.WithPointer(s) } }
func WithHeaderKeymaps() Option       { return func(o *Options) { o.showKeybindHeader = true } }
func WithHeaderKeymapFmt(fn KeymapFormatter) Option {
	return func(o *Options) { o.headerKeymapFmt = fn }
}
func WithHeaderSeparator(sep string) Option { return func(o *Options) { o.headerSeparator = sep } }
func WithHeaderSeparatorFmt(fn SeparatorFormatter) Option {
	return func(o *Options) { o.headerSeparatorFmt = fn }
}
