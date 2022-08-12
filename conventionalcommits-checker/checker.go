package main

import (
	"container/list"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Checker struct {
	Options Option
	Commits []string
	Rules   []Rule
	Result  Result
}

func CreateChecker() Checker {
	checker := Checker{
		Options: Option{
			ExcludeBase: false,
		},
		Rules: []Rule{
			{
				RuleType:  Rule_AlwaysCheck,
				CheckType: Check_ConventionalCommitsBase,
			},
			{
				RuleType:  Rule_AlwaysCheck,
				CheckType: Check_RequireBody,
			},
			{
				RuleType:  Rule_ConditionialCheck,
				RuleData:  "uniontech.com",
				CheckType: Check_RequireTrailers,
				CheckData: "Log",
			},
			{
				RuleType:  Rule_FallbackCheck,
				RuleData:  "Default merge commit caused by workflow sync",
				CheckType: Check_AllowedSubjectPattern,
				// File Sync from originname/branchname (#61)
				CheckData: `File Sync from [\w\-\.]+\/[\w\-\.]+ \(#[\d]+\)`,
			},
			{
				RuleType:  Rule_FallbackCheck,
				RuleData:  "Default GitHub auto merge commit message",
				CheckType: Check_AllowedSubjectPattern,
				// Merge branch 'originname:branchname' into branchname
				// Merge pull request #61 from originname/branchname
				CheckData: `Merge (?:pull request #[\d]+ from [\w\-\.\/]+|branch '[\w\-\.:\/]+' into [\w\-\/])`,
			},
		},
		Result: Result{},
	}

	return checker
}

func (checker *Checker) LogNoticeln(v ...any) {
	if checker.Options.GitHubLogger {
		logv := append([]any{"::notice ::"}, v...)
		fmt.Println(logv...)
	} else {
		log.Println(v...)
	}
}

func (checker *Checker) LogErrorln(title string, message string) {
	if checker.Options.GitHubLogger {
		fmt.Println("::error title=" + title + "::" + message)
	} else {
		log.Println(title, message)
	}
}

func (checker *Checker) LogPanicln(v ...any) {
	if checker.Options.GitHubLogger {
		logv := append([]any{"::error ::"}, v...)
		fmt.Println(logv...)
		panic(v)
	} else {
		log.Panicln(v...)
	}
}

func (checker *Checker) GitCmdOutputWithPipe(gitCommand string) string {
	cmd := "git -C '" + checker.Options.RepoDir + "' " + gitCommand
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		checker.LogPanicln("while try to execute command", cmd)
	}
	return string(out)
}

func (checker *Checker) GitCmdOutput(args ...string) string {
	args = append([]string{"-C", checker.Options.RepoDir}, args...)
	out, err := exec.Command(checker.Options.GitPath, args...).Output()
	if err != nil {
		checker.LogPanicln("while try to execute command", args)
	}
	return string(out)
}

func (checker *Checker) ParseArguments() {
	var flags OptionFlags
	flags.GitHubLogger = flag.Bool("github-logger", true, "Use GitHub workflow commands log style")
	flags.ExcludeBase = flag.Bool("exclude-base", true, "Exclude the base commit from checking")
	flags.BaseSha = flag.String("base-sha", "", "Git commit hash range: start")
	flags.HeadSha = flag.String("head-sha", "", "Git commit hash range: end")
	flags.RepoDir = flag.String("repo-dir", "", "Location of the git repo, use the working dir by default")

	flag.Parse()

	checker.Options.GitHubLogger = *flags.GitHubLogger
	checker.Options.ExcludeBase = *flags.ExcludeBase
	checker.Options.BaseSha = *flags.BaseSha
	checker.Options.HeadSha = *flags.HeadSha
	checker.Options.RepoDir = *flags.RepoDir
}

func (checker *Checker) CheckPrequirements() {
	var err error = nil
	checker.Options.GitPath, err = exec.LookPath("git")
	if err != nil {
		checker.LogPanicln("'git' executable not found: ", err)
	}

	if len(checker.Options.RepoDir) <= 0 {
		pwd, err := os.Getwd()
		if err != nil {
			checker.LogPanicln("Cannot get pwd as repo dir", err)
		}
		checker.Options.RepoDir = pwd
	}

	if len(checker.Options.HeadSha) <= 0 {
		headRev := strings.TrimSpace(checker.GitCmdOutput("log", "--pretty=format:%H", "HEAD", "-n1"))
		checker.Options.HeadSha = headRev
		checker.LogNoticeln("Head SHA not provided, will use current HEAD sha:", headRev)
	}

	if len(checker.Options.BaseSha) <= 0 {
		checker.Options.BaseSha = checker.Options.HeadSha
		checker.LogNoticeln("Base SHA not provided, will also use current HEAD.")
	}

	if !checker.Options.ExcludeBase {
		checker.Commits = []string{checker.Options.BaseSha}
	}

	commitsBetweenRaw := strings.TrimSpace(checker.GitCmdOutput("rev-list", checker.Options.BaseSha+".."+checker.Options.HeadSha))
	commitsBetween := strings.Split(commitsBetweenRaw, "\n")
	checker.Commits = append(checker.Commits, commitsBetween...)
}

func (checker *Checker) CheckCommits() Result {
	result := Result{
		Passed: true,
	}

	for _, commit := range checker.Commits {
		checked := checker.checkSingleCommit(commit)
		if !checked.Passed {
			result.Passed = false
		}
		result.Details = append(result.Details, checked)
	}

	return result
}

func (checker *Checker) PrintResult(result Result) {
	fmt.Println("====================")
	fmt.Println("Check Result Summary")
	fmt.Println("====================")
	if result.Passed {
		fmt.Println("Check: Passed")
	} else {
		checker.LogNoticeln("Check: Failed")
	}
	fmt.Println("Details:")
	for _, singleResult := range result.Details {
		if singleResult.Passed {
			if singleResult.IsFallback {
				checker.LogNoticeln(singleResult.SHA + " : Skipped (Reason: " + singleResult.Reason + ")")
			} else {
				fmt.Println(singleResult.SHA + " : Passed")
			}
		} else {
			checker.LogNoticeln(singleResult.SHA + " : Failed")
			checker.LogErrorln(singleResult.Reason, singleResult.Suggestion)
		}
	}
	fmt.Println("====================")
}

func (checker *Checker) checkSingleCommit(commitSha string) CheckedCommitMessage {
	requiredRules := list.New()
	conditionalRules := list.New()
	fallbackRules := list.New()
	allRequiredCheckPassed := true
	reason, suggestion := "", ""

	commitCheckData := CreateCommitCheckData(
		commitSha,
		strings.TrimSpace(checker.GitCmdOutput("show", "-s", "--format=%ae", commitSha)),
		strings.TrimSpace(checker.GitCmdOutput("show", "-s", "--format=%s", commitSha)),
		strings.TrimSpace(checker.GitCmdOutput("show", "-s", "--format=%b", commitSha)),
		strings.TrimSpace(checker.GitCmdOutputWithPipe("show -s --format=%b "+commitSha+" | git interpret-trailers --parse")),
	)

	for _, rule := range checker.Rules {
		switch rule.RuleType {
		case Rule_AlwaysCheck:
			requiredRules.PushBack(rule)
		case Rule_FallbackCheck:
			fallbackRules.PushBack(rule)
		case Rule_ConditionialCheck:
			conditionalRules.PushBack(rule)
		}
	}

	for e := conditionalRules.Front(); e != nil; e = e.Next() {
		rule := e.Value.(Rule)
		switch rule.ConditionType {
		case Condition_ByEmail:
			if strings.Contains(commitCheckData.AuthorEmail, rule.RuleData) {
				requiredRules.PushBack(rule)
			}
		}
	}

	for e := requiredRules.Front(); e != nil; e = e.Next() {
		rule := e.Value.(Rule)
		checkResult := commitCheckData.Check(rule)
		if !checkResult.Passed {
			allRequiredCheckPassed = false
			reason = checkResult.Reason
			suggestion = checkResult.Suggestion
		}
	}

	if allRequiredCheckPassed {
		return CheckedCommitMessage{
			SHA:        commitSha,
			Passed:     allRequiredCheckPassed,
			IsFallback: false,
		}
	}

	// fallback checks
	for e := fallbackRules.Front(); e != nil; e = e.Next() {
		rule := e.Value.(Rule)
		checkResult := commitCheckData.Check(rule)
		if checkResult.Passed {
			return CheckedCommitMessage{
				SHA:        commitSha,
				Passed:     true,
				IsFallback: true,
				Reason:     rule.RuleData,
			}
		}
	}

	return CheckedCommitMessage{
		SHA:        commitSha,
		Passed:     false,
		IsFallback: false,
		Reason:     reason,
		Suggestion: suggestion,
	}
}
