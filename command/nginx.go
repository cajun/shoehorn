package command

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

// template is template that will be used for the nginx configuration.
// It will look for a file named "config/nginx.template.conf".  If the file
// doesn't exists then it will dup a default configuration in that location.
func templateContent() (content string) {
	bytes, err := ioutil.ReadFile("config/nginx.template.conf")
	if err != nil {
		content = NginxTemplate
		ioutil.WriteFile("config/nginx.template.conf", []byte(content), 0666)
	} else {
		content = string(bytes)
	}

	return
}

type NginxPort struct {
	PublicPort string
}

type NginxConf struct {
	App     string
	Ports   []NginxPort
	Domains string
	Allow   []string
	Pwd     string
}

// UpdateNginxConf will replace vars in the 'config/nginx.template.conf' file.
// The template engine being used is the one packaged with golang. It will
// replace the following vars.
//
// * {{.App}} -> shoehorn.cfg#App
// * {{.Pwd}} -> current working dir
// * {{.Domains}} -> shoehorn.cfg#Domains
// * {{.Allow}} -> shoehorn.cfg#Allow
// * {{.Ports}} -> public ports
//
// Then it will attempt to reload nginx configuration.
// NOTE: this is done via a sudo command.
func UpdateNginxConf() (err error) {
	t := template.New("Nginx Conf")
	t, err = t.Parse(templateContent())

	if err != nil {
		return
	}

	file, err := os.OpenFile("config/nginx.conf", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	t.Execute(file, settings())

	cmd := exec.Command("sudo", "nginx", "-s", "reload")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func settings() NginxConf {
	pwd, _ := os.Getwd()

	return NginxConf{
		App:     cfg.App,
		Ports:   ports(),
		Domains: strings.Join(cfg.Domain, " "),
		Allow:   cfg.Allow,
		Pwd:     pwd,
	}

}

func ports() (tcp []NginxPort) {
	runInstances("pulling ports", func(i int, id string) (err error) {
		port := publicPort(i)
		tcp = append(tcp, NginxPort{PublicPort: port.tcp})
		return
	})
	return
}

const NginxTemplate = `
upstream {{.App}} {
  {{ range .Ports }}
    server 127.0.0.1:{{.PublicPort}} fail_timeout=0;
  {{ end }}
}

server {
  listen       80;
  client_max_body_size 20M;
  server_name  {{.Domains}} ;

  keepalive_timeout 5;

  root {{.Pwd}}/public;

  access_log {{.Pwd}}/log/access.log;
  error_log {{.Pwd}}/log/error.log;

	if ($request_method !~ ^(GET|HEAD|PUT|POST|DELETE|OPTIONS)$ ){
		return 405;
	}

  location ~ "^/assets/(.*/)*.*-[0-9a-f]{32}.*/"  {
		gzip_static on;
		expires     max;
		add_header  Cache-Control public;
	}

	location / {
		try_files $uri/index.html $uri.html $uri @app;
		error_page 404              /404.html;
		error_page 422              /422.html;
		error_page 500 502 503 504  /500.html;
		error_page 403              /403.html;
	}

	location @app {
		proxy_pass http://{{.App}};
	}

	location = /favicon.ico {
		expires    max;
		add_header Cache-Control public;
	}

	location ~ \.php$ {
		deny  all;
	}

  {{range .Allow}}
    allow {{.}};
  {{end}}
  allow 127.0.0.1;

  deny all;
}
`
