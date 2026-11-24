// SPDX-FileCopyrightText: 2026 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type RuleType int

const (
	Rule_AlwaysCheck        RuleType = 1 // all Always checks are required to pass (AND)
	Rule_FallbackCheck      RuleType = 2 // any single Fallback checks pass won't make check fail (OR)
	Rule_ConditionialCheck  RuleType = 3 // additional check applies if condition matches.
)

type ConditionType int

const (
	Condition_ByEmail ConditionType = 1 // if email matches
	Condition_ByType  ConditionType = 2 // if scope matches
)

type CheckType int

const (
	Check_ConventionalCommitsBase CheckType = 0 // require type and description in subject line
	Check_TypeInRange             CheckType = 1 // check type in the given type lists
	Check_RequireScope            CheckType = 2 // require providing a scope
	Check_RequireBody             CheckType = 3 // require body (trailers also count as body)
	Check_RequireTrailers         CheckType = 4 // require given trailer(s)
	Check_AllowedSubjectPattern   CheckType = 5 // fallback allowed pattern if subject line doesn't meet convetional commit guideline
	Check_StrictRequireTrailers   CheckType = 6 // strict require given trailer(s), must contain trailer(s) key and value
)

type Rule struct {
	RuleType      RuleType
	CheckType     CheckType
	ConditionType ConditionType
	RuleData      string
	CheckData     string
}

type CommitCheckData struct {
	CommitSHA      string
	AuthorEmail    string
	SubjectRaw     string
	SubjectParsed  SubjectData
	Body           string
	TrailersRaw    string
	TrailersParsed TrailersData
}

func CreateCommitCheckData(sha string, email string, subject string, body string, trailers string) CommitCheckData {
	return CommitCheckData{
		CommitSHA:      sha,
		AuthorEmail:    email,
		SubjectRaw:     subject,
		SubjectParsed:  ParseSubjectLine(subject),
		Body:           body,
		TrailersRaw:    trailers,
		TrailersParsed: ParseTrailers(trailers),
	}
}

// For the returned value, only `Passed`, `Reason`, and `Suggestion` was used
func (commitCheckData *CommitCheckData) Check(rule Rule) CheckedCommitMessage {
	switch rule.CheckType {
	case Check_ConventionalCommitsBase:
		if commitCheckData.SubjectParsed.Valid != true {
			return CheckedCommitMessage{
				Passed:     false,
				Reason:     "Invalid subject line",
				Suggestion: "check if your subject match the given format: `type: description` or `type(scope): description`",
			}
		}
	case Check_RequireBody:
		if len(commitCheckData.Body) <= 0 {
			return CheckedCommitMessage{
				Passed:     false,
				Reason:     "Body is empty",
				Suggestion: "Body is required, please put some detail description about why this commit is needed",
			}
		}
	case Check_StrictRequireTrailers:
		if !commitCheckData.TrailersParsed.StrictContainsAny(strings.Split(rule.CheckData, "|")) {
			return CheckedCommitMessage{
				Passed:     false,
				Reason:     "Required trailer (any of: " + rule.CheckData + ") or its value not found",
				Suggestion: "You can validate your trailers with `git interpret-trailers`",
			}
		}
	case Check_RequireTrailers:
		if !commitCheckData.TrailersParsed.ContainsAny(strings.Split(rule.CheckData, "|")) {
			return CheckedCommitMessage{
				Passed:     false,
				Reason:     "Required trailer (any of: " + rule.CheckData + ") not found",
				Suggestion: "You can validate your trailers with `git interpret-trailers`",
			}
		}
	case Check_TypeInRange:
		allowedTypes := strings.Split(rule.CheckData, "|")
		for _, t := range allowedTypes {
			if commitCheckData.SubjectParsed.Type == t {
				return CheckedCommitMessage{Passed: true}
			}
		}
		return CheckedCommitMessage{
			Passed:     false,
			Reason:     "Type not in allowed list: " + rule.CheckData,
			Suggestion: "Please use one of the allowed types: " + rule.CheckData,
		}
	case Check_RequireScope:
		if commitCheckData.SubjectParsed.Scope == "" {
			return CheckedCommitMessage{
				Passed:     false,
				Reason:     "Scope is required but not provided",
				Suggestion: "Please add a scope to your commit subject",
			}
		}
		if rule.CheckData != "" {
			allowedScopes := strings.Split(rule.CheckData, "|")
			for _, s := range allowedScopes {
				if commitCheckData.SubjectParsed.Scope == s {
					return CheckedCommitMessage{Passed: true}
				}
			}
			return CheckedCommitMessage{
				Passed:     false,
				Reason:     "Scope not in allowed list: " + rule.CheckData,
				Suggestion: "Please use one of the allowed scopes: " + rule.CheckData,
			}
		}
	case Check_AllowedSubjectPattern:
		r := regexp.MustCompile(rule.CheckData)
		if !r.MatchString(commitCheckData.SubjectRaw) {
			return CheckedCommitMessage{
				Passed:     false,
				Reason:     "Subject doesn't match regex [" + rule.CheckData + "]",
				Suggestion: "Please check your subject line",
			}
		}
	}

	return CheckedCommitMessage{
		Passed: true,
	}
}

// RuleConfig is the JSON representation of a Rule, used for deserialization.
type RuleConfig struct {
	RuleType      string `json:"rule_type"`
	CheckType     string `json:"check_type"`
	ConditionType string `json:"condition_type"`
	RuleData      string `json:"rule_data"`
	CheckData     string `json:"check_data"`
}

func parseRuleType(s string) (RuleType, error) {
	switch s {
	case "always_check":
		return Rule_AlwaysCheck, nil
	case "conditional_check":
		return Rule_ConditionialCheck, nil
	case "fallback_check":
		return Rule_FallbackCheck, nil
	default:
		return 0, fmt.Errorf("invalid rule_type: %q", s)
	}
}

func parseCheckType(s string) (CheckType, error) {
	switch s {
	case "conventional_commits_base":
		return Check_ConventionalCommitsBase, nil
	case "type_in_range":
		return Check_TypeInRange, nil
	case "require_scope":
		return Check_RequireScope, nil
	case "require_body":
		return Check_RequireBody, nil
	case "require_trailers":
		return Check_RequireTrailers, nil
	case "allowed_subject_pattern":
		return Check_AllowedSubjectPattern, nil
	case "strict_require_trailers":
		return Check_StrictRequireTrailers, nil
	default:
		return 0, fmt.Errorf("invalid check_type: %q", s)
	}
}

func parseConditionType(s string) (ConditionType, error) {
	switch s {
	case "by_email":
		return Condition_ByEmail, nil
	case "by_type":
		return Condition_ByType, nil
	default:
		return 0, fmt.Errorf("invalid condition_type: %q", s)
	}
}

// ParseRules deserializes a JSON string into a slice of Rule.
// It validates enum mappings and pre-compiles regex patterns for
// allowed_subject_pattern rules to ensure user-provided regexes are valid.
func ParseRules(jsonStr string) ([]Rule, error) {
	var configs []RuleConfig
	if err := json.Unmarshal([]byte(jsonStr), &configs); err != nil {
		return nil, fmt.Errorf("failed to parse rules JSON: %w", err)
	}

	if configs == nil {
		return nil, fmt.Errorf("rules JSON is null, at least one rule is required")
	}

	rules := make([]Rule, 0, len(configs))
	for i, cfg := range configs {
		var rule Rule
		var err error

		rule.RuleType, err = parseRuleType(cfg.RuleType)
		if err != nil {
			return nil, fmt.Errorf("rule #%d: %w", i+1, err)
		}

		if rule.RuleType == Rule_ConditionialCheck {
			if cfg.ConditionType == "" {
				return nil, fmt.Errorf("rule #%d: conditional_check requires non-empty condition_type", i+1)
			}
			if cfg.RuleData == "" {
				return nil, fmt.Errorf("rule #%d: conditional_check requires non-empty rule_data", i+1)
			}
		}

		rule.CheckType, err = parseCheckType(cfg.CheckType)
		if err != nil {
			return nil, fmt.Errorf("rule #%d: %w", i+1, err)
		}

		if cfg.ConditionType != "" {
			rule.ConditionType, err = parseConditionType(cfg.ConditionType)
			if err != nil {
				return nil, fmt.Errorf("rule #%d: %w", i+1, err)
			}
		}

		rule.RuleData = cfg.RuleData
		rule.CheckData = cfg.CheckData

		if rule.CheckType == Check_AllowedSubjectPattern {
			if rule.CheckData == "" {
				return nil, fmt.Errorf("rule #%d: allowed_subject_pattern requires non-empty check_data", i+1)
			}
			if _, err := regexp.Compile(rule.CheckData); err != nil {
				return nil, fmt.Errorf("rule #%d: invalid regex in check_data: %w", i+1, err)
			}
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

// DefaultRules returns the built-in 6 predefined rules, matching the
// behavior of the original hardcoded rules on master.
func DefaultRules() []Rule {
	return []Rule{
		{
			RuleType:  Rule_AlwaysCheck,
			CheckType: Check_ConventionalCommitsBase,
		},
		{
			RuleType:  Rule_AlwaysCheck,
			CheckType: Check_RequireBody,
		},
		{
			RuleType:      Rule_ConditionialCheck,
			RuleData:      "uniontech.com",
			ConditionType: Condition_ByEmail,
			CheckType:     Check_RequireTrailers,
			CheckData:     "Log",
		},
		{
			RuleType:      Rule_ConditionialCheck,
			RuleData:      "uniontech.com",
			ConditionType: Condition_ByEmail,
			CheckType:     Check_RequireTrailers,
			CheckData:     "Issue|PMS",
		},
		{
			RuleType:  Rule_FallbackCheck,
			RuleData:  "Default merge commit caused by workflow sync",
			CheckType: Check_AllowedSubjectPattern,
			CheckData: `File Sync from [\w\-\.]+\/[\w\-\.]+ \(#[\d]+\)`,
		},
		{
			RuleType:  Rule_FallbackCheck,
			RuleData:  "Default GitHub auto merge commit message",
			CheckType: Check_AllowedSubjectPattern,
			CheckData: `Merge (?:pull request #[\d]+ from [\w\-\.\/]+|branch '[\w\-\.:\/]+' into [\w\-\/])`,
		},
	}
}
