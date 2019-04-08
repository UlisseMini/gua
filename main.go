package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	lua "github.com/yuin/gopher-lua"
)

func LErr(err error) lua.LValue {
	if err == nil {
		return lua.LNil
	}

	return lua.LString(err.Error())
}

func luaErrWrap(L *lua.LState, err error) {
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LNil)
	}
}

type LConn struct{ net.Conn }

// The write function exported to lua. returns nil or error string.
func (c *LConn) LuaWrite(L *lua.LState) int {
	s := L.ToString(1)
	_, err := io.WriteString(c, s)
	luaErrWrap(L, err)
	return 1
}

// The close function exported to lua.
func (c *LConn) LuaClose(L *lua.LState) int {
	err := c.Close()
	luaErrWrap(L, err)
	return 1
}

// The Read function exported to lua.
func (c *LConn) LuaRead(L *lua.LState) int {
	bufSize := L.ToInt(1)
	buf := make([]byte, bufSize)

	n, err := c.Read(buf)
	L.Push(lua.LString(buf[:n]))
	L.Push(LErr(err))
	return 2
}

func Dial(L *lua.LState) int {
	addr := L.ToString(1)
	rawConn, err := net.Dial("tcp", addr)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(LErr(err))
		return 2
	}
	t := L.NewTable()
	conn := LConn{rawConn}
	t.RawSetString("write", L.NewFunction(conn.LuaWrite))
	t.RawSetString("close", L.NewFunction(conn.LuaClose))
	t.RawSetString("read", L.NewFunction(conn.LuaRead))

	L.Push(t)
	L.Push(lua.LNil)
	return 2
}

func NewState() *lua.LState {
	L := lua.NewState()
	L.SetGlobal("dial", L.NewFunction(Dial))
	return L
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <file.lua>\n", os.Args[0])
		return
	}

	L := NewState()
	defer L.Close()

	if err := L.DoFile(os.Args[1]); err != nil {
		log.Print(err)
		return
	}
}
