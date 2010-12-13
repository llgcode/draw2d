package postscript


func ifoperator(interpreter *Interpreter) {
	operator := NewProcedure(interpreter.PopProcedureDefinition())
	condition := interpreter.PopBoolean()
	if condition {
		operator.Execute(interpreter)
	}
}

func ifelse(interpreter *Interpreter) {
	operator2 := NewProcedure(interpreter.PopProcedureDefinition())
	operator1 := NewProcedure(interpreter.PopProcedureDefinition())
	condition := interpreter.PopBoolean()
	if condition {
		operator1.Execute(interpreter)
	} else {
		operator2.Execute(interpreter)
	}
}

func foroperator(interpreter *Interpreter) {
	proc := NewProcedure(interpreter.PopProcedureDefinition())
	limit := interpreter.PopFloat()
	inc := interpreter.PopFloat()
	initial := interpreter.PopFloat()

	for i := initial; i <= limit; i += inc {
		interpreter.Push(i)
		proc.Execute(interpreter)
	}
}



func initControlOperators(interpreter *Interpreter) {
	interpreter.SystemDefine("if", NewOperator(ifoperator))
	interpreter.SystemDefine("ifelse", NewOperator(ifelse))
	interpreter.SystemDefine("for", NewOperator(foroperator))
}