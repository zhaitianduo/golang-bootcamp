package myJson

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var B interface {
	B()
}

type Bb []int

func (b *Bb) B() {
	return
}

type S struct {
	int
	I  int
	S  *string
	T  bool
	A  [3]int
	Ss []S2
	S1
	S2 *S2
	Bb
	Err error
	M   map[string]int
	N   map[int]string
}

type S1 struct {
	F   int
	Sss string
}

type S2 struct {
	I int
}

type Ss struct {
	int
	I  int
	s  *string
	T  bool
	a  [3]int
	Ss []S2
	S1
	s2 *S2
	Bb
	err error
	M   map[string]int
	N   map[int]string
}

func TestPrintJson(t *testing.T) {
	testS := "test"
	m := make(map[string]int)
	m["abc"] = 123
	n := make(map[int]string)
	n[123] = "abc"
	s := Ss{6666, 1, &testS, true, [3]int{3, 2, 1}, []S2{S2{1}, S2{2}, S2{3}}, S1{3, "hello"}, &S2{4}, Bb([]int{333}), errors.New("aa"), m, n}
	re, err := PrintStructJson(s)
	if err != nil {
		fmt.Println(err.Error)
	}
	fmt.Println("result json is: ", re)
	assert.Equal(t, `{"int":6666,"I":1,"s":"test","T":true,"a":[3,2,1],"Ss":[{"I":1},{"I":2},{"I":3}],"F":3,"Sss":"hello","s2":{"I":4},"Bb":[333],"err":{"s":"aa"},"M":{"abc":123},"N":{123:"abc"}}`, re)

	//syBytes, err := json.Marshal(s)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println("system json is: ", string(syBytes))
}

func BenchmarkJson(b *testing.B) {
	testS := "test"
	m := make(map[string]int)
	m["abc"] = 123
	n := make(map[int]string)
	n[123] = "abc"
	s := S{6666, 1, &testS, true, [3]int{3, 2, 1}, []S2{S2{1}, S2{2}, S2{3}}, S1{3, "hello"}, &S2{4}, Bb([]int{333}), errors.New("aa"), m, n}
	for i := 0; i < b.N; i++ {
		PrintStructJson(s)
	}
}

func BenchmarkSystem(b *testing.B) {
	testS := "test"
	m := make(map[string]int)
	m["abc"] = 123
	n := make(map[int]string)
	s := S{6666, 1, &testS, true, [3]int{3, 2, 1}, []S2{S2{1}, S2{2}, S2{3}}, S1{3, "hello"}, &S2{4}, Bb([]int{333}), errors.New("aa"), m, n}
	for i := 0; i < b.N; i++ {
		json.Marshal(s)
	}
}
