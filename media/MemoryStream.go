package media

import (
	"errors"
	"io/ioutil"
	"os"
)

// MemoryStream - is an inteface to a memory stream
type MemoryStream interface {
	GetData() []byte
	GetPosition() int
	SetPosition(position int)
	GetSize() int
	Append(data []byte) (int, error)
	Read(count int) ([]byte, error)
	Write(buffer []byte, count int) (int, error)
	SaveToFile(fileName string) error
	Clear()
}

type memoryStream struct {
	Data     []byte
	Position int
	Size     int
}

// NewEmptyMemoryStream - Creates an inteface to a new empty memoryStream
func NewEmptyMemoryStream() MemoryStream {
	return &memoryStream{
		Data:     make([]byte, 0),
		Position: 0,
		Size:     0,
	}
}

// NewMemoryStreamFromBytes - Creates an interface to a new memoryStream from bytes
func NewMemoryStreamFromBytes(data []byte) MemoryStream {
	return &memoryStream{
		Data:     data,
		Position: 0,
		Size:     len(data),
	}
}

// NewMemoryStreamFromFile - Creates an interface to a new memoryStream from file
func NewMemoryStreamFromFile(fileName string) (MemoryStream, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return &memoryStream{
		Data:     data,
		Position: 0,
		Size:     len(data),
	}, nil
}

func (ms *memoryStream) GetData() []byte {
	return ms.Data
}

func (ms *memoryStream) GetPosition() int {
	return ms.Position
}

func (ms *memoryStream) SetPosition(position int) {
	ms.Position = position
}

func (ms *memoryStream) GetSize() int {
	return ms.Size
}

// Read - Read from MemoryStream into Buffer count bytes
func (ms *memoryStream) Read(count int) ([]byte, error) {
	buffer := make([]byte, count)
	if count+ms.Position > ms.Size {
		return nil, errors.New("ERROR, MemoryStream::Read, count+ms.Position > ms.Size")
	}
	copy(buffer, ms.Data[ms.Position:ms.Position+count])
	ms.Position = ms.Position + count
	return buffer, nil
}

func (ms *memoryStream) Append(data []byte) (int, error) {
	count := len(data)
	if count == 0 {
		return -1, errors.New("ERROR, MemoryStream::Append, nothing to write")
	}

	if ms.Data == nil {
		return -1, errors.New("ERROR, MemoryStream:::Append, no Data to append to")
	}

	ms.Data = append(ms.Data, data...)

	return count, nil
}

// Write - Write from Buffer into MemoryStream count bytes
func (ms *memoryStream) Write(buffer []byte, count int) (int, error) {
	if len(buffer) == 0 {
		return -1, errors.New("ERROR, MemoryStream::Write, nothing to write")
	}

	if ms.Data == nil {
		return -1, errors.New("ERROR, MemoryStream:::Write, no Data to append to")
	}

	if ms.Position >= ms.Size {
		ms.Data = append(ms.Data, buffer...)
		ms.Size = ms.Size + count
	} else {
		temp := ms.Data[:ms.Position]
		temp = append(temp, buffer[:count]...)
		temp = append(temp, ms.Data[ms.Position+count:ms.Size]...)
		copy(ms.Data, temp)
	}
	ms.Position = ms.Position + count
	return count, nil
}

// SaveToFile - Save MemoryStream to File
func (ms *memoryStream) SaveToFile(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return errors.New("ERROR, MemoryStram::SaveToFile, " + err.Error())
	}
	defer file.Close()

	if bs, err := ms.Read(ms.Size); err == nil {
		_, err = file.Write(bs)
		if err != nil {
			return errors.New("ERROR, MemoryStram::SaveToFile, " + err.Error())
		}
		return nil
	}
	return errors.New("ERROR, MemoryStram::SaveToFile, failed to read ms.buffer")
}

// Clear - Clears the memoryStream
func (ms *memoryStream) Clear() {
	ms.Data = ms.Data[:0]
	ms.Position = 0
	ms.Size = 0
}
