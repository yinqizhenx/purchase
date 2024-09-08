package cel

import (
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"testing"
)

func main() {

}

func Test_exprReplacement(t *testing.T) {
	var str = `"Hello world! I'm " + name + "."`
	env, err := cel.NewEnv(
		cel.Variable("name", cel.StringType), // 参数类型绑定
	)
	if err != nil {
		t.Fatal(err)
	}
	ast, iss := env.Compile(str) // 编译，校验，执行 str
	if iss.Err() != nil {
		t.Fatal(iss.Err())
	}
	program, err := env.Program(ast)
	if err != nil {
		t.Fatal(err)
	}
	// 初始化 name 变量的值
	values := map[string]interface{}{"name": "CEL"}
	// 传给内部程序并返回执行的结果
	out, detail, err := program.Eval(values)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(detail)
	fmt.Println(out)
}

func Test_LogicExpr1(t *testing.T) {
	var str = `100 + 200 >= 300`
	env, err := cel.NewEnv()
	if err != nil {
		t.Fatal(err)
	}
	ast, iss := env.Compile(str)
	if iss.Err() != nil {
		t.Fatal(iss.Err())
	}
	prog, err := env.Program(ast)
	if err != nil {
		t.Fatal(err)
	}
	out, detail, err := prog.Eval(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(out)
	fmt.Println(detail)
}

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// 求所有传入参数的和
func Add[T Integer](param1, param2 T, ints ...T) T {
	sum := param1 + param2
	for _, v := range ints {
		sum += v
	}
	return sum
}

func Test_add(t *testing.T) {
	str := `6 == Add(age1, age2, age3)`
	env, _ := cel.NewEnv(
		cel.Variable("age1", cel.IntType),
		cel.Variable("age2", cel.IntType),
		cel.Variable("age3", cel.IntType),
		cel.Function("Add", cel.Overload(
			"Add",
			[]*cel.Type{cel.IntType, cel.IntType, cel.IntType},
			cel.IntType,
			cel.FunctionBinding(func(vals ...ref.Val) ref.Val {
				var xx []int64
				for _, v := range vals {
					xx = append(xx, v.Value().(int64))
				}
				return types.Int(Add[int64](xx[0], xx[1], xx[2:]...))
			}),
		)),
	)
	ast, iss := env.Compile(str)
	if iss.Err() != nil {
		t.Fatal(iss.Err())
	}
	prog, err := env.Program(ast)
	if err != nil {
		t.Fatal(err)
	}
	val, detail, err := prog.Eval(map[string]interface{}{"age1": 1, "age2": 2, "age3": 3})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(detail)
	fmt.Println(val)
}
