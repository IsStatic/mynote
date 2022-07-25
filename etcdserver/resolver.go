package etcdserver

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"log"
	"mynote/pkg/constants"
	"strconv"
	"strings"
)

const schema = constants.Schema

// resolver is the implementaion of grpc.resolve.Builder
// Resolver 实现grpc的grpc.resolve.Builder接口的Build与Scheme方法
type Resolver struct {
	endpoints []string
	service   string
	cli       *clientv3.Client
	cc        resolver.ClientConn
}

// NewResolver return resolver builder
// service is service name
func NewResolver(endpoints []string, service string) resolver.Builder {
	return &Resolver{endpoints: endpoints, service: service}
}

// Scheme return etcd schema
func (r *Resolver) Scheme() string {
	// 最好用这种，因为grpc resolver.Register(r)在注册时，会取scheme，如果一个系统有多个grpc发现，就会覆盖之前注册的
	return schema + "_" + r.service
}

// ResolveNow
func (r *Resolver) ResolveNow(rn resolver.ResolveNowOptions) {
}

// Close
func (r *Resolver) Close() {
	r.cli.Close()
	r.Close()
	log.Println("server stop")
}

// Build to resolver.Resolver
// 实现grpc.resolve.Builder接口的方法
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var err error

	r.cli, err = clientv3.New(clientv3.Config{
		Endpoints: r.endpoints,
	})
	if err != nil {
		return nil, fmt.Errorf("grpclb: create clientv3 client failed: %v", err)
	}

	r.cc = cc

	// go r.watch(fmt.Sprintf("/%s/%s/", schema, r.service))
	go r.watch(fmt.Sprintf(r.service))

	return r, nil
}

type addrAndWeight struct {
	resolver.Address
	servicename string
	weight      int
}

//监听服务的变化
func (r *Resolver) watch(prefix string) {
	addrDict := make(map[string]addrAndWeight)

	update := func() {
		addrList := make([]resolver.Address, 0, len(addrDict))
		for _, v := range addrDict {
			addrList = append(addrList,
				SetAddrInfo(v.Address, AddrInfo{Weight: v.weight}))
			log.Println("put key :", v.servicename, "wieght:", v.weight)
		}
		state := resolver.State{
			Addresses:     addrList,
			ServiceConfig: nil,
			Attributes:    nil,
		}
		r.cc.UpdateState(state)
	}

	resp, err := r.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err == nil {
		for i, kv := range resp.Kvs {
			info := &ServiceInfo{}
			err := json.Unmarshal([]byte(kv.Value), info)
			if err != nil {
				log.Println("反序列化失败")
			}
			//log.Println("获得数据：kvs:",resp.Kvs[i].Key,
			//	"KvsValue:",resp.Kvs[i].Value,"Kvkey:",
			//	kv.Key,"KvValue：",kv.Value)
			nodeweigth, err := strconv.Atoi(info.Weight)
			if err != nil {
				nodeweigth = 1
			}
			addrDict[string(resp.Kvs[i].Value)] = addrAndWeight{
				Address:     resolver.Address{Addr: info.IP},
				servicename: info.Name,
				weight:      nodeweigth,
			}
		}
	}

	update()

	rch := r.cli.Watch(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for n := range rch {
		for _, ev := range n.Events {
			switch ev.Type {
			case mvccpb.PUT:
				info := &ServiceInfo{}
				err := json.Unmarshal([]byte(ev.Kv.Value), info)
				if err != nil {
					log.Println(err)
				} else {
					nodeweigth, err := strconv.Atoi(info.Weight)
					if err != nil {
						nodeweigth = 1
					}
					addrDict[string(ev.Kv.Key)] = addrAndWeight{
						Address:     resolver.Address{Addr: info.IP},
						servicename: info.Name,
						weight:      nodeweigth,
					}
				}
			case mvccpb.DELETE:
				delete(addrDict, string(ev.PrevKv.Key))
			}
		}
		update()
	}
}

//SetServiceList 设置服务地址
func (r *Resolver) SetServiceList(key, val string) {
	//获取服务地址
	addr := resolver.Address{Addr: strings.TrimPrefix(key, r.service)}
	//获取服务地址权重
	nodeWeight, err := strconv.Atoi(val)
	if err != nil {
		//非数字字符默认权重为1
		nodeWeight = 1
	}
	//把服务地址权重存储到resolver.Address的元数据中
	addr = SetAddrInfo(addr, AddrInfo{Weight: nodeWeight})
	log.Println("put key :", key, "wieght:", val)
}
