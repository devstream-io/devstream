package trello

import (
	"fmt"

	trelloCommon "github.com/adlio/trello"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/trello"
)

var _ = Describe("board struct", func() {
	var (
		testBoard  *board
		mockTrello *trello.MockTrelloClient
		errMsg     string
	)
	BeforeEach(func() {
		testBoard = &board{
			Name:        "test_board",
			Description: "test_desc",
		}
	})
	Context("create method", func() {
		When("get board failed", func() {
			BeforeEach(func() {
				errMsg = "get board failed"
				mockTrello = &trello.MockTrelloClient{
					GetError: fmt.Errorf(errMsg),
				}
			})
			It("should return err", func() {
				err := testBoard.create(mockTrello)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(errMsg))
			})
		})
		When("board not exist", func() {
			BeforeEach(func() {
				mockTrello = &trello.MockTrelloClient{
					GetValue: &trelloCommon.Board{
						Name: "test",
					},
				}
			})
			It("should return nil", func() {
				err := testBoard.create(mockTrello)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		When("create failed", func() {
			BeforeEach(func() {
				errMsg = "create board failed"
				mockTrello = &trello.MockTrelloClient{
					CreateError: fmt.Errorf(errMsg),
				}
			})
			It("should return nil", func() {
				err := testBoard.create(mockTrello)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(errMsg))
			})
		})
	})
	Context("delete method", func() {
		When("board not exist", func() {
			BeforeEach(func() {
				mockTrello = &trello.MockTrelloClient{}
			})
			It("should return err", func() {
				err := testBoard.delete(mockTrello)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
	Context("get method", func() {
		When("get board failed", func() {
			BeforeEach(func() {
				errMsg = "get board failed"
				mockTrello = &trello.MockTrelloClient{
					GetError: fmt.Errorf(errMsg),
				}
			})
			It("should return err", func() {
				_, err := testBoard.get(mockTrello)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(errMsg))
			})
		})
	})
})
