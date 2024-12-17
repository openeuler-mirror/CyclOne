package agent

import "os/exec"

const (
	preInstallScript  = "/tmp/preInstall.cmd"
	postInstallScript = "/tmp/postInstall.cmd"
)

//runPreInstall run pre install script
func (agent *Agent) runPreInstall() {
	if agent.PreInstallURL == "" {
		agent.log.Info("Skip the preinstall script")
		return
	}

	agent.log.Debugf("Start fetching(%s) preinstall script", agent.PreInstallURL)
	if err := agent.wgetO(agent.PreInstallURL, preInstallScript, 0755); err != nil {
		return
	}

	agent.log.Debug("Start running pre-install script")
	cmd := exec.Command(preInstallScript) // 不能使用bash执行文件方式，因为目标脚本可能并不是shell。
	output, err := cmd.Output()
	if err != nil {
		agent.log.Errorf("Exec %q err: %s\noutput:\n%s", preInstallScript, err, output)
		return
	}
	agent.log.Infof("%s ==>\n%s", preInstallScript, output)
}

//runPostInstall run post install script
func (agent *Agent) runPostInstall() {
	if agent.PostInstallURL == "" {
		agent.log.Info("Skip the postinstall script")
		return
	}

	agent.log.Debugf("Start fetching(%s) postinstall script", agent.PostInstallURL)
	if err := agent.wgetO(agent.PostInstallURL, postInstallScript, 0755); err != nil {
		return
	}

	agent.log.Debug("Start running post-install script")
	cmd := exec.Command(postInstallScript) // 不能使用bash执行文件方式，因为目标脚本可能并不是shell。
	output, err := cmd.Output()
	if err != nil {
		agent.log.Errorf("Exec %q err: %s\noutput:\n%s", postInstallScript, err, output)
		return
	}
	agent.log.Infof("%s ==>\n%s", postInstallScript, output)
}
