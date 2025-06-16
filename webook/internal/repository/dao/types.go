package dao

import (
	"database/sql"
	"encoding/json"
)

type UserEntity struct {
	ID       int64          `gorm:"primaryKey;autoIncrement"`
	Email    sql.NullString `gorm:"uniqueIndex;type:varchar(50)"`
	Password string         `gorm:"type:varchar(200)"`
	Phone    sql.NullString `gorm:"uniqueIndex;type:varchar(20)"`
	Name     string         `gorm:"type:varchar(10)"`
	Birthday string         `gorm:"type:date;default:NULL"`
	AboutMe  string         `gorm:"varchar(50)"`
	Ctime    int64
	Utime    int64
}

type Oauth2BindingEntity struct {
	ID     int64      `gorm:"primary key;autoIncrement"`
	UserID int64      `gorm:"index"`
	User   UserEntity `gorm:"foreignKey:UserID;references:ID"`

	Provider        sql.NullString `gorm:"type:varchar(20);uniqueIndex:idx_provider_external_id"`
	ExternalID      sql.NullString `gorm:"type:varchar(50);uniqueIndex:idx_provider_external_id"`
	ProviderUnionID sql.NullString

	AccessToken  sql.NullString
	RefreshToken sql.NullString

	Ctime     int64
	Utime     int64
	ExtraInfo json.RawMessage `gorm:"type:json"`
}
