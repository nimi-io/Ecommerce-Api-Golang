package helpers

import (
	"Ecommerce-Api/database"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var secretKey = []byte("secret-key")

func GenerateToken(user database.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
	})
	tokenString, err := token.SignedString(secretKey)
	return tokenString, err
}

// func VerifyToken(tokenString string) (any, error) {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return secretKey, nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	if !token.Valid {
// 		log.Println("invalid token")
// 		return nil, err
// 	}

//		return token, nil
//	}

func VerifyToken(tokenString string) (*jwt.Token, database.User, error) {
	var user database.User

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	claims, _ = token.Claims.(jwt.MapClaims)

	claims["id"] = token.Claims.(jwt.MapClaims)["id"]

	if err != nil || !token.Valid {
		return nil, user, errors.New("invalid token")
	}
	///////////////////////////////////////////////////////

	_, ok := claims["id"].(float64)
	if !ok {
		return nil, user, errors.New("user ID not found in token claims")
	}

	// err = database.DB.QueryRow("SELECT * FROM users WHERE id = ?", int(userID)).Scan(
	// 	&user.ID, &user.Email, &user.Username, &user.Password, &user.CreatedAt,
	// )
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return nil, user, errors.New("user not found")
	// 	}
	// 	return nil, user, err
	// }

	///////////////////////////////////////////////////////

	return token, user, nil
}

// func VerifyToken(tokenString string) (*jwt.Token, database.User, error) {
// 	var user database.User

// 	claims := jwt.MapClaims{}

// 	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 		return secretKey, nil
// 	})

// 	// /
// 	userID, ok := claims["id"].(float64) // Assuming ID is stored as float64 in token claims
// 	if !ok {
// 		return nil, user, errors.New("user ID not found in token claims")
// 	}
// 	err = database.DB.QueryRow("SELECT * FROM users WHERE id = ?", int(userID)).Scan(
// 		&user.ID, &user.Email, &user.Username, &user.Password, &user.CreatedAt,
// 	)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, user, errors.New("user not found")
// 		}
// 		return nil, user, err
// 	}

// 	log.Println("what the hell is this", user)

// 	///
// 	claims, _ = token.Claims.(jwt.MapClaims)

// 	claims["id"] = token.Claims.(jwt.MapClaims)["id"]

// 	if err != nil || !token.Valid {
// 		return nil, user, errors.New("invalid token")
// 	}

// 	return token, user, nil
// }
