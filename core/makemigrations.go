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

func (m *Migrate) migrationsPath() (string, error) {
	if path, err := m.rootPath(); err != nil {
		return "", err
	} else {
		return path + m.ModelsRelativePath + MIGRATIONPath, nil
	}

}

func (m *Migrate) MigrationsInit() error {
	// migrate models
	m.DB.AutoMigrate(&OrmMigrations{})

	// init.go need user to init
	// init.go file for search migrations files
	//if mPath, err := m.migrationsPath(); err != nil {
	//	return err
	//} else {
	//	_, err := os.Stat(mPath + "init.go")
	//	if os.IsNotExist(err) {
	//		return m.write(INITContent, "init")
	//	}
	//	return nil
	//}
	return nil
}

func (m *Migrate) MakeMigrations(migrations interface{}, tables ...interface{}) error {
	tableFromFile, head, err := m.genTablesFromMigrationFiles(migrations)
	if err != nil {
		return err
	}
	tableFromObj, err := m.genTableFromObject(tables...)
	if err != nil {
		return err
	}

	var content string
	pre := m.MigrationsPre()
	end, fn := m.MigrationsEnd(head)
	for _, table := range tables {
		name := m.DB.NewScope(table).TableName()
		content += m.genMigrationFileContent(tableFromFile[name], tableFromObj[name])
	}
	if content == "" {
		fmt.Println("No migrations need to make")
	} else {
		if err := m.write(pre+content+end, fn); err != nil {
			return err
		}
	}
	return nil
}

func (m *Migrate) MigrationsPre() string {
	return fmt.Sprintf("package migrations\n\nimport \"%v\"\n\nfunc (*Migrations)Migration_0001_2020052718060846904000() *core.Operations {\n\tvar ops []*core.Operation\n\tops = append(ops,\n", m.PackagePath)
}

func (m *Migrate) MigrationsEnd(latest string) (end, fn string) {
	if latest == "" {
		latest = INITMigrations
	}
	fn = m.genMigrationFileName(latest)
	end = fmt.Sprintf("\t)\n\treturn &core.Operations{Revision: \"%v\", DownRevision: []string{\"%v\"}, Operations: ops}\n}", fn, latest)
	return
}

func (m *Migrate) diff(exists *Table, target *Table) string {
	// 写到这我发现，其实autoMigrate也挺好。
	// 因为只增加字段符合了大数据量的应用场景：减少全表的结构修改，以减少数据服务不可用。
	// 但是数据量大的话做外键岂不是更难搞。不过也不是所有DB都是大表，emm～怎么衡量，还是要好好思量。

	// todo sth wrong
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
			content += m.migrationAddFieldContent(tableName, n.Name, n.Type)
		}
		if o != nil && n == nil {
			content += m.migrationDeleteFieldContent(tableName, o.Name)
		}
		if o != nil && n != nil {
			content += m.migrationAlterFieldContent(tableName, o.Name, o.Type, n.Type)
		}
	}
	return content
}

func (m *Migrate) migrationAddFieldContent(tableName, fieldName, fileType string) string {
	// todo with quote?
	return fmt.Sprintf("\t\t&core.Operation{Action: core.ADDField, TableName: \"%v\", ColumnName: \"%v\", Type: \"%v\"},\n", tableName, fieldName, fileType)
}

func (m *Migrate) migrationDeleteFieldContent(tableName, fieldName string) string {
	return fmt.Sprintf("\t\t&core.Operation{Action: core.DeleteField, TableName: \"%v\", ColumnName: \"%v\"},\n", tableName, fieldName)
}

func (m *Migrate) migrationAlterFieldContent(tableName, fieldName, old, new string) string {
	return fmt.Sprintf("\t\t&core.Operation{Action: core.AlterField, TableName: \"%v\", ColumnName: \"%v\"}, Type: \"%v\"}, TypeNew: \"%v\"", tableName, fieldName, old, new)
}

func (m *Migrate) genMigrationFileContent(exists *Table, target *Table) string {
	var content string
	if target == nil {
		if exists != nil {
			// 如果创建的表未migrate 那么这段语句不会执行。表现的结果是migrate之后表会创建，再次make migrations会remove。
			content += fmt.Sprintf("&core.Operation{Action: core.DELETETable, TableName: \"%v\"},", exists.Name)
		}
		return content
	}
	if exists == nil {
		content = fmt.Sprintf("\t\t&core.Operation{Action: core.ADDTable, TableName: \"%v\"},\n", target.Name)
		for _, field := range target.Fields {
			content += m.migrationAddFieldContent(target.Name, field.Name, field.Type)
		}
	}
	if exists != nil {
		if !reflect.DeepEqual(exists, target) {
			content = m.diff(exists, target)
		}
	}
	return content
}

func (m *Migrate) GetOperationsTree(migrations interface{}) *OperationsNode {
	var operations []*Operations
	valueOf := reflect.ValueOf(migrations)
	typeOf := reflect.TypeOf(migrations)
	if valueOf.NumMethod() < 1 {
		return nil
	}

	for i := 0; i < typeOf.NumMethod(); i++ {
		method := typeOf.Method(i)
		ops := valueOf.MethodByName(method.Name).Call(nil)[0].Interface().(*Operations)
		operations = append(operations, ops)
	}
	node := GenerateOperationsTree(&operations)
	m.Valid(node)
	return node
}

func (m *Migrate) genTablesFromMigrationFiles(migrations interface{}) (map[string]*Table, string, error) {
	var head string
	node := m.GetOperationsTree(migrations)
	heads := m.HeadToString(node)
	if len(heads) != 0 {
		head = heads[0]
	}
	return node.GetTable(), head, nil
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
				indexes[column] = append(indexes[column], name)
			}
		}

		if name, ok := field.TagSettingsGet("UNIQUE_INDEX"); ok {
			names := strings.Split(name, ",")

			for _, name := range names {
				if name == "UNIQUE_INDEX" || name == "" {
					name = scope.Dialect().BuildKeyName("uix", scope.TableName(), field.DBName)
				}
				name, column := scope.Dialect().NormalizeIndexAndColumn(name, field.DBName)
				uniqueIndexes[column] = append(uniqueIndexes[column], name)
			}
		}
	}
	return indexes, uniqueIndexes
}

func (m *Migrate) genTableFromObject(values ...interface{}) (map[string]*Table, error) {
	if len(values) == 0 {
		return nil, fmt.Errorf("no table specified to make\n")
	}
	ret := make(map[string]*Table)
	for _, value := range values {
		scope := m.DB.NewScope(value)
		indexes, uniqueIndexes := m.indexesAndUniqueIndexes(scope)
		table := &Table{Name: scope.TableName()}
		for _, structField := range scope.GetStructFields() {
			field := Field{
				Name:             structField.DBName,
				Type:             scope.Dialect().DataTypeOf(structField),
				IsPrimary:        structField.IsPrimaryKey,
				IndexNames:       indexes[structField.DBName],
				UniqueIndexNames: uniqueIndexes[structField.DBName],
			}
			table.Fields = append(table.Fields, &field)
		}
		ret[table.Name] = table
	}
	return ret, nil
}

func (m *Migrate) genMigrationFiles() error {
	fmt.Println("writing")
	return nil
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

func (m *Migrate) write(migrationString, fileName string) error {
	migrationsPath, err := m.migrationsPath()
	if err != nil {
		return err
	}
	filePath := migrationsPath + fileName + ".go"
	fmt.Println(filePath)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()

	if err != nil {
		return err
	}
	_, err = f.WriteString(migrationString)
	if err != nil {
		return err
	}
	return nil
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

func (m *Migrate) Merge(root *OperationsNode) {
	heads := m.HeadToString(root)
	if len(heads) > 1 {

		//panic(fmt.Sprintf("multi heads %v", strings.Join(*Map2StringList(heads), " ")))
	}
}
