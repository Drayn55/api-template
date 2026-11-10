package audit

import (
	"reflect"
	"strings"
	"time"
)

func ChangedFields(oldObj, newObj any) []string {
	oldVal := reflect.ValueOf(oldObj)
	newVal := reflect.ValueOf(newObj)
	if !oldVal.IsValid() || !newVal.IsValid() {
		return nil
	}
	return changedFields(reflect.Indirect(oldVal), reflect.Indirect(newVal), "")
}

func changedFields(oldVal, newVal reflect.Value, prefix string) []string {
	if !oldVal.IsValid() || !newVal.IsValid() {
		return nil
	}

	oldVal = reflect.Indirect(oldVal)
	newVal = reflect.Indirect(newVal)

	if oldVal.Kind() != reflect.Struct || newVal.Kind() != reflect.Struct {
		return nil
	}

	var fields []string
	typ := oldVal.Type()

	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if sf.PkgPath != "" {
			continue
		}
		if tag := sf.Tag.Get("audit"); tag == "-" {
			continue
		}

		name := fieldName(sf)
		fullName := name
		if prefix != "" {
			fullName = prefix + "." + name
		}

		oldField := reflect.Indirect(oldVal.Field(i))
		newField := reflect.Indirect(newVal.Field(i))

		if isNestedStruct(sf.Type) {
			fields = append(fields, changedFields(oldField, newField, fullName)...)
			continue
		}

		if !valuesEqual(oldField, newField) {
			fields = append(fields, fullName)
		}
	}

	return fields
}

func isNestedStruct(t reflect.Type) bool {
	t = indirectType(t)
	if t == reflect.TypeOf(time.Time{}) {
		return false
	}
	return t.Kind() == reflect.Struct
}

func indirectType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func fieldName(sf reflect.StructField) string {
	if jsonTag := sf.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
		if name, _, ok := strings.Cut(jsonTag, ","); ok && name != "" {
			return name
		}
		if jsonTag != "" {
			return jsonTag
		}
	}
	return toSnakeCase(sf.Name)
}

func toSnakeCase(s string) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('_')
		}
		b.WriteRune(r)
	}
	return strings.ToLower(b.String())
}

func valuesEqual(a, b reflect.Value) bool {
	if !a.IsValid() && !b.IsValid() {
		return true
	}
	if !a.IsValid() || !b.IsValid() {
		return false
	}
	if a.Type() != b.Type() {
		return reflect.DeepEqual(a.Interface(), b.Interface())
	}
	if a.Kind() == reflect.Float32 || a.Kind() == reflect.Float64 {
		return a.Float() == b.Float()
	}
	return reflect.DeepEqual(a.Interface(), b.Interface())
}
