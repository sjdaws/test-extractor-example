package main

import(
	"testing"
)

// Blah is not a test function so should be ignored
//
// FEATURE(ABT-1233): Test this function is ignored.
func Blah() {
}

// testMeh is not exported and should be ignored
//
// FEATURE(ABT-1233): Test this function is ignored.
func testMeh() {

}

// TestingPayAnyone tests payments to third parties
//
// FEATURE(ABT-123): Create a payment with the downstream services and
// verify the payment has been created correctly and payment status is returned via
// Command Centre.
// BUG(ABT-334): Added additional verification to ensure that payment description allows more than 4 characters
//
// BUG(ABT-9282)Another bug but this one has the description mashed up against the brackets.
// INVALID(ABT-782) This category is invalid and should be treated as just a comment.
//
// This is a trailing comment
// BUG(ABT-1930): This bug is right at the bottom of the docblock
func TestPayAnyone(_ *testing.T) {

}

// TestTransfer tests transfer between accounts
// FEATURE(ABT-909): Create a payment between two accounts owned by the same
// entity.
func TestTransfer(_ *testing.T) {

}

// TestNoTraces has no traces so should be ignored
func TestNoTraces(_ *testing.T) {

}
