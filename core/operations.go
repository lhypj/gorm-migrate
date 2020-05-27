package core

import (
	"fmt"
)

type ActionType uint8

const (
	ADDField ActionType = iota
	DELETEField
	ALTERField
	ADDTable
	DELETETable
	ADDIndex
	DELETEIndex
)

type Operation struct {
	Action     ActionType
	TableName  string
	ColumnName string
	Type       string
	TypeNew    string
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

func searchTree(root *OperationsNode, tableFieldType map[string]map[string]string) {
	if root == nil {
		return
	}
	if root.Ops != nil {
		for _, op := range root.Ops.Operations {
			switch op.Action {
			case ADDTable:
				column := make(map[string]string)
				tableFieldType[op.TableName] = column
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
				fmt.Println("todo add indexes")
			case DELETEIndex:
				fmt.Println("todo delete indexes")
			}
		}
	}
	for _, child := range root.Children {
		searchTree(child, tableFieldType)
	}
}

func (root *OperationsNode) GetTable() map[string]*Table {
	tables := make(map[string]*Table)

	tableFieldType := make(map[string]map[string]string)
	searchTree(root, tableFieldType)

	for tableName := range tableFieldType {
		t := tableFieldType[tableName]
		table := Table{Name: tableName}
		for fs := range t {
			table.Fields = append(table.Fields, &Field{Name: fs, Type: t[fs]})
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
