package etcdserver

import (
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
	"math/rand"
	"sync"
)

//权重
//官方自带的负载均衡只有pick_first和round_robin两种负载均衡策略，所以自己设计一个加权随机策略
//我们需要把服务器地址的权重添加进去，但是地址resolver.Address并没有提供权重的属性。
//官方给的答复是：把权重存储到地址的元数据metadata中。
type attributeKey struct{}

type AddrInfo struct {
	Weight int
}

//通过set方法设置权重
func SetAddrInfo(addr resolver.Address, addrInfo AddrInfo) resolver.Address {
	addr.Attributes = attributes.New()
	addr.Attributes = addr.Attributes.WithValues(attributeKey{}, addrInfo)
	return addr
}

//通过get方法获得权重
func GetAddrInfo(addr resolver.Address) AddrInfo {
	v := addr.Attributes.Value(attributeKey{})
	ai, _ := v.(AddrInfo)
	return ai
}

//实现加权随机策略

const Name = "Weight"

var (
	minWeight = 1
	maxWeight = 5
)

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilderV2(Name, &rrPickerBuilder{}, base.Config{HealthCheck: false})
}

func init() {
	balancer.Register(newBuilder())
}

type rrPickerBuilder struct {
}

//实现V2PickerBuilder的bliud方法
func (*rrPickerBuilder) Build(info base.PickerBuildInfo) balancer.V2Picker {
	grpclog.Infof("weightPicker: newPicker called with info: %v", info)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPickerV2(balancer.ErrNoSubConnAvailable)
	}
	//根据权重获得权重数量的子链接
	//使用空间换时间的方式，把权重转成地址个数（例如addr1的权重是3，
	//那么添加3个子连接到切片中；addr2权重为1，则添加1个子连接；
	//选择子连接时候，按子连接切片长度生成随机数，以随机数作为下标就是选中的子连接），
	//避免重复计算权重。考虑到内存占用，权重定义从1到5权重。
	var scs []balancer.SubConn
	for subConn, addr := range info.ReadySCs {
		node := GetAddrInfo(addr.Address)
		if node.Weight <= 0 {
			node.Weight = minWeight
		} else if node.Weight > 5 {
			node.Weight = maxWeight
		}
		for i := 0; i < node.Weight; i++ {
			scs = append(scs, subConn)
		}
	}
	return &rrPicker{
		subConns: scs,
	}
}

type rrPicker struct {
	subConns []balancer.SubConn
	//加锁防止并发问题
	mu sync.Mutex
}

//实现V2Picker的pick方法
func (p *rrPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.Lock()
	index := rand.Intn(len(p.subConns))
	sc := p.subConns[index]
	p.mu.Unlock()
	return balancer.PickResult{SubConn: sc}, nil
}
