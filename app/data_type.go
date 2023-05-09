package main

var rMap = make(map[string]string, 0)

func RSet(key, value string) {
	rMap[key] = value
}

func RGet(key string) string {
	if val, ok := rMap["foo"]; ok {
		return val
	}
	return "nil"
}