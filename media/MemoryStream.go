package media

import (
	"fmt"
	"os"
)

// MemoryStream is a memory stream
type MemoryStream struct {
	data     []byte
	Position int
	Size     int
}

// Read - Read from MemoryStream into Buffer count bytes
func (ms *MemoryStream) Read(buffer []byte, count int) int {
	if len(buffer) < count {
		return -1
	}
	if count+ms.Position > ms.Size {
		return -1
	}
	copy(buffer, ms.data[ms.Position:ms.Position+count])
	ms.Position = ms.Position + count
	return count
}

// Write - Write from Buffer into MemoryStream count bytes
func (ms *MemoryStream) Write(buffer []byte, count int) int {
	if ms.Position >= ms.Size {
		ms.data = append(ms.data, buffer...)
		ms.Size = ms.Size + count
	} else {
		temp := ms.data[:ms.Position]
		temp = append(temp, buffer[:count]...)
		temp = append(temp, ms.data[ms.Position+count:ms.Size]...)
		copy(ms.data, temp)
	}
	ms.Position = ms.Position + count
	return count
}

// LoadFromFile - Load from File into MemoryStream
func (ms *MemoryStream) LoadFromFile(FileName string) bool {
	flag := false

	file, err := os.Open(FileName)
	if err != nil {
		fmt.Println("ERROR, opening file")
		return flag
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Println("ERROR, getting file Stats")
		return flag
	}

	size := int(stat.Size())
	bs := make([]byte, size)
	_, err = file.Read(bs)
	if err != nil {
		fmt.Println("ERROR, reading file")
		return flag
	}
	ms.Write(bs, size)
	return true
}

// SaveToFile - Save MemoryStream to File
func (ms *MemoryStream) SaveToFile(FileName string) bool {
	flag := false

	file, err := os.Create(FileName)
	if err != nil {
		fmt.Println("ERROR, opening file")
		return flag
	}
	defer file.Close()
	bs := make([]byte, ms.Size)
	if ms.Read(bs, ms.Size) != -1 {
		_, err = file.Write(bs)
		if err != nil {
			fmt.Println("ERROR, writing to file")
			return flag
		}
		return true
	}
	return false
}

func (ms *MemoryStream) Clear() {
	ms.data = ms.data[:0]
	ms.Position = 0
	ms.Size = 0
}
