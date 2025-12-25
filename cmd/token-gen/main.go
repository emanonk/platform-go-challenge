package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	var (
		user     = flag.String("user", "user-1", "User id to set as JWT subject (sub)")
		issuer   = flag.String("iss", "favourites-api", "JWT issuer (iss)")
		audience = flag.String("aud", "web", "JWT audience (aud)")
		minutes  = flag.Int("minutes", 15, "Token validity in minutes")
		keyPath  = flag.String("key", "private.pem", "Path to RSA private key PEM (for signing)")
	)
	flag.Parse()

	fmt.Println("Generating JWT token...", *user, *issuer, *audience, *minutes, *keyPath)

	keyData, err := os.ReadFile(*keyPath)
	if err != nil {
		panic(err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		panic(err)
	}

	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   *user, // ✅ critical: user identity goes here
		Issuer:    *issuer,
		Audience:  jwt.ClaimStrings{*audience},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(*minutes) * time.Minute)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(privKey)
	if err != nil {
		panic(err)
	}

	fmt.Println(signed)
}
