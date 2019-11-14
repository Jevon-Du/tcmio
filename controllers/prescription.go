package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"tcmio/models"
	"tcmio/util"
)

func (this *MainController) ListPrescriptions() {
	var tars []models.Prescription
	var tar models.Prescription
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

/* node id infor
1->Pres:1-1493-------->+10000
2->TCM:1-618---------->+1000
3->Ingredient:1-16437->+100000
4->Ligand:1-126972---->+200000
5->Target:1-400------->+0

*/

func (this *MainController) AnalyzePrescription() {
	fmt.Println("Prescription analysis")
	nameType := strings.TrimSpace(this.GetString("type"))
	presName := strings.TrimSpace(this.GetString("kw"))
	fmt.Println(presName)

	var p models.Prescription
	var pres []models.Prescription
	err := p.Query().Filter(nameType, presName).One(&p)
	if err != nil {
		panic(err)
	}
	pres = append(pres, p)

	var nodes []Node
	var edges []Edge

	var n Node
	var e Edge
	n.Id = p.Id + 10000
	n.Group = "prescriptions"
	n.Label = p.ChineseName

	nodes = append(nodes, n)

	// get TCM for prescription
	var tcms []models.TCM
	for _, s := range pres {
		var pt models.TCM_Prescription
		var pts []models.TCM_Prescription
		_, err = pt.Query().Limit(10000).Filter("pres_name", s.ChineseName).All(&pts)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\t%d\n", s.ChineseName, len(pts))

		for _, k := range pts {
			var t models.TCM
			err = t.Query().Filter("chinese_name", k.TcmName).One(&t)
			if err != nil {
				panic(err)
			}
			tcms = append(tcms, t)

			n.Id = t.Id + 1000
			n.Group = "tcms"
			n.Label = t.PinyinName
			nodes = append(nodes, n)
			e.From = p.Id + 400000
			e.To = n.Id
			edges = append(edges, e)
		}
	}

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
		//fmt.Println(len(its))
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
	util.WriteFile(string(res), "prescription.json")

	this.responseMsg.SuccessMsg("", resultsInfo)
	this.Data["json"] = this.responseMsg
	this.ServeJSON()

}
