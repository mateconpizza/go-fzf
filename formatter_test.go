package menu

import (
	"context"
	"slices"
	"testing"
	"time"
)

func TestFormatItemsPreprocessed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		items          []string
		cancelBefore   bool
		cancelAfter    int
		want           []string
		wantNoMoreData bool
	}{
		{
			name:  "multiple_items",
			items: []string{"one", "two", "three"},
			want:  []string{"one", "two", "three"},
		},
		{
			name:  "single_item",
			items: []string{"one"},
			want:  []string{"one"},
		},
		{
			name:  "empty_slice",
			items: []string{},
			want:  []string{},
		},
		{
			name:  "nil_slice",
			items: nil,
			want:  nil,
		},
		{
			name:  "empty_strings",
			items: []string{"", "", ""},
			want:  []string{"", "", ""},
		},
		{
			name:  "duplicate_items",
			items: []string{"one", "one", "two", "one"},
			want:  []string{"one", "one", "two", "one"},
		},
		{
			name:           "context_canceled_before_start",
			items:          []string{"one", "two"},
			cancelBefore:   true,
			wantNoMoreData: true,
		},
		{
			name:           "context_canceled_after_first_item",
			items:          []string{"one", "two", "three"},
			cancelAfter:    1,
			want:           []string{"one"},
			wantNoMoreData: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			if tt.cancelBefore {
				cancel()
			}

			ch := formatItemsPreprocessed(ctx, tt.items)

			if tt.wantNoMoreData {
				got := make([]string, 0, len(tt.want))

				for range len(tt.want) {
					select {
					case item, ok := <-ch:
						if !ok {
							t.Fatalf("channel closed early: got %v; want %v", got, tt.want)
						}
						got = append(got, item)
					case <-time.After(time.Second):
						t.Fatal("timed out waiting for item")
					}
				}

				if tt.cancelAfter > 0 {
					cancel()
				}

				select {
				case item, ok := <-ch:
					if ok {
						t.Fatalf("received unexpected item after cancellation: %q", item)
					}
				case <-time.After(time.Second):
					t.Fatal("channel was not closed after cancellation")
				}

				if !slices.Equal(got, tt.want) {
					t.Fatalf("got %v; want %v", got, tt.want)
				}

				return
			}

			var got []string

			for {
				select {
				case item, ok := <-ch:
					if !ok {
						if !slices.Equal(got, tt.want) {
							t.Fatalf("got %v; want %v", got, tt.want)
						}
						return
					}
					got = append(got, item)

				case <-time.After(time.Second):
					t.Fatal("timed out waiting for channel to close")
				}
			}
		})
	}
}
