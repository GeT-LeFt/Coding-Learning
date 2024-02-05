package mr

import "fmt"
import "log"
import "net"
import "os"
import "net/rpc"
import "net/http"
import "sync"
import "time"


type Coordinator struct {
	// Your definitions here.
	Workers map[int]bool
	Files []string
}

// Your code here -- RPC handlers for the worker to call.

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)      			// 注册当前对象（协调器）为 RPC 服务对象，使其方法可被远程调用。
	rpc.HandleHTTP()   				// 将RPC服务绑定到HTTP协议上。
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()	// realized in rpc.go 用于生成套接字文件路径
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.


	return ret
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{
		Workers: make(map[int]bool), // 初始化 Workers map
		Files:	 files,
	}

	// Your code here.
	var mu sync.Mutex					// 辅助函数可注释掉，每隔1s查看worker状态
	go printMapContent(c.Workers, &mu)	// 辅助函数可注释掉，每隔1s查看worker状态

	c.server()
	return &c
}

func printMapContent(m map[int]bool, mu *sync.Mutex) {    // 辅助函数可注释掉，每隔1s查看worker状态
    for {
        // 锁定 map，以确保在打印时没有其他线程对其进行修改
        mu.Lock()
        fmt.Println("Workers Status:")
        for key, value := range m {
            fmt.Printf("%d: %v | ", key, value)
        }
        mu.Unlock()
		fmt.Printf("\n")
        time.Sleep(1 * time.Second) // 等待 1 秒
    }
}