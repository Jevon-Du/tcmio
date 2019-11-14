package controllers

import (
	"fmt"
	"tcmio/models"
)

func (this *MainController) ListTargets1() {
	var tars []models.Target
	var tar models.Target
	offset, _ := this.GetInt64("start")
	limit, _ := this.GetInt64("length")
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

	draw, _ := this.GetInt64("draw")
	offset, _ := this.GetInt64("start")
	limit, _ := this.GetInt64("length")
	fmt.Println(offset)
	fmt.Println(limit)
	fmt.Println(this.GetString("order"))
	fmt.Println(this.Ctx.Input)

	/*var maps []orm.Params
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Target)).RelatedSel()
	_, err := qs.Limit(limit, offset).Values(&maps, "name", "uniprot_id", "function", "chembl_id", "gene_name", "pdb", "protein_family", "mass", "length", "ec_number", "kegg")
	//_, err := o.Raw("SELECT name,journal FROM target LEFT JOIN doc ON target.ref_id = doc.id").Values(&tars)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(maps)
	*/
	var tar models.Target
	var tars []models.Target
	_, err := tar.Query().Limit(limit, offset).OrderBy("id").All(&tars)
	if err != nil {
		fmt.Println(err)
	}

	total, _ := tar.Query().Count()

	var data DataList
	data.Draw = draw
	data.RecordsTotal = total
	data.RecordsFiltered = total
	data.Data = tars

	//this.responseMsg.SuccessMsg("", data)
	this.Data["json"] = data
	this.ServeJSON()
	this.TplName = "browse.html"

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
