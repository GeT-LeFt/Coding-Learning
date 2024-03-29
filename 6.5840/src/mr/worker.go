package mr

import "fmt"
import "log"
import "net/rpc"
import "hash/fnv"
import "os"
import "time"


//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}


//
// main/mrworker.go calls this function.
//
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// example
	// CallExample()

	// Your worker implementation here.
	getReduceNumebr()
	// for {						// 循环要任务
		requestTask()
		doMap()
		doReduce()
		time.Sleep(time.Second * 1)
	// }
}

//
// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
//
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}    // Defined in rpc.go

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}    // Defined in rpc.go

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
}

func getReduceNumebr() {

}

func requestTask() {
	args := RequestTaskArgs{}
	args.WorkerId = os.Getpid()		// 把程序ID作为worker的id
	reply := RequestTaskReply{}
	ok := call("Coordinator.RequestTaskReply", &args, &reply)
	if ok {
		fmt.Printf("Worker-%v: request task!\n", args.WorkerId)
		fmt.Printf("reply: %v\n", reply)
	} else {
		fmt.Printf("call failed!\n")
	}
}

func doMap() {
	fmt.Printf("Worker is doing Map\n")
}

func doReduce() {
	fmt.Printf("Worker is doing Reduce\n")
}

//
// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()	// realized in rpc.go
	c, err := rpc.DialHTTP("unix", sockname)	// realized in "net/rpc" and try to create a RPC client connection
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)    // realized in "net/rpc" and making remote procedure calls (RPC) on an established RPC client connection
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
