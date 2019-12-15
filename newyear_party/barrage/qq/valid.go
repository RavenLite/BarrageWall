package qq
import (
	"crypto/sha1"
	"fmt"
	"io"
)

// 使用SHA-1算法散列字符串，返回HEX格式的结果
func sha1String(str string) string  {
	w := sha1.New()
	_, _ = io.WriteString(w, str)
	return fmt.Sprintf("%x", w.Sum(nil))
}

func UserInfoValid(sessionKey, raw , sign string) bool  {
	return sha1String(raw+ sessionKey) == sign
}
