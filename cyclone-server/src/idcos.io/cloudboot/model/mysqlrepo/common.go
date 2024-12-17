package mysqlrepo

import (
	"fmt"
	rawStrings "strings"

	"github.com/jinzhu/gorm"
	"idcos.io/cloudboot/utils/strings"
)

// MultiQuery 多行模糊查询
func MultiQuery(db *gorm.DB, column string, multiStr string) *gorm.DB {
	if multiStr != "" {
		col := rawStrings.Split(column, ".")
		tableName, colName := "", ""
		if len(col) > 1 {
			tableName = col[0] + "."
			colName = col[1]
		} else {
			colName = col[0]
		}
		cs := strings.MultiLines2Slice(multiStr)
		var sb rawStrings.Builder
		for i, c := range cs {
			if i == 0 {
				sb.WriteString(fmt.Sprintf("%s`%s` LIKE '%s' ", tableName, colName, fmt.Sprintf("%%%s%%", c)))
			} else {
				sb.WriteString(fmt.Sprintf("OR %s`%s` LIKE '%s' ", tableName, colName, fmt.Sprintf("%%%s%%", c)))
			}
		}
		if len(cs) > 0 {
			fmt.Print(sb.String())
			db = db.Where(fmt.Sprintf("(%s)", sb.String()))
		}
	}
	return db
}

// MultiMatchQuery 多行精确查询
func MultiMatchQuery(db *gorm.DB, column string, multiStr string) *gorm.DB {
	if multiStr != "" {
		col := rawStrings.Split(column, ".")
		tableName, colName := "", ""
		if len(col) > 1 {
			tableName = col[0] + "."
			colName = col[1]
		} else {
			colName = col[0]
		}
		cs := strings.MultiLines2Slice(multiStr)
		db = db.Where(tableName+"`"+colName+"` IN (?)", cs)
	}
	return db
}

// MultiMatchWithSpaceQuery 多行精确查询(字段包含空格)
func MultiMatchWithSpaceQuery(db *gorm.DB, column string, multiStr string) *gorm.DB {
	if multiStr != "" {
		col := rawStrings.Split(column, ".")
		tableName, colName := "", ""
		if len(col) > 1 {
			tableName = col[0] + "."
			colName = col[1]
		} else {
			colName = col[0]
		}
		cs := strings.MultiLines2SliceWithSpace(multiStr)
		db = db.Where(tableName+"`"+colName+"` IN (?)", cs)
	}
	return db
}

// MultiEnumQuery 多值枚举查询，精确匹配
func MultiEnumQuery(db *gorm.DB, column string, multiStr string) *gorm.DB {
	if multiStr != "" {
		col := rawStrings.Split(column, ".")
		tableName, colName := "", ""
		if len(col) > 1 {
			tableName = col[0] + "."
			colName = col[1]
		} else {
			colName = col[0]
		}
		cs := strings.MultiLines2Slice(multiStr)
		db = db.Where(tableName+"`"+colName+"` IN (?)", cs)
	}
	return db
}

// MultiNumQuery 多数值查询:id=1,2,3,4
func MultiNumQuery(db *gorm.DB, column string, multiVal []uint) *gorm.DB {
	if len(multiVal) != 0 {
		col := rawStrings.Split(column, ".")
		tableName, colName := "", ""
		if len(col) > 1 {
			tableName = col[0] + "."
			colName = col[1]
		} else {
			colName = col[0]
		}
		db = db.Where(tableName+"`"+colName+"` IN (?)", multiVal)
	}
	return db
}

// ConcatColumnSliceStringQuery 多字符串切片
func ConcatColumnSliceStringQuery(db *gorm.DB, column string, multiVal []string) *gorm.DB {
	if len(multiVal) != 0 {
		db = db.Where("CONCAT("+column+") IN (?)", multiVal)
	}
	return db
}
// ConcatColumnStringQuery
func ConcatColumnStringQuery(db *gorm.DB, column string, multiStr string) *gorm.DB {
	if len(multiStr) != 0 {
		db = db.Where("CONCAT("+column+") = (?)", multiStr)
	}
	return db
}