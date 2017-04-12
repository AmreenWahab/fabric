/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package msgstore

import (
	"sync"

	"time"

	"github.com/hyperledger/fabric/gossip/common"
)

var noopLock = func() {}

// invalidationTrigger is invoked on each message that was invalidated because of a message addition
// i.e: if add(0), add(1) was called one after the other, and the store has only {1} after the sequence of invocations
// then the invalidation trigger on 0 was called when 1 was added.
type invalidationTrigger func(message interface{})

// NewMessageStore returns a new MessageStore with the message replacing
// policy and invalidation trigger passed.
func NewMessageStore(pol common.MessageReplacingPolicy, trigger invalidationTrigger) MessageStore {
	return newMsgStore(pol, trigger)
}

// NewMessageStoreExpirable returns a new MessageStore with the message replacing
// policy and invalidation trigger passed. It supports old message expiration after msgTTL, during expiration first external
// lock taken, expiration callback invoked and external lock released. Callback and external lock can be nil.
func NewMessageStoreExpirable(pol common.MessageReplacingPolicy, trigger invalidationTrigger, msgTTL time.Duration, externalLock func(), externalUnlock func(), externalExpire func(interface{})) MessageStore {
	store := newMsgStore(pol, trigger)

	store.expirable = true
	store.msgTTL = msgTTL

	if externalLock != nil {
		store.externalLock = externalLock
	}

	if externalUnlock != nil {
		store.externalUnlock = externalUnlock
	}

	if externalExpire != nil {
		store.expireMsgCallback = externalExpire
	}

	go store.expirationRoutine()
	return store
}

func newMsgStore(pol common.MessageReplacingPolicy, trigger invalidationTrigger) *messageStoreImpl {
	return &messageStoreImpl{
		pol:        pol,
		messages:   make([]*msg, 0),
		invTrigger: trigger,

		expirable:         false,
		externalLock:      noopLock,
		externalUnlock:    noopLock,
		expireMsgCallback: func(m interface{}) {},
		expiredCount:      0,

		doneCh: make(chan struct{}),
	}

}

// MessageStore adds messages to an internal buffer.
// When a message is received, it might:
// 	- Be added to the buffer
// 	- Discarded because of some message already in the buffer (invalidated)
// 	- Make a message already in the buffer to be discarded (invalidates)
// When a message is invalidated, the invalidationTrigger is invoked on that message.
type MessageStore interface {
	// add adds a message to the store
	// returns true or false whether the message was added to the store
	Add(msg interface{}) bool

	// Checks if message is valid for insertion to store
	// returns true or false whether the message can be added to the store
	CheckValid(msg interface{}) bool

	// size returns the amount of messages in the store
	Size() int

	// get returns all messages in the store
	Get() []interface{}

	// Stop all associated go routines
	Stop()
}

type messageStoreImpl struct {
	pol        common.MessageReplacingPolicy
	lock       sync.RWMutex
	messages   []*msg
	invTrigger invalidationTrigger

	expirable    bool
	msgTTL       time.Duration
	expiredCount int

	externalLock      func()
	externalUnlock    func()
	expireMsgCallback func(msg interface{})

	doneCh   chan struct{}
	stopOnce sync.Once
}

type msg struct {
	data    interface{}
	created time.Time
	expired bool
}

// add adds a message to the store
func (s *messageStoreImpl) Add(message interface{}) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	n := len(s.messages)
	for i := 0; i < n; i++ {
		m := s.messages[i]
		switch s.pol(message, m.data) {
		case common.MessageInvalidated:
			return false
		case common.MessageInvalidates:
			s.invTrigger(m.data)
			s.messages = append(s.messages[:i], s.messages[i+1:]...)
			n--
			i--
		}
	}

	s.messages = append(s.messages, &msg{data: message, created: time.Now()})
	return true
}

// Checks if message is valid for insertion to store
func (s *messageStoreImpl) CheckValid(message interface{}) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for _, m := range s.messages {
		if s.pol(message, m.data) == common.MessageInvalidated {
			return false
		}
	}
	return true
}

// size returns the amount of messages in the store
func (s *messageStoreImpl) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.messages) - s.expiredCount
}

// get returns all messages in the store
func (s *messageStoreImpl) Get() []interface{} {
	res := make([]interface{}, 0)

	s.lock.RLock()
	defer s.lock.RUnlock()

	for _, msg := range s.messages {
		if !msg.expired {
			res = append(res, msg.data)
		}
	}
	return res
}

func (s *messageStoreImpl) expireMessages() {
	s.externalLock()
	s.lock.Lock()
	defer s.lock.Unlock()
	defer s.externalUnlock()

	n := len(s.messages)
	for i := 0; i < n; i++ {
		m := s.messages[i]
		if !m.expired {
			if time.Since(m.created) > s.msgTTL {
				m.expired = true
				s.expireMsgCallback(m.data)
				s.expiredCount++
			}
		} else {
			if time.Since(m.created) > (s.msgTTL * 2) {
				s.messages = append(s.messages[:i], s.messages[i+1:]...)
				n--
				i--
				s.expiredCount--
			}

		}
	}
}

func (s *messageStoreImpl) needToExpire() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, msg := range s.messages {
		if !msg.expired && time.Since(msg.created) > s.msgTTL {
			return true
		} else if time.Since(msg.created) > (s.msgTTL * 2) {
			return true
		}
	}
	return false
}

func (s *messageStoreImpl) expirationRoutine() {
	for {
		select {
		case <-s.doneCh:
			return
		case <-time.After(s.expirationCheckInterval()):
			if s.needToExpire() {
				s.expireMessages()
			}
		}
	}
}

func (s *messageStoreImpl) Stop() {
	stopFunc := func() {
		close(s.doneCh)
	}
	s.stopOnce.Do(stopFunc)
}

func (s *messageStoreImpl) expirationCheckInterval() time.Duration {
	return s.msgTTL / 100
}