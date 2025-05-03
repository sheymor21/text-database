package Test

import "strings"

func getId(row string) string {
	split := strings.Split(row, "|")
	return strings.TrimSpace(split[2])
}
