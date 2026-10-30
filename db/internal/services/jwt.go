package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// En un entorno de producción real, esto DEBE venir de variables de entorno (.env)
var jwtKey = []byte("super_secreta_lab3_ucla")

// GenerarToken crea un JWT válido por 24 horas con los datos del usuario
func GenerateToken(id int, email string) (string, error) {
	// Definimos la información ("claims") que viajará dentro del token
	claims := jwt.MapClaims{
		"id_user": id,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Caduca en 24 horas
	}

	// Creamos el token usando el algoritmo de encriptación HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Lo firmamos con nuestra clave secreta
	return token.SignedString(jwtKey)
}
