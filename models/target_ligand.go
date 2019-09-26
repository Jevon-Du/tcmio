package models

import "github.com/astaxie/beego/orm"

type Target_ligand struct {
	Id             int64
	TargetChemblId string
	MolChemblId    string
	Value          float64
	Unit           string
	Type           string
	RefId          int64
}

func (m *Target_ligand) TableName() string {
	return "Target_ligand"
}

func (m *Target_ligand) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Target_ligand) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Target_ligand) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Target_ligand) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Target_ligand) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
