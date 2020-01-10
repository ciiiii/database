package database

type Service struct {
	ID       uint64 `gorm:"type:bigint(20) unsigned auto_increment;not null;primary_key"`
	Name     string `gorm:"type:varchar(255)"`
	Host     string `gorm:"type:varchar(15)"`
	Port     int    `gorm:"type:integer"`
	Priority int    `gorm:"type:integer"`
	Weight   int    `gorm:"type:integer"`
	Text     string `gorm:"type:varchar(255)"`
	Mail     bool   `gorm:"type:boolean"`
	TTL      uint32 `gorm:"type:integer"`
}
