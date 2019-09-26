package models

import (
	"github.com/astaxie/beego/orm"
)

type Docs struct {
	Id        int64
	DocId     int64
	Journal   string
	Year      string
	Volume    string
	Issue     string
	FirstPage string
	LastPage  string
	PubmedId  string
	Doi       string
	Title     string
	Authors   string
	Abstract  string
}

func (m *Docs) TableName() string {
	return "doc"
}

func (m *Docs) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Docs) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Docs) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Docs) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Docs) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
