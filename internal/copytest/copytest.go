// Package copytest provides testing utilites relating to struct copying
package copytest

import "github.com/jinzhu/copier"

// DeepCopy creates deep copies of structs.
// Pass struct pointers in the form <&out1>, <&in1>, <&out2>, <&in2>, ...
func DeepCopy(outIn ...any) (err error) {
	if len(outIn)%2 != 0 {
		panic("need even number of arguments in deepCopy")
	}
	// deep copy will copy values behind pointers instead of the pointers themselves
	copyOption := copier.Option{DeepCopy: true}

	for i := 0; i < len(outIn); i += 2 {
		if err = copier.CopyWithOption(outIn[i], outIn[i+1], copyOption); err != nil {
			return
		}
	}
	return
}
