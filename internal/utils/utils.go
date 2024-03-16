package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"os"
)

type BatchBody struct {
	CorrelationId string `json:"correlation_id"`
	OriginalUrl   string `json:"original_url"`
}

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

func GenerateBatchBody(itemsCnt int) {
	arr := []BatchBody{}

	for i := 0; i <= itemsCnt; i++ {

		arr = append(arr, BatchBody{
			CorrelationId: GetRandomKey(5),
			OriginalUrl:   fmt.Sprintf("http://%s/yandex.ru/%s", GetRandomKey(4), GetRandomKey(7)),
		})
	}

	resp, err := json.MarshalIndent(arr, " ", "    ")
	if err != nil {
		fmt.Sprintln(err)
	}

	file, err := os.OpenFile("/tmp/batch_body.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Sprintln(err)
	}
	defer file.Close()

	file.Write(resp)

	//fmt.Println(string(resp))
}
