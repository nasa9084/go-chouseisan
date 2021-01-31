package chouseisan

import (
	"fmt"
	"testing"
	"time"
)

func TestZeller(t *testing.T) {
	tests := []struct {
		year  int
		month time.Month
		day   int
		want  int
	}{
		{
			year:  1992,
			month: time.July,
			day:   19,
			want:  1, // Sunday
		},
		{
			year:  2021,
			month: time.February,
			day:   1,
			want:  2, // Monday
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d/%02d/%02d", tt.year, tt.month, tt.day), func(t *testing.T) {
			got := zellersCongruence(tt.year, tt.month, tt.day)
			if got != tt.want {
				t.Errorf("%d != %d", got, tt.want)
				return
			}
		})
	}
}
