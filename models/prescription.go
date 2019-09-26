package models

import "github.com/astaxie/beego/orm"

type Prescription struct {
	Id          int64
	ChineseName string
	PinyinName  string
	Ingredients string
	Indication  string
	Effect      string
	RefSource   string
}

func (m *Prescription) TableName() string {
	return "prescription"
}

func (m *Prescription) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Prescription) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Prescription) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Prescription) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Prescription) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
