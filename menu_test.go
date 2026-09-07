package menu

import (
	"errors"
	"reflect"
	"testing"
)

func TestSelect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		items   []any
		opts    []Option
		want    []any
		wantErr error
	}{
		{
			name:    "normal_string_selection",
			items:   []any{"apple", "banana", "cherry"},
			opts:    []Option{WithRunner(&fakeRunner{retcode: 0, output: "banana"})},
			want:    []any{"banana"},
			wantErr: nil,
		},
		{
			name:    "normal_int_selection",
			items:   []any{1, 2, 3, 4, 5},
			opts:    []Option{WithRunner(&fakeRunner{retcode: 0, output: "4"})},
			want:    []any{4},
			wantErr: nil,
		},
		{
			name:    "empty_slice_zero_value",
			items:   []any{},
			opts:    nil,
			want:    nil,
			wantErr: ErrNoItems,
		},
		{
			name:    "nil_slice",
			items:   nil,
			opts:    nil,
			want:    nil,
			wantErr: ErrNoItems,
		},
		{
			name:    "single_item_boundary",
			items:   []any{"only_choice"},
			opts:    []Option{WithRunner(&fakeRunner{retcode: 0, output: "only_choice"})},
			want:    []any{"only_choice"},
			wantErr: nil,
		},
		{
			name:    "error_action_aborted",
			items:   []any{"yes", "no"},
			opts:    []Option{WithRunner(&fakeRunner{retcode: 130, output: ""})},
			want:    nil,
			wantErr: ErrActionAborted,
		},
		{
			name:    "error_no_match",
			items:   []any{"foo", "bar"},
			opts:    []Option{WithRunner(&fakeRunner{retcode: 1, output: ""})},
			want:    nil,
			wantErr: ErrNoMatching,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Select(tt.items, tt.opts...)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("Select() expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Select() expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Select() unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Select() = %v; want %v", got, tt.want)
			}
		})
	}
}
