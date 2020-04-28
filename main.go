package main

import (
	"github.com/NiuStar/protoc-gen-gotag/module"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func main() {
	/*coder1,_ := gocoder.NewCoderWtihFile("api/mdna/EbsService.go")
	coder2,_ := gocoder.NewCoderWtihFile("api2/mdna/EbsService.go")
	coder1.Merge(coder2)
	fmt.Println(coder1.Export())
	return*/
	//pgs.Init(pgs.DebugEnv("GOTAG_DEBUG")).RegisterModule(module.New()).RegisterPostProcessor(pgsgo.GoFmt()).Render()
	pgs.Init(pgs.DebugMode()).RegisterModule(module.New()).RegisterPostProcessor(pgsgo.GoFmt()).Render()
}
