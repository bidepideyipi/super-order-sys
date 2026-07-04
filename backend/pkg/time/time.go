package time

import (
	"database/sql/driver"
	"time"
)

const DateTimeFormat = "2006-01-02 15:04"

// CustomTime 自定义时间类型，JSON序列化时使用指定格式
type CustomTime struct {
	time.Time
}

// MarshalJSON 实现 JSON 序列化
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.IsZero() {
		return []byte(`""`), nil
	}
	formatted := ct.Format(DateTimeFormat)
	return []byte(`"` + formatted + `"`), nil
}

// UnmarshalJSON 实现 JSON 反序列化
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == `""` {
		ct.Time = time.Time{}
		return nil
	}
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	t, err := time.Parse(DateTimeFormat, str)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

// Value 实现 driver.Valuer 接口，用于数据库写入
func (ct CustomTime) Value() (driver.Value, error) {
	return ct.Time, nil
}

// Scan 实现 sql.Scanner 接口，用于数据库读取
func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		ct.Time = time.Time{}
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return nil
	}
	ct.Time = t
	return nil
}

// NewCustomTime 创建 CustomTime
func NewCustomTime(t time.Time) CustomTime {
	return CustomTime{Time: t}
}
