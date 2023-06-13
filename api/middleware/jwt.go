package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		fmt.Println("Token not present in the Header")
		return fmt.Errorf("Unauthorized")
	}

	claims, err := validateToken(token)
	if err != nil {
		return err
	}
	exp, ok := claims["expires"].(float64)
	if !ok {
		fmt.Println("can't covert to time.Time")
	}
	expTime := time.Unix(int64(exp), 0)
	if time.Now().After(expTime) {
		return fmt.Errorf("Token Expired")
	}

	return c.Next()
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Invalid singing method", token.Header["alg"])
			return nil, fmt.Errorf("Unauthorized")

		}

		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("falied to parse JWT token:", err)
		return nil, fmt.Errorf("Unauthorized")
	}

	if !token.Valid {
		fmt.Println("Invlaid token")
		return nil, fmt.Errorf("Unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Unauthorized")
	}
	return claims, nil
}
