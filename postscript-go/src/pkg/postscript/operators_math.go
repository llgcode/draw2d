package postscript

import (
	"math"
)
// begin Primitive Operator implementation


func add(interpreter *Interpreter) {
	f2 := interpreter.PopFloat()
	f1 := interpreter.PopFloat()
	interpreter.Push(f1 + f2)
}

func sub(interpreter *Interpreter) {
	f2 := interpreter.PopFloat()
	f1 := interpreter.PopFloat()
	interpreter.Push(f1 - f2)
}

func div(interpreter *Interpreter) {
	f2 := interpreter.PopFloat()
	f1 := interpreter.PopFloat()
	interpreter.Push(f1 / f2)
}
func mul(interpreter *Interpreter) {
	f2 := interpreter.PopFloat()
	f1 := interpreter.PopFloat()
	interpreter.Push(f1 * f2)
}

func sqrt(interpreter *Interpreter) {
	f := interpreter.PopFloat()
	interpreter.Push(float(math.Sqrt(float64(f))))
}

func atan(interpreter *Interpreter) {
	den := interpreter.PopFloat()
	num := interpreter.PopFloat()
	interpreter.Push(float(math.Atan2(float64(num), float64(den))) * (180.0 / math.Pi))
}

func cos(interpreter *Interpreter) {
	a := interpreter.PopFloat()
	interpreter.Push(float(math.Cos(float64(a))) * (180.0 / math.Pi))
}

func sin(interpreter *Interpreter) {
	a := interpreter.PopFloat()
	interpreter.Push(float(math.Sin(float64(a))) * (180.0 / math.Pi))
}

func neg(interpreter *Interpreter) {
	f := interpreter.PopFloat()
	interpreter.Push(-f)
}

func round(interpreter *Interpreter) {
	f := interpreter.PopFloat()
	interpreter.Push(float(int(f + 0.5)))
}

func abs(interpreter *Interpreter) {
	f := interpreter.PopFloat()
	interpreter.Push(float(math.Fabs(float64(f))))
}



func initMathOperators(interpreter *Interpreter) {
	interpreter.SystemDefine("add", NewOperator(add))
	interpreter.SystemDefine("sub", NewOperator(sub))
	interpreter.SystemDefine("mul", NewOperator(mul))
	interpreter.SystemDefine("div", NewOperator(div))
	interpreter.SystemDefine("sqrt", NewOperator(sqrt))
	interpreter.SystemDefine("cos", NewOperator(cos))
	interpreter.SystemDefine("sin", NewOperator(sin))
	interpreter.SystemDefine("atan", NewOperator(atan))
	interpreter.SystemDefine("neg", NewOperator(neg))
	interpreter.SystemDefine("round", NewOperator(round))
	interpreter.SystemDefine("abs", NewOperator(abs))
}