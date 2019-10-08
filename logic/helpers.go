package logic

import (
	"bufio"
	"regexp"
	"strings"
)

func scanWithPrefix(sc *bufio.Scanner, prefix string) (string, error) {
	if !sc.Scan() {
		return "", errInvalidFormat
	}

	return strings.TrimSpace(strings.TrimPrefix(sc.Text(), prefix)), nil
}

func getContent(rg *regexp.Regexp, desc string) string {
	results := rg.FindStringSubmatch(desc)
	if len(results) != 2 {
		return ""
	}

	return strings.TrimSpace(results[1])
}
