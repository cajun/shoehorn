package command

import (
	"fmt"
	"github.com/cajun/shoehorn/logger"
	"os"
	"strings"

	//"io/ioutil"
	"text/template"
)

type UpstartConf struct {
	App string
	Exe string
	Pwd string
}

func init() {
	available.addInfo("upstart", Executor{
		description: "install upstart files",
		run:         InstallUpstart})
}

func InstallUpstart(args ...string) {
	t := template.New("Upstart Conf")
	t, err := t.Parse(upstartConf)

	if err != nil {
		logger.Log(err.Error() + "\n")
	}

	wd, _ := os.Getwd()
	path := strings.Split(wd, "/")

	name := path[len(path)-1:][0]
	fileName := fmt.Sprintf("/etc/init/%s.conf", name)
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)

	conf := UpstartConf{App: name, Exe: exe(), Pwd: wd}

	t.Execute(file, conf)
}

const upstartConf = `
description ".{{App}} containers"

start on started nginx

exec {{.Exe}} -wait true -root {{.Pwd}} start
post-stop {{.Exe}} -root {{.Pwd}} stop
`
