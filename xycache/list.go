package xycache

type node[kt comparable] struct {
	key  kt
	prev *node[kt]
	next *node[kt]
}

type list[kt comparable] struct {
	first *node[kt]
	last  *node[kt]
}

func (l *list[kt]) append(item *node[kt]) {
	if item == nil {
		return
	}

	item.prev = l.last
	item.next = nil

	if l.first == nil {
		l.first = item
	}

	if l.last != nil {
		l.last.next = item
	}
	l.last = item
}

func (l *list[kt]) remove(n *node[kt]) {
	if n == nil {
		return
	}

	if n.prev != nil {
		n.prev.next = n.next
	}

	if n.next != nil {
		n.next.prev = n.prev
	}

	if n == l.first {
		l.first = n.next
	}

	if n == l.last {
		l.last = n.prev
	}
}
