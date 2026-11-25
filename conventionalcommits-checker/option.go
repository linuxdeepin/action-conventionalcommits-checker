// SPDX-FileCopyrightText: 2026 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

package main

type Option struct {
	GitHubLogger   bool
	ExcludeBase    bool
	GitPath        string
	BaseSha        string
	HeadSha        string
	RepoDir        string
	Rules          string
	ExtraRules     string
	RulesFile      string
	ExtraRulesFile string
}

type OptionFlags struct {
	GitHubLogger   *bool
	ExcludeBase    *bool
	GitPath        *string
	BaseSha        *string
	HeadSha        *string
	RepoDir        *string
	Rules          *string
	ExtraRules     *string
	RulesFile      *string
	ExtraRulesFile *string
}
