package evaluator

import (
	"fmt"

	"github.com/Savvelius/go-interp/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(strs ...object.Object) object.Object {
			if len(strs) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(strs))
			}

			if str, ok := strs[0].(*object.String); ok {
				return &object.Integer{Value: int64(len(str.Value))}
			}

			if arr, ok := strs[0].(*object.Array); ok {
				return &object.Integer{Value: int64(len(arr.Value))}
			}

			if hash, ok := strs[0].(*object.Hash); ok {
				return &object.Integer{Value: int64(len(hash.Pairs))}
			}

			return newError("argument to `len` not supported, got %s", strs[0].Type())
		},
	},
	"typeOf": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			arg := args[0]
			return &object.String{Value: string(arg.Type())}
		},
	},
	"print": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect() + " ")
			}
			fmt.Print("\n")

			return NULL
		},
	},
}
