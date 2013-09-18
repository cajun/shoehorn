package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// template is template that will be used for the nginx configuration.
// It will look for a file named "config/nginx.template.conf".  If the file
// doesn't exists then it will dup a default configuration in that location.
func template() string {
	bytes, err := ioutil.ReadFile("config/nginx.template.conf")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	content := string(bytes)
	return content
}

// UpdateNginxConf will replace vars in the 'config/nginx.template.conf' file.
// It will replace the following vars
//
// * %APP% -> shoehorn.cfg#App
// * %PWD% -> current working dir
// * %DOMAINS% -> shoehorn.cfg#Domains
// * %ALLOW% -> shoehorn.cfg#Allow
//
// Then it will attempt to reload nginx configuration.
// NOTE: this is done via a sudo command.
func UpdateNginxConf() (err error) {
	content := template()

	pwd, err := os.Getwd()

	if err != nil {
		return
	}

	content = strings.Replace(content, "%APP%", cfg.App, -1)
	content = strings.Replace(content, "%PWD%", pwd, -1)
	content = strings.Replace(content, "%PORT%", string(publicPort(0)[0].tcp), -1)
	content = strings.Replace(content, "%DOMAINS%", strings.Join(cfg.Domain, " "), -1)
	content = strings.Replace(content, "%ALLOW%", strings.Join(cfg.Allow, "\n"), -1)

	ioutil.WriteFile("config/nginx.conf", []byte(content), 0644)

	cmd := exec.Command("sudo", "nginx", "-s", "reload")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
