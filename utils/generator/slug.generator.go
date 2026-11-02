package generator

import (
	"crypto/rand"
	"encoding/hex"
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

func SlugCollisionHandler(slug string) string{
	// if ente kaiyil same name ulla vera slug veruvane 

	b:=make([]byte, 2)
	 _, _ = rand.Read(b)
    return slug + "-" + hex.EncodeToString(b)
}