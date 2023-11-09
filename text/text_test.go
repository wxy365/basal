package text

import (
	"fmt"
	"testing"
)

func TestPascal2Snake(t *testing.T) {
	//s := "ChinaIDCard"
	//o := Pascal2Snake(s)
	//if o != "china_id_card" {
	//	t.Fail()
	//}
	//
	//s = "china_id_card"
	//o = Snake2Pascal(s)
	//if o != "ChinaIdCard" {
	//	t.Fail()
	//}

	//sb := &Sub{}
	//fmt.Println(sb.Method1())

	//var a []Super
	//b := Super{}
	//c := &b
	//a = append(a, *c)
	//
	//fmt.Printf("%p,%p", c, &a[0])

	//d := Super{}
	var d Super
	e := &d
	f := d
	g := *e
	h := make([]Super, 1)
	h[0] = f

	fmt.Printf("d:%p, e:%p, f:%p, g:%p, h[0]:%p", &d, e, &f, &g, &h[0])
	fmt.Println()
	fmt.Printf("d:%+v, f: %+v", d, f)
}

type Super struct {
	A string
	B int
}

func (s *Super) Method1() string {
	return s.Method2()
}

func (s *Super) Method2() string {
	return "Super.Method2"
}

type Sub struct {
	Super
}

func (s *Sub) Method2() string {
	return "Sub.Method2"
}
