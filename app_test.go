package main

import (
	"testing"
)

func TestMergeList(t *testing.T) {
	tests := []struct {
		name string
		list []List
		want string
	}{
		{"sample", []List{List{Name: "name"}}, " name"},
		{
			"two",
			[]List{List{Name: "name1"}, List{Name: "name2"}},
			" name1 name2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeList(tt.list); got != tt.want {
				t.Errorf("MergeList() = %v, want %v", got, tt.want)
			}
		})
	}
}
