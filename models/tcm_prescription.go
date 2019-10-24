package models

import (
	"github.com/astaxie/beego/orm"
)

type TCM_Prescription struct {
	Id       int64
	PresId   int64
	TcmId    int64
	PresName string
	TcmName  string
}

func (m *TCM_Prescription) TableName() string {
	return "tcm_prescription"
}

func (m *TCM_Prescription) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *TCM_Prescription) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *TCM_Prescription) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *TCM_Prescription) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *TCM_Prescription) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
