/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package dm

import (
	"io"
)

const (
	READ_LEN = Dm_build_1089
)

type iOffRowBinder interface {
	read(buf *Dm_build_283)
	isReadOver() bool
	getObj() any
}

type offRowBinder struct {
	obj          any
	encoding     string
	readOver     bool
	buffer       *Dm_build_283
	position     int32
	offRow       bool
	targetLength int64
}

func newOffRowBinder(obj any, encoding string, targetLength int64) *offRowBinder {
	return &offRowBinder{
		obj:          obj,
		encoding:     encoding,
		targetLength: targetLength,
		readOver:     false,
		buffer:       Dm_build_287(),
		position:     0,
	}
}

type offRowBytesBinder struct {
	*offRowBinder
}

func newOffRowBytesBinder(obj []byte, encoding string) *offRowBytesBinder {
	var binder = &offRowBytesBinder{
		newOffRowBinder(obj, encoding, int64(IGNORE_TARGET_LENGTH)),
	}
	binder.read(binder.buffer)
	binder.offRow = binder.buffer.Dm_build_288() > Dm_build_1086
	return binder
}

func (b *offRowBytesBinder) read(buf *Dm_build_283) {
	if b.buffer.Dm_build_288() > 0 {
		buf.Dm_build_320(b.buffer)
	} else if !b.readOver {
		var obj = b.obj.([]byte)
		buf.Dm_build_309(obj, 0, len(obj))
		b.readOver = true
	}
}

func (b *offRowBytesBinder) isReadOver() bool {
	return b.readOver
}

func (b *offRowBytesBinder) getObj() any {
	return b.obj
}

type offRowBlobBinder struct {
	*offRowBinder
}

func newOffRowBlobBinder(blob DmBlob, encoding string) *offRowBlobBinder {
	var binder = &offRowBlobBinder{
		newOffRowBinder(blob, encoding, int64(IGNORE_TARGET_LENGTH)),
	}
	binder.read(binder.buffer)
	binder.offRow = binder.buffer.Dm_build_288() > Dm_build_1086
	return binder
}

func (b *offRowBlobBinder) read(buf *Dm_build_283) {
	if b.buffer.Dm_build_288() > 0 {
		buf.Dm_build_320(b.buffer)
	} else if !b.readOver {
		var obj = b.obj.(DmBlob)
		var totalLen, _ = obj.GetLength()
		var leaveLen = totalLen - int64(b.position)
		var readLen = int32(leaveLen)
		if leaveLen > READ_LEN {
			readLen = READ_LEN
		}
		var bytes, _ = obj.getBytes(int64(b.position)+1, readLen)
		b.position += readLen
		if b.position == int32(totalLen) {
			b.readOver = true
		}
		buf.Dm_build_309(bytes, 0, len(bytes))
	}
}

func (b *offRowBlobBinder) isReadOver() bool {
	return b.readOver
}

func (b *offRowBlobBinder) getObj() any {
	return b.obj
}

type offRowClobBinder struct {
	*offRowBinder
}

func newOffRowClobBinder(clob DmClob, encoding string) *offRowClobBinder {
	var binder = &offRowClobBinder{
		newOffRowBinder(clob, encoding, int64(IGNORE_TARGET_LENGTH)),
	}
	binder.read(binder.buffer)
	binder.offRow = binder.buffer.Dm_build_288() > Dm_build_1086
	return binder
}

func (b *offRowClobBinder) read(buf *Dm_build_283) {
	if b.buffer.Dm_build_288() > 0 {
		buf.Dm_build_320(b.buffer)
	} else if !b.readOver {
		var obj = b.obj.(DmClob)
		var totalLen, _ = obj.GetLength()
		var leaveLen = totalLen - int64(b.position)
		var readLen = int32(leaveLen)
		if leaveLen > READ_LEN {
			readLen = READ_LEN
		}
		var str, _ = obj.getSubString(int64(b.position)+1, readLen)
		var bytes = Dm_build_1.Dm_build_217(str, b.encoding, nil)
		b.position += readLen
		if b.position == int32(totalLen) {
			b.readOver = true
		}
		buf.Dm_build_309(bytes, 0, len(bytes))
	}
}

func (b *offRowClobBinder) isReadOver() bool {
	return b.readOver
}

func (b *offRowClobBinder) getObj() any {
	return b.obj
}

type offRowReaderBinder struct {
	*offRowBinder
}

func newOffRowReaderBinder(reader io.Reader, encoding string) *offRowReaderBinder {
	var binder = &offRowReaderBinder{
		newOffRowBinder(reader, encoding, int64(IGNORE_TARGET_LENGTH)),
	}
	binder.read(binder.buffer)
	binder.offRow = binder.buffer.Dm_build_288() > Dm_build_1086
	return binder
}

func (b *offRowReaderBinder) read(buf *Dm_build_283) {
	if b.buffer.Dm_build_288() > 0 {
		buf.Dm_build_320(b.buffer)
	} else if !b.readOver {
		var err error
		var readLen = READ_LEN
		var reader = b.obj.(io.Reader)
		var bytes = make([]byte, readLen)
		readLen, err = reader.Read(bytes)
		if err == io.EOF {
			b.readOver = true
			return
		}
		b.position += int32(readLen)
		if readLen < len(bytes) || b.targetLength != int64(IGNORE_TARGET_LENGTH) && int64(b.position) == b.targetLength {
			b.readOver = true
		}
		buf.Dm_build_309(bytes[0:readLen], 0, readLen)
	}
}

func (b *offRowReaderBinder) readAll() []byte {
	var byteArray = Dm_build_287()
	b.read(byteArray)
	for !b.readOver {
		b.read(byteArray)
	}
	return byteArray.Dm_build_330()
}

func (b *offRowReaderBinder) isReadOver() bool {
	return b.readOver
}

func (b *offRowReaderBinder) getObj() any {
	return b.obj
}
