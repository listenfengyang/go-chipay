package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"testing"
)

func TestBuildSignString(t *testing.T) {
	params := map[string]interface{}{
		"companyId":       "12511234561",
		"companyOrderNum": "M_1636027918971",
		"areaCode":        "86",
		"phone":           "18900000006",
		"totalAmount":     300,
		"asyncUrl":        "http://demo.chippaytest.com/#/v2",
		"syncUrl":         "http://demo.chippaytest.com/#/v2",
	}
	raw := BuildSignString(params)
	expected := "areaCode=86&asyncUrl=http://demo.chippaytest.com/#/v2&companyId=12511234561&companyOrderNum=M_1636027918971&phone=18900000006&syncUrl=http://demo.chippaytest.com/#/v2&totalAmount=300"
	if raw != expected {
		t.Fatalf("unexpected raw sign string: %s", raw)
	}
}

func TestSignAndVerify(t *testing.T) {
	privateKey, publicKey := mustGenerateKeyPair()
	raw := "a=1&b=2"
	sign, err := SignSHA256RSA(raw, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	ok, err := VerifySHA256RSA(raw, sign, publicKey)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("verify failed")
	}
}

func TestVerifyBody(t *testing.T) {
	privateKey, publicKey := mustGenerateKeyPair()
	body := []byte(`{"companyOrderNum":"S_1111111","tradeStatus":"1"}`)
	signRaw, err := BuildSignStringFromBody(body)
	if err != nil {
		t.Fatal(err)
	}
	sign, err := SignSHA256RSA(signRaw, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	withSign := []byte(`{"companyOrderNum":"S_1111111","tradeStatus":"1","sign":"` + sign + `"}`)
	ok, _, err := VerifyBody(withSign, publicKey)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("verify body failed")
	}
}

func mustGenerateKeyPair() (string, string) {
	private, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	privateDer, err := x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		panic(err)
	}
	publicDer, err := x509.MarshalPKIXPublicKey(&private.PublicKey)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(privateDer), base64.StdEncoding.EncodeToString(publicDer)
}
