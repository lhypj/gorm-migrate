package core

import (
	"fmt"
)

type ActionType uint8
type OperationSlice []*Operations

const (
	ADDField ActionType = iota
	DELETEField
	ALTERField
	ADDTable
	DELETETable
	ADDIndex
	DELETEIndex
	ADDUniqueIndex
	DELETEUniqueIndex
)

const (
	ADDIndexStr          = "core.ADDIndex"
	DELETEIndexStr       = "core.DELETEIndex"
	ADDUniqueIndexStr    = "core.ADDUniqueIndex"
	DELETEUniqueIndexStr = "core.DELETEUniqueIndex"
)

type Operation struct {
	Action          ActionType
	TableName       string
	ColumnName      string
	Type            string
	TypeNew         string
	IsPrimary       bool
	IndexName       string
	IndexFieldNames []string
}

type Operations struct {
	Revision     string
	DownRevision []string
	Operations   []*Operation
}

type OperationsNode struct {
	Ops      *Operations
	Children []*OperationsNode
}

func searchTree(node *OperationsNode,
	tableFieldType map[string]map[string]string,
	tableIndexes map[string]map[string][]string,
	tableUniqueIndexes map[string]map[string][]string) {
	if node == nil {
		return
	}
	if node.Ops != nil {
		for _, op := range node.Ops.Operations {
			if tableFieldType[op.TableName] == nil {
				tableFieldType[op.TableName] = make(map[string]string)
			}
			if tableIndexes[op.TableName] == nil {
				tableIndexes[op.TableName] = make(map[string][]string)
			}
			if tableUniqueIndexes[op.TableName] == nil {
				tableUniqueIndexes[op.TableName] = make(map[string][]string)
			}
			switch op.Action {
			case ADDField:
				tableFieldType[op.TableName][op.ColumnName] = op.Type
			case DELETEField:
				if tableFieldType[op.TableName] != nil && tableFieldType[op.TableName][op.ColumnName] != "" {
					delete(tableFieldType[op.TableName], op.ColumnName)
				}
			case ALTERField:
				if tableFieldType[op.TableName][op.ColumnName] == op.TypeNew {
					panic(fmt.Sprintf("Alter Action is err %v => %v", op.Type, op.TypeNew))
				}
				tableFieldType[op.TableName][op.ColumnName] = op.TypeNew
			case DELETETable:
				if tableFieldType[op.TableName] == nil {
					panic(fmt.Sprintf("Table:%v not exists, cannot delete", op.TableName))
				}
				delete(tableFieldType, op.TableName)
			case ADDIndex:
				tableIndexes[op.TableName][op.IndexName] = op.IndexFieldNames
			case DELETEIndex:
				if tableIndexes[op.TableName][op.IndexName] != nil {
					delete(tableIndexes[op.TableName], op.IndexName)
				}
			case ADDUniqueIndex:
				tableUniqueIndexes[op.TableName][op.IndexName] = op.IndexFieldNames
			case DELETEUniqueIndex:
				if tableUniqueIndexes[op.TableName][op.IndexName] != nil {
					delete(tableUniqueIndexes[op.TableName], op.IndexName)
				}
			}
		}
	}
	for _, child := range node.Children {
		searchTree(child, tableFieldType, tableIndexes, tableUniqueIndexes)
	}
}

func (root *OperationsNode) GetTable() map[string]*Table {
	// todo 修改主键，现阶段Table未添加主键的更改信息
	tables := make(map[string]*Table)

	tableFieldType := make(map[string]map[string]string)
	tableIndexes := make(map[string]map[string][]string)
	tableUniqueIndexes := make(map[string]map[string][]string)
	searchTree(root, tableFieldType, tableIndexes, tableUniqueIndexes)

	for tableName := range tableFieldType {
		table := Table{Name: tableName}

		t := tableFieldType[tableName]
		for fs := range t {
			table.Fields = append(table.Fields, &Field{Name: fs, Type: t[fs]})
		}

		for idxName, columns := range tableIndexes[tableName] {
			table.Indexes = append(table.Indexes, &Index{Name: idxName, FieldName: columns})
		}

		for idxName, columns := range tableUniqueIndexes[tableName] {
			table.UniqueIndexes = append(table.UniqueIndexes, &Index{Name: idxName, FieldName: columns})
		}

		tables[tableName] = &table
	}
	return tables
}

func GenerateTree(m map[string]*Operations, children map[string][]*Operations, revision string) *OperationsNode {
	var node OperationsNode
	node.Ops = m[revision]
	if revision != INITMigrations && node.Ops == nil {
		return &node
	}
	for _, ops := range children[revision] {
		if _node := GenerateTree(m, children, ops.Revision); _node != nil {
			node.Children = append(node.Children, _node)
		}
	}
	return &node
}

func PrintTree(root *OperationsNode, blank string) {
	if root == nil {
		return
	}
	fmt.Printf("%v", blank)
	fmt.Println(root.Ops)
	for _, child := range root.Children {
		PrintTree(child, blank+"\t")
	}
}

func ShowTree(root *OperationsNode) {
	fmt.Println("start draw tree")
	PrintTree(root, "")
	fmt.Println("end draw tree")
}

func GenerateOperationsTree(operations *[]*Operations) *OperationsNode {
	childrenRevisionOperations := make(map[string][]*Operations)
	revisionOperations := make(map[string]*Operations)
	revisionOperations[INITMigrations] = nil
	for _, ops := range *operations {
		revisionOperations[ops.Revision] = ops
		for _, downReversion := range ops.DownRevision {
			childrenRevisionOperations[downReversion] = append(childrenRevisionOperations[downReversion], ops)
		}
	}
	node := GenerateTree(revisionOperations, childrenRevisionOperations, "")
	//ShowTree(node)
	return node
}

func (ops OperationSlice) Len() int           {return len(ops)}
func (ops OperationSlice) Swap(i, j int)      {ops[i], ops[j] = ops[j], ops[i]}
func (ops OperationSlice) Less(i, j int) bool {return ops[i].Revision < ops[j].Revision}