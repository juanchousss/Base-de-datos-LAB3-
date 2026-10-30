package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// IMPORTANTE: Esta clave debe ser exactamente LA MISMA que definiste en internal/services/jwt.go
var jwtKey = []byte("super_secreta_lab3_ucla")

// Creamos un tipo personalizado para la llave del contexto (buena práctica en Go para evitar colisiones)
type contextKey string

const UserIDKey contextKey = "userID"

// Autenticacion es el middleware que protege nuestras rutas
func Authentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Extraer el token del Header "Authorization"
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Acceso denegado: Token requerido", http.StatusUnauthorized)
			return
		}

		// 2. Validar que el formato sea "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Acceso denegado: Formato de token inválido", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		// 3. Parsear y validar la firma del token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validar que el algoritmo de encriptación sea el correcto (HMAC)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return jwtKey, nil
		})

		// Si hay un error o el token expiró/fue alterado
		if err != nil || !token.Valid {
			http.Error(w, "Acceso denegado: Token inválido o expirado", http.StatusUnauthorized)
			return
		}

		// 4. Extraer los datos (claims) que guardamos dentro del token
		if claims, ok := token.Claims.(jwt.MapClaims); ok {

			// VALIDACIÓN SEGURA: Preguntamos si existe antes de convertir
			val, existe := claims["id_user"]
			if !existe {
				http.Error(w, "Acceso denegado: El token no contiene id_user", http.StatusUnauthorized)
				return
			}

			// Como ya sabemos que existe, la conversión es 100% segura
			userID := int(val.(float64))

			// 5. Inyectar el ID del usuario en el Contexto de la petición HTTP
			ctx := context.WithValue(r.Context(), UserIDKey, userID)

			// Dejar pasar la petición al siguiente controlador, llevando el contexto actualizado
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Acceso denegado: Error al leer el token", http.StatusUnauthorized)
			return
		}
	}
}
