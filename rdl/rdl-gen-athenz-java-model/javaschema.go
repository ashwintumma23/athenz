// Copyright 2015 Yahoo Inc.
//           2020 Verizon Media Modified to generate java client code for Athenz Clients
// Licensed under the terms of the Apache version 2.0 license. See LICENSE file for terms.

package main

import (
	"fmt"
	"io"
	"sort"
	"text/template"

	"github.com/ardielle/ardielle-go/rdl"
)

func GenerateSchema(schema *rdl.Schema, cName string, writer io.Writer, ns string, banner string) error {
	reg := rdl.NewTypeRegistry(schema)
	funcMap := template.FuncMap{
		"header": func() string {
			s := fmt.Sprintf("//\n// This file generated by %s. Do not modify!\n//\n\n", banner)
			s2 := GenerationPackage(schema, ns)
			if s2 != "" {
				s = s + "package " + s2 + ";\n"
			}
			if ns != "com.yahoo.rdl" {
				s = s + "import com.yahoo.rdl.*;\n"
			}
			return s
		},
		"version": func() string {
			if schema.Version != nil {
				return fmt.Sprintf("        sb.version(%d);", *schema.Version)
			} else {
				return ""
			}
		},
		"name": func() string { return string(schema.Name) },
		"namespace": func() string {
			if schema.Namespace != "" {
				return fmt.Sprintf("\n        sb.namespace(%q);", schema.Namespace)
			} else {
				return ""
			}
		},
		"comment": func() string {
			if schema.Comment != "" {
				return fmt.Sprintf("\n        sb.comment(%q);", schema.Comment)
			} else {
				return ""
			}
		},
		"cname":       func() string { return cName },
		"typeDef":     func(t *rdl.Type) string { return javaGenerateTypeConstructor(reg, t) },
		"resourceDef": func(r *rdl.Resource) string { return javaGenerateResourceConstructor(reg, r) },
	}
	t := template.Must(template.New("util").Funcs(funcMap).Parse(javaSchemaTemplate))
	return t.Execute(writer, schema)
}

func javaGenerateTypeConstructor(reg rdl.TypeRegistry, t *rdl.Type) string {
	switch t.Variant {
	case rdl.TypeVariantBaseType:
		return fmt.Sprintf("    //s.type(BaseType) NYI - %v", t.BaseType)
	case rdl.TypeVariantStructTypeDef:
		return javaGenerateStructTypeConstructor(reg, t.StructTypeDef)
	case rdl.TypeVariantEnumTypeDef:
		return javaGenerateEnumTypeConstructor(t.EnumTypeDef)
	case rdl.TypeVariantStringTypeDef:
		return javaGenerateStringTypeConstructor(t.StringTypeDef)
	case rdl.TypeVariantMapTypeDef:
		return "    //s.type(MapTypeDef) NYI - " + string(t.MapTypeDef.Name)
	case rdl.TypeVariantArrayTypeDef:
		return "    //s.type(ArrayTypeDef) NYI - " + string(t.ArrayTypeDef.Name)
	case rdl.TypeVariantBytesTypeDef:
		return "    //s.type(BytesTypeDef) NYI - " + string(t.BytesTypeDef.Name)
	case rdl.TypeVariantNumberTypeDef:
		return javaGenerateNumberTypeConstructor(t.NumberTypeDef)
	case rdl.TypeVariantUnionTypeDef:
		return javaGenerateUnionTypeConstructor(t.UnionTypeDef)
	case rdl.TypeVariantAliasTypeDef:
		//java doesn't support aliases, all references should have been rewritten
		return ""
	}
	return fmt.Sprintf("    //s.type(%v);", t.Variant)
}

func javaGenerateStringTypeConstructor(t *rdl.StringTypeDef) string {
	s := fmt.Sprintf("    sb.stringType(%q)", t.Name)
	if t.Comment != "" {
		s += fmt.Sprintf("\n            .comment(%q)", t.Comment)
	}
	if t.Pattern != "" {
		s += fmt.Sprintf("\n            .pattern(%q)", t.Pattern)
	}
	if t.MinSize != nil {
		s += fmt.Sprintf("\n            .minSize(%d)", *t.MinSize)
	}
	if t.MaxSize != nil {
		s += fmt.Sprintf("\n            .maxSize(%d)", *t.MaxSize)
	}
	return s + ";"
}

func numberValueString(n *rdl.Number) string {
	switch n.Variant {
	case rdl.NumberVariantInt32:
		return fmt.Sprint(*n.Int32)
	case rdl.NumberVariantInt64:
		return fmt.Sprint(*n.Int64)
	case rdl.NumberVariantFloat64:
		return fmt.Sprint(*n.Float64)
	}
	return "0"
}

func javaGenerateNumberTypeConstructor(t *rdl.NumberTypeDef) string {
	stype := t.Type //fix: fix recursive numeric types.
	s := fmt.Sprintf("    sb.numberType(%q, %q)", t.Name, stype)
	if t.Comment != "" {
		s += fmt.Sprintf("\n            .comment(%q)", t.Comment)
	}
	if t.Min != nil {
		s += fmt.Sprintf("\n            .min(%s)", numberValueString(t.Min))
	}
	if t.Max != nil {
		s += fmt.Sprintf("\n            .max(%s)", numberValueString(t.Max))
	}
	return s + ";"
}

func javaGenerateEnumTypeConstructor(t *rdl.EnumTypeDef) string {
	s := fmt.Sprintf("    sb.enumType(%q)", t.Name)
	if t.Comment != "" {
		s += fmt.Sprintf("\n            .comment(%q)", t.Comment)
	}
	for _, e := range t.Elements {
		s += fmt.Sprintf("\n            .element(%q)", e.Symbol)
	}
	return s + ";"
}

func javaGenerateUnionTypeConstructor(t *rdl.UnionTypeDef) string {
	s := fmt.Sprintf("    sb.unionType(%q)", t.Name)
	if t.Comment != "" {
		s += fmt.Sprintf("\n            .comment(%q)", t.Comment)
	}
	for _, v := range t.Variants {
		s += fmt.Sprintf("\n            .variant(%q)", v)
	}
	return s + ";"
}

func javaLiteral(otype *rdl.Type, oval interface{}) string {
	switch otype.Variant {
	case rdl.TypeVariantStringTypeDef:
		s := fmt.Sprint(oval)
		return fmt.Sprintf("%q", s)
	case rdl.TypeVariantNumberTypeDef:
		return fmt.Sprint(oval)
	case rdl.TypeVariantAliasTypeDef:
		//fix this. Should be passing base type in to handle this.
	case rdl.TypeVariantBaseType:
		switch *otype.BaseType {
		case rdl.BaseTypeString, rdl.BaseTypeTimestamp, rdl.BaseTypeUUID, rdl.BaseTypeSymbol:
			s := fmt.Sprint(oval)
			return fmt.Sprintf("%q", s)
		case rdl.BaseTypeBool:
			return fmt.Sprintf("%v", oval)
		case rdl.BaseTypeInt8, rdl.BaseTypeInt16, rdl.BaseTypeInt32, rdl.BaseTypeInt64:
			return fmt.Sprintf("%v", oval)
		case rdl.BaseTypeFloat32, rdl.BaseTypeFloat64:
			return fmt.Sprintf("%v", oval)
		}
	case rdl.TypeVariantEnumTypeDef:
		return fmt.Sprintf("%v.%v", otype.EnumTypeDef.Name, oval)
	}
	return "null"
}

func javaGenerateStructTypeConstructor(reg rdl.TypeRegistry, t *rdl.StructTypeDef) string {
	s := fmt.Sprintf("    sb.structType(%q", t.Name)
	if t.Type != "Struct" {
		s += fmt.Sprintf(", %q)", t.Type)
	} else {
		s += ")"
	}
	if t.Comment != "" {
		s += fmt.Sprintf("\n            .comment(%q)", t.Comment)
	}
	for _, f := range t.Fields {
		if f.Keys != "" {
			fkeys := string(f.Keys)   //javaType(reg, f.Keys, false, "", "")
			fitems := string(f.Items) //javaType(reg, f.Items, false, "", "")
			s += fmt.Sprintf("\n            .mapField(%q, %q, %q, %v, %q)", f.Name, fkeys, fitems, f.Optional, f.Comment)
		} else if f.Items != "" {
			fitems := string(f.Items) //javaType(reg, f.Items, false, "", "")
			s += fmt.Sprintf("\n            .arrayField(%q, %q, %v, %q)", f.Name, fitems, f.Optional, f.Comment)
		} else {
			ftype := string(f.Type) //javaType(reg, f.Type, f.Optional, "", "")
			if f.Default != nil {
				ft := reg.FindType(f.Type)
				ss := "null"
				if ft != nil {
					ss = javaLiteral(ft, f.Default)
				}
				s += fmt.Sprintf("\n            .field(%q, %q, %v, %q, %s)", f.Name, ftype, f.Optional, f.Comment, ss)
			} else {
				s += fmt.Sprintf("\n            .field(%q, %q, %v, %q)", f.Name, ftype, f.Optional, f.Comment)
			}
		}
	}
	return s + ";"
}

func javaGenerateResourceConstructor(reg rdl.TypeRegistry, rez *rdl.Resource) string {
	rTypeName := rez.Type
	if rez.Method == "PUT" || rez.Method == "POST" {
		for _, ri := range rez.Inputs {
			if !ri.PathParam && ri.QueryParam == "" && ri.Header == "" {
				rTypeName = ri.Type
				break
			}
		}
	}
	s := fmt.Sprintf("    sb.resource(%q, %q, %q)", rTypeName, rez.Method, rez.Path)
	if rez.Comment != "" {
		s += fmt.Sprintf("\n            .comment(%q)", rez.Comment)
	}
	if rez.Name != "" {
		s += fmt.Sprintf("\n            .name(%q)", rez.Name)
	}
	for _, ri := range rez.Inputs {
		def := "null"
		if ri.Default != nil {
			ft := reg.FindType(ri.Type)
			if ft != nil {
				def = javaLiteral(ft, ri.Default)
			}
		}
		if ri.PathParam {
			s += fmt.Sprintf("\n            .pathParam(%q, %q, %q)", ri.Name, ri.Type, ri.Comment)
		} else if ri.QueryParam != "" {
			s += fmt.Sprintf("\n            .queryParam(%q, %q, %q, %s, %q)", ri.QueryParam, ri.Name, ri.Type, def, ri.Comment)
		} else if ri.Header != "" {
			s += fmt.Sprintf("\n            .headerParam(%q, %q, %q, %s, %q)", ri.Header, ri.Name, ri.Type, def, ri.Comment)
		} else {
			s += fmt.Sprintf("\n            .input(%q, %q, %q)", ri.Name, ri.Type, ri.Comment)
		}
	}
	for _, ro := range rez.Outputs {
		s += fmt.Sprintf("\n            .output(%q, %q, %q, %q)", ro.Header, ro.Name, ro.Type, ro.Comment)
	}
	if rez.Auth != nil {
		if rez.Auth.Domain != "" {
			s += fmt.Sprintf("\n            .auth(%q, %q, %v, %q)", rez.Auth.Action, rez.Auth.Resource, rez.Auth.Authenticate, rez.Auth.Domain)
		} else if rez.Auth.Authenticate {
			s += fmt.Sprintf("\n            .auth(%q, %q, true)", rez.Auth.Action, rez.Auth.Resource)
		} else {
			s += fmt.Sprintf("\n            .auth(%q, %q)", rez.Auth.Action, rez.Auth.Resource)
		}
	}
	s += fmt.Sprintf("\n            .expected(%q)", rez.Expected)
	//build a sorted order for the exceptions, to make them predictable. Go randomizes the order otherwise.
	var syms []string
	for sym := range rez.Exceptions {
		syms = append(syms, sym)
	}
	sort.Strings(syms)
	for _, sym := range syms {
		re := rez.Exceptions[sym]
		s += fmt.Sprintf("\n            .exception(%q, %q, %q)\n", sym, re.Type, re.Comment)
	}
	if rez.Async != nil {
		if *rez.Async {
			s += fmt.Sprintf("\n            .async()\n")
		}
	}
	return s + ";"
}

const javaSchemaTemplate = `{{header}}
public class {{cname}} {

    private final static Schema INSTANCE = build();
    public static Schema instance() {
        return INSTANCE;
    }

    private static Schema build() {
        SchemaBuilder sb = new SchemaBuilder("{{name}}");
{{version}}{{namespace}}{{comment}}
{{range .Types}}
    {{typeDef .}}
{{end}}
{{range .Resources}}
    {{resourceDef .}}
{{end}}

        return sb.build();
    }

}
`
