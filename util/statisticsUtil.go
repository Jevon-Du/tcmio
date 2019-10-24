package util

import (
	"fmt"
	"os"
	"tcmio/models"

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
