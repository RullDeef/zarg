// text - Textual inmemory chat implementation
package text

// type client struct {
// 	profile *domain.Profile
// 	channel chan domain.ChatMessage
// }

// type TextChat struct {
// 	clients map[domain.ProfileID]client
// 	mutex   sync.RWMutex
// }

// var _ domain.Chat = &TextChat{}

// func New() *TextChat {
// 	return &TextChat{
// 		clients: make(map[domain.ProfileID]client),
// 	}
// }

// func (c *TextChat) Connect(user *domain.Profile) error {
// 	c.mutex.Lock()
// 	defer c.mutex.Unlock()
// 	if _, ok := c.clients[user.ID]; !ok {
// 		c.clients[user.ID] = client{
// 			profile: user,
// 			channel: make(chan domain.ChatMessage),
// 		}
// 	}
// 	return nil
// }

// func (c *TextChat) Disconnect(user *domain.Profile) error {
// 	c.mutex.Lock()
// 	defer c.mutex.Unlock()

// 	if _, ok := c.clients[user.ID]; ok {
// 		delete(c.clients, user.ID)
// 	}
// 	return nil
// }

// func (c *TextChat) Send(ctx context.Context, msg domain.ChatMessage) error {
// 	panic("TODO")
// }

// func (c *TextChat) ReceiverChan(ctx context.Context) <-chan domain.ChatMessage {
// 	panic("TODO")
// }
