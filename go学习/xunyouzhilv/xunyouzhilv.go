package main

import "fmt"

// 定义bfs函数用来求最短时间
func bfs(n int, k int) int {
    // 定义t和ok判断映射中是否已到达某个位置 _赋值用来防止变量定义却未使用
    var t int = 0
    _ = t
    var ok = true
    time := make(map[int]int)  // time用于存储到达每个位置的步数
    queue := []int{n}          // 队列用于存放走过的路，初始为小青的位置n
    time[n] = 0                // 将小青初始位置的步数设为 0 步
    // 循环直到队列为空
    for len(queue) > 0 {
        // 取出队首元素
        curr := queue[0]
        queue = queue[1:]
        // 如果当前位置是小码的位置则返回步数
        if curr == k {
            return time[curr]
        }
        // 如果当前位置-1未越界且未被寻过则加入队列
        if curr-1 >= 0 {
            if t,ok = time[curr-1];!ok {
                queue = append(queue, curr-1)
                time[curr-1] = time[curr] + 1
            }
        }
        // 如果当前位置+1未越界且未被寻过则加入队列
        if curr+1 <= 100000 {
            if t,ok = time[curr+1];!ok {
                queue = append(queue, curr+1)
                time[curr+1] = time[curr] + 1
            }
        }
        // 如果当前位置*2未越界且未被寻过则加入队列
        if curr*2 <= 100000 {
            if t,ok = time[curr*2];!ok {
                queue = append(queue, curr*2)
                time[curr*2] = time[curr] + 1
            }
        }
    }
    // 如果最终没有找到则返回 -1
    return -1
}


func main() {
    var n, k int
    fmt.Println("请输入小青和小码的位置:")
    fmt.Scan(&n, &k)
    fmt.Println("最短时间是：",bfs(n, k),"分钟")
}
