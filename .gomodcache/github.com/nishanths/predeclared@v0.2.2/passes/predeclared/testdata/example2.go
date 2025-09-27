package print

import cap "foo/capacity"

func make(i *int) T {
	copy := cap.Var(*i)
	return T{copy}
}
