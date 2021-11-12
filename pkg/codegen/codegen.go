package codegen

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

var marshalTemplate = template.Must(template.New("").Parse(`
func (v {{.StructName}}) MarshalABI(e *abi.Encoder) error {
	var err error
	{{- if .UsesExists }}
	var exists bool
	{{- end -}}
	{{- if .UsesLength }}
	var length int
	{{- end -}}
	{{- range .Fields -}}
	{{- .Code -}}
	{{- end }}
	return err
}
`))

var unmarshalTemplate = template.Must(template.New("").Parse(`
func (v *{{.StructName}}) UnmarshalABI(d *abi.Decoder) error {
	var err error
	{{- if .UsesExists }}
	var exists bool
	{{- end -}}
	{{- if .UsesLength }}
	var length uint
	{{- end -}}
	{{- range .Fields -}}
	{{- .Code -}}
	{{- end }}
	return err
}
`))

var writeTemplate = template.Must(template.New("").Parse(`
	if err = e.Write{{ .Method }}(
		{{- if .IsPtr }}(*{{ end -}}
		v.{{ .Name -}}
		{{- if .IsPtr }}){{ end -}}
		{{- if or .IsSlice .IsArray }}[i]{{ end -}}
	); err != nil {
		return err
	}`,
))

var readTemplate = template.Must(template.New("").Parse(`
	{{ if .IsPtr -}}
	var tmp {{ .Type }}
	{{ end -}}
	if {{ if .IsPtr }}tmp{{ else -}}
		v.{{ .Name }}
		{{- if or .IsSlice .IsArray }}[i]{{ end -}}
	{{- end -}}
	, err = d.Read{{ .Method }}(); err != nil {
		return err
	}
	{{- if .IsPtr }}
	v.{{ .Name }}{{ if or .IsSlice .IsArray }}[i]{{ end }} = &tmp
	{{- end }}`,
))

var writeOtherTemplate = template.Must(template.New("").Parse(`
	if err = {{ if and .IsPtr (or .IsSlice .IsArray) }}(*{{ end -}}
		v.{{ .Name -}}
		{{- if and .IsPtr (or .IsSlice .IsArray) }}){{ end -}}
		{{- if or .IsSlice .IsArray }}[i]{{ end -}}
	.MarshalABI(e); err != nil {
		return err
	}`,
))

var readOtherTemplate = template.Must(template.New("").Parse(`
	if err = {{ if and .IsPtr (or .IsSlice .IsArray) }}(*{{ end -}}
		v.{{ .Name -}}
		{{- if and .IsPtr (or .IsSlice .IsArray) }}){{ end -}}
		{{- if or .IsSlice .IsArray }}[i]{{ end -}}
	.UnmarshalABI(d); err != nil {
		return err
	}`,
))

var writeSliceTemplate = template.Must(template.New("").Parse(`
	length = len({{ if .IsPtr }}*{{ end }}v.{{ .Name }})
    if err = e.WriteVaruint(uint(length)); err != nil {
		return err
	}
	for i := 0; i < length; i++ {
		{{- .Code }}
	}`,
))

var readSliceTemplate = template.Must(template.New("").Parse(`
	if length, err = d.ReadVaruint(); err != nil {
		return err
	}
	v.{{ .Name }} = make([]{{ .Type }}, length)
	for i := 0; i < int(length); i++ {
		{{- .Code }}
	}`,
))

var writeArrayTemplate = template.Must(template.New("").Parse(`
	for i := 0; i < {{ .ArrayLen }}; i++ {
		{{- .Code }}
	}`,
))

var readArrayTemplate = template.Must(template.New("").Parse(`
	for i := 0; i < {{ .ArrayLen }}; i++ {
		{{- .Code }}
	}`,
))

var writeOptionalTemplate = template.Must(template.New("").Parse(`
	exists = v.{{ .Name }} != nil
	{{- if .IsOptional }}
	if err = e.WriteBool(exists); err != nil {
		return err
	}
	{{- end }}
	if exists {
		{{- .Code }}
	}
	{{- if not .IsOptional }} else {
		return errors.New("encountered nil for non-optional field: {{ .Name }}")
	}
	{{- end -}}
`))

var readOptionalTemplate = template.Must(template.New("").Parse(`
	if exists, err = d.ReadBool(); err != nil {
		return err
	}
	if exists {
		{{- .Code }}
	}`,
))

var writeMapTemplate = template.Must(template.New("").Parse(`
	length = len({{ if .IsPtr }}*{{ end }}v.{{ .Name }})
    if err = e.WriteVaruint(uint(length)); err != nil {
		return err
	}
	for i := 0; i < length; i++ {
		{{- .Code }}
	}`,
))

type ctx struct {
	StructName string
	UsesExists bool
	UsesLength bool
	Fields     []*field
}

type field struct {
	Name       string
	Type       string
	IsPtr      bool
	IsOptional bool
	IsSlice    bool
	IsArray    bool
	ArrayLen   int
	Code       string
	Method     string
}

func GenUnmarshalFn(structType interface{}) string {
	v := reflect.ValueOf(structType)
	if v.Kind() != reflect.Struct {
		panic("only structs can be synthesized")
	}
	t := v.Type()
	ctx := ctx{
		StructName: t.Name(),
	}
	for i := 0; i < t.NumField(); i++ {
		f := field{
			Name: t.Field(i).Name,
			Type: t.Field(i).Type.String(),
		}
		if t.Field(i).Type.Kind() == reflect.Map {
			// TODO: handle maps
		}
		if f.Type[0] == '*' {
			f.IsPtr = true
			f.Type = f.Type[1:]
		}
		if strings.HasPrefix(f.Type, "[]") {
			f.IsSlice = true
			f.Type = f.Type[2:]
		} else if strings.HasPrefix(f.Type, "[") {
			f.IsArray = true
			for j := 1; j < len(f.Type); j++ {
				if f.Type[j] == ']' {
					f.ArrayLen, _ = strconv.Atoi(f.Type[1:j])
					f.Type = f.Type[j+1:]
					break
				}
			}
		}

		if t.Field(i).Tag.Get("eosio") == "optional" {
			if !f.IsPtr {
				panic("cant generate code for non-pointer optionals")
			}
			f.IsOptional = true
		}
		ctx.Fields = append(ctx.Fields, &f)
	}

	for _, f := range ctx.Fields {
		switch f.Type {
		case "string":
			f.Method = "String"
		case "uint64":
			f.Method = "Uint64"
		case "uint32":
			f.Method = "Uint32"
		case "uint16":
			f.Method = "Uint16"
		case "uint8":
			f.Method = "Uint8"
		case "int64":
			f.Method = "Int64"
		case "int32":
			f.Method = "Int32"
		case "int16":
			f.Method = "Int16"
		case "int8":
			f.Method = "Int8"
		case "bool":
			f.Method = "Bool"
		case "[]uint8", "[]byte":
			f.Method = "Bytes"
		case "int":
			f.Method = "Varint"
		case "uint":
			f.Method = "Varuint"
		case "float32":
			f.Method = "Float32"
		case "float64":
			f.Method = "Float64"
		}
	}

	// marshal pass

	for _, f := range ctx.Fields {
		tpl := writeTemplate
		if f.Method == "" {
			tpl = writeOtherTemplate
		}
		f.Code = render(tpl, f)
		if f.IsSlice {
			ctx.UsesLength = true
			f.Code = indent(f.Code)
			f.Code = render(writeSliceTemplate, f)
		}
		if f.IsArray {
			ctx.UsesLength = true
			f.Code = indent(f.Code)
			f.Code = render(writeArrayTemplate, f)
		}
		if f.IsOptional || f.IsPtr {
			ctx.UsesExists = true
			f.Code = indent(f.Code)
			f.Code = render(writeOptionalTemplate, f)
		}
	}

	rv := render(marshalTemplate, ctx)

	// unmarshal pass

	ctx.UsesLength = false
	ctx.UsesExists = false

	for _, f := range ctx.Fields {
		tpl := readTemplate
		if f.Method == "" {
			tpl = readOtherTemplate
		}
		f.Code = render(tpl, f)
		if f.IsSlice {
			ctx.UsesLength = true
			f.Code = indent(f.Code)
			f.Code = render(readSliceTemplate, f)
		}
		if f.IsArray {
			f.Code = indent(f.Code)
			f.Code = render(readArrayTemplate, f)
		}
		if f.IsOptional {
			ctx.UsesExists = true
			f.Code = indent(f.Code)
			f.Code = render(readOptionalTemplate, f)
		}
	}

	rv += "\n" + render(unmarshalTemplate, ctx)

	return rv

}

func render(t *template.Template, v interface{}) string {
	buf := new(bytes.Buffer)
	err := t.Execute(buf, v)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func indent(s string) string {
	return "	" + strings.Replace(s, "\n", "\n	", -1)
}
