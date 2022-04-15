package transcoder

import (
	"encoding/binary"
	"fmt"
	"strings"
)

func GetUint32(in []byte, length int) uint32 {
	c := make([]byte, length)
	copy(c, in)
	return binary.LittleEndian.Uint32(c)
}

func ReadSegment(in []byte, out []byte, seg_offset uint32, seg_size uint32, i uint32, rawSize uint32) error {
	var count int8
	out_offset := i * rawSize
	in_offset := seg_offset

	for (out_offset - i*rawSize) < rawSize {
		count = int8(in[in_offset])
		in_offset++
		if count >= 0 {
			copy(out[out_offset:out_offset+uint32(count+1)], in[in_offset:in_offset+uint32(count+1)])
			in_offset += uint32(count + 1)
			out_offset += uint32(count + 1)
		} else {
			if (count <= -1) && (count >= -127) {
				newByte := in[in_offset]
				in_offset++
				for j := uint32(0); j < uint32(-count+1); j++ {
					out[j+out_offset] = newByte
				}
				out_offset += uint32(-count + 1)
				if in_offset-seg_offset > seg_size {
					return fmt.Errorf("ERROR, overflow decoding RLE")
				}
			}
		}
	}
	return nil
}

func RLEdecode(in []byte, out []byte, length uint32, size uint32, PhotoInt string) error {
	var segment_count, offset, i uint32
	var segment_offset [15]uint32
	var segment_length [15]uint32

	offset = 0
	for i := 0; i < 15; i++ {
		segment_offset[i] = 0
		segment_length[i] = 0
	}

	segment_count = GetUint32(in, 4)
	for i := 0; i < 15; i++ {
		offset = offset + 4
		segment_offset[i] = GetUint32(in[offset:], 4)
	}

	segment_offset[segment_count] = length
	temp := make([]byte, size)
	for i = 0; i < segment_count; i++ {
		segment_length[i] = segment_offset[i+1] - segment_offset[i]
		ReadSegment(in, temp, segment_offset[i], segment_length[i], i, size/segment_count)
	}

	offset = size / segment_count
	if (strings.Contains(PhotoInt, "MONO")) && (segment_count == 2) {
		for i = 0; i < size/segment_count; i++ {
			out[2*i] = temp[i+offset]
			out[2*i+1] = temp[i]
		}
	} else if (strings.Contains(PhotoInt, "MONO")) && (segment_count == 1) {
		for i = 0; i < size; i++ {
			out[i] = temp[i]
		}
	} else if (PhotoInt == "YBR_FULL") && (segment_count == 3) {
		var Y, Cb, Cr float32
		for i = 0; i < size/segment_count; i++ {
			Y = float32(temp[i])
			Cb = float32(temp[i+offset])
			Cr = float32(temp[i+2*offset])
			out[3*i] = byte(Y + 1.402*(Cr-128.0))
			out[3*i+1] = byte(Y - 0.344136*(Cb-128.0) - 0.714136*(Cr-128.0))
			out[3*i+2] = byte(Y + 1.772*(Cb-128.0))
		}
	} else if (PhotoInt == "RGB") && (segment_count == 3) {
		for i = 0; i < size/segment_count; i++ {
			out[3*i] = temp[i]
			out[3*i+1] = temp[i+offset]
			out[3*i+2] = temp[i+2*offset]
		}
	} else {
		return fmt.Errorf("ERROR, format not supported")
	}
	return nil
}
