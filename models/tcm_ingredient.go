package models

import (
	"github.com/astaxie/beego/orm"
)

type TCM_Ingredient struct {
	Id           int64
	TcmId        int64
	IngredientId int64
}

func (m *TCM_Ingredient) TableName() string {
	return "tcm_ingredient"
}

func (m *TCM_Ingredient) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *TCM_Ingredient) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *TCM_Ingredient) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *TCM_Ingredient) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *TCM_Ingredient) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
