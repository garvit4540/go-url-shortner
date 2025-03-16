package helpers

import (
	"github.com/garvit4540/go-url-shortner/trace"
	"os"
	"strings"
)

// If http is not in front of url, we put it there
func EnforceHttp(url string) string {
	if url[:4] != "http" {
		trace.LogInfo(trace.HttpEnforced, map[string]interface{}{"url": url})
		return "http://" + url
	}
	return url
}

// This function checks if this url is getting directed to localhost:3000
func RemoveDomainError(url string) bool {
	if url == os.Getenv("DOMAIN") {
		return false
	}
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]
	if newURL == os.Getenv("DOMAIN") {
		return false
	}
	return true
}
