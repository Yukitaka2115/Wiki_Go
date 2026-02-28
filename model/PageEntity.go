package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Page struct {
	ID           int        `gorm:"primary_key"`
	Title        string     `json:"Title"`
	Brief        string     `json:"Brief"`
	Background   string     `json:"Background"`
	History      mapHandler `json:"History" gorm:"type:json"`
	Mainchara    Chara
	MaincharaStr string
	Relatives    Chara
	RelativesStr string
	Team         Group
	TeamStr      string
}

type Group struct {
	Chara string
	Grade string
	Group string
}

type Chara struct {
	Chara string
	Cast  string
	Info  string
}

type mapHandler map[string]interface{}

func (mh mapHandler) Value() (driver.Value, error) {
	str, err := json.Marshal(mh)
	if err != nil {
		return nil, err
	}
	return string(str), nil
}

func (mh *mapHandler) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte, value: %v", value)
	}

	return json.Unmarshal(b, &mh)
}
