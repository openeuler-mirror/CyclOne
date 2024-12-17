package build

import (
	"strings"
)

const (
	// VersionNumber 版本号
	VersionNumber = "1.5.5"
)

// The value of variables come from `gb build -ldflags '-X "build.Date=xxxxx" -X "build.CommitID=xxxx"' `
var (
	// Date build time
	Date string
	// Branch current git branch
	Branch string
	// Commit git commit id
	Commit string
)

// Version 生成版本信息
func Version() string {
	var buf strings.Builder
	buf.WriteString(VersionNumber)

	if Date != "" {
		buf.WriteByte('\n')
		buf.WriteString("date: ")
		buf.WriteString(Date)
	}
	if Branch != "" {
		buf.WriteByte('\n')
		buf.WriteString("branch: ")
		buf.WriteString(Branch)
	}
	if Commit != "" {
		buf.WriteByte('\n')
		buf.WriteString("commit: ")
		buf.WriteString(Commit)
	}
	return buf.String()
}
