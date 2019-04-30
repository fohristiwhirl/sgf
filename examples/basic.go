package main

import (
	"fmt"
	sgf ".."
)

func main() {

	// Start a new game tree and get the root node...

	root := sgf.NewTree(19)
	node := root

	// Here we create the ancient Chinese pattern...

	node.AddValue("AB", sgf.Point(3, 3))
	node.AddValue("AB", sgf.Point(15, 15))
	node.AddValue("AW", sgf.Point(15, 3))
	node.AddValue("AW", sgf.Point(3, 15))

	// The normal way to create new nodes is by playing moves.
	// If successful, PlayMove() returns the new node.

	node, err := node.PlayMove("cf")					// "cf" is SGF-speak
	fmt.Printf("%v\n", err)								// Prints nil (no error)

	// We can get an SGF coordinate (e.g. "cf") by calling Point().
	// Note that the coordinate system is zeroth-based, from the top left.

	node, err = node.PlayMove(sgf.Point(2, 5))
	fmt.Printf("%v\n", err)								// Already filled

	// Illegal moves (including suicide and basic ko) will return an error.
	// As a convenience, PlayMove() returns the original node in this case.
	// You may still wish to check for errors...

	node, err = node.PlayMove(sgf.Point(19, 19))
	fmt.Printf("%v\n", err)								// Off-board

	// We can create variations from any node.

	node = node.Parent()
	node.PlayMove(sgf.Point(13, 2))						// Create variation 1
	node.PlayMove(sgf.Point(16, 5))						// Create variation 2

	// Colours are determined intelligently, but we can always force a colour.

	node.PlayMoveColour(sgf.Point(2, 5), sgf.WHITE)		// Create variation 3

	// We can iterate through a node's children.

	for i, child := range node.Children() {
		child.SetValue("C", fmt.Sprintf("Comment %d", i))
	}

	// And we can go down those variations if we wish.
	// (Errors ignored here for simplicity.)

	node, _ = node.PlayMove(sgf.Point(5, 16))			// Create variation 4 and go down it
	node, _ = node.PlayMove(sgf.Point(2, 12))			// ...continue going down it
	node, _ = node.PlayMove(sgf.Point(3, 17))			// ...continue going down it

	// Passes are a thing.
	// Doing the same action on the same node many times just returns the first-created child each time.

	foo := node.Pass()
	bar := node.Pass()									// Does not create a new node
	node = node.Pass()									// Does not create a new node

	fmt.Printf("%v, %v\n", foo == bar, bar == node)		// true, true

	// We can directly manipulate SGF properties...
	// We can also examine the board.

	node.SetValue("C", "White passed. Lets highlight all white stones for some reason...")

	board := node.Board()								// Note that this is a deep copy

	for x := 0; x < board.Size; x++ {
		for y := 0; y < board.Size; y++ {
			if board.State[x][y] == sgf.WHITE {
				node.AddValue("TR", sgf.Point(x, y))
			}
		}
	}

	// It is also possible to directly manage node creation and properties,
	// though this has no legality checks...

	node = sgf.NewNode(node)							// Specify the parent
	node.AddValue("B", "dj")

	// It is possible to edit board-altering properties even if a node has
	// children. All cached boards in descendent nodes will be cleared, and
	// remade as needed.

	root.AddValue("AB", "jj")							// Editing the root.
	node, err = node.PlayMove("jj")						// Trying to play at the current node.
	fmt.Printf("%v\n", err)								// Prints the err: jj is not empty.

	// We can adjust the tree so this branch we are on is the main line...

	node.MakeMainLine()

	// Calling Save() will save the entire tree, regardless of node position.

	node.Save("foo.sgf")

	// We can also load files.

	node, err = sgf.Load("foo.sgf")
}
