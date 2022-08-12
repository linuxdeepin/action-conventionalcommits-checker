package main

type CheckedCommitMessage struct {
	SHA        string
	Passed     bool
	IsFallback bool
	Reason     string
	Suggestion string
}

type Result struct {
	Passed  bool
	Details []CheckedCommitMessage
}
