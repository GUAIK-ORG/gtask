package session

import (
	"gtask/pkg/utils/crypt"
	"strconv"
	"sync"
	"time"
)

type Session struct {
	secretKey string
	token     string
}

var ins *Session
var once sync.Once

func Ins() *Session {
	once.Do(func() {
		ins = &Session{}
	})
	return ins
}

func (s *Session) Init(secretKey string) {
	s.secretKey = secretKey
}

func (s *Session) GetToken(secretKey string) string {
	if secretKey == s.secretKey {
		s.token = crypt.Str2MD5(strconv.FormatInt(time.Now().UnixNano(), 16))
		return s.token
	}
	return ""
}

func (s *Session) CheckToken(token string) bool {
	return token == s.token
}
