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

func (this *MainController) Index() {
	this.TplName = "index.html"
}

func (this *MainController) Help() {
	this.TplName = "help.html"
}

func (this *MainController) Browse() {
	this.TplName = "browse.html"
}

func (this *MainController) Target() {
	this.TplName = "target.html"
}

func (this *MainController) Ligand() {
	this.TplName = "ligand.html"
}

func (this *MainController) Ingredient() {
	this.TplName = "ingredient.html"
}

func (this *MainController) TCM() {
	this.TplName = "tcm.html"
}

func (this *MainController) Prescription() {
	this.TplName = "prescription.html"
}

func (this *MainController) Structure() {
	this.TplName = "structure.html"
}

func (this *MainController) MOA() {
	this.TplName = "moa.html"
}
