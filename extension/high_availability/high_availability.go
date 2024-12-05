package high_availability

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	_net "net"
	"sort"
	"strconv"
	"strings"
	"time"
)

type HighAvailability struct {
	context          iface.Context
	onlineAddrTable  map[string]int64
	offlineAddrTable map[string]int64
	running          bool
}

func GetInstance(globalContext iface.Context) iface.Component {
	return &HighAvailability{
		context:          globalContext,
		onlineAddrTable:  make(map[string]int64),
		offlineAddrTable: make(map[string]int64),
		running:          false,
	}
}

func (h *HighAvailability) GetName() string {
	return "high_availability"
}

func (h *HighAvailability) GetDescribe() string {
	return "高可用模块"
}

func (h *HighAvailability) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.NO_VALUE, cons.LOWER_V, "", "view", false,
		nil, "查看当前的节点在线列表")
	cm.AddParameters(cons.STRING, cons.LOWER_R, "highAvailability.addresses", "run", false,
		nil, "启动高可用模块，多个地址间英文逗号分割")
	cm.AddParameters(cons.STRING, cons.LOWER_T, "", "transport", false,
		nil, "交互节点信息")
}

func (h *HighAvailability) Do(params map[string]any) (resp []byte) {
	_, ok := params["view"]
	if ok {
		if !h.running {
			return []byte("请先启动高可用模块")
		}
		result := make(map[string]map[string]int64)
		result["online"] = h.onlineAddrTable
		result["offline"] = h.offlineAddrTable
		marshal, err := json.Marshal(result)
		if err != nil {
			return []byte("节点信息序列化失败")
		}
		return marshal
	}
	r, ok := params["run"]
	if ok {
		if !h.running {
			h.running = true
			err := h.start(r.(string))
			if err != nil {
				h.running = false
				return []byte("高可用模块启动失败")
			}
			return []byte("高可用模块启动成功")
		} else {
			return []byte("高可用模块已经启动，请不要重复操作")
		}
	}
	t, ok := params["transport"]
	if ok {
		if !h.running {
			return []byte("请先启动高可用模块")
		}
		transport, err := h.acceptTransport(t.(string))
		if err != nil {
			return nil
		}
		return transport
	}
	return []byte("指令解析失败")
}

func (h *HighAvailability) start(adds string) error {
	component := h.context.FindComponent(cons.WEB_SVC, false)
	if component == nil {
		return errors.New("需要前置模块【web_server】")
	}
	serverPort := h.context.GetConfig().GetInt("web_server.port")
	myAddr := getLocalIp() + ":" + strconv.Itoa(serverPort)
	h.onlineAddrTable[myAddr] = time.Now().Unix()
	if adds != myAddr && adds != "localhost:"+strconv.Itoa(serverPort) && adds != "127.0.0.1:"+strconv.Itoa(serverPort) {
		configAddresses := strings.Split(adds, ",")
		for _, c := range configAddresses {
			err := h.transport(c, true)
			if err == nil {
				h.onlineAddrTable[c] = time.Now().Unix()
			} else {
				h.offlineAddrTable[c] = time.Now().Unix()
				return err
			}
		}
	}
	return nil
}

func (h *HighAvailability) transport(addr string, needOffline bool) error {
	param := make(map[string]map[string]int64)
	param["online"] = h.onlineAddrTable
	param["offline"] = h.offlineAddrTable
	marshal, err := json.Marshal(param)
	if err != nil {
		return err
	}
	commands := []string{
		h.GetName(),
		cons.LOWER_T.GetCommandName(),
		string(marshal),
		"--f",
		addr,
	}
	result, err2 := h.context.Servlet(commands, false)
	if err2 != nil {
		if needOffline {
			go h.offline(addr)
		}
		log.Errorf("请求[%s]失败", addr)
		return err2
	}
	onlineAndOfflineAddr := make(map[string]map[string]int64)
	err1 := json.Unmarshal(result, &onlineAndOfflineAddr)
	if err1 != nil {
		if needOffline {
			go h.offline(addr)
		}
		log.Errorf("反序列化[%s]地址列表失败", addr)
		return err1
	}
	receiveOnlineAddr := onlineAndOfflineAddr["online"]
	receiveOfflineAddr := onlineAndOfflineAddr["offline"]
	mergeAddr(receiveOnlineAddr, h.onlineAddrTable)
	mergeAddr(receiveOfflineAddr, h.offlineAddrTable)
	checkOffline(h.onlineAddrTable, h.offlineAddrTable)
	return nil
}

func (h *HighAvailability) acceptTransport(n string) ([]byte, error) {
	returnParam := make(map[string]map[string]int64)
	returnParam["online"] = h.onlineAddrTable
	returnParam["offline"] = h.offlineAddrTable
	marshal, err := json.Marshal(returnParam)
	if err != nil {
		log.Error("序列化地址列表失败")
		return nil, err
	}
	param := make(map[string]map[string]int64)
	err1 := json.Unmarshal([]byte(n), &param)
	if err1 != nil {
		log.Error("反序列化地址列表失败")
		return nil, err1
	}
	receiveOnlineAddrTable := param["online"]
	receiveOfflineAddrTable := param["offline"]
	mergeAddr(receiveOnlineAddrTable, h.onlineAddrTable)
	mergeAddr(receiveOfflineAddrTable, h.offlineAddrTable)
	checkOffline(h.onlineAddrTable, h.offlineAddrTable)
	return marshal, nil
}

func (h *HighAvailability) heartbeat() {
	serverPort := h.context.GetConfig().GetInt("web_server.port")
	myAddr := getLocalIp() + ":" + strconv.Itoa(serverPort)
	h.onlineAddrTable[myAddr] = time.Now().Unix()
	addr := getSendAddr(myAddr, h.onlineAddrTable)
	h.transport(addr, true)
}

func (h *HighAvailability) offline(addr string) {
	h.offlineAddrTable[addr] = time.Now().Unix()
	for k, _ := range h.onlineAddrTable {
		checkOffline(h.onlineAddrTable, h.offlineAddrTable)
		_, b := h.offlineAddrTable[k]
		if b {
			h.transport(k, false)
		}
	}
}

func getLocalIp() string {
	addrs, err := _net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*_net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func mergeAddr(m1 map[string]int64, m2 map[string]int64) {
	for k, v := range m1 {
		i, e := m2[k]
		if e {
			if i < v {
				m2[k] = v
			}
		} else {
			m2[k] = m1[k]
		}
	}
}

func checkOffline(on map[string]int64, off map[string]int64) {
	for k, v := range on {
		i, e := off[k]
		if e {
			if i > v {
				delete(on, k)
			} else {
				delete(off, k)
			}
		}
	}
}

func getSendAddr(myAddr string, on map[string]int64) string {
	index := -1
	addrs := []string{}
	for k, _ := range on {
		addrs = append(addrs, k)
	}
	sort.Strings(addrs)
	for i, v := range addrs {
		if myAddr == v {
			index = i
		}
	}
	if index == -1 || index == len(addrs) {
		index = 0
	}
	return addrs[index]
}
