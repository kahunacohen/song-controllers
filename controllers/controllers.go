package controllers

import "regexp"

func SanitizeInput(inp string) string {
	r := regexp.MustCompile(`<.*?>`)
	return r.ReplaceAllString(inp, "")
}
