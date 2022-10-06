package testdata

import (
	"github.com/axelarnetwork/utils/checks/testdata/imported"
)

type RegularStruct struct {
	Integer int
	Boolean bool
	String  string
	Uint    uint // Deprecated
}

type StructWithPrivateFields struct {
	integer int
	Boolean bool
}

type StructWithEmbedded struct {
	RegularStruct
	SomeString string
}

type StructWithIgnoreTag struct {
	integer int `decl-check:"ignore"`
	boolean bool
	String  string
	Uint    uint `decl-check:"ignore"`
}

func _() {
	_ = RegularStruct{}

	_ = RegularStruct{
		Integer: 0,
		Boolean: false,
		String:  "",
	}

	_ = RegularStruct{
		Integer: 0,
		Boolean: false,
		String:  "",
		Uint:    0,
	}

	_ = RegularStruct{
		0,
		false,
		"",
		0,
	}

	// should fail
	_ = RegularStruct{
		Integer: 0,
		Boolean: false,
	}

	// should fail 2x
	_ = RegularStruct{
		Boolean: false,
	}
}

func _() {
	_ = StructWithPrivateFields{
		integer: 0,
		Boolean: false,
	}

	_ = StructWithPrivateFields{
		0,
		false,
	}

	// should fail
	_ = StructWithPrivateFields{
		Boolean: false,
	}

	// should fail
	_ = StructWithPrivateFields{
		integer: 0,
	}
}

func _() {
	_ = StructWithEmbedded{
		RegularStruct: RegularStruct{},
		SomeString:    "",
	}
	_ = StructWithEmbedded{
		RegularStruct{},
		"",
	}

	// should fail
	_ = StructWithEmbedded{
		SomeString: "",
	}

	// should fail
	_ = StructWithEmbedded{
		RegularStruct: RegularStruct{},
	}

	// should fail
	_ = StructWithEmbedded{
		// should fail
		RegularStruct: RegularStruct{
			Integer: 0,
			String:  "",
		},
	}
}

func _() {
	_ = imported.ImportedStruct{Boolean: false}
}

func _() {
	_ = StructWithIgnoreTag{
		integer: 0,
		boolean: false,
		String:  "",
		Uint:    0,
	}

	_ = StructWithIgnoreTag{
		boolean: false,
		String:  "",
	}

	// should fail 2x
	_ = StructWithIgnoreTag{
		integer: 0,
		Uint:    0,
	}
}
