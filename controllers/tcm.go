package controllers

import (
	"fmt"
	"tcmio/models"
)

func (this *MainController) ListTCMs() {
	var tars []models.TCM
	var tar models.TCM
	offset, _ := this.GetInt64("offset")
	limit, _ := this.GetInt64("limit")
	fmt.Println(offset)
	fmt.Println(limit)
	_, err := tar.Query().Limit(limit, offset).All(&tars)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tars)

	this.responseMsg.SuccessMsg("", tars)
	this.Data["json"] = this.responseMsg
	this.ServeJSON()

}

func (this *MainController) DetailTCM() {
	var tar models.TCM
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
