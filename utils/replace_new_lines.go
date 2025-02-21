package utils

import (
	"html/template"
	"strings"
)

func ReplaceNewlines(s string) template.HTML {
	return template.HTML(strings.ReplaceAll(s, "\n", "<br>"))
}
