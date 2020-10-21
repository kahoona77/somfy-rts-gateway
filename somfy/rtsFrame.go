package somfy

import (
	"encoding/binary"
)

func GetFrame(device *Device, btn Button) []byte {
	data := make([]byte, 7)
	data[0] = device.EncryptionKey
	data[1] = byte(btn) << 4

	rollingCode := make([]byte, 2)
	binary.LittleEndian.PutUint16(rollingCode, device.RollingCode)
	data[2] = rollingCode[1]
	data[3] = rollingCode[0]

	address := make([]byte, 4)
	binary.LittleEndian.PutUint32(address, device.Address)
	data[4] = address[0]
	data[5] = address[1]
	data[6] = address[2]

	//create checksum
	var checksum byte = 0
	data[1] = data[1] & 0xF0
	for i := 0; i < 7; i++ {
		checksum = checksum ^ data[i] ^ (data[i] >> 4)
	}
	checksum = checksum & 0xf
	data[1] = data[1] | checksum

	//obfuscate frame
	for i := 1; i < 7; i++ {
		data[i] = data[i] ^ data[i-1]
	}

	return data
}
