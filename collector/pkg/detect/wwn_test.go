package detect_test

import (
	"fmt"
	"testing"

	"github.com/analogj/scrutiny/collector/pkg/detect"
	"github.com/stretchr/testify/require"
)

func TestWwn_FromStringTable(t *testing.T) {
	// setup
	tests := []struct {
		wwnStr string
		wwn    detect.Wwn
	}{
		{"0x5002538e40a22954", detect.Wwn{Naa: 5, Oui: 9528, Id: 61213911380}}, // sda
		{"0x5000cca264eb01d7", detect.Wwn{Naa: 5, Oui: 3274, Id: 10283057623}}, // sdb
		{"0x5000cca264ec3183", detect.Wwn{Naa: 5, Oui: 3274, Id: 10283135363}}, // sdc
		{"0x5000cca252c859cc", detect.Wwn{Naa: 5, Oui: 3274, Id: 9978796492}},  // sdd
		{"0x50014ee20b2a72a9", detect.Wwn{Naa: 5, Oui: 5358, Id: 8777265833}},  // sde
		{"0x5000cca264ebc248", detect.Wwn{Naa: 5, Oui: 3274, Id: 10283106888}}, // sdf
		{"0x5000c500673e6b5f", detect.Wwn{Naa: 5, Oui: 3152, Id: 1732143967}},  // sdg
	}
	// test
	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.wwnStr)
		t.Run(testname, func(t *testing.T) {
			str := tt.wwn.ToString()
			require.Equal(t, tt.wwnStr, str)
		})
	}
}
