package tools

import "strconv"

// ParseUint is like ParseInt but for unsigned numbers. It's stolen from the
// strconv package and streamlined for uint64s.
func ParseUint(s []byte) (n uint64, err error) {
	const maxUint64 = (1<<64 - 1)
	var cutoff, maxVal uint64

	cutoff = maxUint64/10 + 1
	maxVal = 1<<uint(64) - 1

	for i := 0; i < len(s); i++ {
		var v byte
		d := s[i]
		switch {
		case '0' <= d && d <= '9':
			v = d - '0'
		case 'a' <= d && d <= 'z':
			v = d - 'a' + 10
		case 'A' <= d && d <= 'Z':
			v = d - 'A' + 10
		default:
			n = 0
			err = strconv.ErrSyntax
			goto Error
		}
		if v >= 10 {
			n = 0
			err = strconv.ErrSyntax
			goto Error
		}

		if n >= cutoff {
			// n*base overflows
			n = maxUint64
			err = strconv.ErrRange
			goto Error
		}
		n *= 10

		n1 := n + uint64(v)
		if n1 < n || n1 > maxVal {
			// n+v overflows
			n = maxUint64
			err = strconv.ErrRange
			goto Error
		}
		n = n1
	}

	return n, nil

Error:
	return n, &strconv.NumError{Func: "ParseUint", Num: string(s), Err: err}
}
