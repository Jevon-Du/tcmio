package routers

import (
	"tcmio/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:Index")

	beego.Router("/index", &controllers.MainController{}, "get:Index")
	beego.Router("/help", &controllers.MainController{}, "get:Help")
	beego.Router("/browse", &controllers.MainController{}, "get:Browse")
	beego.Router("/structure", &controllers.MainController{}, "get:Structure")
	beego.Router("/moa", &controllers.MainController{}, "get:MOA")

	beego.Router("/target", &controllers.MainController{}, "get:Target")
	beego.Router("/targets", &controllers.MainController{}, "get:ListTargets")
	beego.Router("/targets/:id([0-9]+)", &controllers.MainController{}, "get:DetailTarget")

	beego.Router("/ligand", &controllers.MainController{}, "get:Ligand")
	beego.Router("/ligands", &controllers.MainController{}, "get:ListLigands")
	beego.Router("/ligands/:id([0-9]+)", &controllers.MainController{}, "get:DetailLigand")
	beego.Router("/ligands/structure/:method", &controllers.MainController{}, "get:SearchLigands")

	beego.Router("/ingredient", &controllers.MainController{}, "get:Ingredient")
	beego.Router("/ingredients", &controllers.MainController{}, "get:ListIngredients")
	beego.Router("/ingredients/:id([0-9]+)", &controllers.MainController{}, "get:DetailIngredient")

	beego.Router("/tcm", &controllers.MainController{}, "get:TCM")
	beego.Router("/tcms", &controllers.MainController{}, "get:ListTCMs")
	beego.Router("/tcms/:id([0-9]+)", &controllers.MainController{}, "get:DetailTCM")

	beego.Router("/prescription", &controllers.MainController{}, "get:Prescription")
	beego.Router("/prescriptions", &controllers.MainController{}, "get:ListPrescriptions")
	beego.Router("/prescriptions/:id([0-9]+)", &controllers.MainController{}, "get:DetailPrescription")

	//beego.Router("/structure/analysis", &controllers.MainController{}, "post:StructureSearch")
	//beego.Router("/moa/analysis", &controllers.MainController{}, "post:MOA")

}
