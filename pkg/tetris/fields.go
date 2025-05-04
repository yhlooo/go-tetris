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

	if !f.IsValid() {
		// 不合法，复原
		f.active.Row -= row
		f.active.Column -= col
		return false
	}

	return true
}

// ChangeActiveBlock 更换活跃方块
//
// 若更换后方块没有超出边界且没有与其他方块重合则更换成功并返回 true ，否则不更换并返回 false
func (f *Field) ChangeActiveBlock(block *Block) bool {
	oldActive := f.active
	f.active = block
	if !f.IsValid() {
		f.active = oldActive
		return false
	}
	return true
}

var tCorners = [4]Location{{0, 0}, {0, 2}, {2, 0}, {2, 2}}

// PinActiveBlock 钉住当前活跃方块清除填满的行然后用新方块替换活跃方块
//
// 若更换方块完后活跃方块没有超出边界且没有与其他方块重合则操作成功并返回 ok=true ，否则不更换方块（但仍执行钉住和清除操作）并返回 ok=false
func (f *Field) PinActiveBlock(newBlock *Block) (tSpin bool, clearLines int, ok bool) {
	// 固定活跃方块
	if f.active != nil {
		if f.active.Type == BlockT {
			// 检查是否 T-Spin
			corners := 0
			for _, cornerLoc := range tCorners {
				row := f.active.Row + cornerLoc.Row()
				col := f.active.Column + cornerLoc.Column()
				if row < 0 || col < 0 || col >= f.cols {
					corners++
					continue
				}
				if cell, _ := f.Block(row, col); cell != BlockNone {
					corners++
					continue
				}
			}
			if corners >= 3 {
				tSpin = true
			}
		}
		for _, cell := range f.active.Cells() {
			_ = f.SetBlock(cell.Row(), cell.Column(), f.active.Type)
		}
	}

	// 清除填满的行
	for i := 0; i < len(f.filled); {
		row := f.filled[i]
		full := true
		for _, cell := range row {
			if cell == BlockNone {
				full = false
				break
			}
		}
		if full {
			f.filled = append(f.filled[:i], f.filled[i+1:]...)
			clearLines++
		} else {
			i++
		}
	}
	for i := 0; i < clearLines; i++ {
		f.filled = append(f.filled, make([]BlockType, f.cols))
	}

	// 更换活跃方块
	return tSpin, clearLines, f.ChangeActiveBlock(newBlock)
}

// IsValid 是否合法
//
// 活跃方块没有超出左右和下边界且不与其他方块重合则返回 true ，否则返回 false
func (f *Field) IsValid() bool {
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
