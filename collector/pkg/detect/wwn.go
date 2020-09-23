package detect

import (
	"fmt"
	"strings"
)

type Wwn struct {
	Naa        uint64 `json:"naa"`
	Oui        uint64 `json:"oui"`
	Id         uint64 `json:"id"`
	VendorCode string `json:"vendor_code"`
}

// this is an incredibly basic converter, that only works for "Registered" IEEE format - NAA5
// https://standards.ieee.org/content/dam/ieee-standards/standards/web/documents/tutorials/fibre.pdf
// references:
// - https://metacpan.org/pod/Device::WWN
// - https://en.wikipedia.org/wiki/World_Wide_Name
// - https://storagemeat.blogspot.com/2012/08/decoding-wwids-or-how-to-tell-whats-what.html
// - https://bryanchain.com/2016/01/20/breaking-down-an-naa-id-world-wide-name/

/*
+----------+---+---+---+---+---+---+---+---+
| Byte/Bit | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
+----------+---+---+---+---+---+---+---+---+
| 0        | NAA (5h)      | (MSB)         |
+----------+---------------+               +
| 1        |                               |
+----------+            IEEE OUI           |
| 2        |                               |
+----------+               +---------------+
| 3        | (LSB)         | (MSB)         |
+----------+---------------+               +
| 4        |                               |
|          |                               |
+----------+                               |
| 5        |            Vendor ID          |
+----------+                               |
| 6        |                               |
+----------+                               |
| 7        |                         (LSB) |
+----------+-------------------------------+


*/

func (wwn *Wwn) ToString() string {

	var wwnBuffer uint64

	wwnBuffer = wwn.Id           //start with vendor ID
	wwnBuffer += (wwn.Oui << 36) //add left-shifted OUI
	wwnBuffer += (wwn.Naa << 60) //NAA is a number from 1-6, so decimal == hex.

	//TODO: may need to support additional versions in the future.

	return strings.ToLower(fmt.Sprintf("%#x", wwnBuffer))
}
