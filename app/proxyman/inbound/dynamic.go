package inbound

import (
	"context"
	"sync"
	"time"

	"v2ray.com/core/app/log"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/proxyman/mux"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
)

type DynamicInboundHandler struct {
	tag            string
	ctx            context.Context
	cancel         context.CancelFunc
	proxyConfig    interface{}
	receiverConfig *proxyman.ReceiverConfig
	portMutex      sync.Mutex
	portsInUse     map[net.Port]bool
	workerMutex    sync.RWMutex
	worker         []worker
	lastRefresh    time.Time
	mux            *mux.Server
}

func NewDynamicInboundHandler(ctx context.Context, tag string, receiverConfig *proxyman.ReceiverConfig, proxyConfig interface{}) (*DynamicInboundHandler, error) {
	ctx, cancel := context.WithCancel(ctx)
	h := &DynamicInboundHandler{
		ctx:            ctx,
		tag:            tag,
		cancel:         cancel,
		proxyConfig:    proxyConfig,
		receiverConfig: receiverConfig,
		portsInUse:     make(map[net.Port]bool),
		mux:            mux.NewServer(ctx),
	}

	return h, nil
}

func (h *DynamicInboundHandler) allocatePort() net.Port {
	from := int(h.receiverConfig.PortRange.From)
	delta := int(h.receiverConfig.PortRange.To) - from + 1

	h.portMutex.Lock()
	defer h.portMutex.Unlock()

	for {
		r := dice.Roll(delta)
		port := net.Port(from + r)
		_, used := h.portsInUse[port]
		if !used {
			h.portsInUse[port] = true
			return port
		}
	}
}

func (h *DynamicInboundHandler) waitAnyCloseWorkers(ctx context.Context, cancel context.CancelFunc, workers []worker, duration time.Duration) {
	time.Sleep(duration)
	cancel()
	ports2Del := make([]net.Port, len(workers))
	for idx, worker := range workers {
		ports2Del[idx] = worker.Port()
		worker.Close()
	}

	h.portMutex.Lock()
	for _, port := range ports2Del {
		delete(h.portsInUse, port)
	}
	h.portMutex.Unlock()
}

func (h *DynamicInboundHandler) refresh() error {
	h.lastRefresh = time.Now()

	timeout := time.Minute * time.Duration(h.receiverConfig.AllocationStrategy.GetRefreshValue()) * 2
	concurrency := h.receiverConfig.AllocationStrategy.GetConcurrencyValue()
	ctx, cancel := context.WithTimeout(h.ctx, timeout)
	workers := make([]worker, 0, concurrency)

	address := h.receiverConfig.Listen.AsAddress()
	if address == nil {
		address = net.AnyIP
	}
	for i := uint32(0); i < concurrency; i++ {
		port := h.allocatePort()
		p, err := proxy.CreateInboundHandler(ctx, h.proxyConfig)
		if err != nil {
			log.Trace(newError("failed to create proxy instance").Base(err).AtWarning())
			continue
		}
		nl := p.Network()
		if nl.HasNetwork(net.Network_TCP) {
			worker := &tcpWorker{
				tag:          h.tag,
				address:      address,
				port:         port,
				proxy:        p,
				stream:       h.receiverConfig.StreamSettings,
				recvOrigDest: h.receiverConfig.ReceiveOriginalDestination,
				dispatcher:   h.mux,
				sniffers:     h.receiverConfig.DomainOverride,
			}
			if err := worker.Start(); err != nil {
				log.Trace(newError("failed to create TCP worker").Base(err).AtWarning())
				continue
			}
			workers = append(workers, worker)
		}

		if nl.HasNetwork(net.Network_UDP) {
			worker := &udpWorker{
				tag:          h.tag,
				proxy:        p,
				address:      address,
				port:         port,
				recvOrigDest: h.receiverConfig.ReceiveOriginalDestination,
				dispatcher:   h.mux,
			}
			if err := worker.Start(); err != nil {
				log.Trace(newError("failed to create UDP worker").Base(err).AtWarning())
				continue
			}
			workers = append(workers, worker)
		}
	}

	h.workerMutex.Lock()
	h.worker = workers
	h.workerMutex.Unlock()

	go h.waitAnyCloseWorkers(ctx, cancel, workers, timeout)

	return nil
}

func (h *DynamicInboundHandler) monitor() {
	timer := time.NewTicker(time.Minute * time.Duration(h.receiverConfig.AllocationStrategy.GetRefreshValue()))
	defer timer.Stop()

	for {
		select {
		case <-h.ctx.Done():
			return
		case <-timer.C:
			h.refresh()
		}
	}
}

func (h *DynamicInboundHandler) Start() error {
	err := h.refresh()
	go h.monitor()
	return err
}

func (h *DynamicInboundHandler) Close() {
	h.cancel()
}

func (h *DynamicInboundHandler) GetRandomInboundProxy() (proxy.Inbound, net.Port, int) {
	h.workerMutex.RLock()
	defer h.workerMutex.RUnlock()

	if len(h.worker) == 0 {
		return nil, 0, 0
	}
	w := h.worker[dice.Roll(len(h.worker))]
	expire := h.receiverConfig.AllocationStrategy.GetRefreshValue() - uint32(time.Since(h.lastRefresh)/time.Minute)
	return w.Proxy(), w.Port(), int(expire)
}
