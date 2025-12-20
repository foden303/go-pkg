package bs

import "unsafe"

// String2Bytes converts byte slice to a string without any memory allocation.
func Bytes2String(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// BytesToString converts []byte to string without allocation.
// func BytesToString(b []byte) string {
// 	return *(*string)(unsafe.Pointer(&b))
// }

// String2Bytes converts string to a byte slice without any memory allocation.
func String2Bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
