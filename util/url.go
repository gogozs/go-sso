package util

import (
	"fmt"
	"golang.org/x/xerrors"
	nurl "net/url"
	"strings"
)

func BuildUrlQuery(url string, m map[string]interface{}) (string, error) {
	if m == nil {
		return url, nil
	}
	u, err := nurl.Parse(url)
	if err != nil {
		return url, xerrors.Errorf("Parse url error: %w", err)
	}
	var s strings.Builder
	for k, v := range m {
		_, err := s.WriteString(fmt.Sprintf("%s=%s", k, v))
		if err != nil {
			return "", xerrors.Errorf("WriteString error: %w", err)
		}
	}
	if u.RawQuery == "" {
		url = fmt.Sprintf("%s?%s", url, s.String())
	} else {
		url = fmt.Sprintf("%s&%s", url, s.String())
	}

	return url, nil
}
