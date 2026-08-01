package set

type Set struct {
	Name string   `json:"name"`
	Args []string `json:"args,omitempty"`
	Body []*Set   `json:"body,omitempty"`
}

func (s *Set) Find(name string) *Set {
	for _, child := range s.Body {
		if child.Name == name {
			return child
		}
	}
	return nil
}

func (s *Set) FindAll(name string) []*Set {
	var result []*Set
	for _, child := range s.Body {
		if child.Name == name {
			result = append(result, child)
		}
	}
	return result
}

func (s *Set) Walk(fn func(*Set)) {
	fn(s)

	for _, child := range s.Body {
		child.Walk(fn)
	}
}
