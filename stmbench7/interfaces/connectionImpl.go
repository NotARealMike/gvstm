package interfaces

type connectionImpl struct {
    typ string
    length int
    from, to AtomicPart
}

func NewConnectionImpl(from, to AtomicPart, typ string, length int) Connection  {
    return &connectionImpl{
        typ:    typ,
        length: length,
        from:   from,
        to:     to,
    }
}

func (c *connectionImpl) GetReversed() Connection {
    return &connectionImpl{
        typ:    c.typ,
        length: c.length,
        from:   c.to,
        to:     c.from,
    }
}

func (c *connectionImpl) GetDestination() AtomicPart {
    return c.to
}

func (c *connectionImpl) GetSource() AtomicPart {
    return c.from
}

func (c *connectionImpl) GetType() string {
    return c.typ
}

func (c *connectionImpl) GetLength() int {
    return c.length
}
