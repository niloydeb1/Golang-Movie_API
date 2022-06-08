package v1

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"github.com/golang-jwt/jwt"
	"github.com/niloydeb1/Golang-Movie_API/config"
	"log"
	"time"
)

// RsaKeys object reference for v1.RsaKeys
var rsaKeys *RsaKeys = nil

type Jwt struct {
}

func (j Jwt) GetRsaKeys() *RsaKeys {
	if rsaKeys == nil {
		rsaKeys = &RsaKeys{
			PrivateKey: j.GetPrivateKey(),
			PublicKey:  j.GetPublicKey(),
		}
	}
	return rsaKeys
}

func (j Jwt) GenerateToken(userUUID string, duration int64, data interface{}) (string, string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims = jwt.MapClaims{
		"exp":  time.Now().UTC().Add(time.Duration(duration) * time.Millisecond).Unix(),
		"iat":  time.Now().UTC().Unix(),
		"sub":  userUUID,
		"data": data,
	}
	tokenString, err := token.SignedString(j.GetRsaKeys().PrivateKey)
	if err != nil {
		return "", "", err
	}
	token.Claims = jwt.MapClaims{
		"exp":  time.Now().UTC().Add(time.Duration(duration+duration/4) * time.Millisecond).Unix(),
		"iat":  time.Now().UTC().Unix(),
		"sub":  userUUID,
		"data": data,
	}
	refreshTokenStr, err := token.SignedString(j.GetRsaKeys().PrivateKey)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenStr, nil
}

func (j Jwt) IsTokenValid(tokenString string) bool {
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return j.GetRsaKeys().PublicKey, nil
	})

	var tm time.Time
	switch iat := claims["exp"].(type) {
	case float64:
		tm = time.Unix(int64(iat), 0)
	case json.Number:
		v, _ := iat.Int64()
		tm = time.Unix(v, 0)
	}
	if time.Now().UTC().After(tm) {
		return false
	}
	return true
}

func (j Jwt) GetPrivateKey() *rsa.PrivateKey {
	block, rest := pem.Decode([]byte(config.PrivateKey))
	if rest != nil {
		log.Print(rest)
	}
	privateKeyImported, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Print(err.Error())
	}
	return privateKeyImported
}

func (j Jwt) GetPublicKey() *rsa.PublicKey {
	block, _ := pem.Decode([]byte(config.Publickey))
	publicKeyImported, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		log.Print(err.Error())
	}
	return publicKeyImported
}