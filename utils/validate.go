package utils

import "regexp"

//验证手机号
func ValidatePhone(phone string) bool {
	regular := "^1[3456789]\\d{9}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(phone)
}

func ValidateEmail(email string) bool {
	//pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func MustAlphabetNumberUnderline(content string) bool {
	expression := `^[a-zA-Z0-9_\-]+$`
	reg := regexp.MustCompile(expression)
	return reg.MatchString(content)
}

func ValidateUrl(content string) bool {
	expression := `^(http|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?$`
	reg := regexp.MustCompile(expression)
	return reg.MatchString(content)
}

// ValidateIDNumber 身份证号
func ValidateIDNumber(content string) bool {
	expression := `(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)`
	reg := regexp.MustCompile(expression)
	return reg.MatchString(content)
}

func ValidateStrongPassword(content string) bool {
	expression := `^.*?[a-z]+(.*?){8,}$`
	expression1 := `^.*?[A-Z]+(.*?){8,}$`
	reg := regexp.MustCompile(expression)
	reg1 := regexp.MustCompile(expression1)
	return reg.MatchString(content) && reg1.MatchString(content) && len(content) >= 8
}
