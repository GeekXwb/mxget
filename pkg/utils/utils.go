package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	re = regexp.MustCompile(`[\\/:*?"<>|]`)
)

func TrimInvalidFilePathChars(path string) string {
	return strings.TrimSpace(re.ReplaceAllString(path, " "))
}

func ToJSON(data interface{}, escapeHTML bool) string {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(escapeHTML)
	enc.SetIndent("", "\t")
	if enc.Encode(data) != nil {
		return "{}"
	}
	return buf.String()
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Input(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", prompt)
	text, _ := reader.ReadString('\n')
	input := strings.TrimSpace(text)
	if input == "" {
		return Input(prompt)
	}
	return input
}
