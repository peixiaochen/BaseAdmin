package lib

import (
	"github.com/gookit/filter"
	"github.com/gookit/goutil/strutil"
)

func SubstrMd5(str string, start int, length int) string {
	return strutil.Md5(filter.Substr(strutil.Md5(str), start, length))
}
