package main

import (
	// "fmt"
)


var Handlers = map[string]func([]Value) Value {
	"PING": ping,
	"SET": set,
	"GET": get,
	// "HGET": hget,
	// "HSET": hset,
}

var SETs = map[string]string{}
var HSETs = map[string]map[string]string{}


func ping(args []Value) Value {
	return Value{typ: "string", str: "PONG"}
}

func set(args []Value) Value{
	SETs[args[0].bulk] = args[1].bulk
	return Value{typ: "string", str: "SET"}
}

func get(args []Value) Value{
	return Value{typ: "bulk", bulk: SETs[args[0].bulk]}
}

// func hget(args []Value) Value{
// 	fmt.Println("HGET")
// }

// func hset(args []Value) Value{
// 	fmt.Println("HSET")
// }



