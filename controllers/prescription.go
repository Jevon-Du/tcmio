package controllers

import (
	"fmt"
	"tcmio/models"
)

func (this *MainController) ListPrescriptions() {
	var tars []models.Prescription
	var tar models.Prescription
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

func (this *MainController) DetailPrescription() {
	var tar models.Prescription
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
