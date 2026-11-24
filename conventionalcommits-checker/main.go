// SPDX-FileCopyrightText: 2022 - 2026 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import "os"

func main() {
	options := ParseArguments()
	checker := CreateChecker(options)
	checker.CheckPrequirements()
	result := checker.CheckCommits()
	checker.PrintResult(result)
	if !result.Passed {
		os.Exit(1)
	}
}
