package encoder

import (
	"net/url"
)

// Конвертируем запрос для использование в качестве части URL
func UrlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
