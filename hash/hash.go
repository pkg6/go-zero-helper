package hash

import "golang.org/x/crypto/bcrypt"

// Make 对明文字符串进行bcrypt哈希加密
// 参数 str：需要加密的明文密码
// 返回值：加密后的哈希字符串
// 注意：使用bcrypt默认强度，内部已自动处理盐值，无需手动加盐
func Make(str string) string {
	// GenerateFromPassword 根据明文和默认加密成本生成哈希字节数组
	hashed, _ := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(hashed)
}

// Check 对比明文密码与哈希密码是否匹配
// 参数 plain：用户输入的明文密码
// 参数 hashed：数据库中存储的哈希密码
// 返回值：匹配返回true，不匹配返回false
func Check(plain string, hashed string) bool {
	// CompareHashAndPassword 对比哈希值与明文，不匹配则返回错误
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}
