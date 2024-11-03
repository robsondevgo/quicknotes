package utils

import "github.com/mojocn/base64Captcha"

func GenerateCaptcha() (id string, value string, err error) {
	driver := base64Captcha.NewDriverDigit(100, 200, 4, 0.7, 100)
	captcha := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	id, value, _, err = captcha.Generate()
	return
}

func ValidateCaptcha(id string, answer string) bool {
	return base64Captcha.DefaultMemStore.Verify(id, answer, true)
}
