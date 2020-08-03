package models

import (
	"errors"
	"fmt"
	"unicode"
)

const MaxMessageLen = 10485760

func validateCreateUserInput(username string) error {
	nameLen := len([]rune(username))
	if nameLen > 255 {
		return errors.New("username length can't be more than 255")
	}
	if nameLen == 0 {
		return errors.New("username can't be empty")
	}
	return nil
}

func validateCreateChatInput(chatName string, userIDs []string) error {
	nameLen := len([]rune(chatName))
	if nameLen > 255 {
		return errors.New("chat name length can't be more than 255")
	}
	if nameLen == 0 {
		return errors.New("chat name can't be empty")
	}

	users := make(map[string]bool, len(userIDs))

	for i := range userIDs {
		if len([]rune(userIDs[i])) != 32 {
			if err := validateUserID(userIDs[i]); err != nil {
				return err
			}
		}

		if users[userIDs[i]] {
			return errors.New("user ids must be unique")
		}

		users[userIDs[i]] = true
	}

	if len(userIDs) == 0 {
		return errors.New("chat must have contain at least one user")
	}
	return nil
}

func validateCreateMessageInput(chatID uint64, authorID string, text string) error {
	if err := validateChatID(chatID); err != nil {
		return err
	}
	if err := validateUserID(authorID); err != nil {
		return err
	}

	if len(text) == 0 {
		return errors.New("message text can't be empty")
	}

	if len([]rune(text)) > MaxMessageLen {
		return errors.New(fmt.Sprintf("messsage len can't exceed %d", MaxMessageLen))
	}
	return nil
}

func validateUserID(ID string) error {
	id := []rune(ID)
	if len(id) != 32 {
		return errors.New("user id len must be 32")
	}

	rts := []*unicode.RangeTable{unicode.Digit, unicode.Latin}
	for _, c := range id {
		if !unicode.IsOneOf(rts, c) {
			return errors.New("user id must contain only latin letters or digits")
		}
	}
	return nil
}

func validateChatID(ID uint64) error {
	if ID <= 0 {
		return errors.New("chat id must be positive integer")
	}
	return nil
}
