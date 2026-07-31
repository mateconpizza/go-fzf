package menu

import (
	"errors"
	"slices"
	"strings"
	"testing"
)

func TestMenu_buildHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func() *Menu[any]
		want       string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "normal_with_visible_keybinds",
			setup: func() *Menu[any] {
				return New[any](
					WithHeaderKeymaps(),
					WithHeaderSeparator(" - "),
					WithKeybinds([]*Keymap{
						{Bind: "a", Desc: "Add", Action: "add-action", Enabled: true},
						{Bind: "x", Desc: "Hidden", Action: "hidden-action", Enabled: true, Hidden: true}, // Skipped (hidden)
						{Bind: "d", Desc: "Disabled", Action: "disabled-action", Enabled: false},          // Skipped (disabled)
						{Bind: "s", Desc: "Save", Action: "save-action", Enabled: true},
					}...),
				)
			},
			want: "a:Add - s:Save",
		},
		{
			name: "no_header_flag_returns_early",
			setup: func() *Menu[any] {
				return New[any](
					WithoutHeader(true),
					WithHeaderKeymaps(),
					WithKeybinds([]*Keymap{
						{Bind: "a", Desc: "Add", Enabled: true},
					}...),
				)
			},
			want: "",
		},
		{
			name: "show_keybind_header_false_joins_existing_header",
			setup: func() *Menu[any] {
				return New[any](WithHeader("custom header"))
			},
			want: "custom header",
		},
		{
			name: "empty_description_defaults_to_question_mark",
			setup: func() *Menu[any] {
				return New[any](
					WithHeaderKeymaps(),
					WithKeybinds([]*Keymap{
						{Bind: "a", Desc: "", Action: "a-action", Enabled: true},
					}...),
				)
			},
			want: "a:?",
		},
		{
			name: "all_keybinds_filtered_returns_empty",
			setup: func() *Menu[any] {
				return New[any](
					WithHeaderKeymaps(),
					WithKeybinds([]*Keymap{
						{Bind: "h", Desc: "Hidden", Enabled: true, Hidden: true},
						{Bind: "d", Desc: "Disabled", Enabled: false},
					}...),
				)
			},
			want: "",
		},
		{
			name: "empty_keymaps_list_returns_empty",
			setup: func() *Menu[any] {
				return New[any](WithHeaderKeymaps())
				// No keymaps registered
			},
			want: "",
		},
		{
			name: "shellwords_parse_error_on_malformed_header",
			setup: func() *Menu[any] {
				return New[any](
					WithHeaderKeymaps(),
					WithKeybinds([]*Keymap{
						{Bind: "e", Desc: "unmatched ' quote", Enabled: true},
					}...),
				)
			},
			want:       "",
			wantErr:    true,
			wantErrMsg: "invalid command line string",
		},
		{
			name: "uses_custom_header_when_provided",
			setup: func() *Menu[any] {
				return New[any](
					WithHeader("custom header"),
					WithKeybinds([]*Keymap{
						{Bind: "a", Desc: "Add", Enabled: true},
						{Bind: "x", Desc: "Hidden", Enabled: true, Hidden: true},
						{Bind: "d", Desc: "Delete", Enabled: true},
					}...),
				)
			},
			want: "custom header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := tt.setup()
			got, err := m.buildHeader()

			if tt.wantErr {
				if err == nil {
					t.Fatalf("buildHeader() expected error containing %q, got nil", tt.wantErrMsg)
				}
				if tt.wantErrMsg != "" && !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Fatalf("buildHeader() expected error containing %q, got %v", tt.wantErrMsg, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("buildHeader() unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("buildHeader() = %q; want %q", got, tt.want)
			}
		})
	}
}

func TestMenu_buildKeybindString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		keybinds []*Keymap
		want     []string
	}{
		{
			name:     "empty_no_keybinds",
			keybinds: nil,
			want:     nil,
		},
		{
			name: "single_enabled",
			keybinds: []*Keymap{
				NewKeymap().WithBind(KeyCtrlE).WithExecute("echo {}"),
			},
			want: []string{"ctrl-e:execute(echo {})"},
		},
		{
			name: "multiple_enabled",
			keybinds: []*Keymap{
				NewKeymap().WithBind(KeyCtrlE).WithExecute("echo {}"),
				NewKeymap().WithBind(KeyCtrlA).WithExecute("cat {}"),
			},
			want: []string{
				"ctrl-e:execute(echo {})",
				"ctrl-a:execute(cat {})",
			},
		},
		{
			name: "disabled_keybind_excluded",
			keybinds: []*Keymap{
				NewKeymap().WithBind(KeyCtrlE).WithExecute("echo {}"),
				NewKeymap().WithBind(KeyCtrlA).WithExecute("cat {}").WithEnabled(false),
			},
			want: []string{"ctrl-e:execute(echo {})"},
		},
		{
			name: "missing_action_excluded",
			keybinds: []*Keymap{
				NewKeymap().WithBind(KeyCtrlE), // no WithExecute -> Action == ""
				NewKeymap().WithBind(KeyCtrlA).WithExecute("cat {}"),
			},
			want: []string{"ctrl-a:execute(cat {})"},
		},
		{
			name: "all_disabled_yields_empty",
			keybinds: []*Keymap{
				NewKeymap().WithBind(KeyCtrlE).WithExecute("echo {}").WithEnabled(false),
				NewKeymap().WithBind(KeyCtrlA).WithExecute("cat {}").WithEnabled(false),
			},
			want: nil,
		},
		{
			name: "all_missing_action_yields_empty",
			keybinds: []*Keymap{
				NewKeymap().WithBind(KeyCtrlE),
				NewKeymap().WithBind(KeyCtrlA),
			},
			want: nil,
		},
		{
			name: "action_with_special_chars",
			keybinds: []*Keymap{
				NewKeymap().WithBind(KeyCtrlR).WithExecute(`reload(find . -name "*.go")`),
			},
			want: []string{`ctrl-r:execute(reload(find . -name "*.go"))`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := New[any](WithKeybinds(tt.keybinds...))
			got := m.buildKeybindString()

			if tt.want == nil {
				if got != "" {
					t.Fatalf("buildKeybindString() = %q; want empty", got)
				}
				return
			}

			gotParts := strings.Split(got, ",")
			slices.Sort(gotParts)
			wantParts := slices.Clone(tt.want)
			slices.Sort(wantParts)

			if !slices.Equal(gotParts, wantParts) {
				t.Fatalf("buildKeybindString() = %q; want (any order) %q", got, tt.want)
			}
		})
	}
}

func TestMenu_buildKeybindArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		keybinds []*Keymap
		wantArgs bool // whether argsBuilder should have a --bind= entry
		wantErr  error
	}{
		{
			name:     "no_keybinds_noop",
			keybinds: nil,
			wantArgs: false,
			wantErr:  nil,
		},
		{
			name: "single_keybind_appends_arg",
			keybinds: []*Keymap{
				NewKeymap().WithBind(KeyCtrlE).WithExecute("echo {}"),
			},
			wantArgs: true,
			wantErr:  nil,
		},
		{
			name: "all_disabled_noop",
			keybinds: []*Keymap{
				NewKeymap().WithBind(KeyCtrlE).WithExecute("echo {}").WithEnabled(false),
			},
			wantArgs: false,
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := New[any](WithKeybinds(tt.keybinds...))
			err := m.buildKeybindArgs()

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("buildKeybindArgs() error = %v; want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("buildKeybindArgs() unexpected error: %v", err)
			}

			got := m.argsBuilder.String()
			hasBind := strings.Contains(got, "--bind=")
			if hasBind != tt.wantArgs {
				t.Fatalf("argsBuilder.String() = %q; wantArgs=%v", got, tt.wantArgs)
			}
		})
	}
}
