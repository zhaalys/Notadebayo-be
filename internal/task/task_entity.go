package task

import "time"

type Task struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primarykey" json:"id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	UserId      string    `gorm:"index" json:"userId"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	Label       string    `gorm:"not null" json:"labelId"`
}
