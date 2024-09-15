// package utils

// import (
//     "time"
//     "github.com/dgrijalva/jwt-go"
//     "file-management/models"
//     "os"
// )

// // GenerateJWT creates a JWT token for a user
// func GenerateJWT(user models.User) (string, error) {
//     jwtSecret := os.Getenv("JWT_SECRET")

//     // Create a new token object
//     token := jwt.New(jwt.SigningMethodHS256)

//     // Set token claims
//     claims := token.Claims.(jwt.MapClaims)
//     claims["authorized"] = true
//     claims["email"] = user.Email
//     claims["exp"] = time.Now().Add(time.Hour * 24).Unix()  // Token expires in 24 hours

//     // Sign the token with the secret key
//     tokenString, err := token.SignedString([]byte(jwtSecret))
//     if err != nil {
//         return "", err
//     }

//     return tokenString, nil
// }

// // VerifyJWT verifies a JWT token and returns claims if valid
// func VerifyJWT(tokenString string) (jwt.MapClaims, error) {
//     jwtSecret := os.Getenv("JWT_SECRET")

//     // Parse the token
//     token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//         return []byte(jwtSecret), nil
//     })

//     if err != nil {
//         return nil, err
//     }

//     if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//         return claims, nil
//     }

//     return nil, err
// }


package utils

import (
    "github.com/dgrijalva/jwt-go"
    "os"
    "time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GenerateJWT generates a JWT token
func GenerateJWT(email string) (string, error) {
    claims := jwt.MapClaims{
        "email": email,
        "exp":   time.Now().Add(time.Hour * 24).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// VerifyJWT verifies the JWT token and returns claims
func VerifyJWT(tokenString string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, err
}