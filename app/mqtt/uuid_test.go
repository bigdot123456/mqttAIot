package main

import (
	"fmt"
	"testing"
)

func TestUUID(t *testing.T) {
	id:=getUUID(123)
	fmt.Println("id is:", id)
	// NewNode创建并返回一个新的snowflake节点，正如前面文章谈到，它是有范围的。在0~1023。
	// func NewNode(node int64) (*Node, error) {
	//      mu.Lock()
	//      ...
	//      if n.node < 0 || n.node > n.nodeMax {
	//          return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
	//      }
	// }

	// Generate会生成并返回一个唯一的snowflake ID，为了保证其唯一性，你需要保证你的系统的时间要精确
	// 生成的ID上有多个类型转化的接口，例如：String()，Int64()等
}


