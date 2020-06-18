package settings

import (
	"errors"
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
)

type readOnlySetting struct {
	baseSetting
	store SettingStore
}

func NewReadOnlySetting(id string, title string, description string, dependsOn string, store SettingStore) Setting {
	return &readOnlySetting{
		baseSetting: baseSetting{
			title:       title,
			description: description,
			id:          id,
			dependsOn:   dependsOn,
		},
		store: store,
	}
}

func (s *readOnlySetting) Get(userID string) (interface{}, error) {
	value, err := s.store.GetSetting(userID, s.id)
	if err != nil {
		return "", err
	}
	stringValue, ok := value.(string)
	if !ok {
		return "", errors.New("current value is not a string")
	}

	return stringValue, nil
}

func (s *readOnlySetting) Set(userID string, value interface{}) error {
	return nil
}

func (s *readOnlySetting) GetSlackAttachments(userID, settingHandler string, disabled bool) (*model.SlackAttachment, error) {
	title := fmt.Sprintf("Setting: %s", s.title)
	currentValueMessage := "Disabled"

	if !disabled {
		currentValue, err := s.Get(userID)
		if err != nil {
			return nil, err
		}
		currentValueMessage = fmt.Sprintf("Current value: %s", currentValue)
	}

	text := fmt.Sprintf("%s\n%s", s.description, currentValueMessage)
	sa := model.SlackAttachment{
		Title: title,
		Text:  text,
	}

	return &sa, nil
}

func (s *readOnlySetting) IsDisabled(foreignValue interface{}) bool {
	return foreignValue == "false"
}
