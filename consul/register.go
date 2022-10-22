package consul

import (
	"asset/sync/model"
	"fmt"
	"strconv"

	"github.com/hashicorp/consul/api"
)

type consulServiceRegistry struct {
	client api.Client
}

func (c *consulServiceRegistry) Register(instance model.Instance) bool {
	registration := new(api.AgentServiceRegistration)
	registration.ID = instance.InstanceName
	registration.Name = "node-exporter"
	registration.Address = instance.InstanceAddress
	registration.Port = 9100

	/* check := new(api.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d", registration.Address, registration.Port)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "100s"
	registration.Check = check */

	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (c *consulServiceRegistry) Lister() (map[string]bool, bool) {
	maps := make(map[string]bool)
	srvs, _, _ := c.client.Health().Service("node-exporter", "", true, nil)
	for _, srv := range srvs {
		fmt.Println(srv.Service.ID, srv.Service.Address)
	}
	for _, i := range srvs {
		maps[i.Service.ID] = true
	}
	return maps, true
}

func (c *consulServiceRegistry) ListService() {
	srvs, _ := c.client.Agent().Services()
	for _, srv := range srvs {
		fmt.Println(srv)
	}
}

func (c *consulServiceRegistry) Deregister(instanceid string) {
	_ = c.client.Agent().ServiceDeregister(instanceid)
}

func DiffRegister(instances []model.Instance, currmaps map[string]bool) []model.Instance{
	var newInstance []model.Instance
	maps := make(map[string]model.Instance)
	for _, i := range instances {
		maps[i.InstanceId] = i
		if _, ok := currmaps[i.InstanceId]; ! ok {
			newInstance = append(newInstance, i)
		}
	}
	return newInstance
}

func NewConsulServiceRegistry(host string, port int) (*consulServiceRegistry, error) {
	config := api.DefaultConfig()
	config.Address = host + ":" + strconv.Itoa(port)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &consulServiceRegistry{client: *client}, nil
}
