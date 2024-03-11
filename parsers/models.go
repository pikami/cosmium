package parsers

type SelectStmt struct {
	SelectItems      []SelectItem
	Table            Table
	Filters          interface{}
	Distinct         bool
	Count            int
	Offset           int
	Parameters       map[string]interface{}
	OrderExpressions []OrderExpression
	GroupBy          []SelectItem
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
	FunctionCallConcat       FunctionCallType = "Concat"
	FunctionCallContains     FunctionCallType = "Contains"
	FunctionCallEndsWith     FunctionCallType = "EndsWith"
	FunctionCallStartsWith   FunctionCallType = "StartsWith"
	FunctionCallIndexOf      FunctionCallType = "IndexOf"
	FunctionCallToString     FunctionCallType = "ToString"
	FunctionCallUpper        FunctionCallType = "Upper"
	FunctionCallLower        FunctionCallType = "Lower"
	FunctionCallLeft         FunctionCallType = "Left"
	FunctionCallLength       FunctionCallType = "Length"
	FunctionCallLTrim        FunctionCallType = "LTrim"
	FunctionCallReplace      FunctionCallType = "Replace"
	FunctionCallReplicate    FunctionCallType = "Replicate"
	FunctionCallReverse      FunctionCallType = "Reverse"
	FunctionCallRight        FunctionCallType = "Right"
	FunctionCallRTrim        FunctionCallType = "RTrim"
	FunctionCallSubstring    FunctionCallType = "Substring"
	FunctionCallTrim         FunctionCallType = "Trim"

	FunctionCallIsDefined      FunctionCallType = "IsDefined"
	FunctionCallIsArray        FunctionCallType = "IsArray"
	FunctionCallIsBool         FunctionCallType = "IsBool"
	FunctionCallIsFiniteNumber FunctionCallType = "IsFiniteNumber"
	FunctionCallIsInteger      FunctionCallType = "IsInteger"
	FunctionCallIsNull         FunctionCallType = "IsNull"
	FunctionCallIsNumber       FunctionCallType = "IsNumber"
	FunctionCallIsObject       FunctionCallType = "IsObject"
	FunctionCallIsPrimitive    FunctionCallType = "IsPrimitive"
	FunctionCallIsString       FunctionCallType = "IsString"

	FunctionCallArrayConcat  FunctionCallType = "ArrayConcat"
	FunctionCallArrayLength  FunctionCallType = "ArrayLength"
	FunctionCallArraySlice   FunctionCallType = "ArraySlice"
	FunctionCallSetIntersect FunctionCallType = "SetIntersect"
	FunctionCallSetUnion     FunctionCallType = "SetUnion"

	FunctionCallAggregateAvg   FunctionCallType = "AggregateAvg"
	FunctionCallAggregateCount FunctionCallType = "AggregateCount"
	FunctionCallAggregateMax   FunctionCallType = "AggregateMax"
	FunctionCallAggregateMin   FunctionCallType = "AggregateMin"
	FunctionCallAggregateSum   FunctionCallType = "AggregateSum"

	FunctionCallIn FunctionCallType = "In"
)

var AggregateFunctions = []FunctionCallType{
	FunctionCallAggregateAvg,
	FunctionCallAggregateCount,
	FunctionCallAggregateMax,
	FunctionCallAggregateMin,
	FunctionCallAggregateSum,
}

type FunctionCall struct {
	Arguments []interface{}
	Type      FunctionCallType
}
