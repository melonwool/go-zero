package gen

import (
	"strings"

	"github.com/melonwool/go-zero/tools/goctl/model/sql/template"
	"github.com/melonwool/go-zero/tools/goctl/util"
	"github.com/melonwool/go-zero/tools/goctl/util/stringx"
)

func genUpdate(table Table, withCache bool) (string, string, error) {
	expressionValues := make([]string, 0)
	for _, field := range table.Fields {
		camel := field.Name.ToCamel()
		if camel == "CreateTime" || camel == "UpdateTime" {
			continue
		}

		if field.Name.Source() == table.PrimaryKey.Name.Source() {
			continue
		}

		expressionValues = append(expressionValues, "data."+camel)
	}

	expressionValues = append(expressionValues, "data."+table.PrimaryKey.Name.ToCamel())
	camelTableName := table.Name.ToCamel()
	text, err := util.LoadTemplate(category, updateTemplateFile, template.Update)
	if err != nil {
		return "", "", err
	}

	output, err := util.With("update").
		Parse(text).
		Execute(map[string]interface{}{
			"withCache":             withCache,
			"upperStartCamelObject": camelTableName,
			"primaryCacheKey":       table.PrimaryCacheKey.DataKeyExpression,
			"primaryKeyVariable":    table.PrimaryCacheKey.KeyLeft,
			"lowerStartCamelObject": stringx.From(camelTableName).Untitle(),
			"originalPrimaryKey":    wrapWithRawString(table.PrimaryKey.Name.Source()),
			"expressionValues":      strings.Join(expressionValues, ", "),
		})
	if err != nil {
		return "", "", nil
	}

	// update interface method
	text, err = util.LoadTemplate(category, updateMethodTemplateFile, template.UpdateMethod)
	if err != nil {
		return "", "", err
	}

	updateMethodOutput, err := util.With("updateMethod").
		Parse(text).
		Execute(map[string]interface{}{
			"upperStartCamelObject": camelTableName,
		})
	if err != nil {
		return "", "", nil
	}

	return output.String(), updateMethodOutput.String(), nil
}
