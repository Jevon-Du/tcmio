package routers

import (
	"tcmio/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:Index")

	beego.Router("/index", &controllers.MainController{}, "get:Index")
	beego.Router("/help", &controllers.MainController{}, "get:Help")
	beego.Router("/browse/:category:string", &controllers.MainController{}, "get:Browse")
	beego.Router("/structure", &controllers.MainController{}, "get:Structure")
	beego.Router("/moa", &controllers.MainController{}, "get:MOA")

	beego.Router("/targets", &controllers.MainController{}, "get:ListTargets")
	beego.Router("/targets/:id([0-9]+)", &controllers.MainController{}, "get:Detail")
	beego.Router("/targets/:id([0-9]+)/json", &controllers.MainController{}, "get:DetailTarget")

	beego.Router("/ligands", &controllers.MainController{}, "get:ListLigands")
	beego.Router("/ligands/:id([0-9]+)", &controllers.MainController{}, "get:Detail")
	beego.Router("/ligands/:id([0-9]+)/json", &controllers.MainController{}, "get:DetailLigand")

	beego.Router("/ingredients", &controllers.MainController{}, "get:ListIngredients")
	beego.Router("/ingredients/:id([0-9]+)", &controllers.MainController{}, "get:Detail")
	beego.Router("/ingredients/:id([0-9]+)/json", &controllers.MainController{}, "get:DetailIngredient")

	beego.Router("/tcms", &controllers.MainController{}, "get:ListTCMs")
	beego.Router("/tcms/:id([0-9]+)", &controllers.MainController{}, "get:Detail")
	beego.Router("/tcms/:id([0-9]+)/json", &controllers.MainController{}, "get:DetailTCM")

	beego.Router("/prescriptions", &controllers.MainController{}, "get:ListPrescriptions")
	beego.Router("/prescriptions/:id([0-9]+)", &controllers.MainController{}, "get:Detail")
	beego.Router("/prescriptions/:id([0-9]+)/json", &controllers.MainController{}, "get:DetailPrescription")

	beego.Router("/structure/ligands", &controllers.MainController{}, "*:StructureSearchLigand")
	beego.Router("/structure/analyze/ligands", &controllers.MainController{}, "*:AnalyzeStructureSearchLigand")
	beego.Router("/structure/ingredients", &controllers.MainController{}, "*:StructureSearchIngredient")
	beego.Router("/structure/analyze/ingredients", &controllers.MainController{}, "*:AnalyzeStructureSearchIngredient")

	beego.Router("/network/tcms", &controllers.MainController{}, "*:AnalyzeTCMs")
	beego.Router("/network/prescriptions", &controllers.MainController{}, "*:AnalyzePrescription")

}
