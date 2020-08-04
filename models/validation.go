package models

import (
	"errors"
	"fmt"
	"unicode"
)

const MaxMessageLen = 10485760

// Errors used by validation
var (
	// ErrUsernameTooLong is returned when username is longer than 255.
	ErrUsernameTooLong = errors.New("username length can't be more than 255")

	// ErrUsernameEmpty is return when username is empty string.
	ErrUsernameEmpty = errors.New("username can't be empty")

	// ErrChatNameTooLong is returned when chat name is longer than 255.
	ErrChatNameTooLong = errors.New("chat name length can't be more than 255")

	// ErrChatNameEmpty is returned when chat name is empty string.
	ErrChatNameEmpty = errors.New("chat name can't be empty")

	// ErrChatUserIDsNotUnique is returned when one of user IDs in chat occurs more than one time.
	ErrChatUserIDsNotUnique = errors.New("user ids must be unique")

	// ErrChatNoUsers is returned when attempting to create chat without users.
	ErrChatNoUsers = errors.New("chat must have contain at least one user")

	// ErrMessageTooLong is returned when message is longer than MaxMessageLen const.
	ErrMessageTooLong = errors.New(fmt.Sprintf("messsage len can't exceed %d", MaxMessageLen))

	// ErrUserIDIncorrectLen is returned when user ID length is not equal to 32.
	ErrUserIDIncorrectLen = errors.New("user id len must be 32")

	// ErrUserIDIncorrectRune is returned when user ID contains nor latin nor digit rune.
	ErrUserIDIncorrectRune = errors.New("user id must contain only latin letters or digits")

	// ErrChatIDNotPositive is returned when ChatID is not positive integer.
	ErrChatIDNotPositive = errors.New("chat id must be positive integer")
)

func validateCreateUserInput(username string) error {
	nameLen := len([]rune(username))
	if nameLen > 255 {
		return ErrUsernameTooLong
	}
	if nameLen == 0 {
		return ErrUsernameEmpty
	}
	return nil
}

func validateCreateChatInput(chatName string, userIDs []string) error {
	nameLen := len([]rune(chatName))
	if nameLen > 255 {
		return ErrChatNameTooLong
	}
	if nameLen == 0 {
		return ErrChatNameEmpty
	}

	// Map is here to find repetitive users. By default value is false, if occur we set it to true.
	// If we see users[userIDs[i]] is true then user is repetitive.
	users := make(map[string]bool, len(userIDs))

	for i := range userIDs {
		if len([]rune(userIDs[i])) != 32 {
			if err := validateUserID(userIDs[i]); err != nil {
				return err
			}
		}

		if users[userIDs[i]] {
			return ErrChatUserIDsNotUnique
		}

		users[userIDs[i]] = true
	}

	if len(userIDs) == 0 {
		return ErrChatNoUsers
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
		return ErrMessageTooLong
	}
	return nil
}

func validateUserID(ID string) error {
	id := []rune(ID)
	if len(id) != 32 {
		return ErrUserIDIncorrectLen
	}

	rts := []*unicode.RangeTable{unicode.Digit, unicode.Latin}
	for _, c := range id {
		if !unicode.IsOneOf(rts, c) {
			return ErrUserIDIncorrectRune
		}
	}
	return nil
}

func validateChatID(ID uint64) error {
	if ID <= 0 {
		return ErrChatIDNotPositive
	}
	return nil
}
