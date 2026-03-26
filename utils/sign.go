package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// BuildSignString 将参数按键名升序排序，并拼接为 k=v&k2=v2 格式的待签名字符串。
func BuildSignString(params map[string]interface{}, ignoreKeys ...string) string {
	ignores := map[string]struct{}{}
	for _, key := range ignoreKeys {
		ignores[key] = struct{}{}
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		if _, ok := ignores[k]; ok {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, key := range keys {
		// 失败或异常情况 部分参数不参与签名
		if tradeStatus, ok := params["tradeStatus"]; ok && tradeStatus != "1" {
			for _, item := range []string{"cancelReason", "coinAmount", "total", "unitPrice"} {
				if item == key {
					continue
				}
			}
		}
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(toSignValue(params[key]))
		sb.WriteString("&")
	}
	str := sb.String()
	str = strings.TrimRight(str, "&")
	// str = "coinAmount=100.00000000&coinSign=5&companyOrderNum=M_1710402289968&errorOrderAmount=6.93480000&errorOrderTotal=50.00000000&intentOrderNo=4628493946291201&successAmount=6.93480000&total=50.00000000&tradeOrderTime=1710402308000&tradeStatus=1&unitPrice=7.21000000"
	//str = "areaCode=90&asyncUrl=https://api-test.logtec.dev/fapi/payment/psp/public/chipay/deposit/back&companyId=4548687456245761&companyOrderNum=202603260623050639&name=赫敏·珍珍·格兰杰 &phone=5300231651&syncUrl=https://uc-test.logtec.dev/&totalAmount=22"
	fmt.Printf("sign raw string: %s\n", str)

	return str
}

// BuildSignStringFromBody 将 JSON Body 解析为 map 后生成待签名字符串。
func BuildSignStringFromBody(body []byte, ignoreKeys ...string) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(body, &params); err != nil {
		return "", err
	}
	return BuildSignString(params, ignoreKeys...), nil
}

// SignSHA256RSA 使用 SHA256withRSA 对原始字符串签名，返回 Base64 签名串。
func SignSHA256RSA(rawString string, privateKey string) (string, error) {
	key, err := parsePKCS8PrivateKey(privateKey)
	if err != nil {
		return "", fmt.Errorf("parse private key failed: %w, privateKeyFingerprint=%s", err, keyFingerprint(privateKey))
	}
	hashed := sha256.Sum256([]byte(rawString))
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("rsa sign failed: %w, rawString=%s, hash=%s, privateKeyFingerprint=%s", err, rawString, hex.EncodeToString(hashed[:]), keyFingerprint(privateKey))
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifySHA256RSA 使用 SHA256withRSA 校验签名，失败时返回带完整上下文的错误信息。
func VerifySHA256RSA(rawString string, sign string, publicKey string) (bool, error) {
	key, err := parseX509PublicKey(publicKey)
	if err != nil {
		return false, fmt.Errorf("parse public key failed: %w, publicKeyFingerprint=%s", err, keyFingerprint(publicKey))
	}
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, fmt.Errorf("decode signature failed: %w, signLen=%d, signPrefix=%s", err, len(sign), shorten(sign, 48))
	}
	hashed := sha256.Sum256([]byte(rawString))
	hashHex := hex.EncodeToString(hashed[:])
	if err = rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], signBytes); err != nil {
		if errors.Is(err, rsa.ErrVerification) {
			return false, fmt.Errorf("rsa verify failed: %w, rawString=%s, hash=%s, signLen=%d, signPrefix=%s, publicKeyFingerprint=%s", err, rawString, hashHex, len(signBytes), shorten(sign, 48), keyFingerprint(publicKey))
		}
		return false, fmt.Errorf("rsa verify unexpected error: %w, rawString=%s, hash=%s, signLen=%d, signPrefix=%s, publicKeyFingerprint=%s", err, rawString, hashHex, len(signBytes), shorten(sign, 48), keyFingerprint(publicKey))
	}
	return true, nil
}

// SignMap 对 map 参数生成待签名串并进行签名。
func SignMap(params map[string]interface{}, privateKey string, ignoreKeys ...string) (string, string, error) {
	raw := BuildSignString(params, ignoreKeys...)
	sign, err := SignSHA256RSA(raw, privateKey)
	if err != nil {
		return "", "", fmt.Errorf("sign map failed: %w, rawString=%s", err, raw)
	}
	return sign, raw, nil
}

// VerifyMap 对 map 参数进行验签，返回是否通过、待验签串和错误。
func VerifyMap(params map[string]interface{}, sign string, publicKey string, ignoreKeys ...string) (bool, string, error) {
	raw := BuildSignString(params, ignoreKeys...)
	ok, err := VerifySHA256RSA(raw, sign, publicKey)
	if err != nil {
		return false, raw, err
	}
	return ok, raw, err
}

// VerifyBody 从完整 JSON Body 中读取 sign 并完成验签。
func VerifyBody(body []byte, publicKey string) (bool, string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(body, &params); err != nil {
		return false, "", fmt.Errorf("unmarshal verify body failed: %w", err)
	}
	signValue, ok := params["sign"]
	if !ok {
		return false, "", fmt.Errorf("missing sign field")
	}
	sign := toSignValue(signValue)
	delete(params, "sign")
	raw := BuildSignString(params)
	flag, err := VerifySHA256RSA(raw, sign, publicKey)
	return flag, raw, err
}

// parsePKCS8PrivateKey 解析私钥，支持 PKCS#8 DER/PEM 与 PKCS#1。
func parsePKCS8PrivateKey(input string) (*rsa.PrivateKey, error) {
	keyBytes, err := decodeKey(input)
	if err != nil {
		return nil, fmt.Errorf("decode private key failed: %w", err)
	}
	if parsed, err := x509.ParsePKCS8PrivateKey(keyBytes); err == nil {
		if rsaKey, ok := parsed.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
		return nil, fmt.Errorf("key is not rsa private key")
	}
	if parsed, err := x509.ParsePKCS1PrivateKey(keyBytes); err == nil {
		return parsed, nil
	}
	return nil, fmt.Errorf("failed to parse private key")
}

// parseX509PublicKey 解析公钥，支持 X509 DER/PEM 与 PKCS#1。
func parseX509PublicKey(input string) (*rsa.PublicKey, error) {
	keyBytes, err := decodeKey(input)
	if err != nil {
		return nil, fmt.Errorf("decode public key failed: %w", err)
	}
	if parsed, err := x509.ParsePKIXPublicKey(keyBytes); err == nil {
		if rsaKey, ok := parsed.(*rsa.PublicKey); ok {
			return rsaKey, nil
		}
		return nil, fmt.Errorf("key is not rsa public key")
	}
	if parsed, err := x509.ParsePKCS1PublicKey(keyBytes); err == nil {
		return parsed, nil
	}
	return nil, fmt.Errorf("failed to parse public key")
}

// decodeKey 支持从 PEM 或 Base64 DER 两种格式读取密钥字节。
func decodeKey(key string) ([]byte, error) {
	trimmed := strings.TrimSpace(key)
	if trimmed == "" {
		return nil, fmt.Errorf("empty key")
	}
	if block, _ := pem.Decode([]byte(trimmed)); block != nil {
		return block.Bytes, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		return nil, fmt.Errorf("decode base64 key failed: %w", err)
	}
	return decoded, nil
}

// toSignValue 将任意参数值转换为签名使用的字符串表示。
func toSignValue(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case json.Number:
		return v.String()
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case int:
		return strconv.Itoa(v)
	case int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// keyFingerprint 生成密钥摘要用于日志定位，避免直接输出密钥本体。
func keyFingerprint(key string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(key)))
	return hex.EncodeToString(sum[:8])
}

// shorten 将超长字符串裁剪为可读日志片段。
func shorten(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max] + "...(truncated)"
}
