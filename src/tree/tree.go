package tree

import "container/list"

type node struct {
	p     int64 // Количество раз, сколько встретился определенный байт
	left  *node
	right *node
	//code
}

type tree struct {
	list.List
}
