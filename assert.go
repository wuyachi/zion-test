package main

type Assertion interface{}

func Assert(result []byte, assertions []Assertion) error {
	return nil
}

/*
1. Euqal: value
2. NotEqual: value
3. contain(list): []value
4. noContain(list): []value
5. elementEqual(struct): field, value
6. elementNotEqual(struct): field, value
*/
