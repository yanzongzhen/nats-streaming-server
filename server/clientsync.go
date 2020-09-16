package server

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type SyncClient struct {
	ClientID string `json:"client_id,omitempty"`
	Topic    string `json:"topic,omitempty"`
	Group    string `json:"group,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
}

type TopicGroup struct {
	Limit    int                        `json:"limit"`
	Clients  map[string][][]*SyncClient `json:"clients"`
	Snapshot map[string][]int           `json:"snapshot"`
}

// to json
func (t *TopicGroup) ToJson() string {
	d, _ := json.Marshal(t)
	return string(d)
}

// judge the client in the array
func IsArray(arr []*SyncClient, cli *SyncClient) (int, bool) {
	for idx, item := range arr {
		if item.ClientID == cli.ClientID && item.Topic == cli.Topic {
			return idx, true
		}
	}
	return -1, false
}

// add client to array
func (t *TopicGroup) Add(client *SyncClient, isAuto bool) {
	if isAuto {
		t.addClientAutoGroup(client)
	} else {
		t.addClientByGroup(client)
	}
}

// delete client from array
func (t *TopicGroup) Delete(client *SyncClient, isAuto bool) {
	if isAuto {
		t.delClientAutoGroup(client)
	} else {
		t.delClientAutoGroup(client)
	}
}

// add client info to array
func (t *TopicGroup) AddClient(clientId, topic, group, host string, port int, isAuto bool) {
	client := &SyncClient{ClientID: clientId, Topic: topic, Group: group, Host: host, Port: port}
	if isAuto {
		t.addClientAutoGroup(client)
	} else {
		t.addClientByGroup(client)
	}
}

// delete client info from array
func (t *TopicGroup) DelClient(clientId, topic string, isAuto bool) {
	client := &SyncClient{ClientID: clientId, Topic: topic, Group: "", Host: "", Port: 0}
	if isAuto {
		t.delClientAutoGroup(client)
	} else {
		t.addClientByGroup(client)
	}
}

// add new client to Topic Group
func (t *TopicGroup) addClientAutoGroup(client *SyncClient) {
	snapshot := t.Snapshot[client.Topic]
	clients := t.Clients[client.Topic]
	// Init Some Info
	if snapshot == nil && clients == nil {
		t.Snapshot = make(map[string][]int)
		t.Clients = make(map[string][][]*SyncClient)
		clients = make([][]*SyncClient, 0)
		snapshot = make([]int, 0)
		t.Snapshot[client.Topic] = snapshot
		t.Clients[client.Topic] = clients
	}
	if len(snapshot) == 0 {
		snapshot = append(snapshot, 1)
		clientTmp := make([]*SyncClient, 0)
		client.Group = "G0"
		clientTmp = append(clientTmp, client)
		clients = append(clients, clientTmp)
	} else {
		for idx, cli := range clients {
			group := fmt.Sprintf("G%v", idx)
			if len(cli) < t.Limit {
				if _, ok := IsArray(cli, client); ok {
					continue
				}
				client.Group = group
				cli = append(cli, client)
				snapshot[idx]++
				clients[idx] = cli
				break
			} else {
				if len(snapshot)-1 >= idx+1 {
					continue
				} else {
					group := fmt.Sprintf("G%v", idx+1)
					client.Group = group
					newCliList := make([]*SyncClient, 0)
					newCliList = append(newCliList, client)
					clients = append(clients, newCliList)
					snapshot = append(snapshot, 1)
				}
			}
		}
	}
	t.Snapshot[client.Topic] = snapshot
	t.Clients[client.Topic] = clients
}

// delete client from Topic Group Map
func (t *TopicGroup) delClientAutoGroup(client *SyncClient) {
	clients := t.Clients[client.Topic]
	snapshot := t.Snapshot[client.Topic]
	for idx, cliList := range clients {
		index, ok := IsArray(cliList, client)
		if ok {
			tmp := append(cliList[:index], cliList[index+1:]...)
			clients[idx] = tmp
			snapshot[idx]--
		}
	}
	t.Snapshot[client.Topic] = snapshot
}

func (t *TopicGroup) GetClient(clientID string, subject string) *SyncClient {
	clients := t.Clients[subject]
	for _, cliList := range clients {
		for _, client := range cliList {
			if client.ClientID == clientID {
				return client
			}
		}
	}
	return nil
}

// add client by input group name
// need implement
func (t *TopicGroup) addClientByGroup(client *SyncClient) {
	panic("addClientByGroup not implement")
}

// delete client by input group name
// need implement
func (t *TopicGroup) delClientByGroup(client *SyncClient) {
	panic("addClientByGroup not implement")
}

// get group client by single client
func (t *TopicGroup) GetClientGroup(client *SyncClient) []*SyncClient {
	clients := t.Clients[client.Topic]
	if client.Group != "" {
		tmp := strings.Split(client.Group, "G")
		index, _ := strconv.Atoi(tmp[len(tmp)-1])
		return clients[index]
	}
	return nil
}

// get group client by group_name and topic
func (t *TopicGroup) GetGroupClientsByGroup(group, subject string) []*SyncClient {
	clients := t.Clients[subject]
	if group != "" {
		tmp := strings.Split(group, "G")
		index, _ := strconv.Atoi(tmp[len(tmp)-1])
		return clients[index]
	}
	return nil
}
