package main

var rMap = make(map[string]string, 0)

func RSet(key, value string) {
	rMap[key] = value
}

func RGet(key string) string {
	if val, ok := rMap[key]; ok {
		return val
	}
	return "nil"
}