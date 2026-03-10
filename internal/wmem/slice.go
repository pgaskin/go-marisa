package wmem

type sliceMemory struct {
	buf []byte
	max int32
}

// SliceMemory allocates a movable slice-backed memory. If allocation fails,
// it panics.
func SliceMemory(cap, max uint32) Memory {
	return &sliceMemory{
		buf: make([]byte, 0, cap),
		max: Pages(max),
	}
}

func (m *sliceMemory) Data() *[]byte {
	return &m.buf
}

func (m *sliceMemory) Grow(delta, _ int32) int32 {
	len := len(m.buf)
	old := int32(len >> PageBits)
	if delta == 0 {
		return old
	}
	new := old + delta
	add := int(new)<<PageBits - len
	if new > m.max || add < 0 {
		return -1
	}
	m.buf = append(m.buf, make([]byte, add)...)
	return old
}

func (m *sliceMemory) Free() {
	m.buf = nil
}
