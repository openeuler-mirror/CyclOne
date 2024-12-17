package strings

import (
	"strconv"
	stdstrings "strings"
)

// DOS2UNIX 将windows字符串转换为Unix字符串
func DOS2UNIX(src string) string {
	return stdstrings.Replace(src, "\r\n", "\n", -1)
}

// UNIX2DOS 将Unix字符串转换为windows字符串
func UNIX2DOS(src string) string {
	return stdstrings.Replace(src, "\n", "\r\n", -1)
}

// MultiLines2Slice 将多行文本换成切片
func MultiLines2Slice(multi string) []string {
	str := DOS2UNIX(stdstrings.TrimSpace(multi))
	str = stdstrings.Replace(str, " ", "\n", -1)
	str = stdstrings.Replace(str, ",", "\n", -1)
	str = stdstrings.Replace(str, ";", "\n", -1)
	return stdstrings.Split(str, "\n")
}

// MultiLines2Slice 将多行文本换成切片（排除空格）
func MultiLines2SliceWithSpace(multi string) []string {
	str := DOS2UNIX(stdstrings.TrimSpace(multi))
	str = stdstrings.Replace(str, ",", "\n", -1)
	str = stdstrings.Replace(str, ";", "\n", -1)
	return stdstrings.Split(str, "\n")
}

// Multi2UintSlice 将逗号分隔的多个数值，转换成[]uint
func Multi2UintSlice(multi string) []uint {
	vals := stdstrings.Split(multi, ",")
	valsUint := make([]uint, 0, len(vals))
	for _, val := range vals {
		valInt, err := strconv.Atoi(val)
		if err != nil {
			return nil
		}
		valsUint = append(valsUint, uint(valInt))
	}
	return valsUint
}
