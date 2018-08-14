package util

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func JsonStringToMap (jsonStr string) map[string]interface{} {
	m := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &m)
	return m
}

func ReadCloserToMap (res io.ReadCloser) map[string]interface{} {
	jsonByte, _ := ioutil.ReadAll(res)
	return JsonStringToMap(string(jsonByte))
}
// map转qs
func MapToQueryString (m map[string]string) string {
	qs := ""
	for k, v := range m {
		qs += "&" + k + "=" + v
	}
	return qs
}