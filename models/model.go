package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
)

func init() {
	dbhost := beego.AppConfig.String("dbhost")
	dbport := beego.AppConfig.String("dbport")
	dbuser := beego.AppConfig.String("dbuser")
	dbpassword := beego.AppConfig.String("dbpassword")
	dbname := beego.AppConfig.String("dbname")
	if dbport == "" {
		dbport = "5432"
	}
	dburl := "user=" + dbuser + " password=" + dbpassword + " dbname=" + dbname + " host=" + dbhost + " port=" + dbport + " sslmode=disable"
	//fmt.Printf(dburl)
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", dburl)

	//orm.RunSyncdb("default", false, true)
	orm.RegisterModel(new(Target))
	orm.RegisterModel(new(Ligand))
	orm.RegisterModel(new(Ingredient))
	orm.RegisterModel(new(TCM))
	orm.RegisterModel(new(TCM_Prescription))
	orm.RegisterModel(new(TCM_Ingredient))
	orm.RegisterModel(new(Ingredient_Target))
	orm.RegisterModel(new(Target_ligand))
	orm.RegisterModel(new(Prescription))
	orm.RegisterModel(new(Docs))
}
