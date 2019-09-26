package main

import (
	_ "tcmio/routers"

	"github.com/astaxie/beego"
)

func main() {
	//util.SDFToDatabase("E:/GIM/TCM/TCMIO/data/all-ligands1.sdf")
	beego.Run()
}
