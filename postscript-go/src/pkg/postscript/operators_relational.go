package postscript


func eq(interpreter *Interpreter) {
	value1 := interpreter.Pop()
	value2 := interpreter.Pop()
	interpreter.Push(value1 == value2)
}

func lt(interpreter *Interpreter) {
	f2 := interpreter.PopFloat()
	f1 := interpreter.PopFloat()
	interpreter.Push(f1 < f2)
}

func initRelationalOperators(interpreter *Interpreter) {
	interpreter.SystemDefine("eq", NewOperator(eq))
	interpreter.SystemDefine("lt", NewOperator(lt))
}