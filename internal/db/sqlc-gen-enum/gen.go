package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var funcMap = template.FuncMap{
	"fmtConstName":  fmtConstName,
	"fmtConstValue": fmtConstValue,
}

func run(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	tmpl, err := template.New("enums.tmpl").Funcs(funcMap).ParseFS(templates, "enums.tmpl")
	if err != nil {
		return nil, err
	}
	tmplData, err := makeTmplData(req)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if err = tmpl.Execute(w, *tmplData); err != nil {
		return nil, err
	}
	w.Flush()
	src, err := format.Source(b.Bytes())
	if err != nil {
		return nil, err
	}

	resp := &plugin.GenerateResponse{
		Files: []*plugin.File{
			{
				Name:     "enums.go",
				Contents: src,
			},
		},
	}

	return resp, nil
}

type tmplData struct {
	Package    string
	TableNames []string
}

func makeTmplData(req *plugin.GenerateRequest) (*tmplData, error) {
	opts, err := parseOptions(req)
	if err != nil {
		return nil, err
	}

	catalog := req.Catalog
	for _, schema := range catalog.Schemas {
		if schema.Name != catalog.DefaultSchema {
			continue
		}
		tableNames := make([]string, len(schema.Tables))
		for i, table := range schema.Tables {
			tableNames[i] = table.Rel.Name
		}

		data := &tmplData{
			Package:    opts.Package,
			TableNames: tableNames,
		}
		return data, nil
	}
	return nil, fmt.Errorf("could not find default schema '%s' in schema list (len=%d)", catalog.DefaultSchema, len(catalog.Schemas))
}

var titleCaser = cases.Title(language.English)

func fmtConstName(tableName string) string {
	return strings.ReplaceAll(fmtConstValue(tableName), " ", "")
}

func fmtConstValue(tableName string) string {
	noUnderscores := strings.ReplaceAll(tableName, "_", " ")
	titled := titleCaser.String(strings.ToLower(noUnderscores))
	return titled
}

type options struct {
	Package string `json:"package" yaml:"package"`
}

func parseOptions(req *plugin.GenerateRequest) (*options, error) {
	opts := &options{}
	if err := json.Unmarshal(req.Settings.Codegen.Options, opts); err != nil {
		return nil, fmt.Errorf("parse codegen options: %w", err)
	}

	return opts, nil
}
