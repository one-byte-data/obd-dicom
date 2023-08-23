package network

import "github.com/one-byte-data/obd-dicom/media"

// AsyncOperationWindow - AsyncOperationWindow
type AsyncOperationWindow interface {
	GetMaxNumberOperationsInvoked() uint16
	GetMaxNumberOperationsPerformed() uint16
	Size() uint16
	Read(ms media.MemoryStream) (err error)
	ReadDynamic(ms media.MemoryStream) (err error)
}

type asyncOperationWindow struct {
	ItemType                     byte //0x53
	Reserved1                    byte
	Length                       uint16
	MaxNumberOperationsInvoked   uint16
	MaxNumberOperationsPerformed uint16
}

// NewAsyncOperationWindow - NewAsyncOperationWindow
func NewAsyncOperationWindow() AsyncOperationWindow {
	return &asyncOperationWindow{
		ItemType: 0x53,
	}
}

func (async *asyncOperationWindow) GetMaxNumberOperationsInvoked() uint16 {
	return async.MaxNumberOperationsInvoked
}

func (async *asyncOperationWindow) GetMaxNumberOperationsPerformed() uint16 {
	return async.MaxNumberOperationsPerformed
}

func (async *asyncOperationWindow) Size() uint16 {
	return async.Length + 4
}

func (async *asyncOperationWindow) Read(ms media.MemoryStream) (err error) {
	if async.ItemType, err = ms.GetByte(); err != nil {
		return err
	}
	return async.ReadDynamic(ms)
}

func (async *asyncOperationWindow) ReadDynamic(ms media.MemoryStream) (err error) {
	if async.Reserved1, err = ms.GetByte(); err != nil {
		return err
	}
	if async.Length, err = ms.GetUint16(); err != nil {
		return err
	}
	if async.MaxNumberOperationsInvoked, err = ms.GetUint16(); err != nil {
		return err
	}
	if async.MaxNumberOperationsPerformed, err = ms.GetUint16(); err != nil {
		return err
	}
	return
}
