package ipmi

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strconv"
	"strings"

	"idcos.io/cloudboot/hardware/oob"
)

// findUserByName 根据带外用户名查找带外用户。
// 若指定名称的用户不存在，则返回nil。
func (worker *worker) findUserByName(name string) (*oob.User, error) {
	users, err := worker.Users()
	if err != nil {
		return nil, err
	}
	for i := range users {
		if users[i].Name == name {
			return &users[i], nil
		}
	}
	return nil, nil
}

// findUserIndexByName 返回指定用户名的用户在切片中的索引号
func (worker *worker) findUserIndexByName(users []oob.User, name string) (index int) {
	for i := range users {
		if users[i].Name == name {
			return i
		}
	}
	return -1
}

const (
	defaultUserID = 2
)

// newUserID 返回新建用户的ID
func (worker *worker) newUserID() (id int, err error) {
	users, err := worker.Users()
	if err != nil {
		return 0, err
	}
	if len(users) <= 0 {
		return defaultUserID, nil
	}

	for i := range users {
		name := strings.TrimSpace(users[i].Name)
		if name == "" || strings.Contains(name, "(Empty") { // 系统预置用户
			return users[i].ID, nil
		}
	}
	return users[len(users)-1].ID + 1, nil
}

// userAccess 返回通道下指定用户的Access信息
func (worker *worker) userAccess(channel, userID int) (*oob.UserAccess, error) {
	output, err := worker.Base.ExecByShell(tool, "channel", "getaccess", strconv.Itoa(channel), strconv.Itoa(userID))
	if err != nil {
		return nil, err
	}
	var access oob.UserAccess
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "User Name") {
			access.UserName = worker.extractValue(line, ":")
		} else if strings.HasPrefix(line, "Access Available") {
			access.AccessAvailable = worker.extractValue(line, ":")
		} else if strings.HasPrefix(line, "Link Authentication") {
			access.LinkAuthentication = worker.extractValue(line, ":")
		} else if strings.HasPrefix(line, "IPMI Messaging") {
			access.IPMIMessaging = worker.extractValue(line, ":")
		} else if strings.HasPrefix(line, "Privilege Level") {
			access.PrivilegeLevel = oob.IntUserLevel(worker.extractValue(line, ":"))
		}
	}
	return &access, nil
}

// extractValue 截取kv对中v的内容。假设，kv内容为"name : voidint"，那么将返回"voidint"。
func (worker *worker) extractValue(kv, sep string) (value string) {
	if !strings.Contains(kv, sep) {
		return kv
	}
	return strings.TrimSpace(strings.SplitN(kv, sep, 2)[1])
}

var (
	numReg = regexp.MustCompile("^\\d+$")
)

// ParseUsers 返回stdout内容中的带外用户信息
func ParseUsers(output []byte) (items []oob.User, err error) {
	var started bool
	rd := bufio.NewReader(bytes.NewBuffer(output))
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "ID") && strings.Contains(line, "Name") {
			started = true
			continue
		}
		if !started {
			continue
		}

		arr := strings.Fields(line)
		if len(arr) < 6 || !numReg.MatchString(arr[0]) || arr[1] == "true" || arr[1] == "false" {
			continue
		}

		id, _ := strconv.Atoi(arr[0])
		items = append(items, oob.User{
			ID:   id,
			Name: arr[1],
		})
	}
	return items, nil
}
