package models

import "time"

type Like struct {
	IDUser    int       `json:"id_user"`
	IDPost    int       `json:"id_post"`
	CreatedAt time.Time `json:"created_at"`
}
