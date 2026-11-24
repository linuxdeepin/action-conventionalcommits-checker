package main

type Option struct {
	GitHubLogger             bool
	ExcludeBase              bool
	GitPath                  string
	BaseSha                  string
	HeadSha                  string
	RepoDir                  string
	CheckConventionalCommits bool
	CheckRequireBody         bool
	CheckRequireLogTrailer   bool
	CheckRequireIssueTrailer bool
	AllowFileSyncCommits     bool
	AllowMergeCommits        bool
}

type OptionFlags struct {
	GitHubLogger             *bool
	ExcludeBase              *bool
	GitPath                  *string
	BaseSha                  *string
	HeadSha                  *string
	RepoDir                  *string
	CheckConventionalCommits *bool
	CheckRequireBody         *bool
	CheckRequireLogTrailer   *bool
	CheckRequireIssueTrailer *bool
	AllowFileSyncCommits     *bool
	AllowMergeCommits        *bool
}
