package models

import "github.com/astaxie/beego/orm"

type TCM struct {
	Id             int64
	ChineseName    string `orm:"unique;size(32);index"`
	PinyinName     string `orm:"unique;type(text);index"`
	EnglishName    string `orm:"unique;type(text);index"`
	UsePart        string `orm:"type(text);null"`
	PropertyFlavor string `orm:"type(text);null"`
	ChannelTropism string `orm:"type(text);null"`
	Effect         string `orm:"type(text);null"`
	Indication     string `orm:"type(text);null"`
	RefSource      string `orm:"type(text);null"`
}

func (m *TCM) TableName() string {
	return "tcm"
}

func (m *TCM) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *TCM) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *TCM) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *TCM) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *TCM) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
