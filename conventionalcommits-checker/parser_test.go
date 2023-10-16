package main

import (
	"reflect"
	"testing"
)

func TestParseSubjectLine(t *testing.T) {
	type args struct {
		subject string
	}
	tests := []struct {
		name string
		args args
		want SubjectData
	}{
		{
			name: "regular commit msg",
			args: args{
				subject: "chore: regular description",
			},
			want: SubjectData{
				Valid:       true,
				Type:        "chore",
				Description: "regular description",
			},
		},
		{
			name: "commit msg contains upper case character",
			args: args{
				subject: "chore: Regular Description",
			},
			want: SubjectData{
				Valid:       true,
				Type:        "chore",
				Description: "Regular Description",
			},
		},
		{
			name: "commit msg contains symbol",
			args: args{
				subject: "chore: description `symbol`...",
			},
			want: SubjectData{
				Valid:       true,
				Type:        "chore",
				Description: "description `symbol`...",
			},
		},
		{
			name: "regular two line commit msg",
			args: args{
				subject: "chore: regular description\nnew line",
			},
			want: SubjectData{
				Valid:       true,
				Type:        "chore",
				Description: "regular description\nnew line",
			},
		},
		{
			name: "regular commit msg but missing space after colon",
			args: args{
				subject: "chore:regular description",
			},
			want: SubjectData{
				Valid: false,
			},
		},
		{
			name: "regular commit msg with scope",
			args: args{
				subject: "chore(doc): regular description",
			},
			want: SubjectData{
				Valid:       true,
				Type:        "chore",
				Scope:       "doc",
				Description: "regular description",
			},
		},
		{
			name: "regular commit msg with scope and skip test",
			args: args{
				subject: "[skip test]chore(doc): regular description",
			},
			want: SubjectData{
				Valid:       true,
				Type:        "chore",
				Scope:       "doc",
				Description: "regular description",
			},
		},
		{
			name: "invalid commit msg",
			args: args{
				subject: "chore regular description",
			},
			want: SubjectData{
				Valid: false,
			},
		},
		{
			name: "commit message with complex scope",
			args: args{
				subject: "chore(tools/some-tools): xxx",
			},
			want: SubjectData{
				Valid:       true,
				Type:        "chore",
				Scope:       "tools/some-tools",
				Description: "xxx",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseSubjectLine(tt.args.subject); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSubjectLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
