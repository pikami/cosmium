package parsers

type LogicalExpressionType int

const (
	LogicalExpressionTypeOr LogicalExpressionType = iota
	LogicalExpressionTypeAnd
)

type ConstantType int

const (
	ConstantTypeString ConstantType = iota
	ConstantTypeInteger
	ConstantTypeFloat
	ConstantTypeBoolean
)

type SelectStmt struct {
	Columns []FieldPath
	Table   Table
	Filters interface{}
}

type Table struct {
	Value string
}

type FieldPath struct {
	Alias string
	Path  []string
}

type LogicalExpression struct {
	Expressions []interface{}
	Operation   LogicalExpressionType
}

type ComparisonExpression struct {
	Left      interface{}
	Right     interface{}
	Operation string
}

type Constant struct {
	Type  ConstantType
	Value interface{}
}
