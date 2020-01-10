package database

type Service struct {
	ID       uint   `gorm:"primary_key"`
	Name     string `gorm:"type:varchar(255)"`
	Host     string `json:"host,omitempty";gorm:"type:varchar(255)"`
	Port     int    `json:"port,omitempty";gorm:"type:integer"`
	Priority int    `json:"priority,omitempty";gorm:"type:integer"`
	Weight   int    `json:"weight,omitempty";gorm:"type:integer"`
	Text     string `json:"text,omitempty";gorm:"type:varchar(255)"`
	Mail     bool   `json:"mail,omitempty";gorm:"type:boolean"`
	TTL      uint32 `json:"mail,omitempty";gorm:"type:integer"`
}