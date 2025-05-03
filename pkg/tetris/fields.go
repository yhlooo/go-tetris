package tetris

// FieldReader 场读出器
type FieldReader interface {
	// Size 获取场大小
	Size() (rows, cols int)
	// Block 获取指定位置已填充的方块类型
	Block(row, col int) (BlockType, bool)
	// BlockWithActiveBlock 获取指定位置的方块类型，包含活跃方块
	BlockWithActiveBlock(row, col int) BlockType
	// ActiveBlock 获取当前活跃方块
	ActiveBlock() *Block
}

// NewField 创建 Field
func NewField(rows, cols int, block *Block) *Field {
	if rows < 2 {
		rows = 2
	}
	if cols < 4 {
		cols = 4
	}

	filled := make([][]BlockType, rows)
	for i := range filled {
		filled[i] = make([]BlockType, cols)
	}
	return &Field{
		rows:   rows,
		cols:   cols,
		active: block,
		filled: filled,
	}
}

// Field 场
//
// 非线程安全
type Field struct {
	// 总行列数
	rows, cols int
	// 当前活跃方块
	active *Block
	// 已填充方块
	filled [][]BlockType
}

var _ FieldReader = &Field{}

// Size 获取场大小
func (f *Field) Size() (rows, cols int) {
	return f.rows, f.cols
}

// BlockWithActiveBlock 获取指定位置的方块类型，包含活跃方块
func (f *Field) BlockWithActiveBlock(row, col int) BlockType {
	ret, _ := f.Block(row, col)
	if ret != BlockNone {
		return ret
	}
	if f.active == nil {
		return BlockNone
	}
	for _, cell := range f.active.Cells() {
		if cell.Row() == row && cell.Column() == col {
			return f.active.Type
		}
	}
	return BlockNone
}

// Block 获取指定位置已填充的方块类型
func (f *Field) Block(row, col int) (BlockType, bool) {
	if row < 0 || len(f.filled) <= row {
		return 0, false
	}
	if col < 0 || len(f.filled[row]) <= col {
		return 0, false
	}
	return f.filled[row][col], true
}

// ActiveBlock 获取当前活跃方块
func (f *Field) ActiveBlock() *Block {
	return f.active
}

// SetBlock 设置指定位置方块类型
func (f *Field) SetBlock(row, col int, blockType BlockType) bool {
	if row < 0 || len(f.filled) <= row {
		return false
	}
	if col < 0 || len(f.filled[row]) <= col {
		return false
	}
	f.filled[row][col] = blockType
	return true
}

// MoveActiveBlock 移动活跃方块
//
// 若移动后方块没有超出边界且没有与其他方块重合则移动成功并返回 true ，否则不移动并返回 false
func (f *Field) MoveActiveBlock(row, col int) bool {
	if f.active == nil {
		return false
	}

	// 移动
	f.active.Row += row
	f.active.Column += col

	if !f.isValid() {
		// 不合法，复原
		f.active.Row -= row
		f.active.Column -= col
		return false
	}

	return true
}

// RotateActiveBlock 旋转活跃方块
//
// 若旋转后方块没有超出边界且没有与其他方块重合则移动成功并返回 true ，否则不旋转并返回 false
//
// TODO: 暂不支持旋转后通过少量平移避开边界或其他方块
func (f *Field) RotateActiveBlock(dir int) bool {
	if f.active == nil {
		return false
	}

	// 旋转
	oldDir := f.active.Dir
	f.active.Dir = BlockDir(int(oldDir)+dir) % 4

	if !f.isValid() {
		// 不合法，复原
		f.active.Dir = oldDir
		return false
	}

	return true
}

// ChangeActiveBlock 更换活跃方块
func (f *Field) ChangeActiveBlock(block *Block) {
	f.active = block
}

// PinActiveBlock 钉住当前活跃方块并用新方块替换
func (f *Field) PinActiveBlock(newBlock *Block) {
	if f.active == nil {
		f.active = newBlock
		return
	}

	for _, cell := range f.active.Cells() {
		_ = f.SetBlock(cell.Row(), cell.Column(), f.active.Type)
	}
	f.active = newBlock
}

// isValid 是否合法
//
// 活跃方块没有超出左右和下边界且不与其他方块重合则返回 true ，否则返回 false
func (f *Field) isValid() bool {
	if f.active == nil {
		return true
	}
	for _, cell := range f.active.Cells() {
		row := cell.Row()
		col := cell.Column()
		// 超出边界
		if row < 0 || col < 0 || col >= f.cols {
			return false
		}
		// 与其它方块重合
		if filled, _ := f.Block(row, col); filled != BlockNone {
			return false
		}
	}
	return true
}
