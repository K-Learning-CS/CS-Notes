package main

import (
	"errors"
	"golang.org/x/exp/constraints"
)

var ErrUnimplemented = errors.New("unimplemented")

// Node 表示树中的一个节点。
type Node[T constraints.Ordered] struct {
	parent *Node[T] // 父节点指针
	left   *Node[T] // 左子节点指针
	right  *Node[T] // 右子节点指针
	extra  int8     // 附加字段
	data   T        // 数据字段
}

// NewNode 创建一个带有给定数据的新节点。
func NewNode[T constraints.Ordered](data T) *Node[T] {
	return &Node[T]{data: data}
}

// Data 返回节点中存储的数据。
func (n *Node[T]) Data() T {
	return n.data
}

// next 根据中序遍历返回树中的下一个节点。
func (n *Node[T]) next() *Node[T] {
	// 如果当前节点的右节点不为空则后继节点为右子树的最小节点
	if n.right != nil {
		return minimum(n.right)
	}
	// 如果右子树为空则后继节点为最近的左祖先
	x := n
	// 是右祖先则往上一层走
	for x == x.parent.right {
		x = x.parent
	}

	return x.parent
}

// setExtra 设置节点的附加字段。
func (n *Node[T]) setExtra(extra int8) {
	AddExtraCounter(1)

	n.extra = extra
}

// minimum 返回以 x 为根的子树中的最小节点。
func minimum[T constraints.Ordered](x *Node[T]) *Node[T] {
	for x.left != nil {
		x = x.left
	}

	return x
}

// maximum 返回以 x 为根的子树中的最大节点。
func maximum[T constraints.Ordered](x *Node[T]) *Node[T] {
	for x.right != nil {
		x = x.right
	}

	return x
}

// transplant 用节点 y 替换节点 x 在树中的位置。
func transplant[T constraints.Ordered](x *Node[T], y *Node[T]) {
	// 如果 x 是父节点的左节点
	if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	if y != nil {
		y.parent = x.parent
	}
}

// Tree 表示一棵树的数据结构。
type Tree[T constraints.Ordered] interface {
	Size() int
	Empty() bool
	Begin() *Node[T]
	End() *Node[T]
	Clear()
	Find(data T) *Node[T]
	Insert(*Node[T])
	Delete(*Node[T])
}

// BaseTree 是 Tree 接口的基本实现。
type BaseTree[T constraints.Ordered] struct {
	sentinel Node[T]
	start    *Node[T]
	size     int
}

// NewBaseTree 创建一个 BaseTree 的新实例。
func NewBaseTree[T constraints.Ordered]() *BaseTree[T] {
	t := new(BaseTree[T])
	t.start = &t.sentinel

	return t
}

// Size 返回树中的节点数量。
func (t *BaseTree[T]) Size() int {
	return t.size
}

// Empty 检查树是否为空。
func (t *BaseTree[T]) Empty() bool {
	return t.Size() == 0
}

// Begin 返回树中的第一个节点。
func (t *BaseTree[T]) Begin() *Node[T] {
	return t.start
}

// End 返回表示树结尾的哨兵节点。
func (t *BaseTree[T]) End() *Node[T] {
	return &t.sentinel
}

// Clear 从树中移除所有节点。
func (t *BaseTree[T]) Clear() {
	t.End().left = nil
	t.start = t.End()
	t.size = 0
}

// Find 在树中搜索具有给定数据的节点。
// 如果找到，返回该节点；否则返回哨兵节点。
func (t *BaseTree[T]) Find(data T) *Node[T] {
	x := t.End()

	for y := x.left; y != nil; {
		AddSearchCounter(1)

		switch {
		case data < y.data:
			y = y.left
		case y.data < data:
			y = y.right
		default:
			return y
		}
	}

	return x
}

// Insert 将节点插入树中。
func (t *BaseTree[T]) Insert(*Node[T]) {
	panic(ErrUnimplemented)
}

// Delete 从树中删除节点。
func (t *BaseTree[T]) Delete(*Node[T]) {
	panic(ErrUnimplemented)
}
