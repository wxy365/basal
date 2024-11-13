package types

import (
	"math/big"
	"strconv"
)

type Number[T NumberUnion] interface {
	String() string
	Parse(str string) error
	Bool() bool // convert number to bool
	Box(T)
	UnBox() T
}

type Int int

func (i *Int) String() string {
	return strconv.Itoa(int(*i))
}

func (i *Int) Parse(str string) error {
	ii, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		return err
	}
	*i = Int(ii)
	return nil
}

func (i *Int) Bool() bool {
	return *i > 0
}

func (i *Int) Box(t int) {
	*i = Int(t)
}

func (i *Int) UnBox() int {
	return int(*i)
}

type I8 int8

func (i *I8) UnBox() int8 {
	return int8(*i)
}

func (i *I8) Box(t int8) {
	*i = I8(t)
}

func (i *I8) Bool() bool {
	return *i > 0
}

func (i *I8) String() string {
	return strconv.Itoa(int(*i))
}

func (i *I8) Parse(str string) error {
	ii, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		return err
	}
	*i = I8(ii)
	return nil
}

type I16 int16

func (i *I16) UnBox() int16 {
	return int16(*i)
}

func (i *I16) Box(t int16) {
	*i = I16(t)
}

func (i *I16) Bool() bool {
	return *i > 0
}

func (i *I16) String() string {
	return strconv.Itoa(int(*i))
}

func (i *I16) Parse(str string) error {
	ii, err := strconv.ParseInt(str, 10, 16)
	if err != nil {
		return err
	}
	*i = I16(ii)
	return nil
}

type I32 int8

func (i *I32) UnBox() int32 {
	return int32(*i)
}

func (i *I32) Box(t int32) {
	*i = I32(t)
}

func (i *I32) Bool() bool {
	return *i > 0
}

func (i *I32) String() string {
	return strconv.Itoa(int(*i))
}

func (i *I32) Parse(str string) error {
	ii, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return err
	}
	*i = I32(ii)
	return nil
}

type I64 int64

func (i *I64) UnBox() int64 {
	return int64(*i)
}

func (i *I64) Box(t int64) {
	*i = I64(t)
}

func (i *I64) Bool() bool {
	return *i > 0
}

func (i *I64) String() string {
	return strconv.Itoa(int(*i))
}

func (i *I64) Parse(str string) error {
	ii, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*i = I64(ii)
	return nil
}

type Uint uint

func (u *Uint) UnBox() uint {
	return uint(*u)
}

func (u *Uint) Box(t uint) {
	*u = Uint(t)
}

func (u *Uint) Bool() bool {
	return *u > 0
}

func (u *Uint) String() string {
	return strconv.Itoa(int(*u))
}

func (u *Uint) Parse(str string) error {
	ii, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		return err
	}
	*u = Uint(ii)
	return nil
}

type U8 uint8

func (u *U8) UnBox() uint8 {
	return uint8(*u)
}

func (u *U8) Box(t uint8) {
	*u = U8(t)
}

func (u *U8) Bool() bool {
	return *u > 0
}

func (u *U8) String() string {
	return strconv.Itoa(int(*u))
}

func (u *U8) Parse(str string) error {
	ii, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		return err
	}
	*u = U8(ii)
	return nil
}

type U16 uint16

func (u *U16) UnBox() uint16 {
	return uint16(*u)
}

func (u *U16) Box(t uint16) {
	*u = U16(t)
}

func (u *U16) Bool() bool {
	return *u > 0
}

func (u *U16) String() string {
	return strconv.Itoa(int(*u))
}

func (u *U16) Parse(str string) error {
	ii, err := strconv.ParseInt(str, 10, 16)
	if err != nil {
		return err
	}
	*u = U16(ii)
	return nil
}

type U32 uint32

func (u *U32) UnBox() uint32 {
	return uint32(*u)
}

func (u *U32) Box(t uint32) {
	*u = U32(t)
}

func (u *U32) Bool() bool {
	return *u > 0
}

func (u *U32) String() string {
	return strconv.Itoa(int(*u))
}

func (u *U32) Parse(str string) error {
	ii, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return err
	}
	*u = U32(ii)
	return nil
}

type U64 uint64

func (u *U64) UnBox() uint64 {
	return uint64(*u)
}

func (u *U64) Box(t uint64) {
	*u = U64(t)
}

func (u *U64) Bool() bool {
	return *u > 0
}

func (u *U64) String() string {
	return strconv.Itoa(int(*u))
}

func (u *U64) Parse(str string) error {
	ii, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*u = U64(ii)
	return nil
}

type F32 float32

func (f *F32) UnBox() float32 {
	return float32(*f)
}

func (f *F32) Box(t float32) {
	*f = F32(t)
}

func (f *F32) Bool() bool {
	return *f > 0
}

func (f *F32) String() string {
	return strconv.FormatFloat(float64(*f), 'f', 9, 32)
}

func (f *F32) Parse(str string) error {
	ii, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return err
	}
	*f = F32(ii)
	return nil
}

type F64 float64

func (f *F64) UnBox() float64 {
	return float64(*f)
}

func (f *F64) Box(t float64) {
	*f = F64(t)
}

func (f *F64) Bool() bool {
	return *f > 0
}

func (f *F64) String() string {
	return strconv.FormatFloat(float64(*f), 'f', 9, 64)
}

func (f *F64) Parse(str string) error {
	ii, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}
	*f = F64(ii)
	return nil
}

type C64 complex64

func (c *C64) UnBox() complex64 {
	return complex64(*c)
}

func (c *C64) Box(t complex64) {
	*c = C64(t)
}

func (c *C64) Bool() bool {
	return real(c.UnBox()) > 0
}

func (c *C64) String() string {
	return strconv.FormatComplex(complex128(c.UnBox()), 'c', 9, 64)
}

func (c *C64) Parse(str string) error {
	ii, err := strconv.ParseComplex(str, 64)
	if err != nil {
		return err
	}
	*c = C64(ii)
	return nil
}

type C128 complex128

func (c *C128) UnBox() complex128 {
	return complex128(*c)
}

func (c *C128) Box(t complex128) {
	*c = C128(t)
}

func (c *C128) Bool() bool {
	return real(c.UnBox()) > 0
}

func (c *C128) String() string {
	return strconv.FormatComplex(c.UnBox(), 'c', 9, 128)
}

func (c *C128) Parse(str string) error {
	ii, err := strconv.ParseComplex(str, 128)
	if err != nil {
		return err
	}
	*c = C128(ii)
	return nil
}

type BigInt big.Int

func (b *BigInt) String() string {
	ii := b.UnBox()
	return ii.String()
}

func (b *BigInt) Parse(str string) error {
	ii := new(I64)
	err := ii.Parse(str)
	if err != nil {
		return err
	}
	bi := big.NewInt(ii.UnBox())
	*b = BigInt(*bi)
	return nil
}

func (b *BigInt) Bool() bool {
	bi := b.UnBox()
	return bi.Cmp(new(big.Int)) > 0
}

func (b *BigInt) Box(t big.Int) {
	*b = BigInt(t)
}

func (b *BigInt) UnBox() big.Int {
	return big.Int(*b)
}

type BigFloat big.Float

func (b *BigFloat) String() string {
	ii := b.UnBox()
	return ii.String()
}

func (b *BigFloat) Parse(str string) error {
	f, _, err := big.ParseFloat(str, 10, 9, big.ToNearestAway)
	if err != nil {
		return err
	}
	*b = BigFloat(*f)
	return nil
}

func (b *BigFloat) Bool() bool {
	bi := b.UnBox()
	return bi.Cmp(new(big.Float)) > 0
}

func (b *BigFloat) Box(t big.Float) {
	*b = BigFloat(t)
}

func (b *BigFloat) UnBox() big.Float {
	return big.Float(*b)
}

type NumberUnion interface {
	BasicNumberUnion | ComplexUnion | BigNumberUnion
}

type BasicNumberUnion interface {
	IntUnion | UintUnion | FloatUnion
}

type IntUnion interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type UintUnion interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type FloatUnion interface {
	~float32 | ~float64
}

type ComplexUnion interface {
	~complex64 | ~complex128
}

type BigNumberUnion interface {
	big.Int | big.Float
}
