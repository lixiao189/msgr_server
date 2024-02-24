package test

import (
	"sync"
	"testing"
	"time"

	"github.com/slainsama/msgr_server/bot/handler"
	"github.com/slainsama/msgr_server/bot/types"
)

func TestNewConversationHandler(t *testing.T) {
	updateChan := make(chan types.TelegramUpdate, 5)
	var wg sync.WaitGroup
	const (
		start = iota
		sendHello
		end
	)

	dispatcher := handler.NewUpdateDispatcher()

	// Initialize conversation handler
	conversationHandler := handler.NewConversationHandler(
		handler.NewLocalPersistence(),
		handler.StateMap{
			start: {
				handler.NewCommandHandler("/startUpload", func(u *types.TelegramUpdate) int {
					t.Log("Start upload")
					return sendHello
				}),
			},
			sendHello: {
				handler.NewCommandHandler("/hello", func(u *types.TelegramUpdate) int {
					t.Log("Hello")
					return end
				}),
			},
			end: {
				handler.NewCommandHandler("/endUpload", func(u *types.TelegramUpdate) int {
					t.Log("End upload")
					return handler.HandleSuccess
				}),
			},
		},
		[]*handler.Handler{},

		time.Minute,
		nil,
		handler.NewKeyLock(),
	)
	dispatcher.AddHandler(conversationHandler)

	// Mock update message channel
	wg.Add(2)
	go func() {
		// Send start message
		updateChan <- newStartUpdate()
		time.Sleep(time.Second)
		// Send message
		updateChan <- newHelloUpdate()
		time.Sleep(time.Second)
		// Send end message
		updateChan <- newEndUpdate()
		time.Sleep(time.Second)

		wg.Done()
	}()

	go func() {
		for i := 0; i < 3; i++ {
			newUpdate := <-updateChan
			dispatcher.Dispatch(&newUpdate)
		}
		wg.Done()
	}()

	wg.Wait()
}

func TestNewConversationHandlerWithMultiChoice(t *testing.T) {
	updateChan := make(chan types.TelegramUpdate, 5)
	var wg sync.WaitGroup
	const (
		start = iota
		sendHello
		end
	)

	dispatcher := handler.NewUpdateDispatcher()

	// Initialize conversation handler
	conversationHandler := handler.NewConversationHandler(
		handler.NewLocalPersistence(),
		handler.StateMap{
			start: {
				handler.NewCommandHandler("/startUpload", func(u *types.TelegramUpdate) int {
					t.Log("Start upload")
					return sendHello
				}),
			},
			sendHello: {
				handler.NewCommandHandler("/hello", func(u *types.TelegramUpdate) int {
					t.Log("Hello")
					return end
				}),
				handler.NewCommandHandler("/another_hello", func(u *types.TelegramUpdate) int {
					t.Log("Another Hello")
					return end
				}),
			},
			end: {
				handler.NewCommandHandler("/endUpload", func(u *types.TelegramUpdate) int {
					t.Log("End upload")
					return handler.HandleSuccess
				}),
			},
		},
		[]*handler.Handler{},

		time.Minute,
		nil,
		handler.NewKeyLock(),
	)
	dispatcher.AddHandler(conversationHandler)

	// Mock update message channel
	wg.Add(2)
	go func() {
		// Send start message
		updateChan <- newStartUpdate()
		time.Sleep(time.Second)
		// Send message
		updateChan <- newAnotherHelloUpdate()
		time.Sleep(time.Second)
		// Send end message
		updateChan <- newEndUpdate()
		time.Sleep(time.Second)
		wg.Done()
	}()

	go func() {
		for i := 0; i < 3; i++ {
			newUpdate := <-updateChan
			dispatcher.Dispatch(&newUpdate)
		}
		wg.Done()
	}()

	wg.Wait()
}

func TestNewConversationHandlerWithoutPermission(t *testing.T) {
	updateChan := make(chan types.TelegramUpdate, 5)
	var wg sync.WaitGroup
	const (
		start = iota
		sendHello
		end
	)
	role := make(map[int64]bool)
	role[123] = true
	role[426] = false

	dispatcher := handler.NewUpdateDispatcher()

	// Initialize conversation handler
	conversationHandler := handler.NewConversationHandler(
		handler.NewLocalPersistence(),
		handler.StateMap{
			start: {
				handler.NewCommandHandler("/startUpload", func(u *types.TelegramUpdate) int {
					t.Log("Start upload")
					return sendHello
				}),
			},
			sendHello: {
				handler.NewCommandHandler("/hello", func(u *types.TelegramUpdate) int {
					userID := u.Message.From.ID
					if role[userID] {
						t.Log("Pass Auth")
						return end
					} else {
						t.Log("No permission")
						return handler.HandleFailed
					}
				}),
			},
			end: {
				handler.NewCommandHandler("/endUpload", func(u *types.TelegramUpdate) int {
					t.Log("End upload")
					return handler.HandleSuccess
				}),
			},
		},
		[]*handler.Handler{},

		time.Minute,
		nil,
		handler.NewKeyLock(),
	)
	dispatcher.AddHandler(conversationHandler)

	// Mock update message channel
	wg.Add(2)
	go func() {
		updateChan <- newStartUpdate()
		updateChan <- newHelloUpdate()
		updateChan <- newEndUpdate()
		wg.Done()
	}()

	go func() {
		for i := 0; i < 3; i++ {
			newUpdate := <-updateChan
			dispatcher.Dispatch(&newUpdate)
		}
		time.Sleep(time.Second * 2)
		wg.Done()
	}()

	wg.Wait()
}

// func TestConversationHandlerWithTimeout(t *testing.T) {
// 	updateChan := make(chan types.TelegramUpdate, 5)
// 	var wg sync.WaitGroup
// 	const (
// 		start = iota
// 		sendHello
// 		end
// 	)

// 	dispatcher := handler.NewUpdateDispatcher()

// 	// Initialize conversation handler
// 	conversationHandler := handler.NewConversationHandler(
// 		"testUpload",
// 		handler.NewCommandHandler("/startUpload", func(u *types.TelegramUpdate) {
// 			t.Log("Start upload")
// 			handler.UpdateState("testUpload", sendHello, u)
// 		}),
// 		handler.HandlerMap{
// 			sendHello: {
// 				handler.NewCommandHandler("/hello", func(u *types.TelegramUpdate) {
// 					t.Log("Hello")
// 					handler.UpdateState("testUpload", end, u)
// 				}),
// 			},
// 		},
// 		handler.NewCommandHandler("/endUpload", func(u *types.TelegramUpdate) {
// 			t.Log("End upload")
// 		}),
// 	)
// 	conversationHandler.SetConversationTimeout(time.Second)
// 	conversationHandler.SetTimeoutTask(func(u *types.TelegramUpdate) { t.Log("Timeout", u.Message.From.ID) })

// 	dispatcher.AddHandler(conversationHandler)

// 	// Mock update message channel
// 	wg.Add(2)
// 	go func() {
// 		// Send start message
// 		updateChan <- newStartUpdate()
// 		time.Sleep(time.Second * 2)
// 		// Send message
// 		updateChan <- newHelloUpdate()
// 		time.Sleep(time.Second)
// 		wg.Done()
// 	}()

// 	go func() {
// 		for i := 0; i < 2; i++ {
// 			newUpdate := <-updateChan
// 			dispatcher.Dispatch(&newUpdate)
// 		}
// 		wg.Done()
// 	}()

// 	wg.Wait()
// }
