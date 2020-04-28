package module

import (
	"github.com/NiuStar/gocoder"
	"os"
)

func merge(file1,content string,log func(...interface{})) error {
	coder2,_ := gocoder.NewCoder(content)
	//return coder2.Save(file1)

	_,err := os.Stat(file1)
	if err != nil {
		log("不存在文件",file1)
		return coder2.Save(file1)
	}
	log("合并代码")
	coder,_ := gocoder.NewCoderWtihFile(file1)
	coder.Merge(coder2)
	return coder.Export()
}
