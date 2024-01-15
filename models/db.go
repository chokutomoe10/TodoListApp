package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type List struct {
	ID          uint      `gorm:"primaryKey" json:"id" form:"id"`
	Title       string    `gorm:"type:varchar(100);not null" json:"title" form:"title"`
	Description string    `gorm:"type:varchar(100);not null" json:"description" form:"description"`
	SubLists    []SubList `gorm:"foreignKey:ListID" json:"subLists" form:"subLists"`
	Files       string    `gorm:"type:varchar(100)" json:"files" form:"files"`
}

type SubList struct {
	ID           uint   `gorm:"primaryKey" json:"id" form:"id"`
	ListID       uint   `gorm:"not null" json:"list_id" form:"list_id"`
	Title        string `gorm:"type:varchar(100);not null" json:"subList_title" form:"subList_title"`
	Description  string `gorm:"type:varchar(100);not null" json:"subList_description" form:"subList_description"`
	SubListFiles string `gorm:"type:varchar(100)" json:"subList_files" form:"subList_files"`
}

var DB *gorm.DB

func DatabaseConnect() {
	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=pgsql089 dbname=todo_list port=5432"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&List{}, &SubList{})

	DB = db
}
