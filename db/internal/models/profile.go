package models

// UserProfile agrupa los datos del usuario, sus estadísticas y su muro de publicaciones
type UserProfile struct {
	IDUser         int    `json:"id_user"`
	Email          string `json:"email"`
	FollowersCount int    `json:"followers_count"`
	FollowingCount int    `json:"following_count"`
	Posts          []Post `json:"posts"`
}
