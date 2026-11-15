package terminal

import "sync"

type Block struct {
	ID         string
	ParentID   string
	View       string
	Meta       map[string]interface{}
}

type BlockManager struct {
	blocks map[string]*Block
	mu     sync.RWMutex
}

func NewBlockManager() *BlockManager {
	return &BlockManager{blocks: make(map[string]*Block)}
}

func (m *BlockManager) CreateBlock(id, parentID, view string) *Block {
	m.mu.Lock()
	defer m.mu.Unlock()
	block := &Block{ID: id, ParentID: parentID, View: view, Meta: make(map[string]interface{})}
	m.blocks[id] = block
	return block
}

func (m *BlockManager) GetBlock(id string) *Block {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.blocks[id]
}
