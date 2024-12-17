package ping

import "os/exec"

// 测试ping联通行性
func PingTest(ip string) error {
	return exec.Command("ping", "-c3", "-i0.5", "-W1", ip).Run()
}

// 测试ping6联通行性
func Ping6Test(ipv6 string) error {
	return exec.Command("ping6", "-c3", "-i0.5", "-W1", ipv6).Run()
}
