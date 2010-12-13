package postscript

type OperatorFunc func(interpreter *Interpreter)

type PrimitiveOperator struct {
	f OperatorFunc
}

func NewOperator(f OperatorFunc) *PrimitiveOperator {
	return &PrimitiveOperator{f}
}

func (o *PrimitiveOperator) Execute(interpreter *Interpreter) {
	o.f(interpreter)
}


func save(interpreter *Interpreter) {
}

func restore(interpreter *Interpreter) {
}



func initSystemOperators(interpreter *Interpreter) {
	interpreter.SystemDefine("save", NewOperator(save))
	interpreter.SystemDefine("restore", NewOperator(restore))
	initStackOperator(interpreter)
	initMathOperators(interpreter)
	initDictionaryOperators(interpreter)
	initRelationalOperators(interpreter)
	initControlOperators(interpreter)
	initDrawingOperators(interpreter)
}
