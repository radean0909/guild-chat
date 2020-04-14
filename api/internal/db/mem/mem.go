package mem

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/radean0909/guild-chat/api/internal/db"
	"github.com/radean0909/guild-chat/api/internal/errors"
	"github.com/radean0909/guild-chat/api/models"
)

type (
	Key struct {
		Sender, Recipient string
	}

	Driver struct {
		mux    sync.RWMutex
		msgs   map[string]*models.Message   // primary key is linked to a single id
		convos map[Key]*models.Conversation // complex primary key
		users  map[string]*models.User      //primary key is a single id
	}
)

var (
	_ db.Driver = new(Driver)
)

// NewDriver - creates a in-memory database driver
func NewDriver() *Driver {
	return &Driver{
		mux:    sync.RWMutex{},
		msgs:   map[string]*models.Message{},
		convos: map[Key]*models.Conversation{},
		users:  map[string]*models.User{},
	}
}

// GetMessage - gets a single message by id
func (d *Driver) GetMessage(id string) (*models.Message, error) {

	if id == "" {
		return nil, errors.ErrBadRequest
	}

	d.mux.RLock()
	defer d.mux.RUnlock()

	msg, ok := d.msgs[id]
	if !ok {
		return nil, errors.ErrNotFound
	}

	return msg, nil
}

// CreateMessage - creates a new message
func (d *Driver) CreateMessage(msg *models.Message) (*models.Message, error) {
	if msg.Recipient == "" || msg.Sender == "" {
		return nil, errors.ErrBadRequest
	}

	d.mux.RLock()
	defer d.mux.RUnlock()

	// if the user or sender id are invalid, throw error - in this case a not found, so as not to tip off malicious attacks of a bad user id
	if _, ok := d.users[msg.Sender]; !ok {
		return nil, errors.ErrNotFound
	}

	if _, ok := d.users[msg.Recipient]; !ok {
		return nil, errors.ErrNotFound
	}

	msg.Date = time.Now()
	msg.ID = uuid.New().String()
	d.msgs[msg.ID] = msg

	// now add to the existing conversation, or create a new one
	if convo, ok := d.convos[Key{msg.Sender, msg.Recipient}]; ok {
		convo.Messages = append(convo.Messages, msg)
		convo.Updated = time.Now()
	} else {
		convo, err := d.CreateConversation(msg.Sender, msg.Recipient)
		if err != nil {
			return nil, err
		}
		convo.Messages = append(convo.Messages, msg)
	}

	return msg, nil
}

// ListMessages - lists messages all messages for a recipient
// from and until times can be passed to further narrow results to conversations that have been updated in the timeframe
// if a 0 time is passed for either of these values, that filtering parameter is ignored
// if limit is 0, it is ignored, otherwise only the most recent messages, up to a count of limit, are returned
func (d *Driver) ListMessages(recipient string, from, until time.Time, limit int) ([]*models.Message, error) {
	if recipient == "" {
		return nil, errors.ErrBadRequest
	}

	if from.After(until) {
		return nil, errors.ErrBadRequest
	}

	d.mux.RLock()
	defer d.mux.RUnlock()

	if _, ok := d.users[recipient]; !ok {
		return nil, errors.ErrNotFound
	}

	msgs := []*models.Message{}
	count := 0
	// only return conversations that match the filters
	// don't need to do anything to 0 time from param to "ignore" it

	// if until is 0 time, set to now
	if until.Equal(time.Time{}) {
		until = time.Now()
	}

	for _, msg := range d.msgs {
		if limit > 0 && count > limit {
			break
		}

		if msg.Recipient == recipient &&
			(msg.Date.After(from) || msg.Date.Equal(from)) &&
			(msg.Date.Before(until) || msg.Date.Equal(until)) {
			msgs = append(msgs, msg)
			count++
		}
	}

	return msgs, nil
}

// GetConversation - gets a conversation between a sender and recipient.
// from and until times can be passed to further narrow results to conversations that have been updated in the timeframe
// if a 0 time is passed for either of these values, that filtering parameter is ignored
func (d *Driver) GetConversation(sender, recipient string, from, until time.Time) (*models.Conversation, error) {
	if sender == "" || recipient == "" {
		return nil, errors.ErrBadRequest
	}

	if from.After(until) {
		return nil, errors.ErrBadRequest
	}

	d.mux.RLock()
	defer d.mux.RUnlock()

	// if the user or sender id are invalid, throw error - in this case a not found, so as not to tip off malicious attacks of a bad user id
	if _, ok := d.users[sender]; !ok {
		return nil, errors.ErrNotFound
	}

	if _, ok := d.users[recipient]; !ok {
		return nil, errors.ErrNotFound
	}

	convo, ok := d.convos[Key{sender, recipient}]
	if !ok {
		return nil, errors.ErrNotFound
	}

	// only return conversations that match the filters
	// don't need to do anything to 0 time from param to "ignore" it

	// if until is 0 time, set to now
	if until.Equal(time.Time{}) {
		until = time.Now()
	}

	if (convo.Updated.After(from) || convo.Updated.Equal(from)) &&
		(convo.Updated.Before(until) || convo.Updated.Equal(until)) {
		return convo, nil
	}

	return nil, errors.ErrNotFound
}

// CreateConversation - creates a conversation between a sender and recipient
func (d *Driver) CreateConversation(sender, recipient string) (*models.Conversation, error) {
	if sender == "" || recipient == "" {
		return nil, errors.ErrBadRequest
	}

	d.mux.RLock()
	defer d.mux.RUnlock()

	// if the user or sender id are invalid, throw error - in this case a not found, so as not to tip off malicious attacks of a bad user id
	if _, ok := d.users[sender]; !ok {
		return nil, errors.ErrNotFound
	}

	if _, ok := d.users[recipient]; !ok {
		return nil, errors.ErrNotFound
	}

	// If a conversation already exists, throw error
	if _, ok := d.convos[Key{sender, recipient}]; ok {
		return nil, errors.ErrBadRequest
	}

	if _, ok := d.convos[Key{recipient, sender}]; ok {
		return nil, errors.ErrBadRequest
	}

	convo := &models.Conversation{
		ID:        uuid.New().String(),
		Sender:    sender,
		Recipient: recipient,
		Updated:   time.Now(),
		Messages:  []*models.Message{},
	}

	// the same conversation should exist, no matter which way we are looking at it
	d.convos[Key{sender, recipient}] = convo
	d.convos[Key{recipient, sender}] = convo

	return convo, nil
}

// ListConversations - lists all conversations between recipient and others.
// from and until times can be passed to further narrow results to conversations that have been updated in the timeframe
// if a 0 time is passed for either of these values, that filtering parameter is ignored
func (d *Driver) ListConversations(recipient string, from, until time.Time) ([]*models.Conversation, error) {
	if recipient == "" {
		return nil, errors.ErrBadRequest
	}

	if from.After(until) {
		return nil, errors.ErrBadRequest
	}

	d.mux.RLock()
	defer d.mux.RUnlock()

	// if the receipient id is invalid, throw error
	if _, ok := d.users[recipient]; !ok {
		return nil, errors.ErrNotFound
	}

	conversations := []*models.Conversation{}

	for key, convo := range d.convos {

		// only return conversations that match the filters
		// don't need to do anything to 0 time from param to "ignore" it

		// if until is 0 time, set to now
		if until.Equal(time.Time{}) {
			until = time.Now()
		}

		if key.Recipient == recipient &&
			(convo.Updated.After(from) || convo.Updated.Equal(from)) &&
			(convo.Updated.Before(until) || convo.Updated.Equal(until)) {
			conversations = append(conversations, convo)
		}
	}

	return conversations, nil
}

// GetUser - gets a single user by id
func (d *Driver) GetUser(id string) (*models.User, error) {

	if id == "" {
		return nil, errors.ErrBadRequest
	}

	d.mux.RLock()
	defer d.mux.RUnlock()

	user, ok := d.users[id]
	if !ok {
		return nil, errors.ErrNotFound
	}

	// if the user has been soft deleted, don't return it
	// note: golang doesn't have null time, so look for 0 time
	if user.ArchivedOn.Before(time.Now()) && !user.ArchivedOn.Equal(time.Time{}) {
		return user, nil
	}

	return nil, errors.ErrNotFound
}

// CreateUser - creates a new user
func (d *Driver) CreateUser(user *models.User) (*models.User, error) {

	if user == nil || user.Email == "" || user.Username == "" {
		return nil, errors.ErrBadRequest
	}

	d.mux.RLock()
	defer d.mux.RUnlock()

	user.ID = uuid.New().String()
	d.users[user.ID] = user
	return user, nil
}

// DeleteUser - soft deletes a user
func (d *Driver) DeleteUser(id string) error {
	if id == "" {
		return errors.ErrBadRequest
	}

	d.mux.RLock()
	defer d.mux.RUnlock()

	user, ok := d.users[id]
	if !ok {
		return errors.ErrNotFound
	}

	user.ArchivedOn = time.Now()

	return nil
}
