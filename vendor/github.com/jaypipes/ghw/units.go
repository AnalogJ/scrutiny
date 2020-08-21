//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

var (
	KB int64 = 1024
	MB       = KB * 1024
	GB       = MB * 1024
	TB       = GB * 1024
	PB       = TB * 1024
	EB       = PB * 1024
)

func unitWithString(size int64) (int64, string) {
	switch {
	case size < MB:
		return KB, "KB"
	case size < GB:
		return MB, "MB"
	case size < TB:
		return GB, "GB"
	case size < PB:
		return TB, "TB"
	case size < EB:
		return PB, "PB"
	default:
		return EB, "EB"
	}
}
