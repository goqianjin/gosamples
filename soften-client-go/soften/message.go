package soften

import (
	"fmt"
	"reflect"

	"github.com/shenqianjin/soften-client-go/soften/checker"
	"github.com/shenqianjin/soften-client-go/soften/message"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/shenqianjin/soften-client-go/soften/internal"
	"github.com/shenqianjin/soften-client-go/soften/internal/strategy"
)

// ------ re-reRouter message ------

type RerouteMessage struct {
	producerMsg pulsar.ProducerMessage
	consumerMsg pulsar.ConsumerMessage
}

// ------ custom consumer message ------

type ConsumerMessage struct {
	pulsar.ConsumerMessage
	StatusMessage
	LeveledMessage
}

// ------ status message interface ------

type StatusMessage interface {
	Status() internal.MessageStatus
}

// ------ leveled message interface ------

type LeveledMessage interface {
	Level() internal.TopicLevel
}

// ---------------------------------------

// ------ status message implementation ------

type statusMessage struct {
	status internal.MessageStatus
}

func (m *statusMessage) Status() internal.MessageStatus {
	return m.status
}

// ------ leveled message implementation ------

type leveledMessage struct {
	level internal.TopicLevel
}

func (m *leveledMessage) Level() internal.TopicLevel {
	return m.level
}

// ------ message receiver help ------

var messageChSelector = &messageChSelectorImpl{}

type messageChSelectorImpl struct {
}

func (mcs *messageChSelectorImpl) receiveAny(chs []<-chan ConsumerMessage) (ConsumerMessage, bool) {
	/*select {
	case msg, ok := <-chs[0]:
		return msg, ok
	case msg, ok := <-chs[1]:
		return msg, ok
	case msg, ok := <-chs[2]:
		return msg, ok
	case msg, ok := <-chs[3]:
		return msg, ok
	}*/
	cases := make([]reflect.SelectCase, len(chs))
	for i, ch := range chs {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	_, value, ok := reflect.Select(cases)
	// ok will be true if the channel has not been closed.
	if rv, valid := value.Interface().(ConsumerMessage); !valid {
		panic(fmt.Sprintf("convert %v to ConsumerMessage failed", value))
	} else {
		return rv, ok
	}
}

func (mcs *messageChSelectorImpl) receiveOneByWeight(chs []<-chan ConsumerMessage, balanceStrategy strategy.IBalanceStrategy, excludedIndexes *[]int) (ConsumerMessage, bool) {
	if len(*excludedIndexes) >= len(chs) {
		excludedIndexes = &[]int{}
		return mcs.receiveAny(chs)
	}
	index := balanceStrategy.Next(*excludedIndexes...)
	select {
	case msg, ok := <-chs[index]:
		return msg, ok
	default:
		*excludedIndexes = append(*excludedIndexes, index)
		return mcs.receiveOneByWeight(chs, balanceStrategy, excludedIndexes)
	}
}

// ------ helper ------

var internalGotoReroute = internal.MessageGoto("Reroute") // for consumer
var internalGotoRoute = internal.MessageGoto("Route")     // for producer

var checkTypeGotoMap = map[internal.CheckType]internal.MessageGoto{
	checker.CheckTypePreDiscard:  message.GotoDiscard,
	checker.CheckTypePreDead:     message.GotoDead,
	checker.CheckTypePreUpgrade:  message.GotoUpgrade,
	checker.CheckTypePreDegrade:  message.GotoDegrade,
	checker.CheckTypePreBlocking: message.GotoBlocking,
	checker.CheckTypePrePending:  message.GotoPending,
	checker.CheckTypePreRetrying: message.GotoRetrying,
	checker.CheckTypePreReroute:  internalGotoReroute,

	checker.CheckTypePostDiscard:  message.GotoDiscard,
	checker.CheckTypePostDead:     message.GotoDead,
	checker.CheckTypePostUpgrade:  message.GotoUpgrade,
	checker.CheckTypePostDegrade:  message.GotoDegrade,
	checker.CheckTypePostBlocking: message.GotoBlocking,
	checker.CheckTypePostPending:  message.GotoPending,
	checker.CheckTypePostRetrying: message.GotoRetrying,
	checker.CheckTypePostReroute:  internalGotoReroute,

	checker.ProduceCheckTypeDiscard:  message.GotoDiscard,
	checker.ProduceCheckTypeDead:     message.GotoDead,
	checker.ProduceCheckTypeUpgrade:  message.GotoUpgrade,
	checker.ProduceCheckTypeDegrade:  message.GotoDegrade,
	checker.ProduceCheckTypeBlocking: message.GotoBlocking,
	checker.ProduceCheckTypePending:  message.GotoPending,
	checker.ProduceCheckTypeRetrying: message.GotoRetrying,
	checker.ProduceCheckTypeRoute:    internalGotoRoute,
}
