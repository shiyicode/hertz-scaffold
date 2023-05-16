package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/three-body/hertz-scaffold/biz/dal"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	outPath := flag.String("o", "biz/dal/query", "output path")
	flag.Parse()

	fmt.Printf("Use outPath: %s\n", *outPath)
	// 全局配置
	g := gen.NewGenerator(gen.Config{
		OutPath:           *outPath,
		Mode:              gen.WithDefaultQuery,
		OutFile:           "",
		ModelPkgPath:      "",
		WithUnitTest:      true,
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: false,
		FieldWithTypeTag:  false,
	})

	// 设置数据库连接
	err := dal.InitMySQL()
	if err != nil {
		panic(err)
	}
	g.UseDB(dal.DB)

	dataMap := map[string]func(columnType gorm.ColumnType) (dataType string){
		"tinyint": func(columnType gorm.ColumnType) (dataType string) {
			fmt.Println(columnType.Name())
			if strings.HasPrefix(columnType.Name(), "tinyint(1)") {
				return "bool"
			}
			return "int"
		},
		"smallint":  func(columnType gorm.ColumnType) (dataType string) { return "int" },
		"mediumint": func(columnType gorm.ColumnType) (dataType string) { return "int" },
		"int":       func(columnType gorm.ColumnType) (dataType string) { return "int" },
		"bigint":    func(columnType gorm.ColumnType) (dataType string) { return "int64" },
	}
	g.WithDataTypeMap(dataMap)

	tables, err := dal.DB.Migrator().GetTables()
	if err != nil {
		panic(err)
	}

	for _, table := range tables {
		fieldOpts := append(AllTablesModelOpt(), GetTableModelOpt(table)...)
		tableModel := g.GenerateModel(table, fieldOpts...)
		g.ApplyBasic(tableModel)
	}

	g.Execute()
}
