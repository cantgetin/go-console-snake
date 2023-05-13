package game

type Stack struct {
	items [][]int
}

func (s *Stack) Push(item []int) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() []int {
	if len(s.items) == 0 {
		return nil
	}
	item := s.items[0]
	s.items = s.items[1:]
	return item
}

func (s *Stack) Size() int {
	return len(s.items)
}

func (s *Stack) Iter() <-chan []int {
	c := make(chan []int)
	go func() {
		for _, item := range s.items {
			c <- item
		}
		close(c)
	}()
	return c
}
