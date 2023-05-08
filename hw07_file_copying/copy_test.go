package main

import (
	"errors"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	cases := []struct {
		name   string
		from   string
		to     string
		offset int64
		limit  int64
		resLen int64
		fail   bool
		err    error
	}{
		{
			name:   "from not exists",
			from:   "/dev/test",
			to:     "/tmp/test.txt",
			offset: 0,
			limit:  0,
			fail:   true,
			err:    os.ErrNotExist,
		}, {
			name:   "unsupported file 1",
			from:   "/dev",
			to:     "/tmp/test.txt",
			offset: 0,
			limit:  10,
			fail:   true,
			err:    ErrUnsupportedFile,
		}, {
			name:   "unsupported file 2",
			from:   "/tmp/empty",
			to:     "/tmp/test.txt",
			offset: 0,
			limit:  10,
			fail:   true,
			err:    ErrUnsupportedFile,
		}, {
			name:   "offset err",
			from:   "testdata/input2.txt",
			to:     "/tmp/test.txt",
			offset: 500,
			limit:  10,
			fail:   true,
			err:    ErrOffsetExceedsFileSize,
		}, {
			name:   "success limit 1",
			from:   "testdata/out_offset0_limit0.txt",
			to:     "/tmp/test.txt",
			offset: 100,
			limit:  1,
			fail:   false,
		}, {
			name:   "success limit 2",
			from:   "testdata/out_offset0_limit0.txt",
			to:     "/tmp/test.txt",
			offset: 101,
			limit:  2,
			fail:   false,
		}, {
			name:   "success limit exceed source",
			from:   "testdata/input2.txt",
			to:     "/tmp/test.txt",
			offset: 0,
			limit:  2000,
			resLen: 4,
			fail:   false,
		}, {
			name:   "success limit 0",
			from:   "testdata/input2.txt",
			to:     "/tmp/test.txt",
			offset: 0,
			limit:  0,
			resLen: 4,
			fail:   false,
		},
	}

	_, _ = os.Create("/tmp/empty")
	buf := make([]byte, int(math.Min(float64(buffSize), float64(1000))))

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := Copy(c.from, c.to, c.offset, c.limit)
			if c.fail != false {
				require.NotNil(t, err)
				require.True(t, errors.Is(err, c.err))
			} else {
				require.Nil(t, err)

				res, _ := os.Open(c.to)
				n, _ := res.Read(buf)

				resLen := c.resLen
				if resLen == 0 {
					resLen = c.limit
				}
				require.True(t, int64(n) == resLen)
				_ = res.Close()
			}
		})
	}
}
