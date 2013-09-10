package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func UpdateNginxConf() (err error) {
	bytes, err := ioutil.ReadFile("conf/nginx.template")
	if err != nil {
		return
	}

	content := string(bytes)

	fmt.Println(content)
	strings.Replace(content, "%APP%", cfg.App, -1)
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	strings.Replace(content, "%PWD%", pwd, -1)
	strings.Replace(content, "%PORT%", string(PublicPort()), -1)
	strings.Replace(content, "%DOMAINS%", strings.Join(cfg.Domain, " "), -1)
	strings.Replace(content, "%ALLOW%", strings.Join(cfg.Allow, "\n"), -1)

	fmt.Println("NGINX Configuration")
	fmt.Println(content)

	ioutil.WriteFile("conf/nginx.conf", []byte(content), 0644)
	return
}
