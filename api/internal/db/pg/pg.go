// package pg

// import (
// 	"sync"

// 	"github.com/radean0909/guild-chat/api/internal/db"
// 	"github.com/radean0909/guild-chat/api/models"
// )

// type (
// 	Key struct {
// 		Sender, Recipient string
// 	}

// 	Driver struct {
// 		mux    sync.RWMutex
// 		msgs   map[string]*models.Message   // primary key is linked to a single id
// 		convos map[Key]*models.Conversation // complex primary key
// 		users  map[string]*models.User      //primary key is a single id
// 	}
// )

// var (
// 	_ db.Driver = new(Driver)
// )

// func (d *Driver) GetMessage(id string) (*models.Message, error) {

// 	sqlQuery := "SELECT * FROM messages WHERE ID = $1"

// 	return nil, nil
// }

// func (d *Driver) CreateMessage(msg *models.Message) (*models.Message, error) {
// 	return nil, nil
// }

// func (d *Driver) GetConversation(sender, recipient string) (*models.Conversation, error) {
// 	return nil, nil
// }

// func (d *Driver) CreateConversation(sender, recipient string) (*models.Conversation, error) {
// 	return nil, nil
// }

// func (d *Driver) ListConversations(recipient string) ([]*models.Conversation, error) {
// 	return nil, nil
// }

// func (d *Driver) GetUser(id string) (*models.User, error) {
// 	return nil, nil
// }

// func (d *Driver) CreateUser(user *models.User) (*models.User, error) {
// 	return nil, nil
// }

// func (d *Driver) DeleteUser(id string) error {
// 	return nil
// }
