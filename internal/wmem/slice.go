package wmem

type SliceMemory struct {
	Buf []byte
	Max int32
}

func (m *SliceMemory) Data() *[]byte {
	return &m.Buf
}

func (m *SliceMemory) Grow(delta, _ int32) int32 {
	len := len(m.Buf)
	old := int32(len >> PageBits)
	if delta == 0 {
		return old
	}
	new := old + delta
	add := int(new)<<PageBits - len
	if new > m.Max || add < 0 {
		return -1
	}
	m.Buf = append(m.Buf, make([]byte, add)...)
	return old
}

func (m *SliceMemory) Close() error {
	m.Buf = nil
	return nil
}
