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
	Workers map[int]bool	// 辅助函数可注释掉，每隔1s查看worker状态
	mu		    sync.Mutex
	nMap		int
	nReduce		int
	mapTasks	[]Task
	reduceTasks []Task
}

type Task struct {
	Type 		string
	Status		string
	Index		int
	FileName	string
	WorkerId	int
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

func (c *Coordinator) RequestTaskReply(args *RequestTaskArgs, reply *RequestTaskReply) error {
	c.mu.Lock()
	c.Workers[args.WorkerId] = true	// 辅助函数可注释掉，每隔1s查看worker状态
	tmpTask := &Task{}
	if c.nMap > 0 {				// 给map任务
		tmpTask = c.selectTask("MapTask", c.mapTasks, args.WorkerId)
	} else if c.nReduce > 0 {	// 给reduce任务
		tmpTask = c.selectTask("ReduceTask", c.reduceTasks, args.WorkerId)
	} else {

	}
	reply.TaskType = tmpTask.Type
	reply.TaskId = args.WorkerId
	reply.TaskFile = tmpTask.FileName
	c.mu.Unlock()
	return nil
}

func (c *Coordinator) selectTask(taskType string, tasks []Task, workerId int) *Task {
	tmpTask := Task{}
	for i := 0; i < len(tasks); i++ {
		if tasks[i].Status == "NotStarted"  && tasks[i].Type == taskType {
			tmpTask = tasks[i]
			tmpTask.Status = "Executing"
			tasks[i].Status = "Executing"
			tmpTask.WorkerId = workerId
			tasks[i].WorkerId = workerId
			return &tmpTask
		}
	}
	return &Task{"NoTask", "Finished", -1, "", -1}	// 全部finish了
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
		Workers: make(map[int]bool), 		// 辅助函数可注释掉，每隔1s查看worker状态
	}

	// Your code here.
	c.nMap = len(files)
	c.nReduce = nReduce
	c.mapTasks = make([]Task, 0, c.nMap)	// 初始化空白map和reduce的task切片
	c.reduceTasks = make([]Task, 0, nReduce)

	for i := 0; i < c.nMap; i++ {			// 分配MapTask
		tmpTask := Task{"MapTask", "NotStarted", i, files[i], -1}
		c.mapTasks = append(c.mapTasks, tmpTask)
	}

	for i := 0; i < c.nReduce; i++ {		// 分配空白ReduceTask
		tmpTask := Task{"ReduceTask", "NotStarted", i, "", -1}
		c.reduceTasks = append(c.reduceTasks, tmpTask)
	}

	// var mu sync.Mutex					// 辅助函数可注释掉，每隔1s查看worker状态
	// go printMapContent(c.Workers, &mu)	// 辅助函数可注释掉，每隔1s查看worker状态
	var mu sync.Mutex						// 辅助函数可注释掉，每隔1s查看task状态
	go checkTaskStatus(c.mapTasks, &mu)		// 辅助函数可注释掉，每隔1s查看task状态
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
func checkTaskStatus(tasks []Task, mu *sync.Mutex) {		// 辅助函数可注释掉，每隔1s查看task状态
	for {
		mu.Lock()
		for _, task := range tasks {
			fmt.Printf("   %-10v | %-10v | %v | %-24v | %v", task.Type, task.Status, task.Index, task.FileName, task.WorkerId)
			fmt.Printf("\n")
		}
		mu.Unlock()
		fmt.Printf("\n")
        time.Sleep(1 * time.Second) // 等待 1 秒
	}
}