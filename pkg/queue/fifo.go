package queue

type FIFO struct {
	queue Queue[any]
}

func (q *FIFO) Enqueue(it any) {

}

func (q *FIFO) Dequeue() any {
	return nil
}

// l2 := list.New()
// u1 := l2.PushFront(1)
// u2 := l2.PushFront(2)
// u3 := l2.PushFront(3)
// u4 := l2.PushFront(4)
// fmt.Println(l2.Remove(u1))
// fmt.Println(l2.Remove(u2))
// fmt.Println(l2.Remove(u3))
// fmt.Println(l2.Remove(u4))
