package helper

import (
	"html"
	"strings"
)

func TrimAndHtmlEscape(str string) string {
	return html.EscapeString(strings.TrimSpace(str))
}
