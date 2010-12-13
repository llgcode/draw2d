package postscript

func pop(interpreter *Interpreter) {
	interpreter.Pop()
}

func dup(interpreter *Interpreter) {
	interpreter.Push(interpreter.Peek())
}

func index(interpreter *Interpreter) {
	f := interpreter.PopInt()
	interpreter.Push(interpreter.Get(int(f)))
}

func roll(interpreter *Interpreter) {
	j := interpreter.PopInt()
	n := interpreter.PopInt()
	values := interpreter.PopValues(n)
	j %= n
	for i := 0; i < n; i++ {
		interpreter.Push(values[(n+i-j)%n])
	}
}

func copystack(interpreter *Interpreter) {
	n := interpreter.PopInt()
	values := interpreter.GetValues(n)
	for _, value := range values {
		interpreter.Push(value)
	}
}

func exch(interpreter *Interpreter) {
	value1 := interpreter.Pop()
	value2 := interpreter.Pop()
	interpreter.Push(value1)
	interpreter.Push(value2)
}

func initStackOperator(interpreter * Interpreter) {
	interpreter.SystemDefine("pop", NewOperator(pop))
	interpreter.SystemDefine("dup", NewOperator(dup))
	interpreter.SystemDefine("index", NewOperator(index))
	interpreter.SystemDefine("copy", NewOperator(copystack))
	interpreter.SystemDefine("roll", NewOperator(roll))
	interpreter.SystemDefine("exch", NewOperator(exch))
}