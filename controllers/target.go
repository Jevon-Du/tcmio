package controllers

import (
	"fmt"
	"tcmio/models"

	"github.com/astaxie/beego/orm"
)

func (this *MainController) ListTargets1() {
	var tars []models.Target
	var tar models.Target
	offset, _ := this.GetInt64("offset")
	limit, _ := this.GetInt64("limit")
	fmt.Println(offset)
	fmt.Println(limit)
	_, err := tar.Query().Limit(limit, offset).OrderBy("name").All(&tars, "name", "uniprot_id")

	if err != nil {
		fmt.Println(err)
		this.responseMsg.ErrorMsg("Not find", "")
		this.Data["json"] = this.responseMsg
		this.ServeJSON()
	}
	fmt.Println(tars)
	this.responseMsg.SuccessMsg("", tars)
	this.Data["json"] = this.responseMsg
	this.ServeJSON()
}

func (this *MainController) ListTargets() {
	var maps []orm.Params
	offset, _ := this.GetInt64("offset")
	limit, _ := this.GetInt64("limit")
	fmt.Println(offset)
	fmt.Println(limit)

	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Target)).RelatedSel()
	_, err := qs.Limit(limit, offset).Values(&maps, "name", "uniprot_id", "function", "chembl_id", "gene_name", "pdb", "protein_family", "mass", "length", "ec_number", "kegg")
	//_, err := o.Raw("SELECT name,journal FROM target LEFT JOIN doc ON target.ref_id = doc.id").Values(&tars)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(maps)
	this.responseMsg.SuccessMsg("", maps)
	this.Data["json"] = this.responseMsg
	this.ServeJSON()
}

func (this *MainController) DetailTarget() {
	var tar models.Target
	id := this.Ctx.Input.Param(":id")
	fmt.Println(id)
	err := tar.Query().Filter("id", id).One(&tar)
	if err != nil {
		fmt.Println(err)
		this.responseMsg.ErrorMsg("Not find", "")
		this.Data["json"] = this.responseMsg
		this.ServeJSON()
	}
	fmt.Println(tar)

	this.responseMsg.SuccessMsg("", tar)
	this.Data["json"] = this.responseMsg
	this.ServeJSON()

}
