package conversion

import (
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

func XssFilter(str string) string {
	p := bluemonday.UGCPolicy()
	ret := p.Sanitize(
		str,
	)

	return strings.TrimSpace(ret)
}
