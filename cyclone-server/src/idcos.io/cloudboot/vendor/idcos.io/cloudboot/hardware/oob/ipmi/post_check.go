package ipmi

import (
	"fmt"
	"strings"

	"idcos.io/cloudboot/hardware"
	"idcos.io/cloudboot/hardware/oob"
)

// checkNetwork 检查实际的OOB网络是否与预期的配置相符
func (w *worker) checkNetwork(sett *oob.NetworkSetting) (items []hardware.CheckingItem) {
	if sett == nil || sett.IPSrc == "" {
		return nil
	}
	network, err := w.Network()
	if err != nil {
		return []hardware.CheckingItem{
			{
				Title:   "Network",
				Matched: hardware.MatchedUnknown,
				Error:   err.Error(),
			},
		}
	}

	items = append(items,
		*hardware.NewCheckingHelper("BootOSIP Source", sett.IPSrc, strings.ToLower(network.IPSrc)).Matcher(hardware.ContainsMatch).Do(),
	)

	if sett.IPSrc == oob.Static {
		items = append(items,
			*hardware.NewCheckingHelper("BootOSIP", sett.StaticIP.IP, network.IP).Do(),
			*hardware.NewCheckingHelper("Netmask", sett.StaticIP.Netmask, network.Netmask).Do(),
			*hardware.NewCheckingHelper("Gateway", sett.StaticIP.Gateway, network.Gateway).Do(),
		)
	}
	return items
}

// checkUser 检查实际的OOB用户是否与预期配置相符
func (w *worker) checkUser(sett *oob.UserSetting) (items []hardware.CheckingItem) {
	if sett == nil {
		return nil
	}
	users, err := w.Users()
	if err != nil {
		return []hardware.CheckingItem{
			{
				Title:   "Users",
				Matched: hardware.MatchedUnknown,
				Error:   err.Error(),
			},
		}
	}

	for _, settUser := range []oob.UserSettingItem(*sett) {
		item := hardware.CheckingItem{
			Title:    "Create User",
			Expected: fmt.Sprintf("%s@%s", settUser.Username, oob.StringUserLevel(settUser.PrivilegeLevel)),
			Matched:  hardware.MatchedNO,
		}
		idx := w.findUserIndexByName(users, settUser.Username)
		// 检查目标用户是否存在
		if idx < 0 {
			item.Actual = "Missing"
			items = append(items, item)
			continue
		}
		// 检查目标用户权限级别
		if users[idx].Access == nil || users[idx].Access.PrivilegeLevel < 0 {
			item.Actual = fmt.Sprintf("%s@unknown", settUser.Username)
			items = append(items, item)
			continue
		}

		item.Actual = fmt.Sprintf("%s@%s", users[idx].Name, oob.StringUserLevel(users[idx].Access.PrivilegeLevel))
		if item.Actual == item.Expected {
			item.Matched = hardware.MatchedYES
		}
		items = append(items, item)
	}
	return items
}
