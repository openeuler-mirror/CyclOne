package bytes

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// Byte2GB 返回字节转化成GB后的数值
func Byte2GB(b Byte) (gb float64) {
	return float64(b) / float64(GB)
}

// Byte2GBRounding 返回字节转化成GB后的整数数值
func Byte2GBRounding(b Byte) (gb int) {
	return int(Byte2GB(b))
}

// Byte2MB 返回字节转化成MB后的数值
func Byte2MB(b Byte) (mb float64) {
	return float64(b) / float64(MB)
}

// Byte2MBRounding 返回字节转化成MB后的整数数值
func Byte2MBRounding(b Byte) (mb int) {
	return int(Byte2MB(b))
}

var (
	// ErrMalformedSizeStringValue 容量值的字符串格式错误
	ErrMalformedSizeStringValue = errors.New("malformed size string value")
	// ErrMalformedUnitStringValue 容量单位的字符串格式错误
	ErrMalformedUnitStringValue = errors.New("malformed unit string value")
)

var (
	numReg = regexp.MustCompile("^\\d+$")
)

// Parse2Byte 将字符串类型容量值和容量单位转化成字节
func Parse2Byte(value, unit string) (size Byte, err error) {
	if !numReg.MatchString(value) {
		return size, ErrMalformedSizeStringValue
	}

	unit = strings.ToUpper(unit)
	val, _ := strconv.Atoi(value)

	switch unit {
	case "B", "BYTE", "BYTES":
		return B * Byte(val), nil
	case "KB":
		return KB * Byte(val), nil
	case "MB":
		return MB * Byte(val), nil
	case "GB":
		return GB * Byte(val), nil
	case "TB":
		return TB * Byte(val), nil
	default:
		return size, ErrMalformedUnitStringValue
	}
}
