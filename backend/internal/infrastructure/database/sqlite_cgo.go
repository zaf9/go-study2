//go:build cgo
// +build cgo

package database

import (
	_ "github.com/mattn/go-sqlite3"
)

// 此文件在使用 CGO 时编译
// mattn/go-sqlite3 是基于 CGO 的 SQLite 驱动
