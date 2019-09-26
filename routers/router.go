package routers

import (
	"tcmio/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	beego.Router("/index", &controllers.MainController{}, "get:Index")
	beego.Router("/help", &controllers.MainController{}, "get:Help")
	beego.Router("/browse", &controllers.MainController{}, "get:Browse")
	beego.Router("/structure", &controllers.MainController{}, "get:Structure")
	beego.Router("/moa", &controllers.MainController{}, "get:MOA")

	beego.Router("/targets", &controllers.MainController{}, "get:ListTargets")
	beego.Router("/targets/:id([0-9]+)", &controllers.MainController{}, "get:DetailTarget")

	beego.Router("/ligands", &controllers.MainController{}, "get:ListLigands")
	beego.Router("/ligands/:id([0-9]+)", &controllers.MainController{}, "get:DetailLigand")
	beego.Router("/ligands/structure/:method", &controllers.MainController{}, "get:SearchLigands")

	beego.Router("/ingredients", &controllers.MainController{}, "get:ListIngredients")
	beego.Router("/ingredients/:id([0-9]+)", &controllers.MainController{}, "get:DetailIngredient")

	beego.Router("/tcms", &controllers.MainController{}, "get:ListTCMs")
	beego.Router("/tcms/:id([0-9]+)", &controllers.MainController{}, "get:DetailTCM")

	beego.Router("/prescriptions", &controllers.MainController{}, "get:ListPrescriptions")
	beego.Router("/prescriptions/:id([0-9]+)", &controllers.MainController{}, "get:DetailPrescription")

	//beego.Router("/structure/analysis", &controllers.MainController{}, "post:StructureSearch")
	//beego.Router("/moa/analysis", &controllers.MainController{}, "post:MOA")

}
