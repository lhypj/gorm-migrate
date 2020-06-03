package core

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	INITMigrations = ""
	MIGRATIONSplit = "_"
	MIGRATIONPath  = "/migrations/"
	INITContent    = "package migrations\n\ntype Migrations struct{}\n"
)
const (
	TIMEFormat       = "20060102150405"
	NORMALTimeFormat = "2006-01-02 15:04:05"
)

func (*Migrate) rootPath() (string, error) {
	return os.Getwd()
}

func (m *Migrate) migrationsPath() string {
	if path, err := m.rootPath(); err != nil {
		panic(err)
	} else {
		return path + m.ModelsRelativePath + MIGRATIONPath
	}

}

func (m *Migrate) MigrationsInit() {
	// migrate models
	m.DB.AutoMigrate(&OrmMigrations{})

	mPath := m.migrationsPath()
	_, err := os.Stat(mPath + "init.go")
	if os.IsNotExist(err) {
		m.write(INITContent, "init")
	}
}

func (m *Migrate) MakeMigrations(tables ...interface{}) {
	defer m.handleErr()()
	var content string
	tableFromFile, head := m.genTablesFromMigrationFiles()
	tableFromObj := m.genTableFromObject(tables...)
	fn := m.genMigrationFileName(head)
	pre := m.MigrationsPre(fn)
	end := m.MigrationsEnd(fn, []string{head})
	for _, table := range tables {
		name := m.DB.NewScope(table).TableName()
		content += m.genMigrationFileContent(tableFromFile[name], tableFromObj[name])
	}
	if content == "" {
		fmt.Println("No migrations need to make")
	} else {
		m.write(pre+content+end, fn)
	}
}

func (m *Migrate) MigrationsPre(fn string) string {
	return fmt.Sprintf("package migrations\n\nimport %v\n\nfunc (*Migrations) Migration_%v() *core.Operations {\n\tvar ops []*core.Operation\n\tops = append(ops,\n",
		m.quoteStrToMigrations(m.PackagePath), fn)
}

func (m *Migrate) MigrationsEnd(fn string, latest []string) string {
	return fmt.Sprintf("\t)\n\treturn &core.Operations{Revision: %v, DownRevision: %v, Operations: ops}\n}",
		m.quoteStrToMigrations(fn), m.quoteStrListToMigrations(latest))
}

func (m *Migrate) diffFields(exists *Table, target *Table) string {
	var content string
	oldField := make(map[string]*Field)
	newField := make(map[string]*Field)
	fields := make(map[string]int)
	tableName := exists.Name
	for _, field := range exists.Fields {
		fields[field.Name] = 1
		oldField[field.Name] = field
	}
	for _, field := range target.Fields {
		fields[field.Name] = 1
		newField[field.Name] = field
	}
	for fn := range fields {
		o := oldField[fn]
		n := newField[fn]
		if o == nil && n != nil {
			content += m.migrationAddFieldContent(tableName, n.Name, n.Type, n.IsPrimary)
		}
		if o != nil && n == nil {
			content += m.migrationDeleteFieldContent(tableName, o.Name, o.Type)
		}
		if o != nil && n != nil && o.Type != n.Type {
			content += m.migrationAlterFieldContent(tableName, o.Name, o.Type, n.Type)
		}
	}
	return content
}

func (m *Migrate) diffIndexes(exists *Table, target *Table, unique bool) string {
	var content string
	oldIndex := make(map[string]*Index)
	newIndex := make(map[string]*Index)
	indexes := make(map[string]int)
	tableName := exists.Name

	eIndexes := exists.Indexes
	tIndexes := target.Indexes
	addAction := ADDIndexStr
	deleteAction := DELETEIndexStr
	if unique {
		eIndexes = exists.UniqueIndexes
		tIndexes = target.UniqueIndexes
		addAction = ADDUniqueIndexStr
		deleteAction = DELETEUniqueIndexStr
	}

	for _, index := range eIndexes {
		indexes[index.Name] = 1
		oldIndex[index.Name] = index
	}
	for _, index := range tIndexes {
		indexes[index.Name] = 1
		newIndex[index.Name] = index
	}
	for idxName := range indexes {
		o := oldIndex[idxName]
		n := newIndex[idxName]
		if o == nil && n != nil {
			content += m.migrationIndexContent(tableName, addAction, n.Name, n.FieldName)
		}
		if o != nil && n == nil {
			content += m.migrationIndexContent(tableName, deleteAction, o.Name, o.FieldName)
		}
		if o != nil && n != nil && !reflect.DeepEqual(o.FieldName, n.FieldName) {
			content += m.migrationIndexContent(tableName, deleteAction, o.Name, o.FieldName)
			content += m.migrationIndexContent(tableName, addAction, n.Name, n.FieldName)
		}
	}
	return content
}

func (m *Migrate) diff(exists *Table, target *Table) string {
	var content string
	content += m.diffFields(exists, target)
	content += m.diffIndexes(exists, target, false)
	content += m.diffIndexes(exists, target, true)
	return content
}

func (m *Migrate) migrationAddFieldContent(tableName, fieldName, fileType string, isPrimary bool) string {
	primaryStr := ""
	if isPrimary {
		primaryStr = ", IsPrimary: true"
	}
	return fmt.Sprintf("\t\t&core.Operation{Action: core.ADDField, TableName: %v, ColumnName: %v, Type: %v%v},\n",
		m.quoteStrToMigrations(tableName), m.quoteStrToMigrations(fieldName), m.quoteStrToMigrations(fileType), primaryStr)
}

func (m *Migrate) migrationDeleteFieldContent(tableName, fieldName, fileType string) string {
	return fmt.Sprintf("\t\t&core.Operation{Action: core.DELETEField, TableName: %v, ColumnName: %v, Type: %v},\n",
		m.quoteStrToMigrations(tableName), m.quoteStrToMigrations(fieldName), m.quoteStrToMigrations(fileType))
}

func (m *Migrate) migrationAlterFieldContent(tableName, fieldName, old, new string) string {
	return fmt.Sprintf("\t\t&core.Operation{Action: core.ALTERField, TableName: %v, ColumnName: %v, Type: %v, TypeNew: %v},\n",
		m.quoteStrToMigrations(tableName), m.quoteStrToMigrations(fieldName), m.quoteStrToMigrations(old), m.quoteStrToMigrations(new))
}

func (m *Migrate) migrationIndexContent(tableName, action, indexName string, indexFields []string) string {
	return fmt.Sprintf("\t\t&core.Operation{TableName: %v, Action: %v, IndexName: %v, IndexFieldNames: %v},\n",
		m.quoteStrToMigrations(tableName), action, m.quoteStrToMigrations(indexName), m.quoteStrListToMigrations(indexFields))
}

func (m *Migrate) quoteStrToMigrations(str string) string {
	return fmt.Sprintf("\"%v\"", str)
}

func (m *Migrate) quoteStrListToMigrations(str []string) string {
	var ret []string
	for _, s := range str {
		ret = append(ret, m.quoteStrToMigrations(s))
	}
	return fmt.Sprintf("[]string{%v}", strings.Join(ret, ", "))
}

func (m *Migrate) genMigrationFileContent(exists *Table, target *Table) string {
	var content string
	if target == nil {
		if exists != nil {
			content += fmt.Sprintf("&core.Operation{Action: core.DELETETable, TableName: %v},",
				m.quoteStrToMigrations(exists.Name))
		}
		return content
	}
	if exists == nil {
		content = fmt.Sprintf("\t\t&core.Operation{Action: core.ADDTable, TableName: %v},\n", m.quoteStrToMigrations(target.Name))
		for _, field := range target.Fields {
			content += m.migrationAddFieldContent(target.Name, field.Name, field.Type, field.IsPrimary)
		}
		for _, index := range target.Indexes {
			indexName := index.Name
			content += fmt.Sprintf(
				"\t\t&core.Operation{Action: core.ADDIndex, TableName: %v, IndexName: %v, IndexFieldNames: %v},\n",
				m.quoteStrToMigrations(target.Name), m.quoteStrToMigrations(indexName), m.quoteStrListToMigrations(index.FieldName))
		}

		for _, index := range target.UniqueIndexes {
			indexName := index.Name
			content += fmt.Sprintf(
				"\t\t&core.Operation{Action: core.ADDUniqueIndex, TableName: %v, UniqueIndexName: %v, UniqueFieldNames: %v},\n",
				m.quoteStrToMigrations(target.Name), m.quoteStrToMigrations(indexName), m.quoteStrListToMigrations(index.FieldName))
		}
	}
	if exists != nil && !reflect.DeepEqual(exists, target){
		content = m.diff(exists, target)
	}
	return content
}

func (m *Migrate) GetOperations() []*Operations {
	var operations []*Operations
	valueOf := reflect.ValueOf(m.Migrations)
	typeOf := reflect.TypeOf(m.Migrations)
	for i := 0; i < typeOf.NumMethod(); i++ {
		operations = append(
			operations,
			valueOf.MethodByName(typeOf.Method(i).Name).Call(nil)[0].Interface().(*Operations),
		)
	}
	return operations
}

func (m *Migrate) GetOperationsTree(withValid bool) *OperationsNode {
	operations := m.GetOperations()
	node := GenerateOperationsTree(&operations)
	if withValid {
		m.Valid(node)
	}
	return node
}

func (m *Migrate) genTablesFromMigrationFiles() (map[string]*Table, string) {
	var head string
	node := m.GetOperationsTree(true)
	heads := m.HeadToString(node)
	if len(heads) != 0 {
		head = heads[0]
	}
	return node.GetTable(), head
}

func (m *Migrate) indexesAndUniqueIndexes(scope *gorm.Scope) (map[string][]string, map[string][]string) {
	var indexes = map[string][]string{}
	var uniqueIndexes = map[string][]string{}

	for _, field := range scope.GetStructFields() {
		if name, ok := field.TagSettingsGet("INDEX"); ok {
			names := strings.Split(name, ",")

			for _, name := range names {
				if name == "INDEX" || name == "" {
					name = scope.Dialect().BuildKeyName("idx", scope.TableName(), field.DBName)
				}
				name, column := scope.Dialect().NormalizeIndexAndColumn(name, field.DBName)
				indexes[name] = append(indexes[name], column)
			}
		}

		if name, ok := field.TagSettingsGet("UNIQUE_INDEX"); ok {
			names := strings.Split(name, ",")

			for _, name := range names {
				if name == "UNIQUE_INDEX" || name == "" {
					name = scope.Dialect().BuildKeyName("uix", scope.TableName(), field.DBName)
				}
				name, column := scope.Dialect().NormalizeIndexAndColumn(name, field.DBName)
				uniqueIndexes[name] = append(uniqueIndexes[name], column)
			}
		}
	}
	return indexes, uniqueIndexes
}

func (m *Migrate) genTableFromObject(values ...interface{}) map[string]*Table {
	if len(values) == 0 {
		panic(fmt.Sprintf("no table specified to make"))
	}
	ret := make(map[string]*Table)
	for _, value := range values {
		scope := m.DB.NewScope(value)
		indexes, uniqueIndexes := m.indexesAndUniqueIndexes(scope)
		table := &Table{Name: scope.TableName()}
		for _, structField := range scope.GetStructFields() {
			field := Field{
				Name:      structField.DBName,
				Type:      scope.Dialect().DataTypeOf(structField),
				IsPrimary: structField.IsPrimaryKey,
			}
			table.Fields = append(table.Fields, &field)
		}
		for indexName, columns := range indexes {
			table.Indexes = append(table.Indexes, &Index{indexName, columns})
		}
		for indexName, columns := range uniqueIndexes {
			table.UniqueIndexes = append(table.UniqueIndexes, &Index{indexName, columns})
		}
		ret[table.Name] = table
	}
	return ret
}

func (*Migrate) genMigrationFileName(latest string) string {
	var ret string
	if latest == INITMigrations {
		ret = "0001" + MIGRATIONSplit
	} else {
		latest = strings.Split(latest, MIGRATIONSplit)[0]
		if i, err := strconv.Atoi(latest); err != nil {
			panic(err)
		} else {
			ret = fmt.Sprintf("%.4d", i+1) + MIGRATIONSplit
		}
	}

	ret += time.Now().Format(TIMEFormat) + strconv.FormatInt(time.Now().UnixNano()%1e8, 10)
	return ret
}

func (m *Migrate) write(migrationString, fileName string) {
	migrationsPath := m.migrationsPath()
	filePath := migrationsPath + fileName + ".go"
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()

	if err != nil {
		panic(err)
	}
	_, err = f.WriteString(migrationString)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v generated success.", filePath)
}

func (m *Migrate) Heads(root *OperationsNode, ret map[string]int) {
	if root == nil {
		return
	}
	if len(root.Children) == 0 {
		ret[root.Ops.Revision] += 1
	}
	for _, child := range root.Children {
		m.Heads(child, ret)
	}
}

func (m *Migrate) HeadToString(root *OperationsNode) []string {
	var ret []string
	heads := make(map[string]int)
	m.Heads(root, heads)
	for k := range heads {
		ret = append(ret, k)
	}
	return ret
}

func (m *Migrate) Valid(root *OperationsNode) {
	heads := m.HeadToString(root)
	if len(heads) > 1 {
		panic(fmt.Sprintf("multi heads %v", strings.Join(heads, " ")))
	}
}
