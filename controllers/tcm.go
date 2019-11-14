package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"tcmio/models"
	"tcmio/util"

	"github.com/astaxie/beego/orm"
)

func (this *MainController) ListTCMs() {
	var tars []models.TCM
	var tar models.TCM
	draw, _ := this.GetInt64("draw")
	offset, _ := this.GetInt64("start")
	limit, _ := this.GetInt64("length")
	fmt.Println(offset)
	fmt.Println(limit)
	_, err := tar.Query().Limit(limit, offset).All(&tars)
	if err != nil {
		fmt.Println(err)
	}

	total, _ := tar.Query().Count()

	var data DataList
	data.Draw = draw
	data.RecordsTotal = total
	data.RecordsFiltered = total
	data.Data = tars

	this.Data["json"] = data
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

func (this *MainController) AnalyzeTCMs() {

	fmt.Println("TCM analysis")
	name_type := this.GetString("type")
	tcmlist := strings.TrimSpace(this.GetString("kw"))
	fmt.Println(tcmlist)
	tcmArray := strings.Split(tcmlist, ",")
	fmt.Println(tcmArray)
	var tcms []models.TCM
	var n Node
	var e Edge
	var nodes []Node
	var edges []Edge
	for _, s := range tcmArray {
		var t models.TCM
		err := t.Query().Filter(name_type, s).One(&t)
		if err != nil {
			panic(err)
		}
		tcms = append(tcms, t)

		n.Id = t.Id + 1000
		n.Group = "tcms"
		n.Label = t.PinyinName
		nodes = append(nodes, n)
	}
	fmt.Printf("TCM:%d\n", len(tcms))
	fmt.Printf("Nodes:%d\n", len(nodes))
	fmt.Printf("Edges:%d\n", len(edges))

	// get ingredinets for tcm
	var ingres []models.Ingredient
	for _, h := range tcms {

		var ti models.TCM_Ingredient
		var tis []models.TCM_Ingredient
		_, err := ti.Query().Filter("tcm_id", h.Id).All(&tis)
		if err != nil {
			panic(err)
		}
		//fmt.Printf("%s\t%d\n", t.ChineseName, len(tis))
		var mols []models.Ingredient
		for _, m := range tis {
			var ing models.Ingredient
			err = ing.Query().Filter("id", m.IngredientId).One(&ing)
			mols = append(mols, ing)
		}

		for _, m1 := range mols {
			e.From = h.Id + 1000
			e.To = m1.Id + 100000
			edges = append(edges, e)
			check := 0
			for _, m2 := range ingres {
				if m1.Id == m2.Id {
					check = 1
					break
				}
			}
			if check == 0 {
				ingres = append(ingres, m1)
				n.Id = m1.Id + 100000
				n.Label = strconv.Itoa(int(m1.Id))
				n.Group = "ingredients"
				nodes = append(nodes, n)
			}
		}
	}

	fmt.Printf("Ingredients:%d\n", len(ingres))
	fmt.Printf("Nodes:%d\n", len(nodes))
	fmt.Printf("Edges:%d\n", len(edges))

	/*
		for _, m := range ingres {
			n.Id = int64(m.Id + 100000)
			str := strconv.Itoa(int(m.Id))
			n.Label = str
			n.Group = 3
			nodes = append(nodes, n)
		}
	*/

	// get target for ingredient
	var tars []models.Target
	for _, i := range ingres {

		var it models.Ingredient_Target
		var its []models.Ingredient_Target
		_, err := it.Query().Filter("ingredient_id", i.Id).All(&its)
		if err != nil {
			//panic(err)
			//fmt.Println(err)
			continue
		}
		for _, x := range its {
			var tr models.Target
			err = tr.Query().Filter("id", x.TargetId).One(&tr)
			if err != nil {
				continue
			}
			e.From = i.Id + 100000
			e.To = tr.Id
			edges = append(edges, e)

			check := 0
			for _, y := range tars {
				if y.Id == tr.Id {
					check = 1
					break
				}
			}
			if check == 0 {
				tars = append(tars, tr)
				n.Id = tr.Id
				n.Label = tr.GeneName
				n.Group = "targets"
				nodes = append(nodes, n)
			}
		}
	}

	fmt.Printf("Targets:%d\n", len(tars))
	fmt.Printf("Nodes:%d\n", len(nodes))
	fmt.Printf("Edges:%d\n", len(edges))

	//enrich := util.DavaidAnalysis(tars)

	resultsInfo := make(map[string]interface{})

	//resultsInfo["pres"] = pres
	//resultsInfo["tcms"] = tcms
	//resultsInfo["ingres"] = ingres
	//resultsInfo["ligs"] = ligs
	//resultsInfo["tars"] = tars
	resultsInfo["nodes"] = nodes
	resultsInfo["edges"] = edges
	//resultsInfo["pathways"] = enrich
	res, _ := json.Marshal(resultsInfo)
	util.WriteFile(string(res), "test6.json")

	this.responseMsg.SuccessMsg("", resultsInfo)
	this.Data["json"] = this.responseMsg
	this.ServeJSON()

}

/*
1->Pres:1-1493->10000-11493
2->TCM:1-618->1000-1618
3->Ingredient:1-16437->100000-116473
4->Ligand:1-126972->200000-326972
5->Target:1-400->1-400

*/

func (this *MainController) AnalyzeTCMs1() {

	fmt.Println("TCM analysis")
	name_type := this.GetString("type")
	tcmlist := strings.TrimSpace(this.GetString("kw"))
	fmt.Println(tcmlist)
	tcmArray := strings.Split(tcmlist, ",")
	fmt.Println(tcmArray)
	var tcms []models.TCM
	var n Node
	var e Edge
	var nodes []Node
	var edges []Edge
	for _, s := range tcmArray {
		var t models.TCM
		err := t.Query().Filter(name_type, s).One(&t)
		if err != nil {
			panic(err)
		}
		tcms = append(tcms, t)

		n.Id = t.Id + 1000
		n.Group = "tcms"
		n.Label = t.PinyinName
		nodes = append(nodes, n)
	}
	fmt.Printf("TCM:%d\n", len(tcms))
	fmt.Printf("Nodes:%d\n", len(nodes))
	fmt.Printf("Edges:%d\n", len(edges))

	// get ingredinets for tcm
	var ingres []models.Ingredient
	for _, h := range tcms {

		var ti models.TCM_Ingredient
		var tis []models.TCM_Ingredient
		_, err := ti.Query().Filter("tcm_id", h.Id).All(&tis)
		if err != nil {
			panic(err)
		}
		//fmt.Printf("%s\t%d\n", t.ChineseName, len(tis))
		var mols []models.Ingredient
		for _, m := range tis {
			var ing models.Ingredient
			err = ing.Query().Filter("id", m.IngredientId).One(&ing)
			mols = append(mols, ing)
		}

		for _, m1 := range mols {
			e.From = h.Id + 1000
			e.To = m1.Id + 100000
			edges = append(edges, e)
			check := 0
			for _, m2 := range ingres {
				if m1.Id == m2.Id {
					check = 1
					break
				}
			}
			if check == 0 {
				ingres = append(ingres, m1)
				n.Id = m1.Id + 100000
				n.Label = strconv.Itoa(int(m1.Id))
				n.Group = "ingredients"
				nodes = append(nodes, n)
			}
		}
	}

	fmt.Printf("Ingredients:%d\n", len(ingres))
	fmt.Printf("Nodes:%d\n", len(nodes))
	fmt.Printf("Edges:%d\n", len(edges))

	/*
		for _, m := range ingres {
			n.Id = int64(m.Id + 100000)
			str := strconv.Itoa(int(m.Id))
			n.Label = str
			n.Group = 3
			nodes = append(nodes, n)
		}
	*/
	//get ligands from ingredients
	var ligs []models.Ligand
	for _, m := range ingres {

		url := "select * from ligand where mol @ (" + "0.95" + ", 1.0, '" + m.Mol + "', 'Tanimoto')::bingo.sim;"
		var mols []models.Ligand
		o := orm.NewOrm()
		_, err := o.Raw(url).QueryRows(&mols)
		if err != nil {
			//fmt.Println(ing.Smiles)
			panic(err)
		}
		//fmt.Printf("%s\t%d\n", ing.Smiles, len(mols))

		for _, mm := range mols {
			e.From = m.Id + 100000
			e.To = mm.Id + 200000
			edges = append(edges, e)
			check := 0
			for _, nn := range ligs {
				if nn.Id == mm.Id {
					check = 1
					break
				}
			}
			if check == 0 {
				ligs = append(ligs, mm)
				n.Id = mm.Id + 200000
				n.Group = "ligands"
				n.Label = mm.ChemblId
				nodes = append(nodes, n)
			}
		}
	}

	fmt.Printf("Ingredients:%d\n", len(ligs))
	fmt.Printf("Nodes:%d\n", len(nodes))
	fmt.Printf("Edges:%d\n", len(edges))

	// get target for ligand
	var tars []models.Target
	for _, lg := range ligs {

		var tl models.Target_ligand
		var tls []models.Target_ligand
		_, err := tl.Query().Filter("mol_chembl_id", lg.ChemblId).All(&tls)
		if err != nil {
			panic(err)
		}
		for _, x := range tls {
			var tr models.Target
			err = tr.Query().Filter("chembl_id", x.TargetChemblId).One(&tr)
			if err != nil {
				continue
			}
			e.From = lg.Id + 200000
			e.To = tr.Id
			edges = append(edges, e)
			//fmt.Printf("%s\t%s\t%s\t%s\n", s.ChineseName, t.ChineseName, ing.Smiles, tr.GeneName)
			check := 0
			for _, y := range tars {
				if y.Id == tr.Id {
					check = 1
					break
				}
			}
			if check == 0 {
				tars = append(tars, tr)
				n.Id = tr.Id
				n.Label = tr.GeneName
				n.Group = "targets"
				nodes = append(nodes, n)
			}
		}
	}

	fmt.Printf("Targets:%d\n", len(tars))
	fmt.Printf("Nodes:%d\n", len(nodes))
	fmt.Printf("Edges:%d\n", len(edges))

	//enrich := util.DavaidAnalysis(tars, j)

	resultsInfo := make(map[string]interface{})

	//resultsInfo["pres"] = pres
	//resultsInfo["tcms"] = tcms
	//resultsInfo["ingres"] = ingres
	//resultsInfo["ligs"] = ligs
	//resultsInfo["tars"] = tars
	resultsInfo["nodes"] = nodes
	resultsInfo["edges"] = edges
	//resultsInfo["pathways"] = enrich
	res, _ := json.Marshal(resultsInfo)
	util.WriteFile(string(res), "test5.json")

	this.responseMsg.SuccessMsg("", resultsInfo)
	this.Data["json"] = this.responseMsg
	this.ServeJSON()

}
