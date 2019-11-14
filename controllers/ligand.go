package controllers

import (
	"fmt"
	"tcmio/models"

	"github.com/astaxie/beego/orm"
)

func (this *MainController) ListLigands() {
	fmt.Println("List ligands")
	var tars []models.Ligand
	var tar models.Ligand
	draw, _ := this.GetInt64("draw")
	offset, _ := this.GetInt64("start")
	limit, _ := this.GetInt64("length")
	fmt.Println(offset)
	fmt.Println(limit)
	_, err := tar.Query().Limit(limit, offset).OrderBy("id").All(&tars)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(tars)

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

func (this *MainController) DetailLigand() {
	fmt.Println("DetailLigand")
	var tar models.Ligand
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

// structure search ligands
func (this *MainController) StructureSearchLigand() {
	fmt.Println("StructureSearchLigand")

	draw, _ := this.GetInt64("draw")
	method := this.GetString("method")
	fmt.Println(method)
	limit := this.GetString("length")
	offset := this.GetString("start")
	if limit == "" {
		limit = "10"
	}
	if offset == "" {
		offset = "0"
	}
	query := this.GetString("query")
	url := ""
	if method == "Similarity" {
		threshold := this.GetString("threshold")
		url = "select * from ligand where mol @ (" + threshold + ", 1.0, '" + query + "', 'Tanimoto')::bingo.sim" + " limit " + limit + " offset " + offset + ";"
	} else if method == "Substructure" {
		url = "select *  from ligand where mol @ ('" + query + "','')::bingo.sub" + " limit " + limit + " offset " + offset + ";"
	} else if method == "Fullstructure" {
		url = "select *  from ligand where mol @ ('" + query + "','')::bingo.exact" + " limit " + limit + " offset " + offset + ";"
	} else {
		this.responseMsg.ErrorMsg("method not support", "")
		this.Data["json"] = this.responseMsg
		this.ServeJSON()
	}

	fmt.Println(url)
	o := orm.NewOrm()
	var hits []models.Ligand
	_, err := o.Raw(url).QueryRows(&hits)
	if err != nil {
		fmt.Println(err)
		this.responseMsg.ErrorMsg("Not find", "")
		this.Data["json"] = this.responseMsg
		this.ServeJSON()
	}
	//fmt.Println(hits)

	var data DataList
	data.Draw = draw
	data.RecordsTotal = 100
	data.RecordsFiltered = 100
	data.Data = hits

	this.Data["json"] = data
	this.ServeJSON()
}

func (this *MainController) AnalyzeStructureSearchLigand() {
	fmt.Println("AnalyzeStructureSearchLigand")
	//resultsInfo := make(map[string]interface{})
	draw, _ := this.GetInt64("draw")
	method := this.GetString("method")
	fmt.Println(method)
	query := this.GetString("query")
	url := ""
	if method == "Similarity" {
		threshold := this.GetString("threshold")
		url = "select * from ligand where mol @ (" + threshold + ", 1.0, '" + query + "', 'Tanimoto')::bingo.sim;"
	} else if method == "Substructure" {
		url = "select *  from ligand where mol @ ('" + query + "','')::bingo.sub;"
	} else if method == "Fullstructure" {
		url = "select *  from ligand where mol @ ('" + query + "','')::bingo.exact;"
	} else {
		this.responseMsg.ErrorMsg("method not support", "")
		this.Data["json"] = this.responseMsg
		this.ServeJSON()
	}

	fmt.Println(url)
	o := orm.NewOrm()
	var hits []models.Ligand
	//_, err := o.Raw("select * from ligand where mol @ (0.5, 1.0, 'Fc1ccccc1Cn2ccc(NC(=O)c3ccc(COc4ccccc4Cl)cc3)n2', 'Tanimoto')::bingo.sim limit 10 offset 0;").QueryRows(&hits)

	//_, err := o.Raw("select * from ligand where mol @ (0.5, 1.0, 'C1=CC=CC=C1', 'Tanimoto')::bingo.sim limit 10 offset 0;").QueryRows(&hits)
	_, err := o.Raw(url).QueryRows(&hits)
	if err != nil {
		fmt.Println(err)
		this.responseMsg.ErrorMsg("Not find", "")
		this.Data["json"] = this.responseMsg
		this.ServeJSON()
	}

	// get activity from ligands

	//var tl models.Target_ligand
	var tlss []models.Target_ligand
	for _, s := range hits {
		var tl models.Target_ligand
		var tls []models.Target_ligand
		//_, err = tl.Query().Filter("mol_chembl_id", s.ChemblId).RelatedSel().All(&tls)
		_, err = tl.Query().Filter("mol_chembl_id", s.ChemblId).All(&tls)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//fmt.Println(tls)
		for _, t := range tls {
			tlss = append(tlss, t)
		}
	}
	//resultsInfo["activity"] = tlss

	// get ingredients from ligands
	// cancer
	/*var ingres []models.Ingredient
	for _, s := range hits {
		fmt.Println(s.ChemblId)
		var ingre models.Ingredient
		err = ingre.Query().Filter("inchikey", s.Inchikey).One(&ingre)
		if err != nil {
			continue
		}
		ingres = append(ingres, ingre)
	}
	resultsInfo["ingredient"] = ingres
	*/
	// convert to json
	//this.responseMsg.SuccessMsg("", resultsInfo)

	var data DataList

	data.Draw = draw
	data.RecordsTotal = int64(len(hits))
	data.RecordsFiltered = int64(len(hits))
	data.Data = tlss

	this.Data["json"] = data
	this.ServeJSON()
}
