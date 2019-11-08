package model

import (
	"database/sql/driver"
	"errors"
	"time"
)

// 通用字段
type BaseModel struct {
	ID        uint      `json:"id" gorm:"primary_key" gorm:"AUTO_INCREMENT"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 挑战状态
type Status string

const (
	FAILD   Status = "失败"
	SUCCESS Status = "挑战成功"
	DOING   Status = "进行中"
)

func (u *Status) Scan(value interface{}) error { *u = Status(value.([]byte)); return nil }
func (u Status) Value() (driver.Value, error)  { return string(u), nil }


type Pagination struct {
	PageNum int `json:"page_num"`
	PageSize int `json:"page_size"`
}

// 自定义 time format
type JsonTime time.Time

const (
	timeFormat = "2006-01-02T15:04:05"
)

func (t *JsonTime) UnmarshalJSON(data []byte) (err error) {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	tmp, err := time.Parse(`"`+timeFormat+`"`, string(data))
	*t = JsonTime(tmp)
	return err
}

func (t JsonTime) MarshalJSON() ([]byte, error) {
	if y := time.Time(t).Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}
	b := make([]byte, 0, len(timeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}