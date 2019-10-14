package main

import (
	_ "tcmio/routers"
	"tcmio/util"

	"github.com/astaxie/beego"
)

func main() {
	//util.SDFToDatabase("E:/GIM/TCM/TCMIO/data/all-ligands1.sdf")
	util.DataProcess()
	beego.Run()
}
