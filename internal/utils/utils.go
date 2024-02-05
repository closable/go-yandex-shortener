package utils

import "net/url"

// check url
func ValidateURL(txtURL string) bool {
	u, err := url.Parse(txtURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}
