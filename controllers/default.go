package controllers

import (
	"tcmio/util"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
	responseMsg util.ResponseMsg
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
