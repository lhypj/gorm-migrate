package core

import (
	"fmt"
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
	m.DB.AutoMigrate(&GOrmMigrations{})

	// init.go file for search migrations files
	if mPath, err := m.migrationsPath(); err != nil {
		return err
	} else {
		_, err := os.Stat(mPath + "init.go")
		if os.IsNotExist(err) {
			return m.write(INITContent, "init")
		}
		return nil
	}
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

func (m *Migrate) genMigrationFileContent(exists *Table, target *Table) string {
	var content string
	if target == nil {
		return content
	}
	if exists == nil {
		content = fmt.Sprintf("\t\t&core.Operation{Action: core.ADDTable, TableName: \"%v\"},\n", target.Name)
		for _, field := range target.Fields {
			content += fmt.Sprintf("\t\t&core.Operation{Action: core.ADDField, TableName: \"%v\", ColumnName: \"%v\", Type: \"%v\"},\n", target.Name, field.Name, field.Type)
		}
	}
	if exists != nil {
		if !reflect.DeepEqual(exists, target) {
			// diff
		}
	}
	return content
}

func (m *Migrate) genTablesFromMigrationFiles(migrations interface{}) (map[string]*Table, string, error) {
	fmt.Println("从已有的migration files(apply, unapply)生成 已生成结构 ing")
	ret := make(map[string]*Table)
	var operations []*Operations
	valueOf := reflect.ValueOf(migrations)
	typeOf := reflect.TypeOf(migrations)
	if valueOf.NumMethod() < 1 {
		return ret, "", nil
	}

	for i := 0; i < typeOf.NumMethod(); i++ {
		method := typeOf.Method(i)
		ops := valueOf.MethodByName(method.Name).Call(nil)[0].Interface().(*Operations)
		operations = append(operations, ops)
	}
	node := GenerateOperationsTree(&operations)
	m.Valid(node)
	heads := *m.HeadToString(node)
	return node.GetTable(), heads[0], nil

}

func (m *Migrate) genTableFromObject(values ...interface{}) (map[string]*Table, error) {
	if len(values) == 0 {
		return nil, fmt.Errorf("no table need to m\n")
	}
	ret := make(map[string]*Table)
	for _, value := range values {
		scope := m.DB.NewScope(value)
		table := &Table{Name: scope.TableName()}
		for _, structField := range scope.GetStructFields() {

			field := Field{Name: structField.Name, Type: scope.Dialect().DataTypeOf(structField)}
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

func (m *Migrate) HeadToString(root *OperationsNode) *[]string {
	var ret []string
	heads := make(map[string]int)
	m.Heads(root, heads)
	for k := range heads {
		ret = append(ret, k)
	}
	return &ret
}

func (m *Migrate) Valid(root *OperationsNode) {
	heads := *m.HeadToString(root)
	if len(heads) > 1 {
		panic(fmt.Sprintf("multi heads %v", strings.Join(heads, " ")))
	}
}

func (m *Migrate) Merge(root *OperationsNode) {
	heads := *m.HeadToString(root)
	if len(heads) > 1 {

		//panic(fmt.Sprintf("multi heads %v", strings.Join(*Map2StringList(heads), " ")))
	}
}
