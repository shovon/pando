package maybeany

type MaybeAny struct {
	node   interface{}
	exists bool
}

func Something(node interface{}) MaybeAny {
	return MaybeAny{
		node:   node,
		exists: true,
	}
}

func Nothing() MaybeAny {
	return MaybeAny{
		node:   nil,
		exists: false,
	}
}

func (m MaybeAny) Get() (interface{}, bool) {
	return m.node, m.exists
}
