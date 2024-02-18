package parsers

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

type SelectItemType int

const (
	SelectItemTypeField SelectItemType = iota
	SelectItemTypeObject
	SelectItemTypeArray
	SelectItemTypeConstant
	SelectItemTypeFunctionCall
)

type SelectItem struct {
	Alias       string
	Path        []string
	SelectItems []SelectItem
	Type        SelectItemType
	Value       interface{}
	IsTopLevel  bool
}

type LogicalExpressionType int

const (
	LogicalExpressionTypeOr LogicalExpressionType = iota
	LogicalExpressionTypeAnd
)

type LogicalExpression struct {
	Expressions []interface{}
	Operation   LogicalExpressionType
}

type ComparisonExpression struct {
	Left      interface{}
	Right     interface{}
	Operation string
}

type ConstantType int

const (
	ConstantTypeString ConstantType = iota
	ConstantTypeInteger
	ConstantTypeFloat
	ConstantTypeBoolean
	ConstantTypeParameterConstant
)

type Constant struct {
	Type  ConstantType
	Value interface{}
}

type OrderDirection int

const (
	OrderDirectionAsc OrderDirection = iota
	OrderDirectionDesc
)

type OrderExpression struct {
	SelectItem SelectItem
	Direction  OrderDirection
}

type FunctionCallType string

const (
	FunctionCallStringEquals FunctionCallType = "StringEquals"
	FunctionCallIsDefined    FunctionCallType = "IsDefined"
)

type FunctionCall struct {
	Arguments []interface{}
	Type      FunctionCallType
}
