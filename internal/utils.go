package internal

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func InsertSorted(strings []string, toInsert string) []string {
	i := sort.SearchStrings(strings, toInsert)

	strings = append(strings, "")
	copy(strings[i+1:], strings[i:])
	strings[i] = toInsert

	return strings
}

func ContainsBinarySearch(strings []string, target string) bool {
	mid := len(strings) / 2

	switch {
	case len(strings) == 0:
		return false
	case strings[mid] > target:
		return ContainsBinarySearch(strings[:mid], target)
	case strings[mid] < target:
		return ContainsBinarySearch(strings[mid+1:], target)
	default:
		return true
	}
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return []string{}, err
	}

	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func WriteLines(path string, lines []string) error {
	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()

	w := bufio.NewWriter(file)

	for _, line := range lines {
		fmt.Fprintln(w, line)
	}

	return w.Flush()
}
