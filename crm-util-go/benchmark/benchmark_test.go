package benchmark

import (
	"strings"
	"testing"
)

var jsonStr = `{"one":"foobar","two":"foobar","three":"foobar","four":"foobar","five":"foobar","six":"foobar","seven":"foobar","eight":"foobar","nine":"foobar","ten":"foobar"}`

func BenchmarkIOReadAll(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		BMIOReadAll(strings.NewReader(jsonStr))
	}
}

func BenchmarkIOCopy(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		BMIOCopy(strings.NewReader(jsonStr))
	}
}

func BenchmarkBMHasStringValueLen(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		BMHasStringValueLen("0123456789abc0123456789abc0123456789abc0123456789abc")
	}
}

func BenchmarkBMHasStringValueEqual(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		BMHasStringValueEqual("0123456789abc0123456789abc0123456789abc0123456789abc")
	}
}
