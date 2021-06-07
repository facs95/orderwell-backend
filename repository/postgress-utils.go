package repository

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var camel = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

func getTypeName(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr || t.Kind() == reflect.Array {
		return underscore(t.Elem().Name())
	} else {
		return underscore(t.Name())
	}
}

func getSchemaTableName(table interface{}, tenantId string) string {
	tableName := getTypeName(table)
	return fmt.Sprintf("%v.%v", tenantId, tableName)
}

func underscore(s string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToLower(strings.Join(a, "_"))
}
