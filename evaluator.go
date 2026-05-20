package main

// Environment: Değişkenleri saklayan yapı
type Environment struct {
	store map[string]int
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]int)}
}

func (e *Environment) Get(name string) (int, bool) {
	val, ok := e.store[name]
	return val, ok
}

func (e *Environment) Set(name string, val int) {
	e.store[name] = val
}

// Evaluate: AST ağacını hafızayı kullanarak çalıştırır
func Evaluate(node Node, env *Environment) int {
	switch n := node.(type) {
	
	case *IntegerLiteral:
		return n.Value

	case *Identifier: // Değişkeni hafızadan çek
		val, ok := env.Get(n.Value)
		if !ok {
			return 0 // Değişken tanımlı değilse 0 dön
		}
		return val

	case *InfixExpression:
		left := Evaluate(n.Left, env)
		right := Evaluate(n.Right, env)
		
		switch n.Operator {
		case "+": return left + right
		case "-": return left - right
		case "*": return left * right
		case "/": return left / right
		}

	case *LetStatement: // Değişken atamasını hafızaya kaydet
		val := Evaluate(n.Value, env)
		env.Set(n.Name.Value, val)
		return val
	}
	
	return 0
}