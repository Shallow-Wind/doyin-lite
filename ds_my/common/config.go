package common

import "time"

const MySqlDSN = "root:cxc666@tcp(127.0.0.1:3306)/byte_dance?charset=utf8mb4"

// MD5Salt MD5加密时的盐
const MD5Salt = "UII34HJ6OIO"

// JWT
const (
	Issuer              = "xhx" // 签发人
	MySecret            = "Fy3Jfa5AD"
	TokenExpirationTime = 14 * 24 * time.Hour * time.Duration(1) // Token过期时间
)

// OSSPreURL OSS前缀
const OSSPreURL = "https://test-vedio-byte.oss-cn-beijing.aliyuncs.com/video/"

// SensitiveWordsPath 敏感词路径
const SensitiveWordsPath = "./utils/SensitiveWords.txt"
