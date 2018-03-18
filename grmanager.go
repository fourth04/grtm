package grtm

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

const (
	// STOP signal
	STOP = "__P"
)

// GoroutineChannel associate name and channel
type GoroutineChannel struct {
	gid   uint64
	name  string
	chMsg chan string
}

// GrManager create a locked map[string]*GoroutineChannel
type GrManager struct {
	mutex      sync.Mutex
	grchannels map[string]*GoroutineChannel
}

func (gm *GrManager) register(name string) error {
	gchannel := &GoroutineChannel{
		gid:  uint64(rand.Int63()),
		name: name,
	}
	gchannel.chMsg = make(chan string)
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	if gm.grchannels == nil {
		gm.grchannels = make(map[string]*GoroutineChannel)
	} else if _, ok := gm.grchannels[gchannel.name]; ok {
		return fmt.Errorf("goroutine channel already defined: %q", gchannel.name)
	}
	gm.grchannels[gchannel.name] = gchannel
	return nil
}

func (gm *GrManager) unregister(name string) error {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	if _, ok := gm.grchannels[name]; !ok {
		return fmt.Errorf("goroutine channel not find: %q", name)
	}
	delete(gm.grchannels, name)
	return nil
}

// NewNormalGoroutine create a normal goroutine including register and unregister, but the register and unregister make no sense, the goroutine will stop itself
func (gm *GrManager) NewNormalGoroutine(name string, fc interface{}, args ...interface{}) {
	go func() {
		//register channel
		err := gm.register(name)
		if err != nil {
			return
		}
		if len(args) > 1 {
			fc.(func(...interface{}))(args)
		} else if len(args) == 1 {
			fc.(func(interface{}))(args[0])
		} else {
			fc.(func())()
		}
		gm.unregister(name)
	}()
}

// NewLoopGoroutine create a loop goroutine including register, it will repeatitively execute the function until a STOP message is sended in the message channel
func (gm *GrManager) NewLoopGoroutine(name string, fc interface{}, args ...interface{}) {
	go func() {
		//register channel
		err := gm.register(name)
		if err != nil {
			return
		}
		for {
			select {
			case info := <-gm.grchannels[name].chMsg:
				taskInfo := strings.Split(info, ":")
				signal, gid := taskInfo[0], taskInfo[1]
				if gid == strconv.Itoa(int(gm.grchannels[name].gid)) {
					if signal == STOP {
						fmt.Printf("[gid: %s, name: %s] quit\n", gid, name)
						gm.unregister(name)
						return
					}
					fmt.Println("unknown signal")
				}
			default:
				fmt.Println("no signal")
				if len(args) > 1 {
					fc.(func(...interface{}))(args)
				} else if len(args) == 1 {
					fc.(func(interface{}))(args[0])
				} else {
					fc.(func())()
				}
			}
		}
	}()
}

// StopGoroutine sends a STOP message to the named message channel to close the associated goroutine
func (gm *GrManager) StopGoroutine(name string) error {
	stopChannel, ok := gm.grchannels[name]
	if !ok {
		return fmt.Errorf("not found goroutine name :" + name)
	}
	gid := strconv.Itoa(int(stopChannel.gid))
	gm.grchannels[name].chMsg <- STOP + ":" + gid
	return nil
}

// NewDiyGoroutine create a diy goroutine, the "diy" means the function will define a chMsg parameter which like a slot, NewDiyGoroutine will register a channel and send it to the function, so the function will use the chMsg to stop
func (gm *GrManager) NewDiyGoroutine(name string, fc interface{}, args ...interface{}) {
	go func() {
		//register channel
		err := gm.register(name)
		if err != nil {
			return
		}
		if len(args) >= 1 {
			fc.(func(chan string, ...interface{}))(gm.grchannels[name].chMsg, args)
		} else if len(args) == 1 {
			fc.(func(chan string, interface{}))(gm.grchannels[name].chMsg, args[0])
		} else {
			fc.(func(chan string))(gm.grchannels[name].chMsg)
		}
	}()
}

// StopDiyGoroutine sends a STOP message to the named message channel to close the associated goroutine, prints quit infomation and unregister the name
func (gm *GrManager) StopDiyGoroutine(name string) error {
	err := gm.StopGoroutine(name)
	if err != nil {
		return err
	}
	gid := strconv.Itoa(int(gm.grchannels[name].gid))
	fmt.Printf("[gid: %s, name: %s] quit\n", gid, name)
	gm.unregister(name)
	return nil
}
