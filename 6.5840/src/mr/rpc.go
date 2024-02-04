package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import "os"
import "strconv"

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

// Add your RPC definitions here.
// MapTask related
type MapTaskArgs struct {    // 可以包含特定于 Map 任务请求的字段
}

type MapTaskReply struct {
	Filename string			// 分配的文件名
	NReduce  int			// Reduce任务数量
	TaskID   int			// Map任务的ID
}

// ReduceTask related
type ReduceTaskArgs struct {// 可以包含特定于 Reduce 任务请求的字段
}

type ReduceTaskReply struct {
	KeyValues []KeyValue	// Map输出的K/V
	TaskID	  int			// reduce任务ID
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {    // 用于生成套接字文件路径
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
