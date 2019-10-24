package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"tcmio/models"

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

	fmt.Println(tars)

	//this.responseMsg.SuccessMsg("", data)
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
	tcmlist := strings.TrimSpace(this.GetString("tcms_name"))
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
		err := t.Query().Filter("chinese_name", s).One(&t)
		if err != nil {
			panic(err)
		}
		tcms = append(tcms, t)

		n.Id = t.Id + 100000
		n.Group = 1
		n.Label = t.PinyinName
		nodes = append(nodes, n)
	}

	nodes = append(nodes, n)

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
			e.From = h.Id + 100000
			e.To = m1.Id
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
			}
		}
	}

	for _, m := range ingres {
		n.Id = int64(m.Id)
		str := strconv.Itoa(int(m.Id))
		n.Label = str
		n.Group = 2
		nodes = append(nodes, n)
	}

	//get ligands from ingredients
	var ligs []models.Ligand
	for _, m := range ingres {

		url := "select * from ligand where mol @ (" + "0.8" + ", 1.0, '" + m.Mol + "', 'Tanimoto')::bingo.sim;"
		var mols []models.Ligand
		o := orm.NewOrm()
		_, err := o.Raw(url).QueryRows(&mols)
		if err != nil {
			//fmt.Println(ing.Smiles)
			panic(err)
		}
		//fmt.Printf("%s\t%d\n", ing.Smiles, len(mols))

		for _, mm := range mols {
			check := 0
			for _, nn := range ligs {
				if nn.Id == mm.Id {
					check = 1
					break
				}
			}
			if check == 0 {
				ligs = append(ligs, mm)
			}
		}
	}

	// get target for ligand
	var tars []models.Target
	for _, n := range ligs {
		var tl models.Target_ligand
		var tls []models.Target_ligand
		_, err := tl.Query().Filter("mol_chembl_id", n.ChemblId).All(&tls)
		if err != nil {
			panic(err)
		}
		for _, x := range tls {
			var tr models.Target
			err = tr.Query().Filter("chembl_id", x.TargetChemblId).One(&tr)
			//fmt.Printf("%s\t%s\t%s\t%s\n", s.ChineseName, t.ChineseName, ing.Smiles, tr.GeneName)

			check := 0
			for _, y := range tars {
				if y.Id == x.Id {
					check = 1
					break
				}
			}
			if check == 0 {
				tars = append(tars, tr)
			}
		}
	}
	//enrich := util.DavaidAnalysis(tars, j)

	resultsInfo := make(map[string]interface{})

	//resultsInfo["pres"] = pres
	resultsInfo["tcms"] = tcms
	resultsInfo["ingres"] = ingres
	resultsInfo["ligs"] = ligs
	resultsInfo["tars"] = tars
	resultsInfo["nodes"] = nodes
	resultsInfo["edges"] = edges
	//resultsInfo["pathways"] = enrich

	this.responseMsg.SuccessMsg("", resultsInfo)
	this.Data["json"] = this.responseMsg
	this.ServeJSON()

}
