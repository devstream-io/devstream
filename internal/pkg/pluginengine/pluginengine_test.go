package pluginengine_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/backend/local"
	"github.com/devstream-io/devstream/internal/pkg/configloader"
	"github.com/devstream-io/devstream/internal/pkg/pluginengine"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

var _ = Describe("Pluginengine", func() {
	var (
		smgr statemanager.Manager
		err  error

		trelloInstance     = "mytrelloboard"
		trelloName         = "trello"
		expectedBoardId    = "1"
		expectedTodoListId = "2"
		trelloKey          = statemanager.StateKey(fmt.Sprintf("%s_%s", trelloName, trelloInstance))
	)

	BeforeEach(func() {
		defer GinkgoRecover()

		stateCfg := configloader.State{
			Backend: "local",
			Options: configloader.StateConfigOptions{
				StateFile: "devstream.state",
			},
		}
		smgr, err = statemanager.NewManager(stateCfg)
		Expect(err).ToNot(HaveOccurred())
		Expect(smgr).NotTo(BeNil())
		_, _ = GinkgoWriter.Write([]byte("new a statemanager"))

		DeferCleanup(func() {
			os.Remove(local.DefaultStateFile)
		})
	})

	It("should be 'one install'", func() {
		instanceID := "a"
		name := "tool-a"

		cfg := &configloader.Config{
			Tools: []configloader.Tool{*getTool(name, instanceID)},
		}
		changes, _ := pluginengine.GetChangesForApply(smgr, cfg)
		GinkgoWriter.Print(changes)
		Expect(len(changes)).To(Equal(1))
		c := changes[0]
		Expect(c.Tool.InstanceID).To(Equal(instanceID))
		Expect(c.ActionName).To(Equal(statemanager.ActionCreate))
	})

	It("should be 'two install'", func() {
		instanceID1, instanceID2 := "a", "b"
		name1, name2 := "tool-a", "too-b"

		cfg := &configloader.Config{
			Tools: []configloader.Tool{*getTool(name1, instanceID1), *getTool(name2, instanceID2)},
		}
		changes, _ := pluginengine.GetChangesForApply(smgr, cfg)

		Expect(len(changes)).To(Equal(2))
		c1 := changes[0]
		Expect(c1.Tool.InstanceID).To(Equal(instanceID1))
		Expect(c1.Tool.Name).To(Equal(name1))
		Expect(c1.ActionName).To(Equal(statemanager.ActionCreate))

		c2 := changes[1]
		Expect(c2.Tool.InstanceID).To(Equal(instanceID2))
		Expect(c2.Tool.Name).To(Equal(name2))
		Expect(c2.ActionName).To(Equal(statemanager.ActionCreate))
	})

	It("should be 1 uninstall when `dtm delete` is triggered against a config with 1 tool and a successful state", func() {
		instanceID := "a"
		name := "tool-a"

		cfg := &configloader.Config{
			Tools: []configloader.Tool{*getTool(name, instanceID)},
		}

		err = smgr.AddState(statemanager.StateKey(fmt.Sprintf("%s_%s", name, instanceID)), statemanager.State{})
		Expect(err).NotTo(HaveOccurred())
		changes, _ := pluginengine.GetChangesForDelete(smgr, cfg, false)

		Expect(len(changes)).To(Equal(1))
		c := changes[0]
		Expect(c.Tool.InstanceID).To(Equal(instanceID))
		Expect(c.Tool.Name).To(Equal(name))
		Expect(c.ActionName).To(Equal(statemanager.ActionDelete))
	})

	It("should handle outputs correctly", func() {
		trelloState := statemanager.State{
			InstanceID: "mytrelloboard",
			Name:       "trello",
			Options:    map[string]interface{}{},
			Resource: map[string]interface{}{
				"outputs": map[string]interface{}{
					"boardId":    expectedBoardId,
					"todoListId": expectedTodoListId,
				},
			},
		}
		err = smgr.AddState(trelloKey, trelloState)
		Expect(err).NotTo(HaveOccurred())

		dependantOptions := map[string]interface{}{
			"boardId":    fmt.Sprintf("${{ %s.%s.outputs.boardId }}", trelloName, trelloInstance),
			"todoListId": fmt.Sprintf("${{ %s.%s.outputs.todoListId }}", trelloName, trelloInstance),
		}
		expectResult := map[string]interface{}{
			"boardId":    expectedBoardId,
			"todoListId": expectedTodoListId,
		}
		errs := pluginengine.HandleOutputsReferences(smgr, dependantOptions)
		Expect(len(errs)).To(BeZero())
		Expect(dependantOptions).To(Equal(expectResult))
	})

	It("should handle output interpolation correctly", func() {
		trelloState := statemanager.State{
			InstanceID: "mytrelloboard",
			Name:       "trello",
			Options:    map[string]interface{}{},
			Resource: map[string]interface{}{
				"outputs": map[string]interface{}{
					"boardId": expectedBoardId,
				},
			},
		}
		err = smgr.AddState(trelloKey, trelloState)
		Expect(err).NotTo(HaveOccurred())

		dependantOptions := map[string]interface{}{
			"boardId": fmt.Sprintf("prefix/${{ %s.%s.outputs.boardId }}/suffix", trelloName, trelloInstance),
		}
		expectResult := map[string]interface{}{
			"boardId": fmt.Sprintf("prefix/%s/suffix", expectedBoardId),
		}
		errs := pluginengine.HandleOutputsReferences(smgr, dependantOptions)
		Expect(len(errs)).To(BeZero())
		Expect(dependantOptions).To(Equal(expectResult))
	})

	It("should give an error when output doesn't exist in the state", func() {
		trelloState := statemanager.State{
			Name:       "trello",
			InstanceID: "mytrelloboard",
			Options:    map[string]interface{}{},
			Resource:   map[string]interface{}{},
		}
		err = smgr.AddState(trelloKey, trelloState)
		Expect(err).NotTo(HaveOccurred())

		dependantOptions := map[string]interface{}{
			"boardId": fmt.Sprintf("${{ %s.%s.outputs.boardId }}", trelloName, trelloInstance),
		}
		expectResult := map[string]interface{}{
			"boardId": fmt.Sprintf("${{ %s.%s.outputs.boardId }}", trelloName, trelloInstance),
		}
		errs := pluginengine.HandleOutputsReferences(smgr, dependantOptions)
		Expect(len(errs)).To(Equal(1))
		Expect(dependantOptions).To(Equal(expectResult))
	})

	It("should give an error when the referred key doesn't exist", func() {
		trelloState := statemanager.State{
			Name:       "trello",
			InstanceID: "mytrelloboard",
			Options:    map[string]interface{}{},
			Resource: map[string]interface{}{
				"outputs": map[string]interface{}{
					"boardId":    expectedBoardId,
					"todoListId": expectedTodoListId,
				},
			},
		}
		err = smgr.AddState(trelloKey, trelloState)
		Expect(err).NotTo(HaveOccurred())

		dependantOptions := map[string]interface{}{
			"boardId":    fmt.Sprintf("${{ %s.%s.outputs.boardId }}", trelloName, trelloInstance),
			"todoListId": fmt.Sprintf("${{ %s.%s.outputs.todoListId }}", trelloName, trelloInstance),
			"someKey":    fmt.Sprintf("${{ %s.%s.outputs.keyNotExist }}", trelloName, trelloInstance),
		}
		expectResult := map[string]interface{}{
			"boardId":    expectedBoardId,
			"todoListId": expectedTodoListId,
			"someKey":    fmt.Sprintf("${{ %s.%s.outputs.keyNotExist }}", trelloName, trelloInstance),
		}
		errs := pluginengine.HandleOutputsReferences(smgr, dependantOptions)
		Expect(len(errs)).To(Equal(1))
		Expect(dependantOptions).To(Equal(expectResult))
	})

	It("should work for nested maps", func() {
		trelloState := statemanager.State{
			Name:       "trello",
			InstanceID: "mytrelloboard",
			Options:    map[string]interface{}{},
			Resource: map[string]interface{}{
				"outputs": map[string]interface{}{
					"boardId": expectedBoardId,
				},
			},
		}
		err = smgr.AddState(trelloKey, trelloState)
		Expect(err).NotTo(HaveOccurred())

		dependantOptions := map[string]interface{}{
			"outerKey": map[string]interface{}{
				"innerKey": fmt.Sprintf("${{ %s.%s.outputs.boardId }}", trelloName, trelloInstance),
			},
		}
		expectResult := map[string]interface{}{
			"outerKey": map[string]interface{}{
				"innerKey": expectedBoardId,
			},
		}
		errs := pluginengine.HandleOutputsReferences(smgr, dependantOptions)
		Expect(len(errs)).To(Equal(0))
		Expect(dependantOptions).To(Equal(expectResult))
	})
})

func getTool(name, instance string) *configloader.Tool {
	return &configloader.Tool{
		Name:       name,
		InstanceID: instance,
		Options:    map[string]interface{}{"key": "value"},
	}
}
