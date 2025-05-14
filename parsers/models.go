package parsers

type SelectStmt struct {
	SelectItems      []SelectItem
	Table            Table
	JoinItems        []JoinItem
	Filters          interface{}
	Exists           bool
	Distinct         bool
	Count            int
	Offset           int
	Parameters       map[string]interface{}
	OrderExpressions []OrderExpression
	GroupBy          []SelectItem
}

type Table struct {
	Value      string
	SelectItem SelectItem
	IsInSelect bool
}

type JoinItem struct {
	Table      Table
	SelectItem SelectItem
}

type SelectItemType int

const (
	SelectItemTypeField SelectItemType = iota
	SelectItemTypeObject
	SelectItemTypeArray
	SelectItemTypeConstant
	SelectItemTypeFunctionCall
	SelectItemTypeSubQuery
)

type SelectItem struct {
	Alias       string
	Path        []string
	SelectItems []SelectItem
	Type        SelectItemType
	Value       interface{}
	Invert      bool
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

	FunctionCallArrayConcat      FunctionCallType = "ArrayConcat"
	FunctionCallArrayContains    FunctionCallType = "ArrayContains"
	FunctionCallArrayContainsAny FunctionCallType = "ArrayContainsAny"
	FunctionCallArrayContainsAll FunctionCallType = "ArrayContainsAll"
	FunctionCallArrayLength      FunctionCallType = "ArrayLength"
	FunctionCallArraySlice       FunctionCallType = "ArraySlice"
	FunctionCallSetIntersect     FunctionCallType = "SetIntersect"
	FunctionCallSetUnion         FunctionCallType = "SetUnion"

	FunctionCallIif FunctionCallType = "Iif"

	FunctionCallMathAbs              FunctionCallType = "MathAbs"
	FunctionCallMathAcos             FunctionCallType = "MathAcos"
	FunctionCallMathAsin             FunctionCallType = "MathAsin"
	FunctionCallMathAtan             FunctionCallType = "MathAtan"
	FunctionCallMathAtn2             FunctionCallType = "MathAtn2"
	FunctionCallMathCeiling          FunctionCallType = "MathCeiling"
	FunctionCallMathCos              FunctionCallType = "MathCos"
	FunctionCallMathCot              FunctionCallType = "MathCot"
	FunctionCallMathDegrees          FunctionCallType = "MathDegrees"
	FunctionCallMathExp              FunctionCallType = "MathExp"
	FunctionCallMathFloor            FunctionCallType = "MathFloor"
	FunctionCallMathIntAdd           FunctionCallType = "MathIntAdd"
	FunctionCallMathIntBitAnd        FunctionCallType = "MathIntBitAnd"
	FunctionCallMathIntBitLeftShift  FunctionCallType = "MathIntBitLeftShift"
	FunctionCallMathIntBitNot        FunctionCallType = "MathIntBitNot"
	FunctionCallMathIntBitOr         FunctionCallType = "MathIntBitOr"
	FunctionCallMathIntBitRightShift FunctionCallType = "MathIntBitRightShift"
	FunctionCallMathIntBitXor        FunctionCallType = "MathIntBitXor"
	FunctionCallMathIntDiv           FunctionCallType = "MathIntDiv"
	FunctionCallMathIntMod           FunctionCallType = "MathIntMod"
	FunctionCallMathIntMul           FunctionCallType = "MathIntMul"
	FunctionCallMathIntSub           FunctionCallType = "MathIntSub"
	FunctionCallMathLog              FunctionCallType = "MathLog"
	FunctionCallMathLog10            FunctionCallType = "MathLog10"
	FunctionCallMathNumberBin        FunctionCallType = "MathNumberBin"
	FunctionCallMathPi               FunctionCallType = "MathPi"
	FunctionCallMathPower            FunctionCallType = "MathPower"
	FunctionCallMathRadians          FunctionCallType = "MathRadians"
	FunctionCallMathRand             FunctionCallType = "MathRand"
	FunctionCallMathRound            FunctionCallType = "MathRound"
	FunctionCallMathSign             FunctionCallType = "MathSign"
	FunctionCallMathSin              FunctionCallType = "MathSin"
	FunctionCallMathSqrt             FunctionCallType = "MathSqrt"
	FunctionCallMathSquare           FunctionCallType = "MathSquare"
	FunctionCallMathTan              FunctionCallType = "MathTan"
	FunctionCallMathTrunc            FunctionCallType = "MathTrunc"

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
