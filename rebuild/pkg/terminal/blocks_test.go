package terminal

import "testing"

func TestBlockManager(t *testing.T) {
	bm := NewBlockManager()
	b := bm.CreateBlock("block1", "tab1", "terminal")

	if b.ID != "block1" {
		t.Errorf("Expected block ID 'block1', got %s", b.ID)
	}

	retrieved := bm.GetBlock("block1")
	if retrieved != b {
		t.Error("Expected to retrieve the same block")
	}

	bm.DeleteBlock("block1")
	if bm.GetBlock("block1") != nil {
		t.Error("Expected block to be deleted")
	}
}
