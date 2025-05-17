package rotationsystems

import (
	"github.com/yhlooo/go-tetris/pkg/tetris/common"
)

// RotationSystem 旋转系统
type RotationSystem interface {
	// RotateRight 将场上活跃方块顺时针旋转 90 度
	RotateRight(field *common.Field) bool
	// RotateLeft 将场上活跃方块逆时针旋转 90 度
	RotateLeft(field *common.Field) bool
}

// SuperRotationSystem 超级旋转系统（ SRS ）
//
// 参考 https://harddrop.com/wiki/SRS
type SuperRotationSystem struct{}

var _ RotationSystem = SuperRotationSystem{}

// srsJLSTZWallKickData SRS 的 J L S T Z 方块踢墙数据
var srsJLSTZWallKickData = map[[2]common.TetrominoDir][]common.Location{
	{common.Dir0, common.DirR}: {{0, 0}, {0, -1}, {+1, -1}, {-2, 0}, {-2, -1}},
	{common.DirR, common.Dir0}: {{0, 0}, {0, +1}, {-1, +1}, {+2, 0}, {+2, +1}},
	{common.DirR, common.Dir2}: {{0, 0}, {0, +1}, {-1, +1}, {+2, 0}, {+2, +1}},
	{common.Dir2, common.DirR}: {{0, 0}, {0, -1}, {+1, -1}, {-2, 0}, {-2, -1}},
	{common.Dir2, common.DirL}: {{0, 0}, {0, +1}, {+1, +1}, {-2, 0}, {-2, +1}},
	{common.DirL, common.Dir2}: {{0, 0}, {0, -1}, {-1, -1}, {+2, 0}, {+2, -1}},
	{common.DirL, common.Dir0}: {{0, 0}, {0, -1}, {-1, -1}, {+2, 0}, {+2, -1}},
	{common.Dir0, common.DirL}: {{0, 0}, {0, +1}, {+1, +1}, {-2, 0}, {-2, +1}},
}

// srsIWallKickData SRS 的 I 方块踢墙数据
var srsIWallKickData = map[[2]common.TetrominoDir][]common.Location{
	{common.Dir0, common.DirR}: {{0, 0}, {0, -2}, {0, +1}, {-1, -2}, {+2, +1}},
	{common.DirR, common.Dir0}: {{0, 0}, {0, +2}, {0, -1}, {+1, +2}, {-2, -1}},
	{common.DirR, common.Dir2}: {{0, 0}, {0, -1}, {0, +2}, {+2, -1}, {-1, +2}},
	{common.Dir2, common.DirR}: {{0, 0}, {0, +1}, {0, -2}, {-2, +1}, {+1, -2}},
	{common.Dir2, common.DirL}: {{0, 0}, {0, +2}, {0, -1}, {+1, +2}, {-2, -1}},
	{common.DirL, common.Dir2}: {{0, 0}, {0, -2}, {0, +1}, {-1, -2}, {+2, +1}},
	{common.DirL, common.Dir0}: {{0, 0}, {0, +1}, {0, -2}, {-2, +1}, {+1, -2}},
	{common.Dir0, common.DirL}: {{0, 0}, {0, -1}, {0, +2}, {+2, -1}, {-1, +2}},
}

// RotateRight 将场上活跃方块顺时针旋转 90 度
func (srs SuperRotationSystem) RotateRight(field *common.Field) bool {
	return srs.rotate(field, 1)
}

// RotateLeft 将场上活跃方块逆时针旋转 90 度
func (srs SuperRotationSystem) RotateLeft(field *common.Field) bool {
	return srs.rotate(field, -1)
}

// rotate 旋转
func (SuperRotationSystem) rotate(field *common.Field, dir int) bool {
	tetromino := field.ActiveTetromino()
	if tetromino == nil {
		return false
	}

	oldDir := tetromino.Dir
	newDir := common.TetrominoDir(int(oldDir)+dir) % 4
	oldRow := tetromino.Row
	oldCol := tetromino.Column

	// 确定踢墙数据
	var wallKickData []common.Location
	switch tetromino.Type {
	case common.J, common.L, common.S, common.T, common.Z:
		wallKickData = srsJLSTZWallKickData[[2]common.TetrominoDir{oldDir, newDir}]
	case common.I:
		wallKickData = srsIWallKickData[[2]common.TetrominoDir{oldDir, newDir}]
	default:
	}
	if wallKickData == nil {
		wallKickData = []common.Location{{0, 0}}
	}

	// 尝试旋转
	tetromino.Dir = newDir
	for _, wallKick := range wallKickData {
		tetromino.Row = oldRow + wallKick.Row()
		tetromino.Column = oldCol + wallKick.Column()
		if field.IsValid() {
			// 旋转成功
			return true
		}
	}

	// 旋转失败，还原
	tetromino.Dir = oldDir
	tetromino.Row = oldRow
	tetromino.Column = oldCol
	return false
}
