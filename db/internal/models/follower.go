package models

import "time"

type Follower struct {
	IDFollower int       `json:"id_follower"`
	IDFollow   int       `json:"id_follow"`
	CreatedAt  time.Time `json:"created_at"`
}
