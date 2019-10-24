package controllers

import (
	"fmt"
	"tcmio/models"

	"github.com/astaxie/beego/orm"
)

func (this *MainController) ListIngredients() {
	var tars []models.Ingredient
	var tar models.Ingredient
	draw, _ := this.GetInt64("draw")
	offset, _ := this.GetInt64("start")
	limit, _ := this.GetInt64("length")
	fmt.Println(offset)
	fmt.Println(limit)
	_, err := tar.Query().Limit(limit, offset).All(&tars)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tars)

	total, _ := tar.Query().Count()

	var data DataList
	data.Draw = draw
	data.RecordsTotal = total
	data.RecordsFiltered = total
	data.Data = tars

	//this.responseMsg.SuccessMsg("", data)
	this.Data["json"] = data
	this.ServeJSON()

}

func (this *MainController) DetailIngredient() {
	var tar models.Ingredient
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

func (this *MainController) SearchIngredient() {
	var tar models.Ingredient
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

func (this *MainController) StructureSearchIngredient() {

	resultsInfo := make(map[string]interface{})
	method := this.Ctx.Input.Param(":method")
	fmt.Println(method)
	query := this.GetString("query")
	url := ""
	if method == "sim" {
		threshold := this.GetString("threshold")
		url = "select * from ligand where mol @ (" + threshold + ", 1.0, '" + query + "', 'Tanimoto')::bingo.sim;"
	} else if method == "sub" {
		url = "select *  from ligand where mol @ ('" + query + "','')::bingo.sub;"
	} else if method == "exact" {
		url = "select *  from ligand where mol @ ('" + query + "','')::bingo.exact;"
	} else {
		this.responseMsg.ErrorMsg("method not support", "")
		this.Data["json"] = this.responseMsg
		this.ServeJSON()
	}

	fmt.Println(url)
	o := orm.NewOrm()
	var hits []models.Ingredient
	_, err := o.Raw(url).QueryRows(&hits)
	if err != nil {
		fmt.Println(err)
		this.responseMsg.ErrorMsg("Not find", "")
		this.Data["json"] = this.responseMsg
		this.ServeJSON()
	}
	fmt.Println(hits)
	resultsInfo["ingredient"] = hits

	// get tcm from ingredient
	o1 := orm.NewOrm()
	var maps []orm.Params
	//url = "select arc.id,arc.title,arc.typeid,art.typename from go_archives as arc left join go_arctype as art on art.id=arc.typeid where arc.typeid=?"
	//url = "select * from ingredient full join tcm_ingredient on ingredient.id = tcm_ingredient.ingredient_id where id = ?"

	o1.Raw("select * from tcm_ingredient as ti join ingredient on ingredient.id=ti.ingredient_id where ingredient.id=?", 39308).Values(&maps)
	o1.Raw("select * from tcm_ingredient join tcm on tcm_ingredient.tcm_id=tcm.id where tcm_ingredient.tcm_id=1").Values(&maps)
	fmt.Println("sdsa")
	fmt.Println(maps)
	resultsInfo["source"] = maps

	// convert to json
	this.responseMsg.SuccessMsg("", resultsInfo)
	this.Data["json"] = this.responseMsg
	this.ServeJSON()
}
