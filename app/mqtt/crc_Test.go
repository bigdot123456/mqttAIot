package main

import (
"fmt"
"github.com/snksoft/crc"
	"testing"
)

func TestCRC(t *testing.T) {
	data := "123456789"
	hash := crc.NewHash(crc.XMODEM)
	xmodemCrc := hash.CalculateCRC([]byte(data))
	fmt.Printf("CRC is 0x%04X\n", xmodemCrc) // prints "CRC is 0x31C3"

	// You can also reuse hash instance for another crc calculation
	// And if data is too big, you may feed it in chunks
	hash.Reset() // Discard crc data accumulated so far
	hash.Update([]byte("123456789")) // feed first chunk
	hash.Update([]byte("01234567890")) // feed next chunk
	xmodemCrc2 := hash.CRC() // gets CRC of whole data "12345678901234567890"
	fmt.Printf("CRC is 0x%04X\n", xmodemCrc2) // prints "CRC is 0x2C89"
}
