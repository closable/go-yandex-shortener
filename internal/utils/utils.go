package utils

import (
	"math/rand"
	"net/url"
)

// check url
func ValidateURL(txtURL string) bool {
	u, err := url.Parse(txtURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func GetRandomKey(length int) string {
	chars := "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	key := ""
	for i := 0; i < length; i++ {
		c := chars[rand.Intn(len(chars))]
		key += string(c)

	}
	return key
}
