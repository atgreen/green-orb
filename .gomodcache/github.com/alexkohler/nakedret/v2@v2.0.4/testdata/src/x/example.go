package x

func justone() (one string) {
	one = `one`
	return // want "naked return in func `justone` with 3 lines of code"
}

func both() (one, two string) {
	one = `one`
	two = `two`
	return // want "naked return in func `both` with 4 lines of code"
}

func three() (one, two string, three int) {
	one = `one`
	two = `two`
	three = 3
	return // want "naked return in func `three` with 5 lines of code"
}

func longFunc() (story string) {
	story = `there once was a man named Michael Finnegan`
	story += `He had whiskers on his chin-negan`
	story += `The wind blew em out and then blew in again`
	story += `Poor old Michael Finnegan Begin again`
	story += `there once was a man named Michael Finnegan`
	story += `He had whiskers on his chin-negan`
	story += `They fell out and grew back in again`
	story += `Poor old Michael Finnegan Begin again`
	story += `there once was a man named Michael Finnegan`
	story += `He had whiskers on his chin-negan`
	story += `The wind blew em out and then blew in again`
	story += `Poor old Michael Finnegan Begin again`
	story += `there once was a man named Michael Finnegan`
	story += `he grew fat and then grew thin again`
	story += `Then he died and had to begin again`
	story += `Poor old Michael Finnegan Begin again`
	story += `there once was a man named Michael Finnegan`
	story += `He had whiskers on his chin-negan`
	story += `The wind blew em out and then blew in again`
	story += `Poor old Michael Finnegan Begin again`
	story += `there once was a man named Michael Finnegan`
	story += `He had whiskers on his chin-negan`
	story += `The wind blew em out and then blew in again`
	story += `Poor old Michael Finnegan Begin again`
	story += `there once was a man named Michael Finnegan`
	story += `He had whiskers on his chin-negan`
	story += `The wind blew em out and then blew in again`
	story += `Poor old Michael Finnegan Begin again`
	story += `there once was a man named Michael Finnegan`
	story += `He had whiskers on his chin-negan`
	story += `The wind blew em out and then blew in again`
	story += `Poor old Michael Finnegan Begin again`
	return // want "naked return in func `longFunc` with 34 lines of code"
}
