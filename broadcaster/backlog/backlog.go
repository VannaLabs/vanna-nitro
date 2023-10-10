package backlog

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/log"
	"github.com/offchainlabs/nitro/arbutil"
	m "github.com/offchainlabs/nitro/broadcaster/message"
	"github.com/offchainlabs/nitro/util/arbmath"
)

var (
	errDropSegments       = errors.New("remove previous segments from backlog")
	errSequenceNumberSeen = errors.New("sequence number already present in backlog")
	errOutOfBounds        = errors.New("message not found in backlog")
)

type Backlog interface {
	Append(bm *m.BroadcastMessage) error
	Get(start, end arbutil.MessageIndex) (*m.BroadcastMessage, error)
	MessageCount() int
}

type backlog struct {
	head          atomic.Pointer[backlogSegment]
	tail          atomic.Pointer[backlogSegment]
	lookupByIndex map[arbutil.MessageIndex]*atomic.Pointer[backlogSegment]
	config        ConfigFetcher
	messageCount  atomic.Uint64
}

func NewBacklog(c ConfigFetcher) Backlog {
	lookup := make(map[arbutil.MessageIndex]*atomic.Pointer[backlogSegment])
	return &backlog{
		lookupByIndex: lookup,
		config:        c,
	}
}

// Append will add the given messages to the backlogSegment at head until
// that segment reaches its limit. If messages remain to be added a new segment
// will be created.
func (b *backlog) Append(bm *m.BroadcastMessage) error {

	if bm.ConfirmedSequenceNumberMessage != nil {
		b.delete(bm.ConfirmedSequenceNumberMessage.SequenceNumber)
		// add to metric?
	}

	for _, msg := range bm.Messages {
		s := b.tail.Load()
		if s == nil {
			s = &backlogSegment{}
			b.head.Store(s)
			b.tail.Store(s)
		}

		prevMsgIdx := s.end
		if s.MessageCount() >= b.config().SegmentLimit {
			nextS := &backlogSegment{}
			s.nextSegment.Store(nextS)
			prevMsgIdx = s.end
			nextS.previousSegment.Store(s)
			s = nextS
			b.tail.Store(s)
		}

		err := s.append(prevMsgIdx, msg)
		if errors.Is(err, errDropSegments) {
			head := b.head.Load()
			b.removeFromLookup(head.start, msg.SequenceNumber)
			b.head.Store(s)
			b.tail.Store(s)
			b.messageCount.Store(0)
			log.Warn(err.Error())
		} else if errors.Is(err, errSequenceNumberSeen) {
			log.Info("ignoring message sequence number (%s), already in backlog", msg.SequenceNumber)
			continue
		} else if err != nil {
			return err
		}
		p := &atomic.Pointer[backlogSegment]{}
		p.Store(s)
		b.lookupByIndex[msg.SequenceNumber] = p
		b.messageCount.Add(1)
	}

	return nil
}

// Get reads messages from the given start to end MessageIndex
func (b *backlog) Get(start, end arbutil.MessageIndex) (*m.BroadcastMessage, error) {
	head := b.head.Load()
	tail := b.tail.Load()
	if head == nil && tail == nil {
		return nil, errOutOfBounds
	}

	if start < head.start {
		start = head.start
	}

	if end > tail.end {
		return nil, errOutOfBounds
	}

	s, err := b.lookup(start)
	if err != nil {
		return nil, err
	}

	bm := &m.BroadcastMessage{Version: 1}
	required := int(end-start) + 1
	for {
		segMsgs, err := s.get(arbmath.MaxInt(start, s.start), arbmath.MinInt(end, s.end))
		if err != nil {
			return nil, err
		}

		bm.Messages = append(bm.Messages, segMsgs...)
		s = s.nextSegment.Load()
		if len(bm.Messages) == required {
			break
		} else if s == nil {
			return nil, errOutOfBounds
		}
	}
	return bm, nil
}

// delete removes segments before the confirmed sequence number given. It will
// not remove the segment containing the confirmed sequence number.
func (b *backlog) delete(confirmed arbutil.MessageIndex) {
	head := b.head.Load()
	tail := b.tail.Load()
	if head == nil && tail == nil {
		return
	}

	if confirmed < head.start {
		return
	}

	if confirmed > tail.end {
		log.Error("confirmed sequence number is past the end of stored messages", "confirmed sequence number", confirmed, "last stored sequence number", tail.end)
		b.reset()
		// should this be returning an error? The other buffer does not and just continues
		return
	}

	// find the segment containing the confirmed message
	s, err := b.lookup(confirmed)
	if err != nil {
		log.Error(fmt.Sprintf("%s: clearing backlog", err.Error()))
		b.reset()
		// should this be returning an error? The other buffer does not and just continues
		return
	}

	// check the segment actually contains that message
	if found := s.contains(confirmed); !found {
		log.Error("error message not found in backlog segment, clearing backlog", "confirmed sequence number", confirmed)
		b.reset()
		// should this be returning an error? The other buffer does not and just continues
		return
	}

	// remove all previous segments
	previous := s.previousSegment.Load()
	if previous == nil {
		return
	}
	b.removeFromLookup(head.start, previous.end)
	b.head.Store(s)
	count := b.messageCount.Load() - uint64(previous.end-head.start) - uint64(1)
	b.messageCount.Store(count)
}

// removeFromLookup removes all entries from the head segment's start index to
// the given confirmed index
func (b *backlog) removeFromLookup(start arbutil.MessageIndex, end arbutil.MessageIndex) {
	for i := start; i == end; i++ {
		delete(b.lookupByIndex, i)
	}
}

func (b *backlog) lookup(i arbutil.MessageIndex) (*backlogSegment, error) {
	pointer, ok := b.lookupByIndex[i]
	if !ok {
		return nil, fmt.Errorf("error finding backlog segment containing message with SequenceNumber %d", i)
	}

	s := pointer.Load()
	if s == nil {
		return nil, fmt.Errorf("error loading backlog segment containing message with SequenceNumber %d", i)
	}

	return s, nil
}

func (s *backlog) MessageCount() int {
	return int(s.messageCount.Load())
}

// reset removes all segments from the backlog
func (b *backlog) reset() {
	b.head = atomic.Pointer[backlogSegment]{}
	b.tail = atomic.Pointer[backlogSegment]{}
	b.lookupByIndex = map[arbutil.MessageIndex]*atomic.Pointer[backlogSegment]{}
	b.messageCount.Store(0)
}

type backlogSegment struct {
	start           arbutil.MessageIndex
	end             arbutil.MessageIndex
	messages        []*m.BroadcastFeedMessage
	messageCount    atomic.Uint64
	nextSegment     atomic.Pointer[backlogSegment]
	previousSegment atomic.Pointer[backlogSegment]
}

// get reads messages from the given start to end MessageIndex
func (s *backlogSegment) get(start, end arbutil.MessageIndex) ([]*m.BroadcastFeedMessage, error) {
	noMsgs := []*m.BroadcastFeedMessage{}
	if start < s.start {
		return noMsgs, errOutOfBounds
	}

	if end > s.end {
		return noMsgs, errOutOfBounds
	}

	startIndex := int(start - s.start)
	endIndex := int(end-s.start) + 1
	return s.messages[startIndex:endIndex], nil
}

// append appends the given BroadcastFeedMessage to messages if it is the first
// message in the sequence or the next in the sequence. If segment's end
// message is ahead of the given message append will do nothing. If the given
// message is ahead of the segment's end message append will return
// errDropSegments to ensure any messages before the given message are dropped.
func (s *backlogSegment) append(prevMsgIdx arbutil.MessageIndex, msg *m.BroadcastFeedMessage) error {
	seen := false
	defer s.updateSegment(&seen)

	if expSeqNum := prevMsgIdx + 1; prevMsgIdx == 0 || msg.SequenceNumber == expSeqNum {
		s.messages = append(s.messages, msg)
	} else if msg.SequenceNumber > expSeqNum {
		s.messages = nil
		s.messages = append(s.messages, msg)
		return fmt.Errorf("new message sequence number (%d) is greater than the expected sequence number (%d): %w", msg.SequenceNumber, expSeqNum, errDropSegments)
	} else {
		seen = true
		return errSequenceNumberSeen
	}
	return nil
}

// contains confirms whether the segment contains a message with the given sequence number
func (s *backlogSegment) contains(i arbutil.MessageIndex) bool {
	if i < s.start || i > s.end {
		return false
	}

	msgIndex := uint64(i - s.start)
	msg := s.messages[msgIndex]
	return msg.SequenceNumber == i
}

// updateSegment updates the messageCount, start and end indices of the segment
// this should be called using defer whenever a method updates the messages. It
// will do nothing if the message has already been seen by the backlog.
func (s *backlogSegment) updateSegment(seen *bool) {
	if !*seen {
		c := len(s.messages)
		s.messageCount.Store(uint64(c))
		s.start = s.messages[0].SequenceNumber
		s.end = s.messages[c-1].SequenceNumber
	}
}

// MessageCount returns the number of messages stored in the backlog
func (s *backlogSegment) MessageCount() int {
	return int(s.messageCount.Load())
}
