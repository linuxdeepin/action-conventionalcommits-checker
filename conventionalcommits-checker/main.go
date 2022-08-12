// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import "os"

func main() {
	// client := github.NewClient(nil)
	// ctx := context.Background()

	// pull, rsp, err := client.PullRequests.Get(ctx, "linuxdeepin", "dtkgui", 40)
	// if err != nil {
	// 	log.Panicln(err)
	// }

	// pull.List
	checker := CreateChecker()
	checker.ParseArguments()
	checker.CheckPrequirements()
	result := checker.CheckCommits()
	checker.PrintResult(result)
	if !result.Passed {
		os.Exit(1)
	}
}
