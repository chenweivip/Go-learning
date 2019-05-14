package base

import (
    "runtime"
    "strings"
)

// GetInvokerLocation 用于获得调用位置。
func GetInvokerLocation(skipNumber int) (funcPath string, fileName string, line int) {
    // runtime.Caller 拿到当前执行的文件名和行号
    // skip如果是0，返回当前调用Caller函数的函数名、文件、程序指针PC
    // runtime.Caller(2)  上2层函数信息
    pc, file, line, ok := runtime.Caller(skipNumber)
    if !ok {
        return "", "", -1
    }
    if index := strings.LastIndex(file, "/"); index > 0{
        fileName = file[index+1 : len(file)]
    }
    funcPtr := runtime.FuncForPC(pc)
    if funcPtr != nil{
        funcPath = funcPtr.Name()
    }
    return funcPath, fileName, line
}
