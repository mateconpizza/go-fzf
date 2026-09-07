package menu

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

type fakeRunnerParse struct {
	retcode  int
	output   string
	parseErr error
}

func (f *fakeRunnerParse) Parse(defaults bool, settings Args) (*RunOptions, error) {
	if f.parseErr != nil {
		return nil, f.parseErr
	}
	return &RunOptions{}, nil
}

func (f *fakeRunnerParse) Run(opts *RunOptions) (int, error) {
	if f.retcode == ExitSuccess {
		opts.Output <- f.output
	}
	return f.retcode, nil
}

func TestSelectFromItems(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		items         []string
		output        string
		retcode       int
		parseErr      error
		nilFormatter  bool
		wantErr       error
		wantErrAny    bool
		wantInterrupt bool
		expected      []string
	}{
		{
			name:     "normal_returns_selected_item",
			items:    []string{"alpha", "beta", "gamma"},
			output:   "beta",
			retcode:  ExitSuccess,
			expected: []string{"beta"},
		},
		{
			name:    "empty_items_returns_error",
			items:   []string{},
			wantErr: ErrNoItems,
		},
		{
			name:         "nil_formatter_defaults_to_defaultPreprocessor",
			items:        []string{"solo"},
			output:       "solo",
			retcode:      ExitSuccess,
			nilFormatter: true,
			expected:     []string{"solo"},
		},
		{
			name:     "single_item_boundary",
			items:    []string{"only"},
			output:   "only",
			retcode:  ExitSuccess,
			expected: []string{"only"},
		},
		{
			name:          "no_match_calls_interrupt",
			items:         []string{"x", "y"},
			retcode:       1,
			wantInterrupt: true,
			wantErr:       ErrNoMatching,
		},
		{
			name:          "generic_fzf_error_calls_interrupt",
			items:         []string{"x", "y"},
			retcode:       2,
			wantInterrupt: true,
			wantErr:       ErrFzf,
		},
		{
			name:          "invalid_shell_command_calls_interrupt",
			items:         []string{"x", "y"},
			retcode:       126,
			wantInterrupt: true,
			wantErr:       ErrInvalidShellCommand,
		},
		{
			name:          "permission_denied_calls_interrupt",
			items:         []string{"x", "y"},
			retcode:       127,
			wantInterrupt: true,
			wantErr:       ErrPermissionDenied,
		},
		{
			name:          "action_aborted_calls_interrupt",
			items:         []string{"x", "y"},
			retcode:       130,
			wantInterrupt: true,
			wantErr:       ErrActionAborted,
		},
		{
			name:          "unrecognized_retcode_calls_interrupt_with_nil_err",
			items:         []string{"x", "y"},
			retcode:       99, // not in handleFzfErr's switch -> nil error
			wantInterrupt: true,
			wantErr:       nil, // note: retcode != ExitSuccess but resulting err is nil
			wantErrAny:    false,
		},
		{
			name:       "parse_error_is_wrapped",
			items:      []string{"a"},
			parseErr:   os.ErrInvalid,
			wantErrAny: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var interruptCalled bool

			m := &Menu[string]{
				Options: Options{
					argsBuilder: &ArgsBuilder{},
					runner: &fakeRunnerParse{
						retcode:  tt.retcode,
						output:   tt.output,
						parseErr: tt.parseErr,
					},
					interruptFn: func(error) { interruptCalled = true },
				},
			}
			if !tt.nilFormatter {
				m.Formatter = func(s string) string { return s }
			}

			result, err := selectFromItems(m, tt.items)

			switch {
			case tt.wantErrAny:
				if err == nil {
					t.Fatalf("selectFromItems() expected a non-nil error, got nil")
				}
			case tt.name == "unrecognized_retcode_calls_interrupt_with_nil_err":
				// handleFzfErr returns nil for an unmapped retcode, and
				// selectFromItems returns that nil straight through as the
				// function's error, even though retcode != ExitSuccess.
				if err != nil {
					t.Fatalf("selectFromItems() expected nil error for unmapped retcode, got %v", err)
				}
			case tt.wantErr != nil:
				if err == nil {
					t.Fatalf("selectFromItems() expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("selectFromItems() expected error %v, got %v", tt.wantErr, err)
				}
			default:
				if err != nil {
					t.Fatalf("selectFromItems() unexpected error: %v", err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Fatalf("selectFromItems() = %v, want %v", result, tt.expected)
				}
			}

			if tt.wantInterrupt && !interruptCalled {
				t.Errorf("selectFromItems() expected interruptFn to be called, but it wasn't")
			}
			if !tt.wantInterrupt && interruptCalled {
				t.Errorf("selectFromItems() did not expect interruptFn to be called, but it was")
			}
		})
	}
}
