package main

import "regexp"

type RuleType int

const (
	Rule_AlwaysCheck       RuleType = 1 // all Always checks are required to pass (AND)
	Rule_FallbackCheck     RuleType = 2 // any single Fallback checks pass won't make check fail (OR)
	Rule_ConditionialCheck RuleType = 3 // additional check applies if condition matches.
)

type ConditionType int

const (
	Condition_ByEmail ConditionType = 1 // if email matches
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
		if !commitCheckData.TrailersParsed.StrictContains(rule.CheckData) {
			return CheckedCommitMessage{
				Passed:     false,
				Reason:     "Required trailer [" + rule.CheckData + "] or [" + rule.CheckData + "]'s value not found",
				Suggestion: "You can validate your trailers with `git interpret-trailers`",
			}
		}
	case Check_RequireTrailers:
		if !commitCheckData.TrailersParsed.Contains(rule.CheckData) {
			return CheckedCommitMessage{
				Passed:     false,
				Reason:     "Required trailer [" + rule.CheckData + "] not found",
				Suggestion: "You can validate your trailers with `git interpret-trailers`",
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
