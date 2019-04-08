package main

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestLuaDial(t *testing.T) {
	const addr = "127.0.0.1:2903"
	var gotConn bool
	go func() {
		l, err := net.Listen("tcp4", addr)
		if err != nil {
			t.Fatal(err)
		}

		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(conn, "OK")
		conn.Close()
		gotConn = true
	}()

	L := NewState()
	if err := L.DoString(fmt.Sprintf(`local conn, err = dial('%s')
							if err ~= nil then error(err) end
							local dat, err = conn.read(128)
							if err ~= nil then error(err) end
							if dat ~= "OK" then
								error(("want %%q; got %%q"):format("OK", dat))
							end
							conn.close()
	`, addr)); err != nil {
		t.Fatalf("lua error: %v", err)
	}

	if !gotConn {
		t.Fatal("Did not get the connection")
	}
}
