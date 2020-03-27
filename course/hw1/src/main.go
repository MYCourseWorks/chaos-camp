package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Node struct {
	Value int
	Prev  *Node
	Next  *Node
}

type CircularLinkedList struct {
	Head  *Node
	Tail  *Node
	Count int
}

func (l *CircularLinkedList) AddLast(v int) *Node {
	n := &Node{Value: v}
	if l.Count == 0 {
		l.Head = n
		l.Tail = n
		l.Head.Next = n
		l.Head.Prev = n
		l.Tail.Next = n
		l.Tail.Prev = n
	} else if l.Count == 1 {
		l.Tail = n
		l.Head.Next = n
		l.Head.Prev = n
		l.Tail.Next = l.Head
		l.Tail.Prev = l.Head
	} else {
		l.Head.Prev = n
		l.Tail.Next = n
		n.Next = l.Head
		n.Prev = l.Tail
		l.Tail = n
	}

	l.Count++
	return n
}

// NOTE:
// This functions is dangerous, if n is not in the CircularLinkedList behavior is undefined.
// We can't protect for that and it should probably not be exposed to the client.
// Offcourse the advantage is that we can remove nodes in O(1) time.
func (l *CircularLinkedList) Remove(n *Node) int {
	if n == l.Head {
		return l.RemoveFirst()
	} else if n == l.Tail {
		return l.RemoveLast()
	} else {
		// Most common case
		item := n.Value
		prevNode := n.Prev
		nextNode := n.Next
		prevNode.Next = n.Next
		nextNode.Prev = n.Prev
		l.Count--
		return item
	}
}

func (l *CircularLinkedList) RemoveFirst() int {
	if l.Count == 0 {
		panic("No elements in CircularLinkedList.")
	}

	l.Count--
	val := l.Head.Value
	if l.Head == l.Tail {
		l.Head = nil
		l.Tail = nil
	} else {
		newHead := l.Head.Next
		l.Head.Next = nil
		l.Head.Prev = nil
		newHead.Prev = l.Tail
		l.Tail.Next = newHead
		l.Head = newHead
	}

	return val
}

func (l *CircularLinkedList) RemoveLast() int {
	if l.Count == 0 {
		panic("No elements in List.")
	}

	l.Count--
	elem := l.Tail.Value
	if l.Head == l.Tail {
		l.Head = nil
		l.Tail = nil
	} else {
		newTail := l.Tail.Prev
		l.Tail.Next = nil
		l.Tail.Prev = nil
		newTail.Next = l.Head
		l.Head.Prev = newTail
		l.Tail = newTail
	}

	return elem
}

// TODO: If length is negative we should move backwords ??
func MoveForward(n *Node, length int) *Node {
	for i := 0; i < length; i++ {
		if n == nil {
			panic("Null reference. Linked List Chain is broken.")
		}
		n = n.Next
	}
	return n
}

func (l *CircularLinkedList) Traverse(cb func(n *Node)) {
	curr := l.Head
	for i := 0; i < l.Count; i++ {
		cb(curr)
		curr = curr.Next
	}
}

func (l *CircularLinkedList) TraverseBack(cb func(n *Node)) {
	curr := l.Tail
	for i := 0; i < l.Count; i++ {
		cb(curr)
		curr = curr.Prev
	}
}

func (l *CircularLinkedList) PrintList() {
	if l.Count == 0 {
		fmt.Println("(Empty)")
		return
	}

	l.Traverse(func(n *Node) {
		fmt.Printf("%d ", n.Value)
	})
	fmt.Print(" --> ")
	l.TraverseBack(func(n *Node) {
		fmt.Printf("%d ", n.Next.Value)
	})
	fmt.Println()
}

func findWinner(n int, m int) int {
	list := new(CircularLinkedList)

	for i := 0; i < n; i++ {
		list.AddLast(i + 1)
	}

	curr := list.Head
	for list.Count > 1 {
		curr = MoveForward(curr, m-1) // m - 1, because we count the current
		tmp := curr.Next
		// fmt.Println(list.Remove(curr))
		list.Remove(curr)
		curr = tmp
	}

	return list.Head.Value
}

func main() {
	// TODO: Better error checking for user input
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("N = ")
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())

	fmt.Print("M = ")
	scanner.Scan()
	m, _ := strconv.Atoi(scanner.Text())

	fmt.Print("P = ")
	fmt.Println(findWinner(n, m))
}
