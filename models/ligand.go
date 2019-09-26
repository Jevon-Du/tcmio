package models

import (
	"github.com/astaxie/beego/orm"
)

type Ligand struct {
	Id           int64
	ChemblId     string
	Name         string
	MolWeight    float64
	Formula      string
	Mol          string
	Smiles       string
	Inchi        string
	Inchikey     string
	Hba          int64
	Hbd          int64
	Rtb          int64
	Alogp        float64
	IngredientId int64
}

func (m *Ligand) TableName() string {
	return "ligand"
}

func (m *Ligand) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Ligand) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Ligand) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Ligand) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Ligand) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
