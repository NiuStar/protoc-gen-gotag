package module

import (
	"bytes"
	"errors"
	pgs "github.com/lyft/protoc-gen-star"
	"unicode"
	"go/format"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func createInitgo(targets map[string]pgs.File,outputPath,projectName string,log func(...interface{})) (error) {

	log("outputPath",outputPath)
	log("projectName",projectName)
	if len(targets) == 0 {
		return nil
	}
	//outputPath := context.OutputPath(targetFirst).String()
	//outputPath = path.Join(path.Dir(outputPath),"../api")
	fileName := path.Join(outputPath,"init.go")

	fileName,err := filepath.Abs(fileName)
	if err != nil {
		return err
	}
	log("fileName:",fileName)
	var codeBuf bytes.Buffer
	codeBuf.WriteString(`package api

`)

	for _,target := range targets {


		var service pgs.Service
		services := target.Services()
		if len(services) > 0 {
			service = services[0]
		} else {
			log("target:",target.Name(),"不是service")
			continue
		}
		packageName := service.Package().ProtoName().String()
		codeBuf.WriteString(`import _ "`)
		codeBuf.WriteString(projectName)
		codeBuf.WriteByte('/')
		codeBuf.WriteString(outputPath)
		codeBuf.WriteByte('/')
		codeBuf.WriteString(packageName)
		codeBuf.WriteString(`"
`)
	}


	return merge(fileName,codeBuf.String(),log)
}

func getLowerServiceName(name string) string {
	if unicode.IsUpper([]rune(name)[0]) {
		buf := bytes.Buffer{}
		buf.WriteByte(name[0] + 32)
		buf.WriteString(name[1:])
		return buf.String()
	}
	return name
}

func getUpServiceName(name string) string {
	if unicode.IsUpper([]rune(name)[0]) {
		return name
	}
	buf := bytes.Buffer{}
	buf.WriteByte(name[0] - 32)
	buf.WriteString(name[1:])
	return buf.String()
}
func createServiceCode(service pgs.Service,outputPath,golangPath,projectName string,log func( ...interface{})) (error) {
	packageName := service.Package().ProtoName().String()
	dirName,_ := filepath.Abs(path.Join(outputPath,packageName))
	_,err := os.Stat(dirName)
	if err != nil {
		if err = os.Mkdir(dirName,0777);err != nil {
			return err
		}
		err := os.Chmod(dirName, 0777)
		if err != err {
			return err
		}
		_,err = os.Stat(dirName)
		if err != err {
			return err
		}
	}
	serviceName := service.Name().String()
	fileName := path.Join(dirName,serviceName)
	var tval string
	tval = service.SourceCodeInfo().TrailingComments()// Extension(tagger.E_Tags, &tval)
	if len(tval) == 0 {
		tval = service.SourceCodeInfo().LeadingComments()
	}
	if len(tval) == 0 {
		tval = serviceName
	}
	tval = strings.TrimSuffix(tval,"\n")

	var codeBuf bytes.Buffer
	codeBuf.WriteString("package ")
	codeBuf.WriteString(packageName)
	codeBuf.WriteByte('\n')
	codeBuf.WriteString(`import (
	"context"
	"code.aliyun.com/new_backend/scodi_nqc/grpcserver/service"
`)
	codeBuf.WriteString("\"")
	codeBuf.WriteString(projectName)
	codeBuf.WriteByte('/')
	codeBuf.WriteString(golangPath)
	codeBuf.WriteByte('/')
	codeBuf.WriteString(packageName)
	codeBuf.WriteString(`"
)

func init() {
`)
	codeBuf.WriteString(getLowerServiceName(serviceName))
	codeBuf.WriteString(":=&")
	codeBuf.WriteString(serviceName)
	codeBuf.WriteString(`{}
`)
	codeBuf.WriteString(`service.ShareServiceManagerInstance().RegisterService(`)
	codeBuf.WriteString(getLowerServiceName(serviceName))
	codeBuf.WriteString(",(*")
	codeBuf.WriteString(packageName)
	codeBuf.WriteByte('.')
	codeBuf.WriteString(serviceName)
	codeBuf.WriteString("Server)(nil),\"")
	codeBuf.WriteString(tval)
	codeBuf.WriteString(`")
}

`)
	codeBuf.WriteString(`type `)
	codeBuf.WriteString(serviceName)
	codeBuf.WriteString(` struct {
}`)
	for _,method := range service.Methods() {
		methodName := getUpServiceName(method.Name().String())
		var mTval string
		mTval = method.SourceCodeInfo().LeadingComments()// Extension(tagger.E_Tags, &tval)
		if len(mTval) == 0 {
			mTval = methodName
		}
		mTval = strings.TrimSuffix(mTval,"\n")
		requestName := method.Input().Name().String()
		responseName := method.Output().Name().String()

		codeBuf.WriteString(`
//`)
		codeBuf.WriteString(mTval)
		codeBuf.WriteString(`
func (`)
		codeBuf.WriteString(getLowerServiceName(serviceName))
		codeBuf.WriteString(" *")
		codeBuf.WriteString(serviceName)
		codeBuf.WriteByte(')')
		codeBuf.WriteString(methodName)
		codeBuf.WriteString("(context context.Context,request *")
		codeBuf.WriteString(packageName)
		codeBuf.WriteByte('.')
		codeBuf.WriteString(requestName)
		codeBuf.WriteString(") (response *")
		codeBuf.WriteString(packageName)
		codeBuf.WriteByte('.')
		codeBuf.WriteString(responseName)
		codeBuf.WriteString(`, err error){
`)
		codeBuf.WriteString(`response = &`)
		codeBuf.WriteString(packageName)
		codeBuf.WriteByte('.')
		codeBuf.WriteString(responseName)
		codeBuf.WriteString(`{}
return
}

`)
	}

	data,err := format.Source(codeBuf.Bytes())

	if err != nil {
		return errors.New(codeBuf.String())
	}
	return merge(fileName + ".go",string(data),log)
}
