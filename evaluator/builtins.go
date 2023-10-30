package evaluator

import "github.com/Savvelius/go-interp/object"

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(strs ...object.Object) object.Object {
			if len(strs) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(strs))
			}

			str, ok := strs[0].(*object.String)
			if !ok {
				return newError("argument to `len` not supported, got %s", strs[0].Type())
			}

			return &object.Integer{Value: int64(len(str.Value))}
		},
	},
}
