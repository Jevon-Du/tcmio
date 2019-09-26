package models

import (
	"github.com/astaxie/beego/orm"
)

type Target struct {
	Id            int64
	Name          string
	GeneName      string
	ProteinFamily string
	Mass          int64
	Length        int64
	EcNumber      string
	UniprotId     string
	ChemblId      string
	Drug          string
	Kegg          string
	Pdb           string
	Function      string
	Sequence      string
}

func (m *Target) TableName() string {
	return "target"
}

func (m *Target) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Target) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Target) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Target) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Target) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
