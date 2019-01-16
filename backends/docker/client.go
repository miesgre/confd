package docker

import (
	"fmt"
	//"os"
	"strings"
	"context"
	"encoding/json"

	"github.com/kelseyhightower/confd/log"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/filters"
)


// Client provides a shell for the env client
type Client struct{
	client *docker.Client
}

// NewDockerClient returns a new client
func NewDockerClient() (*Client, error) {
	cli, err := docker.NewEnvClient()

	client := &Client{
		client:		cli,
	}

	return client, err
}
// GetValues queries the environment for keys
func (c *Client) GetValues(keys []string) (map[string]string, error) {
	result := map[string]string{}
	for _, key := range keys {
		r := make(map[string]string)
		if (key == "/") {
			// Returns eveything
			r, _ = c.GetValues([]string{
				"/containers",
				"/services",
				"/networks",
				"/tasks",
				"/volumes",
				"/nodes",
			});							
		} else {
			parts := strings.Split(strings.Trim(key, "/"), "/")
			namespace := parts[0]
			
			if (namespace == "containers") {
				r, _ = c.GetContainersValues(key);
			} else if (namespace == "services") {
				r, _ = c.GetServicesValues(key);
			} else if (namespace == "networks") {
				r, _ = c.GetNetworksValues(key);
			} else if (namespace == "tasks") {
				r, _ = c.GetTasksValues(key);
			} else if (namespace == "volumes") {
				r, _ = c.GetVolumesValues(key);
			} else if (namespace == "nodes") {
				r, _ = c.GetNodesValues(key);
			}
		}
		for k, v := range r { result[k] = v }
	}

	return result, nil	
}


/**********
 *
 * /networks/byId/<id>
 * /networks/byName/<name>
 *
 **********/
func (c *Client) GetNetworksValues(key string) (map[string]string, error) {
	result := map[string]string{}

	parts := strings.Split(strings.Trim(key, "/"), "/")
	if len(parts) == 1 {
		return c.GetValues([]string{"/networks/byId", "/networks/byName"})
	} else if len(parts) == 2 {
		byKey := parts[1]
		items, _ := c.client.NetworkList(context.Background(), types.NetworkListOptions{})		
		for _, i := range items {
			jsonString, _ := json.Marshal(i)
			if (byKey == "byId") {							
				result[key+"/"+i.ID] = string(jsonString)
			} else if (byKey == "byName") {
				result[key+"/"+i.Name] = string(jsonString)
			}
		}
	} else if len(parts) == 3 {
		byKey := parts[1]
		id := parts[2]
		filter := filters.NewArgs()
		if (byKey == "byId") {
			filter.Add("id", id)		
		} else if (byKey == "byName") {
			filter.Add("name", id)	
		}
		items, _ := c.client.NetworkList(context.Background(), types.NetworkListOptions{filter})
		jsonString, _ := json.Marshal(items[0])
		result[key] = string(jsonString)
	} 
	return result, nil
}


/**********
 *
 * /services/byId/<sid>
 * /services/byName/<sname>
 *
 **********/
func (c *Client) GetServicesValues(key string) (map[string]string, error) {
	result := map[string]string{}

	parts := strings.Split(strings.Trim(key, "/"), "/")
	if len(parts) == 1 {
		return c.GetValues([]string{"/services/byId", "/services/byName"})
	} else if len(parts) == 2 {
		byKey := parts[1]
		services, _ := c.client.ServiceList(context.Background(), types.ServiceListOptions{})		
		for _, s := range services {
			jsonString, _ := json.Marshal(s)
			if (byKey == "byId") {							
				result[key+"/"+s.ID] = string(jsonString)
				names, err := c.GetServicesValues("/services/byId/"+s.ID+"/tasks")
				if err != nil { return nil, err }
				for k, v := range names { result[k] = v }		
			} else if (byKey == "byName") {
				result[key+"/"+s.Spec.Name] = string(jsonString)
			}
		}
	} else if len(parts) == 3 {
		byKey := parts[1]
		id := parts[2]
		filter := filters.NewArgs()
		if (byKey == "byId") {
			filter.Add("id", id)		
		} else if (byKey == "byName") {
			filter.Add("name", id)	
		}
		services, _ := c.client.ServiceList(context.Background(), types.ServiceListOptions{filter})
		jsonString, _ := json.Marshal(services[0])
		result[key] = string(jsonString)
	}

	return result, nil
}


/**********
 *
 * /tasks/byId/<taskid>
 * /tasks/byService/<serviceid>/<taskid>
 *
 **********/
 func (c *Client) GetTasksValues(key string) (map[string]string, error) {
	result := map[string]string{}

	parts := strings.Split(strings.Trim(key, "/"), "/")
	if len(parts) == 1 {
		return c.GetValues([]string{"/tasks/byId", "/tasks/byService"})
	} else if len(parts) == 2 {
		byKey := parts[1]
		networks, _ := c.client.TaskList(context.Background(), types.TaskListOptions{})		
		for _, n := range networks {
			jsonString, _ := json.Marshal(n)
			if (byKey == "byId") {							
				result[key+"/"+n.ID] = string(jsonString)
			} else if (byKey == "byName") {
				result[key+"/"+n.Name] = string(jsonString)
			} else if (byKey == "byService") {
				result[key+"/"+n.ServiceID+"/"+n.ID] = string(jsonString)
			}
		}
	} else if len(parts) == 3 {
		byKey := parts[1]
		id := parts[2]
		filter := filters.NewArgs()
		if (byKey == "byId") {
			filter.Add("id", id)		
		} else if (byKey == "byName") {
			filter.Add("name", id)	
		}
		networks, _ := c.client.NetworkList(context.Background(), types.NetworkListOptions{filter})
		jsonString, _ := json.Marshal(networks[0])
		result[key] = string(jsonString)
	} 
	return result, nil
}


/**********
 *
 * /volumes/<id>
 *
 **********/
 func (c *Client) GetVolumesValues(key string) (map[string]string, error) {
	result := map[string]string{}

	parts := strings.Split(strings.Trim(key, "/"), "/")
	if len(parts) == 1 {
		items, _ := c.client.VolumeList(context.Background(), filters.Args{})
		for _, i := range items.Volumes {
			jsonString, _ := json.Marshal(i)
			result[key+"/"+i.Name] = string(jsonString)
		}
	} else if len(parts) == 2 {
		id := parts[1]
		_, vol, _ := c.client.VolumeInspectWithRaw(context.Background(), id)		
		result[key] = string(vol)
	} 
	return result, nil
}

/**********
 *
 * /nodes/<id>
 *
 **********/
 func (c *Client) GetNodesValues(key string) (map[string]string, error) {
	result := map[string]string{}
	parts := strings.Split(strings.Trim(key, "/"), "/")
	if len(parts) == 1 {
		items, _ := c.client.NodeList(context.Background(), types.NodeListOptions{})		
		for _, i := range items {
			jsonString, _ := json.Marshal(i)
			result[key+"/"+i.ID] = string(jsonString)
		}
	} else if len(parts) == 2 {
		id := parts[1]
		_, item, _ := c.client.NodeInspectWithRaw(context.Background(), id)
		result[key] = string(item)
	} 
	return result, nil
}

/**********
 *
 * /containers/<id>
 *
 **********/
 func (c *Client) GetContainersValues(key string) (map[string]string, error) {
	result := map[string]string{}
	parts := strings.Split(strings.Trim(key, "/"), "/")
	if len(parts) == 1 {
		items, _ := c.client.ContainerList(context.Background(), types.ContainerListOptions{})		
		for _, i := range items {
			jsonString, _ := json.Marshal(i)
			result[key+"/"+i.ID] = string(jsonString)
		}
	} else if len(parts) == 2 {
		id := parts[1]
		_, item, _ := c.client.ContainerInspectWithRaw(context.Background(), id, true)
		result[key] = string(item)
	} 
	return result, nil
}









func (c *Client) WatchPrefix(prefix string, keys []string, waitIndex uint64, stopChan chan bool) (uint64, error) {
	if waitIndex == 0 {
		return 1, nil
	}	

	events, err := c.client.Events(context.Background(), types.EventsOptions{})
	for {
		select {
		case event := <-events:
			log.Debug(fmt.Sprintf("docker event: %s %#v", event.Type, event));
			return 1, nil			
		case err := <-err:
			return 0, err
		case <-stopChan:			
			return 0, nil
		}
	}
	
	return waitIndex, nil
}

