package main

import (
	"golang.org/x/exp/constraints"
)

const (
	avlLeftHeavy  = -1
	avlBalanced   = 0
	avlRightHeavy = +1
)

type AVLTree[T constraints.Ordered] struct {
	*BaseTree[T]
}

func NewAVLTree[T constraints.Ordered]() *AVLTree[T] {
	return &AVLTree[T]{BaseTree: NewBaseTree[T]()}
}

func (t *AVLTree[T]) Insert(z *Node[T]) {
	// 默认为平衡节点
	z.extra = avlBalanced
	//z.parent, z.left, z.right = nil, nil, nil
	// 获取哨兵节点 这里的哨兵节点为头节点
	x, childIsLeft := t.End(), true

	// 如果不是空树则遍历至插入节点的前序节点
	for y := x.left; y != nil; {
		AddSearchCounter(1)

		x, childIsLeft = y, z.data < y.data

		if childIsLeft {
			y = y.left
		} else {
			y = y.right
		}
	}

	// 设置 z 的父节点为前序节点，空树时前序节点为哨兵节点
	z.parent = x

	// 将z绑定至前序节点，空树时绑定至哨兵的左子节点
	if childIsLeft {
		x.left = z
	} else {
		x.right = z
	}

	// 将start设置为最左边的叶子节点
	if t.start.left != nil {
		t.start = t.start.left
	}

	// 平衡树
	t.balanceAfterInsert(x, childIsLeft)
	// 树大小加一
	t.size++
}

// balanceAfterInsert 在 AVL 树中插入新节点后对树进行平衡调整
func (t *AVLTree[T]) balanceAfterInsert(x *Node[T], childIsLeft bool) {
	// 从插入位置开始向根节点遍历
	for ; x != t.End(); x = x.parent {
		// 记录平衡调整的次数
		AddFixupCounter(1)

		// 根据新节点的位置(左子树或右子树)进行不同的平衡调整
		if !childIsLeft {
			// 新节点在右子树
			switch x.extra {
			case avlLeftHeavy:
				// 如果当前节点的平衡因子是左侧重,则将其设为平衡
				x.setExtra(avlBalanced)
				return
			case avlRightHeavy:
				// 如果当前节点的平衡因子是右侧重,则需要进行旋转操作
				if x.right.extra == avlLeftHeavy {
					// 先进行右旋,再进行左旋
					avlRotateRightLeft(x)
				} else {
					// 直接进行左旋
					avlRotateLeft(x)
				}
				return
			default:
				// 如果当前节点的平衡因子是平衡,则将其设为右侧重
				x.setExtra(avlRightHeavy)
			}
		} else {
			// 新节点在左子树
			switch x.extra {
			case avlRightHeavy:
				// 如果当前节点的平衡因子是右侧重,则将其设为平衡
				x.setExtra(avlBalanced)
				return
			case avlLeftHeavy:
				// 如果当前节点的平衡因子是左侧重,则需要进行旋转操作
				if x.left.extra == avlRightHeavy {
					// 先进行左旋,再进行右旋
					avlRotateLeftRight(x)
				} else {
					// 直接进行右旋
					avlRotateRight(x)
				}
				return
			default:
				// 如果当前节点的平衡因子是平衡,则将其设为左侧重
				x.setExtra(avlLeftHeavy)
			}
		}

		// 更新 childIsLeft 的值,准备继续向根节点遍历
		childIsLeft = x == x.parent.left
	}
}

// Delete 从 AVL 树中删除给定的节点
func (t *AVLTree[T]) Delete(z *Node[T]) {
	// 如果被删除节点是根节点,则将根节点指向该节点的后继节点
	if t.start == z {
		t.start = z.next()
	}

	// 获取被删除节点 z 的父节点 x,以及 z 是否是其父节点的左子节点
	x, childIsLeft := z.parent, z == z.parent.left

	// 根据被删除节点 z 的子节点数量,执行不同的删除操作
	switch {
	case z.left == nil:
		// 如果 z 只有右子节点,则用右子节点替换 z
		transplant(z, z.right)
	case z.right == nil:
		// 如果 z 只有左子节点,则用左子节点替换 z
		transplant(z, z.left)
	default:
		// 如果 z 有左右两个子节点,则找到 z 的中序后继节点 y 来替换 z
		if z.extra == avlRightHeavy {
			// 如果 z 的平衡因子是右重,则找到 z 右子树的最小节点 y 来替换 z
			y := minimum(z.right)
			x, childIsLeft = y, y == y.parent.left

			if y.parent != z {
				// 如果 y 不是 z 的直接子节点,则需要先将 y 从其原位置删除
				x = y.parent
				transplant(y, y.right)
				y.right = z.right
				y.right.parent = y
			}

			// 用 y 替换 z
			transplant(z, y)
			y.left = z.left
			y.left.parent = y
			y.extra = z.extra
		} else {
			// 如果 z 的平衡因子是左重,则找到 z 左子树的最大节点 y 来替换 z
			y := maximum(z.left)
			x, childIsLeft = y, y == y.parent.left

			if y.parent != z {
				// 如果 y 不是 z 的直接子节点,则需要先将 y 从其原位置删除
				x = y.parent
				transplant(y, y.left)
				y.left = z.left
				y.left.parent = y
			}

			// 用 y 替换 z
			transplant(z, y)
			y.right = z.right
			y.right.parent = y
			y.extra = z.extra
		}
	}

	// 调整 AVL 树的平衡性
	t.balanceAfterDelete(x, childIsLeft)
	// 树的大小减 1
	t.size--
}

// balanceAfterDelete 在删除节点后调整 AVL 树的平衡性
func (t *AVLTree[T]) balanceAfterDelete(x *Node[T], childIsLeft bool) {
	// 从被删除节点的父节点开始,沿着树的上溯路径调整平衡性
	for ; x != t.End(); x = x.parent {
		// 记录平衡调整的次数
		AddFixupCounter(1)

		if childIsLeft {
			// 如果被删除节点是其父节点的左子节点
			switch x.extra {
			case avlBalanced:
				// 如果父节点是平衡的,则将其设为右重
				x.setExtra(avlRightHeavy)
				return
			case avlRightHeavy:
				// 如果父节点是右重,则需要进行旋转调整
				b := x.right.extra
				if b == avlLeftHeavy {
					// 如果右子节点是左重,则进行左右旋
					avlRotateRightLeft(x)
				} else {
					// 否则进行左旋
					avlRotateLeft(x)
				}
				// 如果旋转后父节点变平衡,则调整结束
				if b == avlBalanced {
					return
				}
				// 否则继续向上调整
				x = x.parent
			default:
				// 如果父节点是左重,则将其设为平衡
				x.setExtra(avlBalanced)
			}
		} else {
			// 如果被删除节点是其父节点的右子节点
			switch x.extra {
			case avlBalanced:
				// 如果父节点是平衡的,则将其设为左重
				x.setExtra(avlLeftHeavy)
				return
			case avlLeftHeavy:
				// 如果父节点是左重,则需要进行旋转调整
				b := x.left.extra
				if b == avlRightHeavy {
					// 如果左子节点是右重,则进行左右旋
					avlRotateLeftRight(x)
				} else {
					// 否则进行右旋
					avlRotateRight(x)
				}
				// 如果旋转后父节点变平衡,则调整结束
				if b == avlBalanced {
					return
				}
				// 否则继续向上调整
				x = x.parent
			default:
				// 如果父节点是右重,则将其设为平衡
				x.setExtra(avlBalanced)
			}
		}

		// 更新被删除节点在其父节点中的位置
		childIsLeft = x == x.parent.left
	}
}

// avlRotateLeft 对 AVL 树中的一个节点进行左旋操作
func avlRotateLeft[T constraints.Ordered](x *Node[T]) {
	// 记录左旋操作的次数
	AddRotateCounter(1)

	// 这里不考虑有左子节点的情况，因为有左子节点时该节点是平衡的

	// 获取当前节点 x 的右子节点 z
	z := x.right
	// 将 x 的右子节点指向 z 的左子节点
	x.right = z.left

	// 如果 z 的左子节点不为空,则将其父节点指向 x
	if z.left != nil {
		z.left.parent = x
	}

	// 将 z 的父节点指向 x 的父节点
	z.parent = x.parent

	// 如果 x 是其父节点的左子节点,则将父节点的左子节点指向 z
	// 否则将父节点的右子节点指向 z
	if x == x.parent.left {
		x.parent.left = z
	} else {
		x.parent.right = z
	}

	// 将 x 设为 z 的左子节点
	z.left = x
	// 将 x 的父节点指向 z
	x.parent = z

	// 根据 z 的平衡因子调整 x 和 z 的平衡因子
	if z.extra == avlBalanced {
		x.setExtra(avlRightHeavy)
		z.setExtra(avlLeftHeavy)
	} else {
		x.setExtra(avlBalanced)
		z.setExtra(avlBalanced)
	}
}

// avlRotateRight 对 AVL 树中的一个节点进行右旋操作
func avlRotateRight[T constraints.Ordered](x *Node[T]) {
	// 记录右旋操作的次数
	AddRotateCounter(1)

	// 这里不考虑有右子节点的情况，因为有右子节点时该节点是平衡的

	// 获取当前节点 x 的左子节点 z
	z := x.left
	// 将 x 的左子节点指向 z 的右子节点
	x.left = z.right

	// 如果 z 的右子节点不为空,则将其父节点指向 x
	if z.right != nil {
		z.right.parent = x
	}

	// 将 z 的父节点指向 x 的父节点
	z.parent = x.parent

	// 如果 x 是其父节点的右子节点,则将父节点的右子节点指向 z
	// 否则将父节点的左子节点指向 z
	if x == x.parent.right {
		x.parent.right = z
	} else {
		x.parent.left = z
	}

	// 将 x 设为 z 的右子节点
	z.right = x
	// 将 x 的父节点指向 z
	x.parent = z

	// 根据 z 的平衡因子调整 x 和 z 的平衡因子
	if z.extra == avlBalanced {
		x.setExtra(avlLeftHeavy)
		z.setExtra(avlRightHeavy)
	} else {
		x.setExtra(avlBalanced)
		z.setExtra(avlBalanced)
	}
}

// avlRotateRightLeft 对 AVL 树中的一个节点执行右左旋操作
func avlRotateRightLeft[T constraints.Ordered](x *Node[T]) {
	// 记录右左旋操作的次数
	AddRotateCounter(2)

	// 获取当前节点 x 的右子节点 z
	z := x.right
	// 获取 z 的左子节点 y
	y := z.left
	// 将 z 的左子节点指向 y 的右子节点
	z.left = y.right

	// 如果 y 的右子节点不为空,则将其父节点指向 z
	if y.right != nil {
		y.right.parent = z
	}

	// 将 y 的右子节点指向 z
	y.right = z
	// 将 z 的父节点指向 y
	z.parent = y
	// 将 x 的右子节点指向 y 的左子节点
	x.right = y.left

	// 如果 y 的左子节点不为空,则将其父节点指向 x
	if y.left != nil {
		y.left.parent = x
	}

	// 将 y 的父节点指向 x 的父节点
	y.parent = x.parent

	// 如果 x 是其父节点的左子节点,则将父节点的左子节点指向 y
	// 否则将父节点的右子节点指向 y
	if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	// 将 y 的左子节点指向 x
	y.left = x
	// 将 x 的父节点指向 y
	x.parent = y

	// 根据 y 的平衡因子调整 x、y 和 z 的平衡因子
	switch y.extra {
	case avlRightHeavy:
		x.setExtra(avlLeftHeavy)
		y.setExtra(avlBalanced)
		z.setExtra(avlBalanced)
	case avlLeftHeavy:
		x.setExtra(avlBalanced)
		y.setExtra(avlBalanced)
		z.setExtra(avlRightHeavy)
	default:
		x.setExtra(avlBalanced)
		z.setExtra(avlBalanced)
	}
}

// avlRotateLeftRight 对 AVL 树中的一个节点执行左右旋操作
func avlRotateLeftRight[T constraints.Ordered](x *Node[T]) {
	// 记录左右旋操作的次数
	AddRotateCounter(2)

	// 开始左旋转
	// 获取当前节点 x 的左子节点 z
	z := x.left
	// 获取 z 的右子节点 y
	y := z.right
	// 将 z 的右子节点指向 y 的左子节点
	z.right = y.left

	// 如果 y 的左子节点不为空,则将其父节点指向 z
	if y.left != nil {
		y.left.parent = z
	}

	// 将 y 的左子节点指向 z
	y.left = z
	// 将 z 的父节点指向 y
	z.parent = y

	// 开始右旋转
	// 将 x 的左子节点指向 y 的右子节点
	x.left = y.right

	// 如果 y 的右子节点不为空,则将其父节点指向 x
	if y.right != nil {
		y.right.parent = x
	}

	// 将 y 的父节点指向 x 的父节点
	y.parent = x.parent

	// 如果 x 是其父节点的右子节点,则将父节点的右子节点指向 y
	// 否则将父节点的左子节点指向 y
	if x == x.parent.right {
		x.parent.right = y
	} else {
		x.parent.left = y
	}

	// 将 y 的右子节点指向 x
	y.right = x
	// 将 x 的父节点指向 y
	x.parent = y

	// 根据 y 的平衡因子调整 x、y 和 z 的平衡因子
	switch y.extra {
	case avlLeftHeavy:
		x.setExtra(avlRightHeavy)
		y.setExtra(avlBalanced)
		z.setExtra(avlBalanced)
	case avlRightHeavy:
		x.setExtra(avlBalanced)
		y.setExtra(avlBalanced)
		z.setExtra(avlLeftHeavy)
	default:
		x.setExtra(avlBalanced)
		z.setExtra(avlBalanced)
	}
}
