package models

type Item struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"type:varchar(255);not null" validate:"required"`
	Description string `json:"description" gorm:"type:text;not null" validate:"required"`
	Icon        string `json:"icon" gorm:"type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null" validate:"required"`
	Color       string `json:"color" gorm:"type:varchar(255);not null" validate:"required"`
	Time        int    `json:"time" gorm:"type:integer;not null" validate:"required"`
}
