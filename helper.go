package main

import "strings"

func removeNextLine(line string) string {
	return strings.TrimSuffix(line, "\n")
}