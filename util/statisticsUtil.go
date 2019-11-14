package util

import (
	"fmt"
	"os"
	"tcmio/models"

	"github.com/tealeg/xlsx"

	"github.com/astaxie/beego/orm"
)

// count ligands for target and targets for ligand
func Statistics() {
	var tar models.Target
	var tars []models.Target

	cnt, err := tar.Query().All(&tars)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cnt)

	fi, err := os.Create("target-ligands.txt")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	for _, s := range tars {

		if s.ChemblId != "" {
			var tl models.Target_ligand
			a, _ := tl.Query().Filter("target_chembl_id", s.ChemblId).Count()
			fmt.Printf("%s\t%d\n", s.ChemblId, a)
			fmt.Fprintf(fi, "%s\t%d\n", s.ChemblId, a)
			//fmt.Printf("%s\t%d\n", s.ChemblId, len(tls))
		}
	}

	fi2, err := os.Create("target-ligands-2.txt")
	if err != nil {
		panic(err)
	}
	defer fi2.Close()
	var lig models.Ligand
	var ligs []models.Ligand

	cnt, err = lig.Query().Limit(1000000).All(&ligs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cnt)

	for _, s := range ligs {

		if s.ChemblId != "" {
			var tl models.Target_ligand
			a, _ := tl.Query().Filter("mol_chembl_id", s.ChemblId).Count()
			fmt.Printf("%s\t%d\n", s.ChemblId, a)
			fmt.Fprintf(fi2, "%s\t%d\n", s.ChemblId, a)
			//fmt.Printf("%s\t%d\n", s.ChemblId, len(tls))
		}
	}

}

// similarity search ligands and ingredients
func SimComp() {
	var ingre models.Ingredient
	var ingres []models.Ingredient

	_, err := ingre.Query().Limit(1000000).All(&ingres)
	if err != nil {
		panic(err)
	}
	o := orm.NewOrm()
	fi, err := os.Create("similarityComp-1.0.txt")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	for _, s := range ingres {
		//fmt.Println(s.Smiles)
		url := "select * from ligand where mol @ (" + "0.7" + ", 1.0, '" + s.Mol + "', 'Tanimoto')::bingo.sim;"
		var hits []models.Ligand
		cnt, err := o.Raw(url).QueryRows(&hits)
		if err != nil {
			//panic(err)
			//fmt.Println(s.Smiles)
			//continue
			cnt = 0
		}
		fmt.Printf("%s\t%d\n", s.Inchikey, cnt)
		fmt.Fprintf(fi, "%s\t%d\n", s.Inchikey, cnt)
	}
}

// generate and statistics network for prescription

func PrescriptionNetwork() {

	var p models.Prescription
	var ps []models.Prescription
	_, err := p.Query().Limit(10000).All(&ps)
	if err != nil {
		panic(err)
	}
	//fmt.Println(len(ps))
	for _, s := range ps {
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
			var ti models.TCM_Ingredient
			var tis []models.TCM_Ingredient
			_, err = ti.Query().Filter("tcm_id", t.Id).All(&tis)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s\t%d\n", t.ChineseName, len(tis))

			for _, m := range tis {
				var ing models.Ingredient
				err = ing.Query().Filter("id", m.IngredientId).One(&ing)

				url := "select * from ligand where mol @ (" + "0.8" + ", 1.0, '" + ing.Mol + "', 'Tanimoto')::bingo.sim;"
				var ligs []models.Ligand
				o := orm.NewOrm()
				_, err := o.Raw(url).QueryRows(&ligs)
				if err != nil {
					fmt.Println(ing.Smiles)
					panic(err)
				}
				fmt.Printf("%s\t%d\n", ing.Smiles, len(ligs))

				for _, n := range ligs {
					var tl models.Target_ligand
					var tls []models.Target_ligand
					_, err = tl.Query().Filter("mol_chembl_id", n.ChemblId).All(&tls)
					if err != nil {
						panic(err)
					}
					for _, x := range tls {
						var tr models.Target
						err = tr.Query().Filter("chembl_id", x.TargetChemblId).One(&tr)
						fmt.Printf("%s\t%s\t%s\t%s\n", s.ChineseName, t.ChineseName, ing.Smiles, tr.GeneName)
					}
				}

			}
		}
	}

}

func CleanData() {
	var ing models.TCM_Ingredient
	var ings []models.TCM_Ingredient
	_, err := ing.Query().Limit(10000000).All(&ings)
	if err != nil {
		fmt.Println(err)
	}

	var uni []models.TCM_Ingredient
	for _, s := range ings {
		check := 0
		for _, t := range uni {
			if s.TcmId == t.TcmId && s.IngredientId == t.IngredientId {
				fmt.Println(s)
				check = 1
				var a models.TCM_Ingredient

				err := a.Query().Filter("id", s.Id).One(&a)
				if err != nil {
					panic(err)
				}
				err = a.Delete()
				if err != nil {
					panic(err)
				}
				break
			}
		}
		if check == 0 {
			uni = append(uni, s)
		}

	}

	fmt.Println(len(ings))
	fmt.Println(len(uni))

}

type NP_Target struct {
	Tcmid    string
	Inchikey string
	GeneName string
}

type UTarget struct {
	Target   string
	GeneName string
	ChemblId string
	Uniprot  string
	PrefName string
	Tid      string
	Organism string
}

type TCM_NP struct {
	Herb     string
	Pinyin   string
	Inchikey string
	Tcmid    string
}

type TCM struct {
	ChineseName string
	Pinyin      string
}

func GetNPTarget() {

	// get tcm-np infor
	xlFile, err := xlsx.OpenFile("TCM_TARGET_part1.xlsx")
	if err != nil {
		fmt.Println("err in open file")
	}
	var nts []NP_Target
	var ts []UTarget
	fmt.Println("start")
	for j, sheet := range xlFile.Sheets {
		if j == 0 {
			for i, row := range sheet.Rows {
				if i < 1 {
					continue
				}
				if row.Cells[0].String() == "" {
					continue
				}
				var nt NP_Target
				nt.Tcmid = row.Cells[0].String()
				nt.Inchikey = row.Cells[1].String()
				nt.GeneName = row.Cells[7].String()
				nts = append(nts, nt)
				check := 0
				for _, s := range ts {
					if s.GeneName == nt.GeneName {
						check = 1
						break
					}
				}
				if check == 0 {
					var t UTarget
					t.GeneName = nt.GeneName
					ts = append(ts, t)
				}
			}
		}
	}

	xlFile, err = xlsx.OpenFile("E:/GIM/TCM/TCMIO/processed_2_bSDTNBI_KR_0.1_0.1_-0.5/TCM_TARGET_part2.xlsx")
	if err != nil {
		fmt.Println("err in open file")
	}
	for j, sheet := range xlFile.Sheets {
		if j == 0 {
			for i, row := range sheet.Rows {
				if i < 1 {
					continue
				}
				if row.Cells[0].String() == "" {
					continue
				}
				var nt NP_Target
				nt.Tcmid = row.Cells[0].String()
				nt.Inchikey = row.Cells[1].String()
				nt.GeneName = row.Cells[7].String()
				nts = append(nts, nt)
				check := 0
				for _, s := range ts {
					if s.GeneName == nt.GeneName {
						check = 1
						break
					}
				}
				if check == 0 {
					var t UTarget
					t.GeneName = nt.GeneName
					ts = append(ts, t)
				}
			}
		}
	}
	fmt.Printf("np-target relations:%d\n", len(nts))
	fmt.Printf("unique target:%d\n", len(ts))

	var tars []models.Target
	o := orm.NewOrm()
	_, err = o.Raw("select * from target").QueryRows(&tars)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(tars))
	cnt := 0
	for _, s := range tars {
		for _, t := range ts {
			if s.GeneName == t.GeneName {
				fmt.Printf("%s\t%s\n", s.GeneName, t.GeneName)
				cnt++
				break
			}
		}
	}
	fmt.Println(cnt)
	cnt = 0
	var snts []NP_Target
	fi, err := os.Create("ingredient-target.txt")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fmt.Fprintf(fi, "ingredient_id\ttarget_id\ttype\n")
	for _, s := range nts {
		var ing models.Ingredient
		err = ing.Query().Filter("inchikey", s.Inchikey).One(&ing)
		if err != nil {
			//fmt.Println(err)
			continue
		}
		var t models.Target
		err = t.Query().Filter("gene_name", s.GeneName).One(&t)
		if err != nil {
			//fmt.Println(err)
			continue
		}
		fmt.Println(s)
		fmt.Fprintf(fi, "%d\t%d\t%s\n", ing.Id, t.Id, "Network-based prediction")
		snts = append(snts, s)
		cnt++
	}
	fmt.Println(cnt)
	fmt.Println(len(snts))

}

func AnalyzeNetworkBasedNPTargetPrediction() {
	fmt.Println("Start")
	GetNPTarget()
	//GetIngredientTargetRelation()
}

func GetIngredientTargetRelation() {

	var ing models.Ingredient
	var ings []models.Ingredient

	_, err := ing.Query().All(&ings)
	if err != nil {
		panic(err)
	}
	cnt := 0
	for _, s := range ings {
		var lig models.Ligand
		err = lig.Query().Filter("inchikey", s.Inchikey).One(&lig)
		if err != nil {
			continue
		}

		var lt models.Target_ligand
		var lts []models.Target_ligand

		_, err = lt.Query().Filter("mol_chembl_id", lig.ChemblId).All(&lts)
		if err != nil {
			continue
		}

		for _, t := range lts {
			//var it models.Ingredient_Target
			//it.IngredientId = s.Id
			//it.TargetId = t.Id
			//it.Insert()
			fmt.Println(t)
			cnt++
		}
	}
	fmt.Println(cnt)

}
