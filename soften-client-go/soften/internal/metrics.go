package internal

import "github.com/prometheus/client_golang/prometheus"

// message consume stages: published -> received -> listened -> before_checked -> handled -> after_checked -> decided

type MetricsProvider struct {
	metricsLevel int

	// client labels: {url=}
	clientsOpened *prometheus.GaugeVec // gauge of opened clients

	// producer labels: {topic=}
	producersOpened        *prometheus.GaugeVec     // gauge of producers
	producerPublishLatency *prometheus.HistogramVec // producers publish latency
	producerDecideLatency  *prometheus.HistogramVec // route latency

	// producer check typed label: {topic=, check_type=}
	producerTypedCheckerOpened *prometheus.GaugeVec   // gauge of produce checkers
	producerTypedCheckPassed   *prometheus.CounterVec // counter of produce check passed
	producerTypedCheckRejected *prometheus.CounterVec // counter of produce check rejected
	producerTypedCheckLatency  *prometheus.HistogramVec

	// producer decider labels: {topic=, goto=}
	producerDeciderGotoOpened *prometheus.GaugeVec     // gauge of routers
	producerDecideGotoSuccess *prometheus.CounterVec   // counter of producers decide messages
	producerDecideGotoFailed  *prometheus.CounterVec   // counter of producers decide messages
	producerDecideGotoLatency *prometheus.HistogramVec // route latency

	// producer decider (router) labels: {topic=, goto=, route_topic}
	producerRouteDeciderOpened *prometheus.GaugeVec     // gauge of routers
	producerRouteDecideSuccess *prometheus.CounterVec   // counter of producers decide messages
	producerRouteDecideFailed  *prometheus.CounterVec   // counter of producers decide messages
	produceRouteDecideLatency  *prometheus.HistogramVec // route latency

	// listener labels: {topics, levels}
	listenersOpened        *prometheus.GaugeVec     // gauge of opened listeners
	listenersRunning       *prometheus.GaugeVec     // gauge of running listeners
	listenerReceiveLatency *prometheus.HistogramVec // labels: topic, level, status
	listenerListenLatency  *prometheus.HistogramVec // labels: topic, level, status
	listenerCheckLatency   *prometheus.HistogramVec
	listenerHandleLatency  *prometheus.HistogramVec
	listenerDecideLatency  *prometheus.HistogramVec

	// listen typed checker labels: {topics, levels, check_type}
	listenerTypedCheckersEnabled *prometheus.GaugeVec   // gauge of checkers
	listenerTypedCheckPassed     *prometheus.CounterVec // counter of produce check passed
	listenerTypedCheckRejected   *prometheus.CounterVec // counter of produce check rejected
	listenerTypedCheckLatency    *prometheus.HistogramVec

	// listener handle labels: {topics=, levels=, goto=}
	listenerHandleGotoLatency *prometheus.HistogramVec
	listenerDecideGotoLatency *prometheus.HistogramVec

	// consumer labels: {topics, levels, topic, level, status}
	consumersOpened        *prometheus.GaugeVec     // gauge of consumers
	consumerReceiveLatency *prometheus.HistogramVec // labels: topic, level, status
	consumerListenLatency  *prometheus.HistogramVec // labels: topic, level, status
	consumerCheckLatency   *prometheus.HistogramVec // labels: topic, level, status
	consumerHandleLatency  *prometheus.HistogramVec // labels: topic, level, status
	consumerDecideLatency  *prometheus.HistogramVec // labels: topic, level, status
	consumerConsumeLatency *prometheus.HistogramVec // including receive, listen, check, handle and decide

	// consumer handle labels: {topics, levels, topic, level, status, goto}
	consumerHandleGoto             *prometheus.CounterVec   // labels: topics, levels, topic, level, status, goto
	consumerHandleGotoLatency      *prometheus.HistogramVec // labels: topics, levels, topic, level, status, goto, reason: {over_reconsume_times, internal_error}
	consumerHandleGotoConsumeTimes *prometheus.HistogramVec // labels: topics, levels, topic, level, status, goto, reason: {over_reconsume_times, internal_error}
	consumerDeciderOpened          *prometheus.GaugeVec     // gauge of handlers
	consumerDecideGotoSuccess      *prometheus.CounterVec
	consumerDecideGotoFailed       *prometheus.CounterVec
	consumerDecideGotoLatency      *prometheus.HistogramVec // labels: topic, level, status

	// messages terminate labels: {topic=, level=, original_topic=, original_level=,}
	messagesAcks              *prometheus.CounterVec
	messagesNacks             *prometheus.CounterVec
	messagesTerminatedLatency *prometheus.HistogramVec // whole latency on message lifecycle, from published to done/discard/dead. labels: topic, level, status[done, discard, dead]

}

func NewMetricsProvider(metricsLevel int, userDefinedLabels map[string]string) *MetricsProvider {
	metrics := &MetricsProvider{}
	return metrics

}

func (v *MetricsProvider) GetClientMetrics(url string) *ClientMetrics {
	labels := prometheus.Labels{"url": url}
	metrics := &ClientMetrics{
		ClientsOpened: v.clientsOpened.With(labels),
	}
	return metrics
}

func (v *MetricsProvider) GetProducerMetrics(topic string) *ProducerMetrics {
	labels := prometheus.Labels{"topic": topic}
	metrics := &ProducerMetrics{
		ProducersOpened: v.producersOpened.With(labels),
		PublishLatency:  v.producerPublishLatency.With(labels),
	}
	return metrics
}

func (v *MetricsProvider) GetProducerTypedCheckMetrics(topic string, checkType CheckType) *TypedCheckMetrics {
	labels := prometheus.Labels{"topic": topic, "type": checkType.String()}
	metrics := &TypedCheckMetrics{
		CheckersOpened: v.producerTypedCheckerOpened.With(labels),
		CheckPassed:    v.producerTypedCheckPassed.With(labels),
		CheckRejected:  v.producerTypedCheckRejected.With(labels),
		CheckLatency:   v.producerTypedCheckLatency.With(labels),
	}
	return metrics
}

func (v *MetricsProvider) GetProducerDecideGotoMetrics(topic, gotoAction string) *DecideGotoMetrics {
	labels := prometheus.Labels{"topic": topic, "goto": gotoAction}
	metrics := &DecideGotoMetrics{
		DecidersOpened: v.producerDeciderGotoOpened.With(labels),
		DecideSuccess:  v.producerDecideGotoSuccess.With(labels),
		DecideFailed:   v.producerDecideGotoFailed.With(labels),
		DecideLatency:  v.producerDecideGotoLatency.With(labels),
	}
	return metrics
}

func (v *MetricsProvider) GetProducerDecideRouteMetrics(topic, gotoAction string, routeTopic string) *DecideGotoMetrics {
	labels := prometheus.Labels{"topic": topic, "goto": gotoAction, "route_topic": routeTopic}
	metrics := &DecideGotoMetrics{
		DecidersOpened: v.producerDeciderGotoOpened.With(labels),
		DecideSuccess:  v.producerDecideGotoSuccess.With(labels),
		DecideFailed:   v.producerDecideGotoFailed.With(labels),
		DecideLatency:  v.producerDecideGotoLatency.With(labels),
	}
	return metrics
}

func (v *MetricsProvider) GetListenMetrics(topics, levels string) *ListenMetrics {
	labels := prometheus.Labels{"topics": topics, "levels": levels}
	metrics := &ListenMetrics{
		ListenersOpened:  v.listenersOpened.With(labels),
		ListenersRunning: v.listenersRunning.With(labels),
		ReceiveLatency:   v.listenerReceiveLatency.With(labels),
		ListenLatency:    v.listenerListenLatency.With(labels),
		CheckLatency:     v.listenerCheckLatency.With(labels),
		HandleLatency:    v.listenerHandleLatency.With(labels),
		DecideLatency:    v.listenerDecideLatency.With(labels),
	}
	return metrics
}

func (v *MetricsProvider) GetListenerTypedCheckMetrics(topics, levels string, checkType CheckType) *TypedCheckMetrics {
	labels := prometheus.Labels{"topics": topics, "levels": levels, "check_type": checkType.String()}
	metrics := &TypedCheckMetrics{
		CheckersOpened: v.listenerTypedCheckersEnabled.With(labels),
		CheckPassed:    v.listenerTypedCheckPassed.With(labels),
		CheckRejected:  v.listenerTypedCheckRejected.With(labels),
		CheckLatency:   v.listenerTypedCheckLatency.With(labels),
	}
	return metrics
}

func (v *MetricsProvider) GetListenerHandleMetrics(topics, levels string, msgGoto MessageGoto) *ListenerHandleMetrics {
	labels := prometheus.Labels{"topics": topics, "levels": levels, "goto": msgGoto.String()}
	metrics := &ListenerHandleMetrics{
		HandleGotoLatency: v.listenerHandleGotoLatency.With(labels),
		DecideGotoLatency: v.listenerDecideGotoLatency.With(labels),
	}
	return metrics
}

func (v *MetricsProvider) GetConsumerMetrics(topics, levels, topic string, level TopicLevel, status MessageStatus) *ConsumerMetrics {
	labels := prometheus.Labels{"topics": topics, "levels": levels, "topic": topic, "level": level.String(), "status": status.String()}
	metrics := &ConsumerMetrics{
		ConsumersOpened: v.consumersOpened.With(labels),
		ReceiveLatency:  v.consumerReceiveLatency.With(labels),
		ListenLatency:   v.consumerListenLatency.With(labels),
		CheckLatency:    v.consumerCheckLatency.With(labels),
		HandleLatency:   v.consumerHandleLatency.With(labels),
		DecideLatency:   v.consumerDecideLatency.With(labels),
		ConsumeLatency:  v.consumerConsumeLatency.With(labels),
	}
	return metrics

}

func (v *MetricsProvider) GetConsumerHandleGotoMetrics(topics, levels, topic string, level TopicLevel, status MessageStatus, msgGoto MessageGoto) *ConsumerHandleGotoMetrics {
	labels := prometheus.Labels{"topics": topics, "levels": levels, "topic": topic, "level": level.String(),
		"status": status.String(), "goto": msgGoto.String()}
	metrics := &DecideGotoMetrics{
		DecidersOpened: v.consumerDeciderOpened.With(labels),
		DecideSuccess:  v.consumerDecideGotoSuccess.With(labels),
		DecideFailed:   v.consumerDecideGotoFailed.With(labels),
		DecideLatency:  v.consumerDecideGotoLatency.With(labels),
	}
	handledMetrics := &ConsumerHandleGotoMetrics{
		DecideGotoMetrics:      metrics,
		HandleGoto:             v.consumerHandleGoto.With(labels),
		HandleGotoLatency:      v.consumerHandleGotoLatency.With(labels),
		HandleGotoConsumeTimes: v.consumerHandleGotoConsumeTimes.With(labels),
	}
	return handledMetrics
}

func (v *MetricsProvider) GetMessageMetrics(topic, gotoAction string, routeTopic string) *MessageMetrics {
	labels := prometheus.Labels{"topic": topic, "goto": gotoAction, "route_topic": routeTopic}
	metrics := &MessageMetrics{
		Acks:              v.messagesAcks.MustCurryWith(labels),
		Nacks:             v.messagesNacks.MustCurryWith(labels),
		TerminatedLatency: v.messagesTerminatedLatency.MustCurryWith(labels),
	}
	return metrics
}

// ------ helper ------

// ClientMetrics composes these metrics related to client with {url=} label:
// soften_clients_opened {url=};
type ClientMetrics struct {
	ClientsOpened prometheus.Gauge
}

// ProducerMetrics composes these metrics generated during producing with {topic=} label:
// soften_client_producers_opened {topic=};
// soften_client_producer_messages_published {topic=};
// soften_client_producer_publish_latency {topic=};
type ProducerMetrics struct {
	ProducersOpened   prometheus.Gauge    // gauge of producers
	MessagesPublished prometheus.Counter  // counter of producers publish messages
	PublishLatency    prometheus.Observer // producers publish latency
}

type TypedCheckMetrics struct {
	CheckersOpened prometheus.Gauge   // gauge of produce checkers
	CheckPassed    prometheus.Counter // counter of produce check passed
	CheckRejected  prometheus.Counter // counter of produce check rejected
	CheckLatency   prometheus.Observer
}

// DecideGotoMetrics composes these metrics generated during routing with {topic=, route_topic=} label:
// soften_client_routers_opened {topic=, route_topic=};
// soften_client_messages_routed {topic=, route_topic=};
// soften_client_messages_publish_latency {topic=, route_topic};
type DecideGotoMetrics struct {
	DecidersOpened prometheus.Gauge    // gauge of routers
	DecideSuccess  prometheus.Counter  // counter of producers decide messages
	DecideFailed   prometheus.Counter  // counter of producers decide messages
	DecideLatency  prometheus.Observer // route latency
}

// ListenMetrics composed these metrics related multi-level framework, including listen, check, handler, reroute:
// soften_client_listeners_running {topics=,levels=};
// soften_client_listeners_closed {topics=,levels=};
type ListenMetrics struct {
	ListenersOpened  prometheus.Gauge    // gauge of opened listeners
	ListenersRunning prometheus.Gauge    // gauge of running listeners
	ReceiveLatency   prometheus.Observer // labels: topic, level, status
	ListenLatency    prometheus.Observer // labels: topic, level, status
	CheckLatency     prometheus.Observer
	HandleLatency    prometheus.Observer
	DecideLatency    prometheus.Observer
}

// ListenerHandleMetrics composed these metrics related to handle process:
// soften_client_handlers_active {topics=,levels=};
type ListenerHandleMetrics struct {
	HandleGotoLatency prometheus.Observer
	DecideGotoLatency prometheus.Observer
}

// ConsumerMetrics composes these metrics generated during consume with {topic=, level=, status=} label:
// soften_client_consumers_opened {topic=, level=, status=};
// soften_client_messages_received
// soften_client_messages_listened
// soften_client_messages_processed
type ConsumerMetrics struct {
	ConsumersOpened prometheus.Gauge    // gauge of consumers
	ReceiveLatency  prometheus.Observer // labels: topic, level, status
	ListenLatency   prometheus.Observer // labels: topic, level, status
	CheckLatency    prometheus.Observer // labels: topic, level, status
	HandleLatency   prometheus.Observer // labels: topic, level, status
	DecideLatency   prometheus.Observer // labels: topic, level, status
	ConsumeLatency  prometheus.Observer // latency on client consumer consume from received to done/transferred. labels: topic, level, status

}

type ConsumerHandleGotoMetrics struct {
	*DecideGotoMetrics
	HandleGoto             prometheus.Counter  // labels: topics, levels, topic, level, status, goto
	HandleGotoLatency      prometheus.Observer // labels: topics, levels, topic, level, statu
	HandleGotoConsumeTimes prometheus.Observer
}

// MessageMetrics composed these metrics related to transfer process:
// soften_client_transfer_opened {topics=,levels=};
type MessageMetrics struct {
	Acks              *prometheus.CounterVec
	Nacks             *prometheus.CounterVec
	TerminatedLatency prometheus.ObserverVec
}
