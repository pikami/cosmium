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
	ConstantTypeParameterConstant
)

type SelectItemType int

const (
	SelectItemTypeField SelectItemType = iota
	SelectItemTypeObject
	SelectItemTypeArray
)

type OrderDirection int

const (
	OrderDirectionAsc OrderDirection = iota
	OrderDirectionDesc
)

type SelectStmt struct {
	SelectItems      []SelectItem
	Table            Table
	Filters          interface{}
	Count            int
	Parameters       map[string]interface{}
	OrderExpressions []OrderExpression
}

type Table struct {
	Value string
}

type SelectItem struct {
	Alias       string
	Path        []string
	SelectItems []SelectItem
	Type        SelectItemType
	IsTopLevel  bool
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

type OrderExpression struct {
	SelectItem SelectItem
	Direction  OrderDirection
}
