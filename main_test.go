package gochain

import (
	"errors"
	"testing"

	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
)

func TestChain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chain Suite")
}

type ChainContext struct {
	FirstName string
	LastName  string
}

type ChainResult struct {
	Value int
}

var _ = Describe("Chain", func() {
	const FIRST_NAME = "John"
	const LAST_NAME = "Doe"

	var chain *Chain[ChainContext, ChainResult]

	BeforeEach(func() {
		chain = NewChain[ChainContext, ChainResult]()
	})

	It("should the change context", func() {
		err := chain.
			Add(func(
				chain *Chain[ChainContext, ChainResult],
				_ StopFunc,
			) error {
				chain.UpdateContext("FirstName", FIRST_NAME)
				return nil
			}).
			Add(func(
				chain *Chain[ChainContext, ChainResult],
				_ StopFunc,
			) error {
				chain.UpdateContext("LastName", LAST_NAME)
				return nil
			}).
			Run()

		Expect(err).Should(BeNil())

		ctx := chain.GetContext()

		Expect(ctx.FirstName).Should(Equal(FIRST_NAME))
		Expect(ctx.LastName).Should(Equal(LAST_NAME))
	})

	It("should return error", func() {
		unknownError := errors.New("unknown error")

		err := chain.
			Add(func(
				chain *Chain[ChainContext, ChainResult],
				_ StopFunc,
			) error {
				return unknownError
			}).
			Add(func(
				chain *Chain[ChainContext, ChainResult],
				_ StopFunc,
			) error {
				chain.UpdateContext("LastName", "Doe")
				return nil
			}).
			Run()

		Expect(err).Should(MatchError(unknownError))
	})

	It("should assign the value -100 to the ChainResult.Value", func() {
		const EXPECTED_VALUE = -100

		err := chain.
			Add(func(
				chain *Chain[ChainContext, ChainResult],
				_ StopFunc,
			) error {
				return chain.UpdateResult("Value", EXPECTED_VALUE)
			}).
			Run()

		Expect(err).Should(BeNil())

		resultValue := chain.GetResult().Value

		Expect(resultValue).Should(Equal(EXPECTED_VALUE))
	})

	It("should stop abruptly", func() {
		err := chain.
			Add(func(
				chain *Chain[ChainContext, ChainResult],
				_ StopFunc,
			) error {
				chain.UpdateContext("FirstName", FIRST_NAME)
				return nil
			}).
			Add(func(
				chain *Chain[ChainContext, ChainResult],
				stop StopFunc,
			) error {
				stop()
				return nil
			}).
			Add(func(
				chain *Chain[ChainContext, ChainResult],
				_ StopFunc,
			) error {
				chain.UpdateContext("LastName", LAST_NAME)
				return nil
			}).
			Run()

		Expect(err).Should(BeNil())

		ctx := chain.GetContext()

		Expect(ctx.FirstName).Should(Equal(FIRST_NAME))
		Expect(ctx.LastName).Should(BeEmpty())
	})
})
