package command

import (
	"fmt"
	"github.com/cajun/shoehorn/logger"
	"os"
	"os/exec"
	"strings"

	//"io/ioutil"
	"text/template"
)

type SystemdConf struct {
	Exe string
	App string
	Pwd string
}

func init() {
	available.addInfo("systemd", Executor{
		description: "install systemd files",
		run:         InstallSystemd})
}

func exe() string {
	out, err := exec.Command("which shoehorn").Output()
	if err != nil {
		logger.Log(err.Error())
		out = []byte("/usr/bin/shoehorn")
	}
	return string(out)
}

func InstallSystemd(args ...string) {
	t := template.New("Systemd Conf")
	t, err := t.Parse(systemdConf)

	if err != nil {
		logger.Log(err.Error())
	}

	wd, _ := os.Getwd()
	path := strings.Split(wd, "/")

	name := path[len(path)-1:][0]
	shoe := exe()
	logger.Log(shoe + "\n")
	fileName := fmt.Sprintf("/usr/lib/systemd/system/%s.service", name)

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	conf := SystemdConf{
		Exe: shoe,
		App: name,
		Pwd: wd}

	t.Execute(file, conf)

}

const systemdConf = `
[Unit]
Description={{.App}}
After=syslog.target nginx.service

[Service]
Type=oneshot
RemainAfterExit=yes
ExecStart={{.Exe}} --root {{.Pwd}} start
ExecStop={{.Exe}} --root {{.Pwd}} stop

[Install]
WantedBy=multi-user.target
`
