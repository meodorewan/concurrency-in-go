# Chapter 5: Concurrency at scale 
This chapter briefs techniques which might be helpful in big project 

### 1. Error Propagation
It is good to consider **error handling** as first class problem when designing the system.

Good error should contains:

*What happened*: to tell exactly what happened: i.e "disk full", "network timeout", etc.
	
*When and where it occured*: Trace stack is neccessary but not in error message. Can store it in another different field of error struct. error should contains context it's running within. And it's good to have timestamp.
	
*User-friendly*: error message should be clear and meaningful.
	
*Identification*: RequestId, etc.

~~~
type MyError struct {
	Inner error
	Message string
	StackTrace string
	Misc map[string]interface{}
}

func wrapError(err error, messagef string, msgArgs ...interface{}) MyError {
	return MyError{
		Inner: err,
		Message: fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),
		Misc: make(map[string]interface{}),
	}
}

// MyError implements error interface
func (err MyError) Error() string { 
	return err.Message
}

// handle error
func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))
	log.Printf("%#v", err)
	fmt.Printf("[%v] %v", key, message)
}

// result
[logID: 1]: 21:46:07 main.LowLevelErr{error:main.MyError{Inner:(*os.PathError)(0xc4200123f0),
Message:"stat /bad/job/binary: no such file or directory",
StackTrace:"goroutine 1 [running]:
runtime/debug.Stack(0xc420012420, 0x2f, 0xc420045d80)
/home/kate/.guix-profile/src/runtime/debug/stack.go:24 +0x79
main.wrapError(0x530200, 0xc4200123f0, 0xc420012420, 0x2f, 0x0, 0x0,
0x0, 0x0, 0x0, 0x0, ...)
/tmp/babel-79540aE/go-src-7954NTK.go:22 +0x62
main.isGloballyExec(0x4d1313, 0xf, 0xc420045eb8, 0x487649, 0xc420056050)
/tmp/babel-79540aE/go-src-7954NTK.go:37 +0xaa
main.runJob(0x4cfada, 0x1, 0x4d4c35, 0x22)
/tmp/babel-79540aE/go-src-7954NTK.go:47 +0x48
main.main()
/tmp/babel-79540aE/go-src-7954NTK.go:67 +0x63
", Misc:map[string]interface {}{}}}

And message is printed to stdout

[1] There was an unexpected issue; please report this as a bug.
~~~

### 2. Timeout and cancellation
There are some reasons for timeout and cancellation

1. Stability
	
	Don't want a request taking long time, waiting computer resources => use timeout
	 
2. Stale data
	
	Some logics require data should be processed before new data coming => use timeout and cancel
	
3. Deadlock prevention

Go `context` package is recommended.

### 3. Heartbeats
### 4. Replicated request
### 5. Rate limiting
### 6. Healing unhealthy goroutines


