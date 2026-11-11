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

func TestTrailersDataContainsAny(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		keys   []string
		want   bool
	}{
		{
			name:  "single key found",
			input: "Log: some log\n",
			keys:  []string{"Log"},
			want:  true,
		},
		{
			name:  "single key not found",
			input: "Log: some log\n",
			keys:  []string{"Issue"},
			want:  false,
		},
		{
			name:  "multiple keys first match",
			input: "Issue: 123\n",
			keys:  []string{"Issue", "PMS"},
			want:  true,
		},
		{
			name:  "multiple keys last match",
			input: "PMS: 456\n",
			keys:  []string{"Issue", "PMS"},
			want:  true,
		},
		{
			name:  "multiple keys no match",
			input: "Log: some log\n",
			keys:  []string{"Issue", "PMS"},
			want:  false,
		},
		{
			name:  "empty keys slice",
			input: "Log: some log\n",
			keys:  []string{},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := ParseTrailers(tt.input)
			if got := td.ContainsAny(tt.keys); got != tt.want {
				t.Errorf("ContainsAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrailersDataStrictContainsAny(t *testing.T) {
	tests := []struct {
		name  string
		input string
		keys  []string
		want  bool
	}{
		{
			name:  "single key with value",
			input: "Log: some log\n",
			keys:  []string{"Log"},
			want:  true,
		},
		{
			name:  "single key without value",
			input: "Log:\n",
			keys:  []string{"Log"},
			want:  false,
		},
		{
			name:  "multiple keys first match with value",
			input: "Issue: 123\n",
			keys:  []string{"Issue", "PMS"},
			want:  true,
		},
		{
			name:  "multiple keys last match with value",
			input: "PMS: 456\n",
			keys:  []string{"Issue", "PMS"},
			want:  true,
		},
		{
			name:  "multiple keys no match",
			input: "Log: some log\n",
			keys:  []string{"Issue", "PMS"},
			want:  false,
		},
		{
			name:  "empty keys slice",
			input: "Log: some log\n",
			keys:  []string{},
			want:  false,
		},
		{
			name:  "key present but empty value, other key with value",
			input: "Issue:\nPMS: 456\n",
			keys:  []string{"Issue", "PMS"},
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := ParseTrailers(tt.input)
			if got := td.StrictContainsAny(tt.keys); got != tt.want {
				t.Errorf("StrictContainsAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckRequireTrailersWithPipe(t *testing.T) {
	tests := []struct {
		name      string
		commit    CommitCheckData
		rule      Rule
		wantPass  bool
		wantEmpty bool
	}{
		{
			name: "pipe CheckData with Issue present",
			commit: CreateCommitCheckData("sha", "dev@uniontech.com", "feat: test", "", "Issue: 123\n"),
			rule:  Rule{RuleType: Rule_ConditionialCheck, RuleData: "uniontech.com", CheckType: Check_RequireTrailers, CheckData: "Issue|PMS"},
			wantPass:  true,
			wantEmpty: true,
		},
		{
			name: "pipe CheckData with PMS present",
			commit: CreateCommitCheckData("sha", "dev@uniontech.com", "feat: test", "", "PMS: 456\n"),
			rule:  Rule{RuleType: Rule_ConditionialCheck, RuleData: "uniontech.com", CheckType: Check_RequireTrailers, CheckData: "Issue|PMS"},
			wantPass:  true,
			wantEmpty: true,
		},
		{
			name: "pipe CheckData with neither present",
			commit: CreateCommitCheckData("sha", "dev@uniontech.com", "feat: test", "", "Log: some log\n"),
			rule:  Rule{RuleType: Rule_ConditionialCheck, RuleData: "uniontech.com", CheckType: Check_RequireTrailers, CheckData: "Issue|PMS"},
			wantPass:  false,
			wantEmpty: false,
		},
		{
			name: "single key no pipe regression pass",
			commit: CreateCommitCheckData("sha", "dev@uniontech.com", "feat: test", "", "Log: some log\n"),
			rule:  Rule{RuleType: Rule_ConditionialCheck, RuleData: "uniontech.com", CheckType: Check_RequireTrailers, CheckData: "Log"},
			wantPass:  true,
			wantEmpty: true,
		},
		{
			name: "single key no pipe regression fail",
			commit: CreateCommitCheckData("sha", "dev@uniontech.com", "feat: test", "", "Issue: 123\n"),
			rule:  Rule{RuleType: Rule_ConditionialCheck, RuleData: "uniontech.com", CheckType: Check_RequireTrailers, CheckData: "Log"},
			wantPass:  false,
			wantEmpty: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.commit.Check(tt.rule)
			if result.Passed != tt.wantPass {
				t.Errorf("Check() Passed = %v, want %v, Reason: %s", result.Passed, tt.wantPass, result.Reason)
			}
			if !tt.wantEmpty && result.Reason == "" {
				t.Errorf("Check() expected non-empty Reason but got empty")
			}
		})
	}
}

func TestCheckStrictRequireTrailersWithPipe(t *testing.T) {
	tests := []struct {
		name      string
		commit    CommitCheckData
		rule      Rule
		wantPass  bool
		wantEmpty bool
	}{
		{
			name: "pipe CheckData with Issue present and has value",
			commit: CreateCommitCheckData("sha", "dev@uniontech.com", "feat: test", "", "Issue: 123\n"),
			rule:  Rule{RuleType: Rule_ConditionialCheck, RuleData: "uniontech.com", CheckType: Check_StrictRequireTrailers, CheckData: "Issue|PMS"},
			wantPass:  true,
			wantEmpty: true,
		},
		{
			name: "pipe CheckData with PMS present and has value",
			commit: CreateCommitCheckData("sha", "dev@uniontech.com", "feat: test", "", "PMS: 456\n"),
			rule:  Rule{RuleType: Rule_ConditionialCheck, RuleData: "uniontech.com", CheckType: Check_StrictRequireTrailers, CheckData: "Issue|PMS"},
			wantPass:  true,
			wantEmpty: true,
		},
		{
			name: "pipe CheckData with neither present",
			commit: CreateCommitCheckData("sha", "dev@uniontech.com", "feat: test", "", "Log: some log\n"),
			rule:  Rule{RuleType: Rule_ConditionialCheck, RuleData: "uniontech.com", CheckType: Check_StrictRequireTrailers, CheckData: "Issue|PMS"},
			wantPass:  false,
			wantEmpty: false,
		},
		{
			name: "pipe CheckData key present but empty value",
			commit: CreateCommitCheckData("sha", "dev@uniontech.com", "feat: test", "", "Issue:\n"),
			rule:  Rule{RuleType: Rule_ConditionialCheck, RuleData: "uniontech.com", CheckType: Check_StrictRequireTrailers, CheckData: "Issue|PMS"},
			wantPass:  false,
			wantEmpty: false,
		},
		{
			name: "single key no pipe regression pass",
			commit: CreateCommitCheckData("sha", "dev@uniontech.com", "feat: test", "", "Log: some log\n"),
			rule:  Rule{RuleType: Rule_ConditionialCheck, RuleData: "uniontech.com", CheckType: Check_StrictRequireTrailers, CheckData: "Log"},
			wantPass:  true,
			wantEmpty: true,
		},
		{
			name: "single key no pipe regression fail",
			commit: CreateCommitCheckData("sha", "dev@uniontech.com", "feat: test", "", "Issue: 123\n"),
			rule:  Rule{RuleType: Rule_ConditionialCheck, RuleData: "uniontech.com", CheckType: Check_StrictRequireTrailers, CheckData: "Log"},
			wantPass:  false,
			wantEmpty: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.commit.Check(tt.rule)
			if result.Passed != tt.wantPass {
				t.Errorf("Check() Passed = %v, want %v, Reason: %s", result.Passed, tt.wantPass, result.Reason)
			}
			if !tt.wantEmpty && result.Reason == "" {
				t.Errorf("Check() expected non-empty Reason but got empty")
			}
		})
	}
}
