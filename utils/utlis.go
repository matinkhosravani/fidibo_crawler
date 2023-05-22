package utils

import (
	"fmt"
	"regexp"
)

func ExtractBoodIDFromURL(url string) (string, error) {
	re := regexp.MustCompile(`model/(\d+)-`)
	match := re.FindStringSubmatch(url)
	if len(match) > 1 {
		return match[1], nil
	}

	return "", fmt.Errorf("cannot extarct Id From %v", url)
}
