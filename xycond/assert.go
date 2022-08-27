package xycond

import "reflect"

// AssertEqual panics if a is different from b.
func AssertEqual(a, b any) {
	ExpectEqual(a, b).assert()
}

// AssertNotEqual panics if a is equal to b.
func AssertNotEqual(a, b any) {
	ExpectNotEqual(a, b).assert()
}

// AssertLessThan panics if a is not less than b.
func AssertLessThan[t number](a, b t) {
	ExpectLessThan(a, b).assert()
}

// AssertNotLessThan panics if a is less than b.
func AssertNotLessThan[t number](a, b t) {
	ExpectNotLessThan(a, b).assert()
}

// AssertGreaterThann panics if a is not greater than b.
func AssertGreaterThan[t number](a, b t) {
	ExpectGreaterThan(a, b).assert()
}

// AssertNotGreaterThan panics if a is greater than b.
func AssertNotGreaterThan[t number](a, b t) {
	ExpectNotGreaterThan(a, b).assert()
}

// AssertPanic panics if the function doesn't panic.
func AssertPanic(f func()) {
	ExpectPanic(f).assert()
}

// AssertNotPanic panics if the function panics.
func AssertNotPanic(f func()) {
	ExpectNotPanic(f).assert()
}

// AssertZero panics if the parameter is not zero.
func AssertZero[t number](a t) {
	ExpectZero(a).assert()
}

// AssertNotZero panics if the parameter is zero.
func AssertNotZero[t number](a t) {
	ExpectNotZero(a).assert()
}

// AssertNil panics if the parameter is not nil.
func AssertNil(a any) {
	ExpectNil(a).assert()
}

// AssertNotNil panics if the parameter is nil.
func AssertNotNil(a any) {
	ExpectNotNil(a).assert()
}

// AssertEmpty panics if the parameter is not empty.
func AssertEmpty(a any) {
	ExpectEmpty(a).assert()
}

// AssertNotEmpty panics if the parameter is empty.
func AssertNotEmpty(a any) {
	ExpectNotEmpty(a).assert()
}

// AssertIs panics if value doesn't belongs to any passed kinds.
func AssertIs(v any, kinds ...reflect.Kind) {
	ExpectIs(v, kinds...).assert()
}

// AssertIsNot panics if value belongs to one of passed kinds.
func AssertIsNot(v any, kinds ...reflect.Kind) {
	ExpectIsNot(v, kinds...).assert()
}

// AssertSame panics if there is at least value' type different from the rest.
func AssertSame(v ...any) {
	ExpectSame(v...).assert()
}

// AssertNotSame panics if all values' type are the same.
func AssertNotSame(v ...any) {
	ExpectNotSame(v...).assert()
}

// AssertWritable panics if the parameter is not a writable channel.
func AssertWritable(c any) {
	ExpectWritable(c).assert()
}

// AssertNotWritable panics if the parameter is a writable channel.
func AssertNotWritable(c any) {
	ExpectNotWritable(c).assert()
}

// AssertReadable panics if the parameter is not a readable channel.
func AssertReadable(c any) {
	ExpectReadable(c).assert()
}

// AssertNotReadable panics if the parameter is a readable channel.
func AssertNotReadable(c any) {
	ExpectNotReadable(c).assert()
}

// AssertError panics if the err doesn't belong to any targets.
func AssertError(err error, targets ...error) {
	ExpectError(err, targets...).assert()
}

// AssertErrorNot panics if the err belongs to one of targets.
func AssertErrorNot(err error, targets ...error) {
	ExpectErrorNot(err, targets...).assert()
}

// AssertTrue panics if the condition is false.
func AssertTrue(b bool) {
	ExpectTrue(b).assert()
}

// AssertFalse panics if the condition is true.
func AssertFalse(b bool) {
	ExpectFalse(b).assert()
}
