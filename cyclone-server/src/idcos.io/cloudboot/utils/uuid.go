package utils

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// UUID 返回UUID字符串
func UUID() string {
	return fmt.Sprintf("%x", uuid.Must(uuid.NewV4()).Bytes())
}
