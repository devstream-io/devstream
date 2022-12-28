package statemanager

import (
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/mapz/concurrentmap"
)

var _ = Describe("Statemanager.state", func() {
	Describe("GenerateStateKeyByToolNameAndInstanceID func", func() {
		It("should ouput state key base on toolName and plugin kind", func() {
			var testCases = []struct {
				toolName       string
				plugKind       string
				expectStateKey StateKey
			}{
				{"test_tool", "test_kind", "test_tool_test_kind"},
				{"123", "1", "123_1"},
			}
			for _, t := range testCases {
				funcResult := GenerateStateKeyByToolNameAndInstanceID(t.toolName, t.plugKind)
				Expect(funcResult).Should(Equal(t.expectStateKey))
			}
		})
	})

	Describe("NewStatesMap func", func() {
		It("should return normal statesMap", func() {
			newStatesMap := NewStatesMap()
			Expect(newStatesMap).ShouldNot(BeNil())
		})
	})

	Describe("StatesMap struct", func() {
		var (
			testMap      StatesMap
			testStateKey StateKey
			testStateVal State
		)

		BeforeEach(func() {
			testStateKey = StateKey("test_key")
			testStateVal = State{
				Name:       "test_tool",
				InstanceID: "test_instance",
			}
			testMap = StatesMap{
				ConcurrentMap: concurrentmap.NewConcurrentMap(
					StateKey(""), State{},
				),
			}
		})

		Context("DeepCopy method", func() {
			It("should return not same map", func() {
				newTestMap := testMap.DeepCopy()
				oldMapAddress := reflect.ValueOf(&testMap)
				newMapAddress := reflect.ValueOf(&newTestMap)
				Expect(oldMapAddress).ShouldNot(Equal(newMapAddress))
				valEqual := reflect.DeepEqual(newTestMap, testMap)
				Expect(valEqual).Should(BeTrue())
			})
		})

		Describe("ToList method", func() {
			It("should return not empty list if have key", func() {
				testMap.Store(testStateKey, testStateVal)
				tList := testMap.ToList()
				Expect(len(tList)).Should(Equal(1))
				Expect(tList[0].Name).Should(Equal(testStateVal.Name))
				Expect(tList[0].InstanceID).Should(Equal(testStateVal.InstanceID))
			})
		})

		Describe("Format method", func() {
			It("should return formated info for map", func() {
				testMap.Store(testStateKey, testStateVal)
				formatedInfo := testMap.Format()
				formatResult := fmt.Sprintf(
					"test_key:\n  name: %s\n  instanceID: %s\n  dependsOn: []\n  options: {}\n  resourceStatus: {}\n",
					testStateVal.Name, testStateVal.InstanceID,
				)
				Expect(string(formatedInfo)).Should(Equal(formatResult))
			})
		})
	})
})
