package controllers

import "regexp"

func sanitizeInput(inp string) string {
	r := regexp.MustCompile(`<.*?>`)
	return r.ReplaceAllString(inp, "")
}
