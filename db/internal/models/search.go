package models

// UserSearch es una versión ligera del usuario para los resultados de búsqueda
type UserSearch struct {
	IDUser   int    `json:"id_user"`
	Username string `json:"username"`
}

// SearchResponse es el paquete maestro que devolverá el buscador
type SearchResponse struct {
	Users []UserSearch `json:"users"`
	Posts []Post       `json:"posts"`
}
