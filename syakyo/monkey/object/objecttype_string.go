// Code generated by "stringer -type ObjectType -linecomment"; DO NOT EDIT.

package object

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[IntegerObjectType-0]
	_ = x[BooleanObjectType-1]
	_ = x[NullObjectType-2]
	_ = x[ReturnValueObjectType-3]
	_ = x[ErrorObjectType-4]
}

const _ObjectType_name = "INTEGERBOOLEANNULLRETURN_VALUEERROR"

var _ObjectType_index = [...]uint8{0, 7, 14, 18, 30, 35}

func (i ObjectType) String() string {
	if i < 0 || i >= ObjectType(len(_ObjectType_index)-1) {
		return "ObjectType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ObjectType_name[_ObjectType_index[i]:_ObjectType_index[i+1]]
}
