package utils

import (
	"crypto/rand"
	"math/big"
)

const inviteCodeChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// GenerateInviteCode 生成 6 位随机大写字母+数字邀请码
func GenerateInviteCode() string {
	code := make([]byte, 6)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(inviteCodeChars))))
		code[i] = inviteCodeChars[n.Int64()]
	}
	return string(code)
}
