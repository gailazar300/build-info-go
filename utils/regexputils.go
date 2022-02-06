package utils

import (
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"regexp"
)

const CredentialsInUrlRegexp = `(http|https|git)://.+@`

func GetRegExp(regex string) (*regexp.Regexp, error) {
	regExp, err := regexp.Compile(regex)
	if errorutils.CheckError(err) != nil {
		return nil, err
	}
	return regExp, nil
}
