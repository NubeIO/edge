// Code generated by "stringer -type=ProductType"; DO NOT EDIT.

package product

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[RubixCompute-0]
	_ = x[RubixComputeIO-1]
	_ = x[RubixCompute5-2]
	_ = x[Edge28-3]
	_ = x[Nuc-4]
	_ = x[AllLinux-5]
	_ = x[None-6]
}

const _ProductType_name = "RubixComputeRubixComputeIORubixCompute5Edge28NucAllLinuxNone"

var _ProductType_index = [...]uint8{0, 12, 26, 39, 45, 48, 56, 60}

func (i ProductType) String() string {
	if i < 0 || i >= ProductType(len(_ProductType_index)-1) {
		return "ProductType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ProductType_name[_ProductType_index[i]:_ProductType_index[i+1]]
}
