package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/urfave/cli"

	"idcos.io/cloudboot/agent"
	"idcos.io/cloudboot/build"
	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/logger"
)

func main() {
	app := cli.NewApp()
	app.Name = "cloudboot-agent"
	app.Description = "cloudboot agent"
	app.Version = build.Version()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-level",
			Usage: "log level, optional values: debug|info|warn|error",
			Value: "debug",
		},
		cli.StringFlag{
			Name:  "log-dir",
			Value: "/var/log/cloudboot",
			Usage: "log file directory",
		},
	}

	app.Action = func(c *cli.Context) error {
		return runAgent(c)
	}

	// TODO 以下实现待优化
	//catch signal, then do sth like recycle resources, kill child process, etc.
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()
	go func() {
		app.Run(os.Args)
		done <- true
	}()

	<-done
	_ = doBeforeExit()
}

func runAgent(ctx *cli.Context) error {
	log := logger.NewBeeLogger(&config.Logger{
		Level:          ctx.String("log-level"),
		LogFile:        filepath.Join(ctx.String("log-dir"), "agent.log"),
		ConsoleEnabled: false,
	})

	agent, err := agent.New(log)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if err = agent.Run(); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func doBeforeExit() error {
	return exec.Command("pkill", "hw-server").Run()
}
