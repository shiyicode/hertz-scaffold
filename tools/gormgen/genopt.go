package main

import (
	"gorm.io/gen"
	"gorm.io/gen/field"
)

func AllTablesModelOpt() []gen.ModelOpt {
	autoCreatedField := gen.FieldGORMTag("created_at", func(tag field.GormTag) field.GormTag {
		tag.Set("autoCreateTime", "")
		return tag
	})
	autoUpdatedField := gen.FieldGORMTag("updated_at", func(tag field.GormTag) field.GormTag {
		tag.Set("autoUpdateTime", "")
		return tag
	})
	softDeleteFieldType := gen.FieldType("deleted_at", "soft_delete.DeletedAt")
	softDeleteFieldTag := gen.FieldGORMTag("deleted_at", func(tag field.GormTag) field.GormTag {
		tag.Set("softDelete", "milli")
		return tag
	})

	return []gen.ModelOpt{autoCreatedField, autoUpdatedField, softDeleteFieldType, softDeleteFieldTag}
}

func GetTableModelOpt(table string) []gen.ModelOpt {
	switch table {
	case "user":
		extraField := gen.FieldType("extra", "datatypes.JSONType[map[string]string]")
		return []gen.ModelOpt{extraField}
	}
	return []gen.ModelOpt{}
}
