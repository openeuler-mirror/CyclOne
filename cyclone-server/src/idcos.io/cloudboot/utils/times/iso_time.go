package times

import (
	"fmt"
	"strings"
	"time"
)

// ISOTime 标准库Time的别名类型。用于格式化为ISO8601标准的时间（形如'2015-09-09 09:09:09'）。
type ISOTime time.Time

const (
	// DateTimeLayout 日期时间转换模式
	DateTimeLayout = "2006-01-02 15:04:05"
	// DateLayout 日期转换模式
	DateLayout = "2006-01-02"
	// DateLayout2
	DateLayout2 = "01-02-06"
)

// MarshalJSON 序列化，实现json.Marshaller接口
func (t ISOTime) MarshalJSON() ([]byte, error) {
	if time.Time(t).IsZero() {
		return []byte(fmt.Sprintf("%q", "")), nil
	}
	// 注意：序列化成的时间字符串必须包含在双引号当中
	return []byte(fmt.Sprintf("%q", time.Time(t).Format(DateTimeLayout))), nil
}

// UnmarshalJSON 反序列化,实现json.Unmarshaler接口
func (t *ISOTime) UnmarshalJSON(b []byte) error {
	sTime := fmt.Sprintf("%s", string(b))

	if strings.HasPrefix(sTime, "\"") {
		sTime = sTime[1:]
	}

	if strings.HasSuffix(sTime, "\"") {
		sTime = sTime[0 : len(sTime)-1]
	}

	tmpT, err := time.Parse(DateTimeLayout, sTime)
	if err != nil {
		return err
	}

	*t = ISOTime(tmpT)
	return nil
}

// MarshalYAML 序列化成YAML
func (t ISOTime) MarshalYAML() (interface{}, error) {
	return time.Time(t).Format(DateTimeLayout), nil
}

// UnmarshalYAML YAML反序列化
// func (t ISOTime) UnmarshalYAML(unmarshal func(interface{}) error) error {
// 	return nil
// }

//String 时间转换
func (t ISOTime) String() string {
	return time.Time(t).Format(DateTimeLayout)
}

// ToDateStr 格式化成日期
func (t ISOTime) ToDateStr() string {
	if time.Time(t).IsZero() {
		return ""
	}
	return time.Time(t).Format(DateLayout)
}

// ToDateStr 格式化成日期
func (t ISOTime) ToTimeStr() string {
	if time.Time(t).IsZero() {
		return ""
	}
	return time.Time(t).Format(DateTimeLayout)
}

// UnixSecToISOTime UNIX秒转化成ISOTime类型
func UnixSecToISOTime(unixSecond int64) ISOTime {
	return ISOTime(time.Unix(unixSecond, 0))
}
