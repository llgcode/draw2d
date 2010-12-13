package postscript


import (
	"log"
)


func get(interpreter *Interpreter) {
	name := interpreter.PopName()
	dictionary := interpreter.Pop().(Dictionary)
	interpreter.Push(dictionary[name])
}

func def(interpreter *Interpreter) {
	value := interpreter.Pop()
	name := interpreter.PopName()
	if p, ok := value.(*ProcedureDefinition); ok {
		value = NewProcedure(p)
	}
	interpreter.Define(name, value)
}

func bind(interpreter *Interpreter) {
	pdef := interpreter.PopProcedureDefinition()
	values := make([]Value, len(pdef.Values))
	for i, value := range pdef.Values {
		if s, ok := value.(string); ok {
			firstChar := s[0]
			if firstChar != '(' && firstChar != '/' {
				v, _ := interpreter.FindValueInDictionaries(s)
				operator, isOperator := v.(Operator)
				if isOperator {
					values[i] = operator
				} else {
					values[i] = value
				}
			} else {
				values[i] = value
			}
		} else {
			values[i] = value
		}
	}
	pdef.Values = values
	interpreter.Push(pdef)
}

func load(interpreter *Interpreter) {
	name := interpreter.PopName()
	value, _ := interpreter.FindValueInDictionaries(name)
	if value == nil {
		log.Printf("Can't find value %s\n", name)
	}
	interpreter.Push(value)
}

func begin(interpreter *Interpreter) {
	interpreter.PushDictionary(interpreter.Pop().(Dictionary))
}

func end(interpreter *Interpreter) {
	interpreter.PopDictionary()
}

func dict(interpreter *Interpreter) {
	interpreter.Push(NewDictionary(interpreter.PopInt()))
}

func where(interpreter *Interpreter) {
	key := interpreter.PopName()
	_, dictionary := interpreter.FindValueInDictionaries(key)
	if dictionary == nil {
		interpreter.Push(false)
	} else {
		interpreter.Push(dictionary)
		interpreter.Push(true)
	}
}

func initDictionaryOperators(interpreter *Interpreter) {
	interpreter.SystemDefine("get", NewOperator(get))
	interpreter.SystemDefine("def", NewOperator(def))
	interpreter.SystemDefine("load", NewOperator(load))
	interpreter.SystemDefine("bind", NewOperator(bind))
	interpreter.SystemDefine("begin", NewOperator(begin))
	interpreter.SystemDefine("end", NewOperator(end))
	interpreter.SystemDefine("dict", NewOperator(dict))
	interpreter.SystemDefine("where", NewOperator(where))
}