package codegen

import (
	"bufio"
	"bytes"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	v3helper "github.com/unionj-cloud/go-doudou/cmd/internal/openapi/v3"
	"github.com/unionj-cloud/go-doudou/toolkit/copier"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var httpHandlerImpl = `package httpsrv

import (
	v3 "github.com/unionj-cloud/go-doudou/toolkit/openapi/v3"
	{{.ServiceAlias}} "{{.ServicePackage}}"
	"net/http"
)

type {{.Meta.Name}}HandlerImpl struct{
	{{.Meta.Name | toLowerCamel}} {{.ServiceAlias}}.{{.Meta.Name}}
}

{{- range $m := .Meta.Methods }}
	func (receiver *{{$.Meta.Name}}HandlerImpl) {{$m.Name}}(w http.ResponseWriter, r *http.Request) {
    	panic("implement me")
    }
{{- end }}

func New{{.Meta.Name}}Handler({{.Meta.Name | toLowerCamel}} {{.ServiceAlias}}.{{.Meta.Name}}) {{.Meta.Name}}Handler {
	return &{{.Meta.Name}}HandlerImpl{
		{{.Meta.Name | toLowerCamel}},
	}
}
`

// GenHttpHandlerImpl generates http handler implementations
func GenHttpHandlerImpl(dir string, ic astutils.InterfaceCollector) {
	var (
		err             error
		modfile         string
		modName         string
		firstLine       string
		handlerimplfile string
		f               *os.File
		tpl             *template.Template
		source          string
		buf             bytes.Buffer
		httpDir         string
	)
	httpDir = filepath.Join(dir, "transport/httpsrv")
	if err = os.MkdirAll(httpDir, os.ModePerm); err != nil {
		panic(err)
	}

	handlerimplfile = filepath.Join(httpDir, "handlerimpl.go")
	if _, err = os.Stat(handlerimplfile); os.IsNotExist(err) {
		modfile = filepath.Join(dir, "go.mod")
		if f, err = os.Open(modfile); err != nil {
			panic(err)
		}
		reader := bufio.NewReader(f)
		firstLine, _ = reader.ReadString('\n')
		modName = strings.TrimSpace(strings.TrimPrefix(firstLine, "module"))

		if f, err = os.Create(handlerimplfile); err != nil {
			panic(err)
		}
		defer f.Close()

		funcMap := make(map[string]interface{})
		funcMap["toLowerCamel"] = strcase.ToLowerCamel
		funcMap["toCamel"] = strcase.ToCamel
		if tpl, err = template.New("handlerimpl.go.tmpl").Funcs(funcMap).Parse(httpHandlerImpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(&buf, struct {
			ServicePackage string
			ServiceAlias   string
			VoPackage      string
			Meta           astutils.InterfaceMeta
		}{
			ServicePackage: modName,
			ServiceAlias:   ic.Package.Name,
			VoPackage:      modName + "/vo",
			Meta:           ic.Interfaces[0],
		}); err != nil {
			panic(err)
		}

		source = strings.TrimSpace(buf.String())
		astutils.FixImport([]byte(source), handlerimplfile)
	} else {
		logrus.Warnf("file %s already exists.", handlerimplfile)
	}
}

var appendHttpHandlerImplTmpl = `
{{- range $m := .Meta.Methods }}
	func (receiver *{{$.Meta.Name}}HandlerImpl) {{$m.Name}}(_writer http.ResponseWriter, _req *http.Request) {
    	var (
			{{- range $p := $m.Params }}
			{{- if isVarargs $p.Type }}
			{{ $p.Name }} = new({{ $p.Type | toSlice }})
			{{- else }}
			{{ $p.Name }} {{ $p.Type }}
			{{- end }}
			{{- end }}
			{{- range $r := $m.Results }}
			{{ $r.Name }} {{ $r.Type }}
			{{- end }}
		)
		{{- $multipartFormParsed := false }}
		{{- $formParsed := false }}
		{{- range $p := $m.Params }}
		{{- if or (eq $p.Type "*multipart.FileHeader") (eq $p.Type "[]*multipart.FileHeader") }}
		{{- if not $multipartFormParsed }}
		if _err := _req.ParseMultipartForm(32 << 20); _err != nil {
			http.Error(_writer, _err.Error(), http.StatusBadRequest)
			return
		}
		{{- $multipartFormParsed = true }}
		{{- end }}
		{{- if contains $p.Type "["}}
		{{$p.Name}} = _req.MultipartForm.File["{{$p.Name}}"]
		{{- else}}
		{{$p.Name}}Files := _req.MultipartForm.File["{{$p.Name}}"]
		if len({{$p.Name}}Files) > 0 {
			{{$p.Name}} = {{$p.Name}}Files[0]
		}
		{{- end}}
		{{- else if or (eq $p.Type "v3.FileModel") (eq $p.Type "*v3.FileModel") (eq $p.Type "[]v3.FileModel") (eq $p.Type "*[]v3.FileModel") (eq $p.Type "...v3.FileModel") }}
		{{- if not $multipartFormParsed }}
		if _err := _req.ParseMultipartForm(32 << 20); _err != nil {
			http.Error(_writer, _err.Error(), http.StatusBadRequest)
			return
		}
		{{- $multipartFormParsed = true }}
		{{- end }}
		{{$p.Name}}FileHeaders, exists := _req.MultipartForm.File["{{$p.Name}}"]
		if exists {
			{{- if not (isOptional $p.Type) }}
			if len({{$p.Name}}FileHeaders) == 0 {
				http.Error(_writer, "no file uploaded for parameter {{$p.Name}}", http.StatusBadRequest)
				return
			}
			{{- end }}
			{{- if isSlice $p.Type }}
			{{- if isOptional $p.Type }}
			if {{$p.Name}} == nil && len({{$p.Name}}FileHeaders) > 0 {
				{{$p.Name}} = new([]v3.FileModel)
			}
			{{- end }}
			for _, _fh :=range {{$p.Name}}FileHeaders {
				_f, _err := _fh.Open()
				if _err != nil {
					http.Error(_writer, _err.Error(), http.StatusBadRequest)
					return
				}
				{{- if isOptional $p.Type }}
				*{{$p.Name}} = append(*{{$p.Name}}, v3.FileModel{
					Filename: _fh.Filename,
					Reader: _f,
				})
				{{- else }}
				{{$p.Name}} = append({{$p.Name}}, v3.FileModel{
					Filename: _fh.Filename,
					Reader: _f,
				})
				{{- end }}
			}
			{{- else}}
			if len({{$p.Name}}FileHeaders) > 0 {
				_fh := {{$p.Name}}FileHeaders[0]
				_f, _err := _fh.Open()
				if _err != nil {
					http.Error(_writer, _err.Error(), http.StatusBadRequest)
					return
				}
				{{- if isOptional $p.Type }}
				{{$p.Name}} = &v3.FileModel{
					Filename: _fh.Filename,
					Reader: _f,
				}
				{{- else }}
				{{$p.Name}} = v3.FileModel{
					Filename: _fh.Filename,
					Reader: _f,
				}
				{{- end }}
			}
			{{- end}}
		}{{- if not (isOptional $p.Type) }} else {
			http.Error(_writer, "missing parameter {{$p.Name}}", http.StatusBadRequest)
			return
		}{{- end }}
		{{- else if eq $p.Type "context.Context" }}
		{{$p.Name}} = _req.Context()
		{{- else if not (isBuiltin $p)}}
		{{- if isOptional $p.Type }}
		if _req.Body != nil {
			if _err := json.NewDecoder(_req.Body).Decode(&{{$p.Name}}); _err != nil {
				if _err != io.EOF {
					http.Error(_writer, _err.Error(), http.StatusBadRequest)
					return				
				}
			}
		}
		{{- else }}
		if _req.Body == nil {
			http.Error(_writer, "missing request body", http.StatusBadRequest)
			return
		} else {
			if _err := json.NewDecoder(_req.Body).Decode(&{{$p.Name}}); _err != nil {
				http.Error(_writer, _err.Error(), http.StatusBadRequest)
				return
			}
		}
		{{- end }}
		{{- else if isSlice $p.Type }}
		{{- if not $formParsed }}
		if _err := _req.ParseForm(); _err != nil {
			http.Error(_writer, _err.Error(), http.StatusBadRequest)
			return
		}
		{{- $formParsed = true }}
		{{- end }}
		if _, exists := _req.Form["{{$p.Name}}"]; exists {
			{{- if IsEnum $p }}
			{{- if isOptional $p.Type }}
			{{- if not (isVarargs $p.Type) }}
			{{$p.Name}} = new({{ TrimPrefix $p.Type "*"}})
			{{- end }}
			{{- end }}
			for _, item := range _req.Form["{{$p.Name}}"] {
				var _{{ $p.Name }} {{ ElementType $p.Type }}
				_{{ $p.Name }}.StringSetter(item)
				{{- if isOptional $p.Type }}
				*{{ $p.Name }} = append(*{{ $p.Name }}, _{{ $p.Name }})
				{{- else }}
				{{ $p.Name }} = append({{ $p.Name }}, _{{ $p.Name }})
				{{- end }}
			}
			{{- else if $p.Type | isSupport }}
			if casted, err := cast.{{$p.Type | castFunc}}E(_req.Form["{{$p.Name}}"]); err != nil {
				http.Error(_writer, err.Error(), http.StatusBadRequest)
				return
			} else {
				{{- if isOptional $p.Type }}
				{{$p.Name}} = &casted
				{{- else }}
				{{$p.Name}} = casted
				{{- end }}
			}
			{{- else }}
			{{- if isOptional $p.Type }}
			_{{$p.Name}} := _req.Form["{{$p.Name}}"]
			{{$p.Name}} = &_{{$p.Name}}
			{{- else }}
			{{$p.Name}} = _req.Form["{{$p.Name}}"]
			{{- end }}
			{{- end }}
		} else {
			if _, exists := _req.Form["{{$p.Name}}[]"]; exists {
				{{- if IsEnum $p }}
				{{- if isOptional $p.Type }}
				{{- if not (isVarargs $p.Type) }}
				{{$p.Name}} = new({{ TrimPrefix $p.Type "*"}})
				{{- end }}
				{{- end }}
				for _, item := range _req.Form["{{$p.Name}}[]"] {
					var _{{ $p.Name }} {{ ElementType $p.Type }}
					_{{ $p.Name }}.StringSetter(item)
					{{- if isOptional $p.Type }}
					*{{ $p.Name }} = append(*{{ $p.Name }}, _{{ $p.Name }})
					{{- else }}
					{{ $p.Name }} = append({{ $p.Name }}, _{{ $p.Name }})
					{{- end }}
				}
				{{- else if $p.Type | isSupport }}
				if casted, err := cast.{{$p.Type | castFunc}}E(_req.Form["{{$p.Name}}[]"]); err != nil {
					http.Error(_writer, err.Error(), http.StatusBadRequest)
					return
				} else {
					{{- if isOptional $p.Type }}
					{{$p.Name}} = &casted
					{{- else }}
					{{$p.Name}} = casted
					{{- end }}
				}
				{{- else }}
				{{- if isOptional $p.Type }}
				_{{$p.Name}} := _req.Form["{{$p.Name}}[]"]
				{{$p.Name}} = &_{{$p.Name}}
				{{- else }}
				{{$p.Name}} = _req.Form["{{$p.Name}}[]"]
				{{- end }}
				{{- end }}
			}{{- if not (isOptional $p.Type) }} else {
				http.Error(_writer, "missing parameter {{$p.Name}}", http.StatusBadRequest)
				return
			}{{- end }}
		}
		{{- else }}
		{{- if not $formParsed }}
		if _err := _req.ParseForm(); _err != nil {
			http.Error(_writer, _err.Error(), http.StatusBadRequest)
			return
		}
		{{- $formParsed = true }}
		{{- end }}
		if _, exists := _req.Form["{{$p.Name}}"]; exists {
			{{- if IsEnum $p }}
			{{- if isOptional $p.Type }}
			{{$p.Name}} = new({{ TrimPrefix $p.Type "*"}})
			{{- end }}
			{{ $p.Name }}.StringSetter(_req.FormValue("{{$p.Name}}"))
			{{- else if $p.Type | isSupport }}
			if casted, err := cast.{{$p.Type | castFunc}}E(_req.FormValue("{{$p.Name}}")); err != nil {
				http.Error(_writer, err.Error(), http.StatusBadRequest)
				return
			} else {
				{{- if isOptional $p.Type }}
				{{$p.Name}} = &casted
				{{- else }}
				{{$p.Name}} = casted
				{{- end }}
			}
			{{- else }}
			{{- if isOptional $p.Type }}
			_{{$p.Name}} := _req.FormValue("{{$p.Name}}")
			{{$p.Name}} = &_{{$p.Name}}
			{{- else }}
			{{$p.Name}} = _req.FormValue("{{$p.Name}}")
			{{- end }}
			{{- end }}
		}{{- if not (isOptional $p.Type) }} else {
			http.Error(_writer, "missing parameter {{$p.Name}}", http.StatusBadRequest)
			return
		}{{- end }}
		{{- end }}
		{{- end }}
		{{ range $i, $r := $m.Results }}{{- if $i}},{{- end}}{{- $r.Name }}{{- end }} = receiver.{{$.Meta.Name | toLowerCamel}}.{{$m.Name}}(
			{{- range $p := $m.Params }}
			{{- if isVarargs $p.Type }}
			*{{ $p.Name }}...,
			{{- else }}
			{{ $p.Name }},
			{{- end }}
			{{- end }}
		)
		{{- range $r := $m.Results }}
			{{- if eq $r.Type "error" }}
				if {{ $r.Name }} != nil {
					if errors.Is({{ $r.Name }}, context.Canceled) {
						http.Error(_writer, {{ $r.Name }}.Error(), http.StatusBadRequest)
					} else if _err, ok := {{ $r.Name }}.(*ddhttp.BizError); ok {
						http.Error(_writer, _err.Error(), _err.StatusCode)
					} else {
						http.Error(_writer, {{ $r.Name }}.Error(), http.StatusInternalServerError)
					}
					return
				}
			{{- end }}
		{{- end }}
		{{- $done := false }}
		{{- range $r := $m.Results }}
			{{- if eq $r.Type "*os.File" }}
				if {{$r.Name}} == nil {
					http.Error(_writer, "No file returned", http.StatusInternalServerError)
					return
				}
				defer {{$r.Name}}.Close()
				var _fi os.FileInfo
				_fi, _err := {{$r.Name}}.Stat()
				if _err != nil {
					http.Error(_writer, _err.Error(), http.StatusInternalServerError)
					return
				}
				_writer.Header().Set("Content-Disposition", "attachment; filename="+_fi.Name())
				_writer.Header().Set("Content-Type", "application/octet-stream")
				_writer.Header().Set("Content-Length", fmt.Sprintf("%d", _fi.Size()))
				io.Copy(_writer, {{$r.Name}})
				{{- $done = true }}	
			{{- end }}
		{{- end }}
		{{- if not $done }}
			if _err := json.NewEncoder(_writer).Encode(struct{
				{{- range $r := $m.Results }}
				{{- if ne $r.Type "error" }}
				{{ $r.Name | toCamel }} {{ $r.Type }} ` + "`" + `json:"{{ $r.Name | convertCase }}{{if $.Omitempty}},omitempty{{end}}"` + "`" + `
				{{- end }}
				{{- end }}
			}{
				{{- range $r := $m.Results }}
				{{- if ne $r.Type "error" }}
				{{ $r.Name | toCamel }}: {{ $r.Name }},
				{{- end }}
				{{- end }}
			}); _err != nil {
				http.Error(_writer, _err.Error(), http.StatusInternalServerError)
				return
			}
		{{- end }}
    }
{{- end }}
`

var importTmpl = `
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	v3 "github.com/unionj-cloud/go-doudou/toolkit/openapi/v3"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"
	{{.ServiceAlias}} "{{.ServicePackage}}"
	"net/http"
	"{{.VoPackage}}"
	"github.com/pkg/errors"
`

var initHttpHandlerImplTmpl = `package httpsrv

import ()

type {{.Meta.Name}}HandlerImpl struct{
	{{.Meta.Name | toLowerCamel}} {{.ServiceAlias}}.{{.Meta.Name}}
}

` + appendHttpHandlerImplTmpl + `

func New{{.Meta.Name}}Handler({{.Meta.Name | toLowerCamel}} {{.ServiceAlias}}.{{.Meta.Name}}) {{.Meta.Name}}Handler {
	return &{{.Meta.Name}}HandlerImpl{
		{{.Meta.Name | toLowerCamel}},
	}
}
`

// GenHttpHandlerImplWithImpl generates http handler implementation
// Parsed value from query string parameters or application/x-www-form-urlencoded form will be string type.
// You may need to convert the type by yourself.
func GenHttpHandlerImplWithImpl(dir string, ic astutils.InterfaceCollector, omitempty bool, caseconvertor func(string) string) {
	var (
		err             error
		modfile         string
		modName         string
		firstLine       string
		handlerimplfile string
		f               *os.File
		modf            *os.File
		tpl             *template.Template
		buf             bytes.Buffer
		httpDir         string
		fi              os.FileInfo
		tmpl            string
		meta            astutils.InterfaceMeta
		importBuf       bytes.Buffer
	)
	httpDir = filepath.Join(dir, "transport/httpsrv")
	if err = os.MkdirAll(httpDir, os.ModePerm); err != nil {
		panic(err)
	}

	handlerimplfile = filepath.Join(httpDir, "handlerimpl.go")
	fi, err = os.Stat(handlerimplfile)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	err = copier.DeepCopy(ic.Interfaces[0], &meta)
	if err != nil {
		panic(err)
	}
	if fi != nil {
		logrus.Warningln("New content will be append to handlerimpl.go file")
		if f, err = os.OpenFile(handlerimplfile, os.O_APPEND, os.ModePerm); err != nil {
			panic(err)
		}
		defer f.Close()
		tmpl = appendHttpHandlerImplTmpl
		unimplementedMethods(&meta, handlerimplfile)
	} else {
		if f, err = os.Create(handlerimplfile); err != nil {
			panic(err)
		}
		defer f.Close()
		tmpl = initHttpHandlerImplTmpl
	}

	modfile = filepath.Join(dir, "go.mod")
	if modf, err = os.Open(modfile); err != nil {
		panic(err)
	}
	reader := bufio.NewReader(modf)
	firstLine, _ = reader.ReadString('\n')
	modName = strings.TrimSpace(strings.TrimPrefix(firstLine, "module"))

	funcMap := make(map[string]interface{})
	funcMap["toLowerCamel"] = strcase.ToLowerCamel
	funcMap["toCamel"] = strcase.ToCamel
	funcMap["contains"] = strings.Contains
	funcMap["isBuiltin"] = v3helper.IsBuiltin
	funcMap["isSupport"] = v3helper.IsSupport
	funcMap["isOptional"] = v3helper.IsOptional
	funcMap["castFunc"] = v3helper.CastFunc
	funcMap["convertCase"] = caseconvertor
	funcMap["isVarargs"] = v3helper.IsVarargs
	funcMap["toSlice"] = v3helper.ToSlice
	funcMap["isSlice"] = v3helper.IsSlice
	funcMap["IsEnum"] = v3helper.IsEnum
	funcMap["TrimPrefix"] = strings.TrimPrefix
	funcMap["ElementType"] = v3helper.ElementType
	if tpl, err = template.New("handlerimpl.go.tmpl").Funcs(funcMap).Parse(tmpl); err != nil {
		panic(err)
	}
	if err = tpl.Execute(&buf, struct {
		ServicePackage string
		ServiceAlias   string
		VoPackage      string
		Meta           astutils.InterfaceMeta
		Omitempty      bool
	}{
		ServicePackage: modName,
		ServiceAlias:   ic.Package.Name,
		VoPackage:      modName + "/vo",
		Meta:           meta,
		Omitempty:      omitempty,
	}); err != nil {
		panic(err)
	}
	original, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	original = append(original, buf.Bytes()...)
	if tpl, err = template.New("himportimpl.go.tmpl").Parse(importTmpl); err != nil {
		panic(err)
	}
	if err = tpl.Execute(&importBuf, struct {
		ServicePackage string
		ServiceAlias   string
		VoPackage      string
	}{
		ServicePackage: modName,
		ServiceAlias:   ic.Package.Name,
		VoPackage:      modName + "/vo",
	}); err != nil {
		panic(err)
	}
	original = astutils.AppendImportStatements(original, importBuf.Bytes())
	//fmt.Println(string(original))
	astutils.FixImport(original, handlerimplfile)
}

func unimplementedMethods(meta *astutils.InterfaceMeta, handlerimplfile string) {
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, handlerimplfile, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := astutils.NewStructCollector(astutils.ExprString)
	ast.Walk(sc, root)
	if handlers, exists := sc.Methods[meta.Name+"HandlerImpl"]; exists {
		var notimplemented []astutils.MethodMeta
		for _, item := range meta.Methods {
			for _, handler := range handlers {
				if item.Name == handler.Name {
					goto L
				}
			}
			notimplemented = append(notimplemented, item)

		L:
		}

		meta.Methods = notimplemented
	}
}
