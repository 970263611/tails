package server_on_off

type QueryInstanceRespEntry struct {
	Code    int
	Message string
	Data    Data
}

type Data struct {
	Instances []Instances
	Count     int
}

type Instances struct {
	InstanceId  string  `json:"instanceId"`
	Ip          string  `json:"ip"`
	Port        int     `json:"port"`
	Weight      float64 `json:"weight"`
	Healthy     bool    `json:"healthy"`
	Enabled     bool    `json:"enabled"`
	Ephemeral   bool    `json:"ephemeral"`
	ClusterName string  `json:"clusterName"`
	ServiceName string  `json:"serviceName"`
}
