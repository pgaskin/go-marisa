package wmem

type sliceMemory struct {
	buf []byte
	max int64
}

// SliceMemory allocates a movable slice-backed memory. If allocation fails,
// it panics.
func SliceMemory(cap, max uint64) Memory {
	return &sliceMemory{
		buf: make([]byte, 0, cap),
		max: int64((int64(max) + PageSize - 1) >> PageBits),
	}
}

func (m *sliceMemory) Slice() *[]byte {
	return &m.buf
}

func (m *sliceMemory) Grow(delta, _ int64) int64 {
	blen := len(m.buf)
	old := int64(blen >> PageBits)
	if delta == 0 {
		return old
	}
	new := old + delta
	add := new<<PageBits - int64(blen)
	if new > m.max || add < 0 || add > int64(^uint(0)>>1) {
		return -1
	}
	m.buf = append(m.buf, make([]byte, int(add))...)
	return old
}

func (m *sliceMemory) Free() {
	m.buf = nil
}
