//go:build windows

package style

import "testing"

func TestLooksLikeMSYSPtyName(t *testing.T) {
	cases := []struct {
		name string
		want bool
	}{
		{`\msys-1888ae32e00d56aa-pty0-to-master`, true},
		{`\msys-1888ae32e00d56aa-pty0-from-master`, true},
		{`\cygwin-1234567890abcdef-pty3-to-master`, true},
		{`\MSYS-DD50A72AB4668B33-PTY1-FROM-MASTER`, true},
		{`\Device\NamedPipe\anonymous`, false},
		{`C:\Users\foo\output.txt`, false},
		{``, false},
	}

	for _, c := range cases {
		if got := looksLikeMSYSPtyName(c.name); got != c.want {
			t.Errorf("looksLikeMSYSPtyName(%q) = %v, want %v", c.name, got, c.want)
		}
	}
}
