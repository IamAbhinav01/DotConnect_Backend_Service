package generator

import (
	"regexp"
	"strings"
)

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

func Slugify(name string) string{
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug,"'","")
	slug = nonAlphanumeric.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")


	if slug == ""{
		slug = "workspace"
	}

	return slug
}