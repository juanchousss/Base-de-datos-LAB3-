package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Driver puro de PostgreSQL
)

// InitDB inicializa y retorna la conexión a la base de datos
func InitDB(user, password, dbname, host, port string) *sql.DB {
	// Cadena de conexión
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=America/Caracas",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error abriendo la base de datos: %v", err)
	}

	// Comprobamos que la conexión sea exitosa
	if err = db.Ping(); err != nil {
		log.Fatalf("Error haciendo ping a la base de datos: %v", err)
	}

	log.Println("Conexión a PostgreSQL establecida con éxito")
	return db
}
