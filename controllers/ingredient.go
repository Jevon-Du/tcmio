package controllers

import (
	"fmt"
	"tcmio/models"
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

	this.responseMsg.SuccessMsg("", data)
	this.Data["json"] = this.responseMsg
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
