package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func UpdateNginxConf() (err error) {
	fmt.Println("In conf")
	bytes, err := ioutil.ReadFile("config/nginx.template.conf")
	if err != nil {
		fmt.Println(err)
		return
	}

	content := string(bytes)

	fmt.Println(content)
	content = strings.Replace(content, "%APP%", cfg.App, -1)
	pwd, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
		return
	}

	content = strings.Replace(content, "%PWD%", pwd, -1)
	content = strings.Replace(content, "%PORT%", string(publicPort()), -1)
	content = strings.Replace(content, "%DOMAINS%", strings.Join(cfg.Domain, " "), -1)
	content = strings.Replace(content, "%ALLOW%", strings.Join(cfg.Allow, "\n"), -1)

	fmt.Println("NGINX Configuration")
	fmt.Println(content)

	ioutil.WriteFile("config/nginx.conf", []byte(content), 0644)

	cmd := exec.Command("sudo", "nginx", "-s", "reload")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
