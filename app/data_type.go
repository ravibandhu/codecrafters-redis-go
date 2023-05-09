package main

import (
	"fmt"
	"strconv"
	"time"
)

type V struct {
	value    string
	expireIn *time.Time
}

type Store struct {
	storage map[string]V
}

func NewStore() *Store {
	return &Store{
		storage: make(map[string]V),
	}
}

func (s *Store) Set(args []Value) (string, error) {
	key, value, expiry := "", "", ""
	expireIn := time.Duration(0)
	if len(args) == 4 {
		key, value, expiry = args[0].String(), args[1].String(), args[3].String()
		expiringIN, err := strconv.Atoi(expiry)
		if err != nil {
			return "-ERR invalid expire time in SET\r\n", err
		}
		expireIn = time.Duration(expiringIN) * time.Millisecond
	} else {
		key, value = args[0].String(), args[1].String()
	}
	if expireIn == time.Duration(0) {
		s.storage[key] = V{value: value}
		return "OK", nil
	}
	e := time.Now().Add(expireIn)
	s.storage[key] = V{
		value:    value,
		expireIn: &e,
	}
	return "OK", nil
}

func (s *Store) Get(key string) string {
	if v, ok := s.storage[key]; ok {
		if v.expireIn == nil || v.expireIn.After(time.Now()) {
			return v.value
		} else {
			delete(s.storage, key)
		}
	}
	return ""
}

func prepareStringResp(s string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", s))
}
func prepareStringRespWithLength(arg string) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
}
func prepareArrayResp(args []Value) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(args[0].String()), args[0].String()))
}
func prepareErrResp() []byte {
	return []byte("$-1\r\n")
}
