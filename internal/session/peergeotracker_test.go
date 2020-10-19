package session

import (
	"fmt"
	"testing"

	"github.com/ipfs/go-bitswap/internal/testutil"
)

func TestPeerGeoTrackerInit(t *testing.T) {
	peers := testutil.GeneratePeers(2)
	pgt := newPeerGeoTracker()
	fmt.Println(peers[0])
	fmt.Println(pgt.getGeoInfo(peers[0]))
	pgt.setGeoInfo(peers[0], 111.11, 22, 22)
	fmt.Println(pgt.getGeoInfo(peers[0]))

}
