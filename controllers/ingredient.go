package controllers

import (
	"fmt"
	"sortutil"
	"strconv"
	"tcmio/models"

	"github.com/astaxie/beego/orm"
)

func (this *MainController) ListIngredients() {
	fmt.Println("List ingredients")
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

func (this *MainController) DetailIngredient() {
	fmt.Println("Detail ingredients")
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
	fmt.Println("StructureSearchIngredient")

	draw, _ := this.GetInt64("draw")
	method := this.GetString("method")

	limit := this.GetString("length")
	offset := this.GetString("start")
	query := this.GetString("query")

	fmt.Println(method)
	fmt.Println(limit)
	fmt.Println(offset)
	fmt.Println(query)
	if limit == "" {
		limit = "10"
	}
	if offset == "" {
		offset = "0"
	}

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
	var hits []models.Ingredient
	_, err := o.Raw(url).QueryRows(&hits)
	if err != nil {
		fmt.Println(err)
		this.responseMsg.ErrorMsg("Not find", "")
		this.Data["json"] = this.responseMsg
		this.ServeJSON()
	}

	var data DataList
	data.Draw = draw
	data.RecordsTotal = 100
	data.RecordsFiltered = 100
	data.Data = hits

	this.Data["json"] = data
	this.ServeJSON()
}

func (this *MainController) AnalyzeStructureSearchIngredient() {
	fmt.Println("AnalyzeStructureSearchIngredient")
	draw, _ := this.GetInt64("draw")
	method := this.GetString("method")
	query := this.GetString("query")
	fmt.Println(method)
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
	var hits []models.Ingredient
	_, err := o.Raw(url).QueryRows(&hits)
	if err != nil {
		fmt.Println(err)
		this.responseMsg.ErrorMsg("Not find", "")
		this.Data["json"] = this.responseMsg
		this.ServeJSON()
	}
	fmt.Println(len(hits))

	// get tcm from ingredient
	var iss []IngredientSource
	for _, s := range hits {
		fmt.Println(s.Id)
		fmt.Println(s.Synonyms)
		var ing models.TCM_Ingredient
		var ings []models.TCM_Ingredient
		_, err = ing.Query().Filter("ingredient_id", s.Id).All(&ings)
		if err != nil {
			continue
		}

		for _, t := range ings {
			var tcm models.TCM
			err := tcm.Query().Filter("id", t.TcmId).One(&tcm)
			if err != nil {
				continue
			}
			check := 0
			for i, x := range iss {
				if x.ChineseName == tcm.ChineseName {
					check = 1
					iss[i].IngredientId += (";" + strconv.FormatInt(s.Id, 10))
					iss[i].Count++
					break
				}
			}
			if check == 0 {
				var is IngredientSource
				is.ChineseName = tcm.ChineseName
				is.EnglishName = tcm.EnglishName
				is.IngredientId = strconv.FormatInt(s.Id, 10)
				is.PinyinName = tcm.PinyinName
				is.Count = 1
				iss = append(iss, is)
			}
		}
	}
	fmt.Println(iss)
	sortutil.DescByField(iss, "Count")
	/*
		o1 := orm.NewOrm()
		var maps []orm.Params
		//url = "select arc.id,arc.title,arc.typeid,art.typename from go_archives as arc left join go_arctype as art on art.id=arc.typeid where arc.typeid=?"
		//url = "select * from ingredient full join tcm_ingredient on ingredient.id = tcm_ingredient.ingredient_id where id = ?"

		//o1.Raw("select * from tcm_ingredient as ti join ingredient on ingredient.id=ti.ingredient_id where ingredient.id=?", 39308).Values(&maps)
		o1.Raw("select chinese_name from tcm_ingredient join tcm on tcm_ingredient.tcm_id=tcm.id where tcm_ingredient.tcm_id=1").Values(&maps)
		fmt.Println("sdsa")
		fmt.Println(maps)
		//resultsInfo["source"] = maps

	*/
	// convert to json

	var data DataList
	data.Draw = draw
	data.RecordsTotal = 100
	data.RecordsFiltered = 100
	data.Data = iss

	this.Data["json"] = data
	this.ServeJSON()
}

type IngredientSource struct {
	EnglishName  string `json:"EnglishName"`
	PinyinName   string `json:"EnglishName"`
	ChineseName  string `json:"ChineseName"`
	IngredientId string `json:"IngredientID"`
	Count        int64
}
