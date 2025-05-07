package utils

import (
	"net/url"
	"strings"
)

func DecodeURL(url string) string {
	url = strings.ReplaceAll(url, "%257B", "{")
	url = strings.ReplaceAll(url, "%2522", "\"")
	url = strings.ReplaceAll(url, "%253A", ":")
	url = strings.ReplaceAll(url, "%252C", ",")
	url = strings.ReplaceAll(url, "%257D", "}")
	url = strings.ReplaceAll(url, "%C3%A7", "ç")
	url = strings.ReplaceAll(url, "%C3%A3", "ã")
	url = strings.ReplaceAll(url, "%C3%B5", "õ")
	url = strings.ReplaceAll(url, "%C3%A0", "à")
	url = strings.ReplaceAll(url, "%C3%A1", "á")
	url = strings.ReplaceAll(url, "%C3%AD", "í")
	url = strings.ReplaceAll(url, "%C3%B3", "ó")
	url = strings.ReplaceAll(url, "%C3%A9", "é")
	url = strings.ReplaceAll(url, "%C3%89", "É")
	url = strings.ReplaceAll(url, "%C3%AA", "ê")
	return url
}

func EncodeURL(input string) (string, error) {
	parsedURL, err := url.Parse(input)
	if err != nil {
		return "", err
	}

	encodedQuery := parsedURL.Query()
	for key, values := range encodedQuery {
		for i, value := range values {
			encodedValue := url.QueryEscape(value)
			values[i] = encodedValue
		}
		encodedQuery[key] = values
	}
	parsedURL.RawQuery = encodedQuery.Encode()

	return parsedURL.String(), nil
}
