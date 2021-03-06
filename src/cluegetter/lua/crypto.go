// ClueGetter - Does things with mail
//
// Copyright 2016 Dolf Schimmel, Freeaqingme.
//
// This Source Code Form is subject to the terms of the Apache License, Version 2.0.
// For its contents, please refer to the LICENSE file.
//
package lua

import (
	"github.com/yuin/gopher-lua"

	"crypto/sha256"
	"encoding/hex"
)

func CryptoLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), cryptoExports)
	L.Push(mod)
	return 1
}

var cryptoExports = map[string]lua.LGFunction{
	"sha256": cryptoSha256,
}

func cryptoSha256(L *lua.LState) int {
	data := L.ToString(2)

	h := sha256.New()
	h.Write([]byte(data))

	L.Push(lua.LString(hex.EncodeToString(h.Sum(nil))))
	return 1
}
