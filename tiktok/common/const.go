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

	DSN           = "douyin:123456@tcp(localhost:3306)/douyin?charset=utf8&parseTime=True&loc=Local"
	DRIVER        = "mysql"
	ParamsErrMsg  = "Params is invalid "
	SqlFindErrMsg = "Sql find is err"
	Url           = "http://127.0.0.1:8080"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

func MD5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
