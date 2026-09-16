# skiplist

Go 实现的跳表（p = 1/4，maxLevel = 32）：`Set` / `Get` / `Delete` / 范围扫描。期望 O(log n) 查找，第 0 层是有序链表。

非并发安全；多 goroutine 请外加 `RWMutex`，或只插不删时再考虑无锁。

配套文章：[Go 跳表的原理与实现](https://seamonw.github.io/blog/2026/09/15/golang-skiplist/)

```bash
go get github.com/seamonw/skiplist
```

需要 Go 1.23（`iter.Seq2`）。

## 快速开始

```go
l := skiplist.New[int, string]()
l.Set(2, "b")
l.Set(1, "a")
v, ok := l.Get(1)

l.Range(1, func(k int, v string) bool {
    fmt.Println(k, v)
    return true // false 停止
})

for k, v := range l.From(1) {
    _, _ = k, v
}
```

`Set` 是 upsert：key 已存在只改 value，不改节点高度。
