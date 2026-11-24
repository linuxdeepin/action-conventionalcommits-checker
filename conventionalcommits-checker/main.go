// SPDX-FileCopyrightText: 2022 - 2026 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import (
	"fmt"
	"os"
)

func main() {
	// client := github.NewClient(nil)
	// ctx := context.Background()

	// pull, rsp, err := client.PullRequests.Get(ctx, "linuxdeepin", "dtkgui", 40)
	// if err != nil {
	// 	log.Panicln(err)
	// }

	// pull.List
	options := ParseArguments()
	checker, err := CreateChecker(options)
	if err != nil {
		fmt.Println("::error title=Configuration Error::" + err.Error())
		os.Exit(1)
	}
	checker.CheckPrequirements()
	result := checker.CheckCommits()
	checker.PrintResult(result)
	if !result.Passed {
		os.Exit(1)
	}
}
