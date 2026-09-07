package menu

import (
	"errors"
	"reflect"
	"testing"
)

type fakeRunner struct {
	retcode int
	output  string
}

func (f *fakeRunner) Parse(defaults bool, settings Args) (*RunOptions, error) {
	return &RunOptions{}, nil
}

func (f *fakeRunner) Run(opts *RunOptions) (int, error) {
	opts.Output <- f.output
	return f.retcode, nil
}

func TestSelectReturnsSelectedItem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		items    []any
		output   string
		expected []any
		recode   int
		err      error
	}{
		{
			name:     "strings",
			items:    []any{"x", "y", "z"},
			output:   "y",
			expected: []any{"y"},
			recode:   0,
		},
		{
			name:     "ints",
			items:    []any{1, 2, 3},
			output:   "3",
			expected: []any{3},
			recode:   0,
		},
		{
			name:     "no items",
			items:    []any{},
			output:   "",
			expected: nil,
			recode:   1,
			err:      ErrNoItems,
		},
		{
			name:     "no match",
			items:    []any{1, 2},
			output:   "",
			expected: nil,
			recode:   1,
			err:      ErrNoMatching,
		},
		{
			name:     "action aborted",
			items:    []any{1, 2},
			output:   "",
			expected: nil,
			recode:   130,
			err:      ErrActionAborted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := make([]any, len(tt.items))
			copy(s, tt.items)

			r := &fakeRunner{
				output:  tt.output,
				retcode: tt.recode,
			}

			m := New[any](WithRunner(r))
			m.SetFormatter(defaultPreprocessor)

			result, err := m.Select(s)

			if tt.recode == 0 {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("expected result %v, got %v", tt.expected, result)
				}
				return
			}

			if result != nil {
				t.Errorf("expected nil result, got: %v", result)
			}
			if err == nil {
				t.Errorf("expected error %v, got nil", tt.err)
			} else if !errors.Is(err, tt.err) {
				t.Errorf("expected error %v, got %v", tt.err, err)
			}
		})
	}
}
