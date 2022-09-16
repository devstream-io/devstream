package jenkins

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/devstream-io/devstream/pkg/util/jenkins/dingtalk"
)

// ApplyDingTalkBot will override the setting of DingTalk bot list in Jenkins.
// The name of this plugin is "dingding-notifications".
func (j *jenkins) ApplyDingTalkBot(config dingtalk.BotConfig) error {
	// build config
	innerConfig, err := dingtalk.BuildDingTalkConfig(config)
	if err != nil {
		return err
	}

	// build request
	configJson, err := json.Marshal(innerConfig)
	if err != nil {
		return fmt.Errorf("marshal dingtalk config failed: %v", err)
	}
	params := map[string]string{
		"json": string(configJson),
	}

	// send request
	ctx := context.Background()
	res, err := j.Requester.Post(ctx, "/manage/dingtalk/configure", nil, nil, params)

	if err != nil {
		return fmt.Errorf("apply dingtalk bot failed: %v", err)
	}
	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("apply dingtalk bot failed: status code: %d, response body: %v", res.StatusCode, string(body))
	}
	return nil
}
