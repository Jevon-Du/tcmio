package util

import (
	"fmt"
	"tcmio/models"

	"github.com/astaxie/beego/orm"
)

func GetTargets() {
	var tars []models.Target
	o := orm.NewOrm()
	_, err := o.Raw("select * from target limit 5 offset 5").QueryRows(&tars)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tars)
}
