package main

type Option struct {
	GitHubLogger bool
	ExcludeBase  bool
	GitPath      string
	BaseSha      string
	HeadSha      string
	RepoDir      string
}

type OptionFlags struct {
	GitHubLogger *bool
	ExcludeBase  *bool
	GitPath      *string
	BaseSha      *string
	HeadSha      *string
	RepoDir      *string
}
