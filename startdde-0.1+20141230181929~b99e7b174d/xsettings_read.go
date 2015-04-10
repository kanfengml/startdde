/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
	"bytes"
	"encoding/binary"
	"github.com/BurntSushi/xgb/xproto"
	"io"
)

func readInterger(buf io.Reader) uint32 {
	body := uint32(0)
	binary.Read(buf, byteOrder, &body)

	return body
}

func readColor(buf io.Reader) []uint16 {
	ret := []uint16{}
	var r uint16

	binary.Read(buf, byteOrder, &r)
	ret = append(ret, r)
	binary.Read(buf, byteOrder, &r)
	ret = append(ret, r)
	binary.Read(buf, byteOrder, &r)
	ret = append(ret, r)
	binary.Read(buf, byteOrder, &r)
	ret = append(ret, r)

	return ret
}

func readString(buf io.Reader) string {
	var nameLen uint32
	binary.Read(buf, byteOrder, &nameLen)
	if nameLen > 1000 {
		logger.Info("name len to long:", nameLen)
		panic("name len to long")
	}

	nameBuf := make([]byte, nameLen)
	binary.Read(buf, byteOrder, &nameBuf)

	leftPad := 3 - (nameLen+3)%4
	buf.Read(make([]byte, leftPad))

	return string(nameBuf)
}

func readString2(buf io.Reader) (string, uint16) {
	var nameLen uint16
	binary.Read(buf, byteOrder, &nameLen)

	nameBuf := make([]byte, nameLen)
	binary.Read(buf, byteOrder, &nameBuf)

	leftPad := 3 - (nameLen+3)%4
	buf.Read(make([]byte, leftPad))

	return string(nameBuf), nameLen
}

func readHeader(buf io.Reader) (byte, uint16, string, uint32) {
	var sType byte
	binary.Read(buf, byteOrder, &sType)
	buf.Read(make([]byte, 1))

	name, nameLen := readString2(buf)
	lastSerial := readInterger(buf)

	return sType, nameLen, name, lastSerial
}

func readXSettings() []*HeaderInfo {
	reply, err := xproto.GetProperty(X, false,
		sReply.Owner,
		getAtom(X, XSETTINGS_SETTINGS),
		getAtom(X, XSETTINGS_SETTINGS),
		0, 10240).Reply()
	if err != nil {
		logger.Info("Get Property Failed:", err)
		panic(err)
	}

	if reply.ValueLen <= 0 {
		return nil
	}

	infos := []*HeaderInfo{}
	data := reply.Value[:reply.ValueLen]
	if data[0] == 1 {
		byteOrder = binary.BigEndian
	} else {
		byteOrder = binary.LittleEndian
	}

	buf := bytes.NewReader(data[4:])

	serial := readInterger(buf)
	numSettings := readInterger(buf)

	logger.Infof("serial: %d, numSettings: %d",
		serial, numSettings)

	for i := uint32(0); i < numSettings; i++ {
		sType, nameLen, name, lastSerial := readHeader(buf)
		info := &HeaderInfo{}
		info.vType = sType
		info.nameLen = nameLen
		info.name = name
		info.lastSerial = lastSerial
		switch sType {
		case XSETTINGS_INTERGER:
			v := readInterger(buf)
			//logger.Infof("%s = %d", name, v)
			//logger.Debugf("\"%s\": \"%d;0\",\n",
			//name, v)
			info.value = v
		case XSETTINGS_STRING:
			v := readString(buf)
			//logger.Debugf("\"%s\": \"%s;1\",\n",
			//name, v)
			info.value = v
		case XSETTINGS_COLOR:
			v := readColor(buf)
			//logger.Debugf("\"%s\": \"%d,%d,%d,%d;2\",\n",
			//name, v[0], v[1], v[2], v[3])
			info.value = v
		}
		infos = append(infos, info)
	}
	print("\n")

	return infos
}
