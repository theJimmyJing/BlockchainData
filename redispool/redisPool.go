// package redispool

// import (
// 	"context"
// 	"net"
// 	"sync"
// 	"time"
// )

// type Options struct {
// 	Dialer  func(context.Context) (net.Conn, error) // 新建连接的工厂函数
// 	OnClose func(*Conn) error                       // 关闭连接的回调函数，在连接关闭的时候的回调

// 	PoolSize           int           // 连接池大小，连接池中的连接的最大数量
// 	MinIdleConns       int           // 最小空闲连接数
// 	MaxConnAge         time.Duration // 从连接池获取连接的超时时间
// 	PoolTimeout        time.Duration // 空闲连接的超时时间
// 	IdleTimeout        time.Duration // 空闲连接的超时时间
// 	IdleCheckFrequency time.Duration // 检查空闲连接频率（超时空闲连接清理的间隔时间）
// }

// type Conn struct {
// 	// 包装net.conn
// 	netConn net.Conn // tcp 连接

// 	rd *proto.Reader // 根据 Redis 通信协议实现的 Reader
// 	wr *proto.Writer // 根据 Redis 通信协议实现的 Writer

// 	Inited    bool      // 是否完成初始化（该连接是否初始化，比如如果需要执行命令之前需要执行的auth,select db 等的标识，代表已经auth,select过）
// 	pooled    bool      // 是否放进连接池的标志，有些场景产生的连接是不需要放入连接池的
// 	createdAt time.Time // 创建时间（超过maxconnage的连接需要淘汰）
// 	usedAt    int64     // 使用时间（该连接执行命令的时候所记录的时间，就是上次用过这个连接的时间点）
// }

// type ConnPool struct {
// 	opt *Options // 连接池参数配置，如上

// 	dialErrorsNum uint32 // 连接失败的错误次数，atomic

// 	lastDialErrorMu sync.RWMutex // 上一次连接错误锁，读写锁
// 	lastDialError   error        // 上一次连接错误（保存了最近一次的连接错误）

// 	queue chan struct{} // 轮转队列，是一个 channel 结构（一个带容量（poolsize）的阻塞队列）

// 	connsMu   sync.Mutex // 连接队列锁
// 	conns     []*Conn    // 连接队列（连接队列，维护了未被删除所有连接，即连接池中所有的连接）
// 	idleConns []*Conn    // 空闲连接队列（空闲连接队列，维护了所有的空闲连接）

// 	poolSize     int // 连接池大小，注意如果没有可用的话要等待
// 	idleConnsLen int // 连接池空闲连接队列长度

// 	stats Stats // 连接池统计的结构体（包含了使用数据）

// 	_closed  uint32        // 连接池关闭标志，atomic
// 	closedCh chan struct{} // 通知连接池关闭通道（用于主协程通知子协程的常用方法）
// }

// // Stats contains pool state information and accumulated stats.
// type Stats struct {
// 	Hits     uint32 // number of times free connection was found in the pool
// 	Misses   uint32 // number of times free connection was NOT found in the pool
// 	Timeouts uint32 // number of times a wait timeout occurred

// 	TotalConns uint32 // number of total connections in the pool
// 	IdleConns  uint32 // number of idle connections in the pool
// 	StaleConns uint32 // number of stale connections removed from the pool
// }

// // 连接池的接口
// type Pooler interface {
// 	NewConn(context.Context) (*Conn, error) // 创建连接
// 	CloseConn(*Conn) error                  // 关闭连接

// 	Get(context.Context) (*Conn, error) // 获取连接
// 	Put(*Conn)                          // 放回连接
// 	Remove(*Conn, error)                // 移除连接

// 	Len() int      // 连接池长度
// 	IdleLen() int  // 空闲连接数量
// 	Stats() *Stats // 连接池统计

// 	Close() error // 关闭连接池
// }

// func init() {

// }
