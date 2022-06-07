package report

import "testing"

const testSubject = "Test subject"
const testBody = "Test body"

func TestEmail(t *testing.T) {
	// 请填入测试目标邮箱
	SendMail(testSubject, testBody, "test@test.com")
}
