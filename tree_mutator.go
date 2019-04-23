package sgf

// MutFunc is the type of function accepted by MutateTree.
type MutFunc func(original *Node, boardsize int) *Node

// MutateTree creates a new tree that is isomorphic to the original tree. The
// only argument is a function which examines each node and returns a new node
// to place in the new tree. The node returned by the mutator function must have
// no parent and no children. The node returned by MutateTree is the one in the
// new tree which is equivalent to the node on which MutateTree was called.
func (self *Node) MutateTree(mutator MutFunc) *Node {

	root := self.GetRoot()
	boardsize := root.RootBoardSize()

	// We mutate the entire tree but we want to return the node that's equivalent to self.
	// To accomplish this, mutate_recursive() gets a pointer to a pointer which it can set
	// when it sees that it is mutating self, which is the initial value of that pointer.

	foo := self

	mutate_recursive(root, boardsize, mutator, &foo)

	if foo == self {
		panic("Node.MutateTree(): failed to set equivalent node, this is normally impossible")
	}

	return foo
}

func mutate_recursive(node *Node, boardsize int, mutator MutFunc, foo **Node) *Node {

	mutant := mutator(node, boardsize)

	if mutant == nil || mutant == node || mutant.parent != nil || len(mutant.children) > 0 {
		panic("mutate_recursive(): the mutator function returned an improper node")
	}

	// foo starts off as the node whose mutant we ultimately want to return at the top level.
	// When we actually see that node, we set foo to be the mutant. See note in MutateTree().
	// This is a slightly-too-cute way of doing it.

	if node == *foo {
		*foo = mutant
	}

	for _, child := range(node.children) {
		mutant_child := mutate_recursive(child, boardsize, mutator, foo)
		mutant_child.parent = mutant
		mutant.children = append(mutant.children, mutant_child)
	}

	return mutant
}
