package util

import (
	"net/url"
)

//
//func GetEncodedUrl(str string) string {
//	if str == "" {
//		return str
//	}
//	u, err := url.Parse(str)
//	if err != nil {
//		logger.Error(err)
//		return str
//	}
//	result := ""
//	if len(u.Query()) > 0 {
//		result = fmt.Sprintf(`%s://%s%s?%s`, u.Scheme, u.Host, u.EscapedPath(), u.Query().Encode())
//	} else {
//		result = fmt.Sprintf(`%s://%s%s`, u.Scheme, u.Host, u.EscapedPath())
//	}
//	return result
//}
//
func GetEncodedUrl(su string) string {
	pu, err := url.Parse(su)
	if err != nil {
		return su
	}

	if len(pu.RawQuery) >= 0 {
		tq, err := url.ParseQuery(pu.RawQuery)
		if err != nil {
			return su
		}
		pu.RawQuery = tq.Encode()
	}

	return pu.String()
}
