package tetris

// RotationSystem 旋转系统
type RotationSystem interface {
	// RotateRight 将场上活跃方块顺时针旋转 90 度
	RotateRight(field *Field) bool
	// RotateLeft 将场上活跃方块逆时针旋转 90 度
	RotateLeft(field *Field) bool
}

// SuperRotationSystem 超级旋转系统（ SRS ）
//
// 参考 https://harddrop.com/wiki/SRS
type SuperRotationSystem struct{}

var _ RotationSystem = SuperRotationSystem{}

// srsJLSTZWallKickData SRS 的 J L S T Z 方块踢墙数据
var srsJLSTZWallKickData = map[[2]BlockDir][]Location{
	{Dir0, DirR}: {{0, 0}, {0, -1}, {+1, -1}, {-2, 0}, {-2, -1}},
	{DirR, Dir0}: {{0, 0}, {0, +1}, {-1, +1}, {+2, 0}, {+2, +1}},
	{DirR, Dir2}: {{0, 0}, {0, +1}, {-1, +1}, {+2, 0}, {+2, +1}},
	{Dir2, DirR}: {{0, 0}, {0, -1}, {+1, -1}, {-2, 0}, {-2, -1}},
	{Dir2, DirL}: {{0, 0}, {0, +1}, {+1, +1}, {-2, 0}, {-2, +1}},
	{DirL, Dir2}: {{0, 0}, {0, -1}, {-1, -1}, {+2, 0}, {+2, -1}},
	{DirL, Dir0}: {{0, 0}, {0, -1}, {-1, -1}, {+2, 0}, {+2, -1}},
	{Dir0, DirL}: {{0, 0}, {0, +1}, {+1, +1}, {-2, 0}, {-2, +1}},
}

// srsIWallKickData SRS 的 I 方块踢墙数据
var srsIWallKickData = map[[2]BlockDir][]Location{
	{Dir0, DirR}: {{0, 0}, {0, -2}, {0, +1}, {-1, -2}, {+2, +1}},
	{DirR, Dir0}: {{0, 0}, {0, +2}, {0, -1}, {+1, +2}, {-2, -1}},
	{DirR, Dir2}: {{0, 0}, {0, -1}, {0, +2}, {+2, -1}, {-1, +2}},
	{Dir2, DirR}: {{0, 0}, {0, +1}, {0, -2}, {-2, +1}, {+1, -2}},
	{Dir2, DirL}: {{0, 0}, {0, +2}, {0, -1}, {+1, +2}, {-2, -1}},
	{DirL, Dir2}: {{0, 0}, {0, -2}, {0, +1}, {-1, -2}, {+2, +1}},
	{DirL, Dir0}: {{0, 0}, {0, +1}, {0, -2}, {-2, +1}, {+1, -2}},
	{Dir0, DirL}: {{0, 0}, {0, -1}, {0, +2}, {+2, -1}, {-1, +2}},
}

// RotateRight 将场上活跃方块顺时针旋转 90 度
func (srs SuperRotationSystem) RotateRight(field *Field) bool {
	return srs.rotate(field, 1)
}

// RotateLeft 将场上活跃方块逆时针旋转 90 度
func (srs SuperRotationSystem) RotateLeft(field *Field) bool {
	return srs.rotate(field, -1)
}

// rotate 旋转
func (SuperRotationSystem) rotate(field *Field, dir int) bool {
	block := field.ActiveBlock()
	if block == nil {
		return false
	}

	oldDir := block.Dir
	newDir := BlockDir(int(oldDir)+dir) % 4
	oldRow := block.Row
	oldCol := block.Column

	// 确定踢墙数据
	var wallKickData []Location
	switch block.Type {
	case BlockJ, BlockL, BlockS, BlockT, BlockZ:
		wallKickData = srsJLSTZWallKickData[[2]BlockDir{oldDir, newDir}]
	case BlockI:
		wallKickData = srsIWallKickData[[2]BlockDir{oldDir, newDir}]
	default:
	}
	if wallKickData == nil {
		wallKickData = []Location{{0, 0}}
	}

	// 尝试旋转
	block.Dir = newDir
	for _, wallKick := range wallKickData {
		block.Row = oldRow + wallKick.Row()
		block.Column = oldCol + wallKick.Column()
		if field.IsValid() {
			// 旋转成功
			return true
		}
	}

	// 旋转失败，还原
	block.Dir = oldDir
	block.Row = oldRow
	block.Column = oldCol
	return false
}
