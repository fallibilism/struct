package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"v/pkg/config"

	protocol "github.com/fallibilism/protocol/go_protocol"
)

func ClaimsToJWT(claims *protocol.LtiAuthClaims) (string, error) {
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: []byte(config.Conf.JWTSecret)},
		(&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", err
	}

	cl := jwt.Claims{
		Issuer:    config.Conf.JWTIssuer,
		NotBefore: jwt.NewNumericDate(time.Now()),
		Expiry:    jwt.NewNumericDate(time.Now().Add(time.Hour * 2)), // valid for 2 hours
		Subject:   claims.UserId,
	}

	return jwt.Signed(sig).Claims(cl).Claims(claims).CompactSerialize()
}

// to hash id using *sha1*
func GenHash(id string) (hash string) {
	hasher := sha1.New()
	hasher.Write([]byte(id))
	hash = hex.EncodeToString(hasher.Sum(nil))

	return hash
}
