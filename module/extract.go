package module

import (
	"github.com/fatih/structtag"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"net/url"
	"strings"
)

type tagExtractor struct {
	pgs.Visitor
	pgs.DebuggerCommon
	pgsgo.Context
	outputPath string
	projectname string
	proto_golang string
	services []pgs.Service
	tags map[string]map[string]*structtag.Tags
}

func newTagExtractor(d pgs.DebuggerCommon, ctx pgsgo.Context) *tagExtractor {
	v := &tagExtractor{DebuggerCommon: d, Context: ctx}
	v.Visitor = pgs.PassThroughVisitor(v)
	return v
}

func (v *tagExtractor) VisitOneOf(o pgs.OneOf) (pgs.Visitor, error) {
	var tval string
	tval = o.SourceCodeInfo().TrailingComments()// Extension(tagger.E_Tags, &tval)
	if len(tval) == 0 {
		tval = o.SourceCodeInfo().LeadingComments()
	}

	msgName := v.Context.Name(o.Message()).String()

	if v.tags[msgName] == nil {
		v.tags[msgName] = map[string]*structtag.Tags{}
	}

	if len(tval) == 0 {
		return v, nil
	}

	tags := parseTags(tval)

	v.tags[msgName][v.Context.Name(o).String()] = tags

	return v, nil
}

func parseTags(tval string) *structtag.Tags {
	tags := &structtag.Tags{}
	list := strings.Split(strings.TrimPrefix(tval,"//")," ")
	for _,l := range list {
		list2 := strings.Split(l,":")
		if len(list2) > 1 {
			values := strings.Split(list2[1],"\"")
			if len(values) == 3 {
				tags.Set(&structtag.Tag{Key: list2[0], Name: values[1]})
			}
		} else {
			tags.Set(&structtag.Tag{Key: "comment", Name: strings.TrimSuffix(l,"\r\n")})
		}
	}
	return tags
}

func (v *tagExtractor) VisitField(f pgs.Field) (pgs.Visitor, error) {
	var tval string
	tval = f.SourceCodeInfo().TrailingComments()// Extension(tagger.E_Tags, &tval)
	if len(tval) == 0 {
		tval = f.SourceCodeInfo().LeadingComments()
	}

	msgName := v.Context.Name(f.Message()).String()

	if f.InOneOf() {
		msgName = f.Message().Name().UpperCamelCase().String() + "_" + f.Name().UpperCamelCase().String()
	}

	if v.tags[msgName] == nil {
		v.tags[msgName] = map[string]*structtag.Tags{}
	}

	if len(tval) == 0 {
		return v, nil
	}


	tags := parseTags(tval)

	v.tags[msgName][v.Context.Name(f).String()] = tags

	return v, nil
}

func (v *tagExtractor) VisitService(service pgs.Service) (pgs.Visitor,error) {
	err := createServiceCode(service,v.outputPath,v.proto_golang,v.projectname,v.Log)
	if err != nil {
		return nil,err
	}
	return v,nil
}
func (v *tagExtractor) Extract(f pgs.File,parameters url.Values) StructTags {
	v.tags = map[string]map[string]*structtag.Tags{}
	v.outputPath = v.OutputPath(f).String()
	v.outputPath = parameters["path"][0]
	v.projectname = parameters["name"][0]
	v.proto_golang = parameters["proto_golang"][0]

	v.services = f.Services()
	v.CheckErr(pgs.Walk(v, f))
	return v.tags
}
