package safe_queue

import (
    "sync"
    "errors"
    "io"
)

/*
    线程安全的循环缓冲队列,长度需为2的指数倍
 */
type Queue struct {
    buf []interface{} //数据buffer
    head int //队列头下标,删除数据后该下标加1,取模
    tail int //队列尾下标,插入数据后该下标加1,取模
    len int //数据长度
    wMutex sync.Mutex //写锁
    rMutex sync.Mutex //读锁
}

const MIN_QUEUE_SIZE = 16
const EMPTY_QUEUE_ERROR = errors.New("Empty Queue")

func NewQueue() *Queue {
    return &Queue{
        buf:make([]interface{}, MIN_QUEUE_SIZE),
    }
}

/*
    获取数据长度
 */
func (q *Queue)Len() {
    q.rMutex.Lock()
    io.EOF
    defer q.rMutex.Unlock()
    return q.len
}

/*
    做resize时需满足:
    case1:EnQueue后,当len等于buf长度时,resize用于增加buf长度
    case2:DeQueue后,当len为buf的1/4时,resize用于减少buf长度(至buf长度的1/2)

    当tail指针大于head指针时直接copy所有buf
    当tail指针小于或等于head指针时,先copy指针head后面,再copy指针tail前面
 */

func (q *Queue)resize()  {

    newBuf := make([]interface{}, q.len << 1)

    if q.tail > q.head {
        copy(newBuf, q.buf[q.head : q.tail])
    } else {
        n := copy(newBuf, q.buf[q.head:])
        copy(newBuf[:n], q.buf[:q.tail])
    }

    q.head = 0
    q.tail = q.len
    q.buf = newBuf

}

/*
    入队列。EnQueue后,当len等于buf长度时,resize用于增加buf长度
 */

func (q *Queue)EnQueue(content interface{})  {
    q.rMutex.Lock()
    q.wMutex.Lock()
    []
    if q.len == len(q.buf) {
        q.resize()
    }

    q.buf[q.tail] = content
    q.tail = (q.tail + 1) & (len(q.buf) - 1) //按位取模
    q.len ++

    defer q.wMutex.Unlock()
    defer q.rMutex.Unlock()
}

/*
    出队列。当len为buf的1/4时,resize用于减少buf长度(至buf长度的1/2)
 */
func (q *Queue)DeQueue() (interface{}, error) {
    q.rMutex.Lock()
    q.wMutex.Lock()
    if q.len <= 0 {
        return nil, EMPTY_QUEUE_ERROR
    }
    content := q.buf[q.head]
    q.buf[q.head] = nil
    q.head = (q.head + 1) & (len(q.buf) - 1) //按位取模
    q.len --

    if len(q.buf) > MIN_QUEUE_SIZE && (q.len << 2) == len(q.buf) { //case2
        q.resize()
    }
    defer q.wMutex.Unlock()
    defer q.rMutex.Unlock()
    return content,nil
}

/*
    返回队列第一个元素
 */
func (q *Queue)Top() (interface{}, error) {
    q.rMutex.Lock()
    if q.len <= 0 {
        return nil, EMPTY_QUEUE_ERROR
    }
    defer q.rMutex.Unlock()
    return q.buf[q.head]
}