package datacenter

const BufferSize = 512

// Buffer a buffer of fixed size.
type Buffer struct {
	Data 			[]byte
	maxSize			int
	BytesWritten	int			// how many bytes have been written into the current buffer
}

func NewBuffer(maxSize int) *Buffer {
	return &Buffer{
		Data:         make([]byte, maxSize),
		maxSize:      maxSize,
		BytesWritten: 0,
	}
}

// Bytes makes a copy and returns the bytes in the buffer.
func (buffer *Buffer) Bytes() []byte {
	r := make([]byte, buffer.BytesWritten)
	copy(r, buffer.Data)
	return r
}

// Clear clears the data in the buffer; note it does not erase the data immediately; it simply resets the pointer.
func (buffer *Buffer) Clear() {
	buffer.BytesWritten = 0
}

// ReadTillFull reads a stream of bytes into the buffer until the buffer is full, and returns the number of bytes read.
func (buffer *Buffer) ReadTillFull(data []byte) int {
	l := len(data)
	read := 0
	for read < l && read < len(data) && buffer.BytesWritten < buffer.maxSize {
		buffer.Data[buffer.BytesWritten] = data[read]
		read ++
		buffer.BytesWritten ++
	}
	return read
}

// IsFull checks if the buffer is full.
func (buffer *Buffer) IsFull() bool {
	return buffer.BytesWritten == buffer.maxSize
}
