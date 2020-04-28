package module

import (
	"go/parser"
	"go/printer"
	"go/token"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	//"fmt"
	"bytes"
	"github.com/fatih/structtag"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type mod struct {
	*pgs.ModuleBase
	pgsgo.Context
}

func New() pgs.Module {
	return &mod{ModuleBase: &pgs.ModuleBase{}}
}

func (m *mod) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.Context = pgsgo.InitContext(c.Parameters())
}

func (mod) Name() string {
	return "gotag"
}

func (m mod)parseValues() (url.Values,error) {
	values,err := url.ParseQuery(m.Params().String())
	if err != nil {
		return nil,err
	}

	if len(values["name"]) == 0 {
		dir,_ := filepath.Abs(".")
		projectName := filepath.Base(dir)
		values.Add("name",projectName)
	}
	if len(values["path"]) == 0 {
		values.Add("path","api")
	}
	if len(values["proto_path"]) == 0 {
		values.Add("proto_path",".")
	}
	return values,nil
}
func (m mod) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {

	//fmt.Println("targets",targets)
	xtv := m.Parameters().Str("xxx")

	xtv = strings.Replace(xtv, "+", ":", -1)

	xt, err := structtag.Parse(xtv)
	m.CheckErr(err)

	xt.Set(&structtag.Tag{Key:"testkey",Name:"testvalue"})
	extractor := newTagExtractor(m, m.Context)

	m.Log("params:",m.Params().String())

	params,err := m.parseValues()
	m.CheckErr(err)

	err = createInitgo(targets,params["path"][0],params["name"][0],m.Log)

	m.Log(m.Parameters())
	m.Log(m.Params())
	m.CheckErr(err)
	for _, f := range targets {
		m.Log("outputPath:",m.OutputPath(f).String())
		logBuffer := bytes.Buffer{}
		gfname := path.Join(params["proto_golang"][0],m.Context.OutputPath(f).SetExt(".go").String())
		logBuffer.WriteString(gfname+"\n")

		tags := extractor.Extract(f,params)


		fs := token.NewFileSet()
		fn, err := parser.ParseFile(fs, gfname, nil, parser.ParseComments)
		m.CheckErr(err)
		//fmt.Println("tag",len(tags))
		m.CheckErr(Retag(fn, tags))

		var buf strings.Builder
		m.CheckErr(printer.Fprint(&buf, fs, fn))
		//m.OverwriteCustomFile(
		//	"./log.txt",logBuffer.String(),0777)
		m.OverwriteGeneratorFile(gfname, buf.String())
	}

	return m.Artifacts()
}
