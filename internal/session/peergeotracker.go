// @lry sessionwantsender.go 中用于trace peer的地理位置信息结构
package session

import (
	bsmsg "github.com/ipfs/go-bitswap/message"
	"github.com/ipfs/go-bitswap/priority_queue"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

// var log = logging.Logger("peergeotracker")

// 地理位置距离层级划分
const (
	L1 float64 = 1000000.00 //(m)
	L2 float64 = 3000000.00
	L3 float64 = 5000000.00
)

type GeoNode struct {
	pid      peer.ID
	distance float64
}

func (this *GeoNode) Less(other interface{}) bool {
	return this.distance < other.(*GeoNode).distance
}

type geoInfo struct {
	peerDistance float64
	peerLocation bsmsg.Location
}

// 分层信息 把距离本节点地理位置分三层 进行节点选择算法
// (msg中包含最近处理发送块数作为请求依据？)
// 用于追踪peer距离本节点距离 peer地理位置
type peerGeoTracker struct {
	pgs map[peer.ID]geoInfo
}

func newPeerGeoTracker() *peerGeoTracker {
	return &peerGeoTracker{
		pgs: make(map[peer.ID]geoInfo),
	}
}
func (pgt *peerGeoTracker) setGeoInfo(from peer.ID, distance float64, lat float64, lng float64) {

	pgt.pgs[from] = geoInfo{
		peerDistance: distance,
		peerLocation: bsmsg.Location{
			Lat: lat,
			Lng: lng,
		},
	}
}

// 在peers中选择最近区域的
// 进一步在prt中选择发送block最多的
func (pgt *peerGeoTracker) chooseNearestPeer(peers []peer.ID, prt *peerResponseTracker) peer.ID {
	// fmt.Printf("@lry_debug in peergeotracker.go chooseNearestPeer \n")
	// 将peers 拆分到 peers_l1 peers_l2 peers_l3 中
	var bestPeers []peer.ID
	peers_l1 := priority_queue.New()
	peers_l2 := priority_queue.New()
	peers_l3 := priority_queue.New()
	// 两种想法 1.分层 从最低层开始选  2. 直接遍历找到个第一层的就返回 同时记录最近第二层的与最近第三层 若遍历到最后都没有第一层的就发最低第二层/最低第三层
	for _, p := range peers {
		info, ok := pgt.pgs[p]
		if ok {
			if info.peerDistance <= L1 {
				peers_l1.Push(&GeoNode{
					pid:      p,
					distance: info.peerDistance,
				})
			} else if info.peerDistance > L1 && info.peerDistance <= L2 {
				peers_l2.Push(&GeoNode{
					pid:      p,
					distance: info.peerDistance,
				})
			} else if info.peerDistance > L2 {
				peers_l3.Push(&GeoNode{
					pid:      p,
					distance: info.peerDistance,
				})
			}
		} else {
			// 没有地理位置信息
			log.Debugw("@lry_debug in peergeotracker.go chooseNearestPeer peer ", p, "geoinfo not exist")
			// fmt.Printf("@lry_debug in peergeotracker.go chooseNearestPeer peer %s geoinfo not exist\n", p)
		}
	}
	if peers_l1.Len() > 0 {
		log.Debugw("@lry_debug in peergeotracker.go chooseNearestPeer Level 1 has ", peers_l1.Len(), " peers")
		// fmt.Printf("@lry_debug in peergeotracker.go chooseNearestPeer Level 1 has %d peers \n", peers_l1.Len())
		for i := 0; i < min(peers_l1.Len(), 5); i++ {
			x := peers_l1.Pop().(*GeoNode)
			bestPeers = append(bestPeers, x.pid)
		}
		return prt.choose(bestPeers)
	}
	if peers_l2.Len() > 0 {
		log.Debugw("@lry_debug in peergeotracker.go chooseNearestPeer Level 2 has ", peers_l2.Len(), " peers")
		// fmt.Printf("@lry_debug in peergeotracker.go chooseNearestPeer Level 2 has %d peers \n", peers_l2.Len())
		for i := 0; i < min(peers_l2.Len(), 5); i++ {
			x := peers_l2.Pop().(*GeoNode)
			bestPeers = append(bestPeers, x.pid)
		}
		return prt.choose(bestPeers)
	}
	if peers_l3.Len() > 0 {
		log.Debugw("@lry_debug in peergeotracker.go chooseNearestPeer Level 3 has ", peers_l3.Len(), " peers")
		// fmt.Printf("@lry_debug in peergeotracker.go chooseNearestPeer Level 3 has %d peers \n", peers_l3.Len())
		for i := 0; i < min(peers_l3.Len(), 5); i++ {
			x := peers_l3.Pop().(*GeoNode)
			bestPeers = append(bestPeers, x.pid)
		}
		return prt.choose(bestPeers)
	}
	// 如果所有节点都没有地理位置信息
	// fmt.Printf("@lry_debug in peergeotracker.go chooseNearestPeer no peers has geoinfo \n")
	log.Debugw("@lry_debug in peergeotracker.go chooseNearestPeer no peers has geoinfo")
	return prt.choose(peers)

}

func (pgt *peerGeoTracker) getGeoInfo(p peer.ID) *geoInfo {
	info, ok := pgt.pgs[p]
	if ok {
		return &info
	}
	return nil
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}

}
