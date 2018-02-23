package models

func checkPassword(pw string) bool {
	if len(pw) >= 6 {
		return true
	} else {
		return false
	}
}

func checkBasicString(s string) bool {
	if len(s) <= 128 {
		return true
	} else {
		return false
	}
}
