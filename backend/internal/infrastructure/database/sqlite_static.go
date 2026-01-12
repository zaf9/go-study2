//go:build !cgo
// +build !cgo

package database

// 此文件在不使用 CGO 时编译
// GoFrame 的 contrib/drivers/sqlite/v2 会自动选择合适的驱动
// 在 CGO 不可用时，它会使用纯 Go 实现的 SQLite（如 modernc.org/sqlite）
