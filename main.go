package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/icattlecoder/godaemon"
	"github.com/robfig/cron"
	"gopkg.in/macaron.v1"
)

type serverInfo struct {
	Server   string
	Port     string
	Password string
	Method   string
	Qrstring string
}

var servers []serverInfo

const (
	shadowsocksEntry0 = "http://www.ishadowsocks.org/"
)

func main() {

	m := macaron.Classic()

	//	默认为静态 dir 为 public 目录
	m.Use(macaron.Static("public"))
	//	设置模板目录,模板后缀
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory:  "templates",
		Extensions: []string{".tmpl", ".html"},
	}))

	m.Get("/", func(ctx *macaron.Context) {

		ctx.Data["servers"] = servers

		ctx.HTML(200, "index")

	})

	{
		//强制运行一次
		getServerSlice()

		c := cron.New()

		// c.AddFunc("0/5 * * * * ?", getServerSlice)
		c.AddFunc("0 10 * * * ?", getServerSlice)
		// c.AddFunc("@hourly", getServerSlice)

		c.Start()
	}

	m.Run(8080)

}

func getServerSlice() {

	servers = nil

	doc, e := goquery.NewDocument(shadowsocksEntry0)
	if e != nil {
		fmt.Println(e)
	}
	c := doc.Find("#free").Find("div.col-sm-4.text-center")
	c.Each(func(i int, content *goquery.Selection) {
		sectionInfo := []string{}
		// serverInfo := serverDetail{}
		content.Find("h4").Each(func(i int, content *goquery.Selection) {
			if i < 4 {
				sectionInfo = append(sectionInfo, strings.Split(content.Text(), ":")[1])
			}
		})
		// 如果密码为空, 不加入到 servers slice 中

		// ss: //method:password@hostname:port
		if sectionInfo[2] != "" {
			si := serverInfo{
				Server:   strings.ToUpper(sectionInfo[0]),
				Port:     sectionInfo[1],
				Password: sectionInfo[2],
				Method:   strings.ToUpper(sectionInfo[3]),
				Qrstring: base64.StdEncoding.EncodeToString([]byte(strings.ToUpper(sectionInfo[3]) + ":" + sectionInfo[2] + "@" + strings.ToUpper(sectionInfo[0]) + ":" + sectionInfo[1])),
			}
			servers = append(servers, si)
		}
	})
}
