package util

import (
	"github.com/tj/assert"
	"testing"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		items string
		want  []string
	}{
		{"test1,test2,test3", []string{"test1", "test2", "test3"}},
		{"Fulton,West,Rod,Farley", []string{"Fulton", "West", "Rod", "Farley"}},
	}

	for _, tt := range tests {
		t.Run("testing split", func(t *testing.T) {
			assert.Equal(t, tt.want, Split(tt.items))
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		items []string
		want  string
	}{
		{[]string{"test1", "test2", "test3"}, "test1,test2,test3"},
		{[]string{"Fulton", "West", "Rod", "Farley"}, "Fulton,West,Rod,Farley"},
	}

	for _, tt := range tests {
		t.Run("testing join", func(t *testing.T) {
			assert.Equal(t, tt.want, Join(tt.items))
		})
	}
}
