package models

import (
	"github.com/astaxie/beego/orm"
)

type Ingredient struct {
	Id       int64
	Name     string
	Synonyms string
	Mol      string
	Smiles   string
	Inchi    string
	Inchikey string
	//MolWeight float64
	//Formula   string
	LigandId int64
}

func (m *Ingredient) TableName() string {
	return "ingredient"
}

func (m *Ingredient) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Ingredient) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Ingredient) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Ingredient) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Ingredient) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
