package testutils

import "github.com/pikami/cosmium/parsers"

func SelectItem_Constant_String(value string) parsers.SelectItem {
	return parsers.SelectItem{
		Type: parsers.SelectItemTypeConstant,
		Value: parsers.Constant{
			Type:  parsers.ConstantTypeString,
			Value: value,
		},
	}
}

func SelectItem_Constant_Int(value int) parsers.SelectItem {
	return parsers.SelectItem{
		Type: parsers.SelectItemTypeConstant,
		Value: parsers.Constant{
			Type:  parsers.ConstantTypeInteger,
			Value: value,
		},
	}
}

func SelectItem_Constant_Float(value float64) parsers.SelectItem {
	return parsers.SelectItem{
		Type: parsers.SelectItemTypeConstant,
		Value: parsers.Constant{
			Type:  parsers.ConstantTypeFloat,
			Value: value,
		},
	}
}

func SelectItem_Constant_Bool(value bool) parsers.SelectItem {
	return parsers.SelectItem{
		Type: parsers.SelectItemTypeConstant,
		Value: parsers.Constant{
			Type:  parsers.ConstantTypeBoolean,
			Value: value,
		},
	}
}

func SelectItem_Constant_Parameter(name string) parsers.SelectItem {
	return parsers.SelectItem{
		Type: parsers.SelectItemTypeConstant,
		Value: parsers.Constant{
			Type:  parsers.ConstantTypeParameterConstant,
			Value: name,
		},
	}
}

func SelectItem_Path(path ...string) parsers.SelectItem {
	return parsers.SelectItem{
		Path: path,
	}
}
