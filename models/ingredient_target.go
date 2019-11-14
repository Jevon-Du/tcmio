package models

import (
	"github.com/astaxie/beego/orm"
)

type Ingredient_Target struct {
	Id           int64
	TargetId     int64
	IngredientId int64
	Type         string
}

func (m *Ingredient_Target) TableName() string {
	return "ingredient_target"
}

func (m *Ingredient_Target) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Ingredient_Target) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Ingredient_Target) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Ingredient_Target) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Ingredient_Target) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
