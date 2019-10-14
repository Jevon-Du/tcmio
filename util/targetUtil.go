package util

import (
	"fmt"
	"strings"
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

func DataProcess() {
	fmt.Println("Processing ...")

	var tar models.Target
	var tars []models.Target
	_, err := tar.Query().All(&tars)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(len(tars))

	for _, t := range tars {
		//fmt.Println(t)
		t.Kegg = strings.TrimRight(t.Kegg, ";")
		t.Drug = strings.TrimRight(t.Drug, ";")
		t.Pdb = strings.TrimRight(t.Pdb, ";")
		t.ChemblId = strings.TrimRight(t.ChemblId, ";")

		err := t.Update()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

}
