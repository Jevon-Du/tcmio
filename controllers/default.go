package controllers

import (
	"tcmio/util"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
	responseMsg util.ResponseMsg
}

type DataList struct {
	Draw            int64       `json:"draw"`
	RecordsTotal    int64       `json:"recordsTotal"`
	RecordsFiltered int64       `json:"recordsFiltered"`
	Data            interface{} `json:"data"`
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func (this *MainController) Index() {
	this.TplName = "index.html"
}

func (this *MainController) Help() {
	this.TplName = "help.html"
}

func (this *MainController) Browse() {
	this.TplName = "browse.html"
}

func (this *MainController) Structure() {
	this.TplName = "structure.html"
}

func (this *MainController) MOA() {
	this.TplName = "moa.html"
}
