package terminal

import "testing"

func TestTabManager(t *testing.T) {
	tm := NewTabManager()
	tab := tm.AddTab("test", nil)
	if tab.Title != "test" { t.Error("title mismatch") }
	if tm.GetActiveTab() != tab { t.Error("active tab mismatch") }
}

func TestBlockManager(t *testing.T) {
	bm := NewBlockManager()
	b := bm.CreateBlock("b1", "t1", "view")
	if bm.GetBlock("b1") != b { t.Error("block retrieval fail") }
}
