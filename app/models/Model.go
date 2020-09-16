package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Model interface {
	Insert() (Id uint64, err error)
	Update(Id uint) (Rows int64, err error)
	Delete(Ids string) (Rows int64, err error)
	GetAll(RequestData map[string]interface{}, Offset uint64, Limit uint8) (Data interface{}, err error)
	GetOne(Id uint) (Data interface{}, err error)
}
type Mysql struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt JsonTime
	UpdatedAt JsonTime
	DeletedAt *time.Time `sql:"index"`
}

type JsonTime struct {
	time.Time
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t JsonTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
func (t JsonTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *JsonTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JsonTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
