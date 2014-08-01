package ldap

import (
	"testing"

	"github.com/vanackere/asn1-ber"
)

type compileTest struct {
	filterStr  string
	filterType int
}

var testFilters = []compileTest{
	compileTest{filterStr: "(&(sn=Miller)(givenName=Bob))", filterType: FilterAnd},
	compileTest{filterStr: "(|(sn=Miller)(givenName=Bob))", filterType: FilterOr},
	compileTest{filterStr: "(!(sn=Miller))", filterType: FilterNot},
	compileTest{filterStr: "(sn=Miller)", filterType: FilterEqualityMatch},
	compileTest{filterStr: "(sn=Mill*)", filterType: FilterSubstrings},
	compileTest{filterStr: "(sn=*Mill)", filterType: FilterSubstrings},
	compileTest{filterStr: "(sn=*Mill*)", filterType: FilterSubstrings},
	compileTest{filterStr: "(sn>=Miller)", filterType: FilterGreaterOrEqual},
	compileTest{filterStr: "(sn<=Miller)", filterType: FilterLessOrEqual},
	compileTest{filterStr: "(sn=*)", filterType: FilterPresent},
	compileTest{filterStr: "(sn~=Miller)", filterType: FilterApproxMatch},
	// compileTest{ filterStr: "()", filterType: FilterExtensibleMatch },
}

func TestFilter(t *testing.T) {
	// Test Compiler and Decompiler
	for _, i := range testFilters {
		filter, err := CompileFilter(i.filterStr)
		if err != nil {
			t.Errorf("Problem compiling %s - %s", i.filterStr, err.Error())
		} else if filter.Tag != uint8(i.filterType) {
			t.Errorf("%q Expected %q got %q", i.filterStr, FilterMap[uint64(i.filterType)], FilterMap[uint64(filter.Tag)])
		} else {
			o, err := DecompileFilter(filter)
			if err != nil {
				t.Errorf("Problem compiling %s - %s", i.filterStr, err.Error())
			} else if i.filterStr != o {
				t.Errorf("%q expected, got %q", i.filterStr, o)
			}
		}
	}
}

func TestDecodePresentFilter(t *testing.T) {

	bs := []byte{0x87, 0x06, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72} // ~ (member=*)

	p := ber.DecodePacket(bs)
	f, err := DecompileFilter(p)
	if err != nil {
		t.Errorf("--- CAN'T DECODE & DECOMPILE FILTER")
		return
	}
	if f != "(member=*)" {
		t.Errorf("expected (member=*), got %s", f)
	}

}

func TestEncodeDecodePresentFilter(t *testing.T) {
	f := "(member=*)"
	p, err := CompileFilter(f)
	if err != nil {
		t.Errorf("cant compile filter")
		return
	}
	bytes := p.Bytes()

	p2 := ber.DecodePacket(bytes)
	f2, err := DecompileFilter(p2)
	if err != nil {
		t.Errorf("cant decompile filter")
		return
	}
	if f != f2 {
		t.Errorf("encode/decode changed filter")
	}
}

func BenchmarkFilterCompile(b *testing.B) {
	b.StopTimer()
	filters := make([]string, len(testFilters))

	// Test Compiler and Decompiler
	for idx, i := range testFilters {
		filters[idx] = i.filterStr
	}

	maxIdx := len(filters)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		CompileFilter(filters[i%maxIdx])
	}
}

func BenchmarkFilterDecompile(b *testing.B) {
	b.StopTimer()
	filters := make([]*ber.Packet, len(testFilters))

	// Test Compiler and Decompiler
	for idx, i := range testFilters {
		filters[idx], _ = CompileFilter(i.filterStr)
	}

	maxIdx := len(filters)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		DecompileFilter(filters[i%maxIdx])
	}
}
