package utils

import (
	"net/url"
	"strings"

	"schej.it/server/logger"
)

func ParseArrayQueryParam(s string) []string {
	decoded, err := url.QueryUnescape(s)
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	arr := strings.Split(decoded, ",")
	return arr
}
