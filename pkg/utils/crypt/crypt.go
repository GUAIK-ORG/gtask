package crypt

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

func Str2MD5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func URLSafeBase64(d []byte) string {
	bytearr := base64.StdEncoding.EncodeToString(d)
	safeurl := strings.Replace(string(bytearr), "/", "_", -1)
	safeurl = strings.Replace(safeurl, "+", "-", -1)
	return safeurl
}
