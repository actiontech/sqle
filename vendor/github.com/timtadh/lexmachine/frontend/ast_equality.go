package frontend

// Equals checks deep equality of the two trees
func (c *Concat) Equals(o AST) bool {
	if x, is := o.(*Concat); is {
		if len(c.Items) != len(x.Items) {
			return false
		}
		for i := range c.Items {
			if !c.Items[i].Equals(x.Items[i]) {
				return false
			}
		}
		return true
	}
	return false
}

// Equals checks deep equality of the two trees
func (a *AltMatch) Equals(o AST) bool {
	if x, is := o.(*AltMatch); is {
		return a.A.Equals(x.A) && a.B.Equals(x.B)
	}
	return false
}

// Equals checks deep equality of the two trees
func (a *Alternation) Equals(o AST) bool {
	if x, is := o.(*Alternation); is {
		return a.A.Equals(x.A) && a.B.Equals(x.B)
	}
	return false
}

// Equals checks deep equality of the two trees
func (m *Match) Equals(o AST) bool {
	if x, is := o.(*Match); is {
		return m.AST.Equals(x.AST)
	}
	return false
}

// Equals checks deep equality of the two trees
func (s *Star) Equals(o AST) bool {
	if x, is := o.(*Star); is {
		return s.AST.Equals(x.AST)
	}
	return false
}

// Equals checks deep equality of the two trees
func (p *Plus) Equals(o AST) bool {
	if x, is := o.(*Plus); is {
		return p.AST.Equals(x.AST)
	}
	return false
}

// Equals checks deep equality of the two trees
func (m *Maybe) Equals(o AST) bool {
	if x, is := o.(*Maybe); is {
		return m.AST.Equals(x.AST)
	}
	return false
}

// Equals checks deep equality of the two trees
func (c *Character) Equals(o AST) bool {
	if x, is := o.(*Character); is {
		return *c == *x
	}
	return false
}

// Equals checks deep equality of the two trees
func (r *Range) Equals(o AST) bool {
	if x, is := o.(*Range); is {
		return *r == *x
	}
	return false
}

// Equals checks deep equality of the two trees
func (e *EOS) Equals(o AST) bool {
	if x, is := o.(*EOS); is {
		return *e == *x
	}
	return false
}
