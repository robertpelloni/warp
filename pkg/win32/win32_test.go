//go:build windows

package win32

import (
	"testing"
	"unsafe"
)

func TestRichEditConstantsAndStructs(t *testing.T) {
	cf := CHARFORMAT2W{}
	cf.CbSize = uint32(unsafe.Sizeof(cf))
	if cf.CbSize == 0 {
		t.Error("expected non-zero size for CHARFORMAT2W")
	}

	cr := CHARRANGE{CpMin: 0, CpMax: 10}
	if cr.CpMin != 0 || cr.CpMax != 10 {
		t.Errorf("unexpected CHARRANGE values: %+v", cr)
	}

	if MSFTEDIT_CLASS != "RICHEDIT50W" {
		t.Errorf("unexpected MSFTEDIT_CLASS: %s", MSFTEDIT_CLASS)
	}
}
