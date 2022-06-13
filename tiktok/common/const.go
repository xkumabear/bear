package common

import (
	"crypto/md5"
	"fmt"
	"io"
)

const (
	SuccessCode    = 0
	ParamsErr      = 101
	ParamsErrExist = 102
	SqlFindErr     = 202
	SqlAddErr      = 203
	SqlDelErr      = 204
	DSN            = "root:root@tcp(localhost:3306)/tiktok?charset=utf8&parseTime=True&loc=Local"
	DRIVER         = "mysql"
	ParamsErrMsg   = "Params is invalid "
	SqlFindErrMsg  = "Sql find is err"
	SqlAddErrMsg   = "Sql add is err"
	SqlDelErrMsg   = "Sql del is err"
	Url            = "http://172.27.7.47:8080"
	TokenKey       = "bear_tiktok"
	Salt           = "bear_tiktok"
)

func MD5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
