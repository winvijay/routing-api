package models

import (
	"fmt"
	"time"

	"github.com/nu7hatch/gouuid"
)

type TcpRouteMapping struct {
	Model
	ExpiresAt time.Time
	TcpRouteMappingEntity
}

type TcpRouteMappingEntity struct {
	TcpRoute
	HostPort        uint16 `gorm:"not null; unique_index:idx_tcp_route; type:int" json:"backend_port"`
	HostIP          string `gorm:"not null; unique_index:idx_tcp_route" json:"backend_ip"`
	ModificationTag `json:"modification_tag"`
	TTL             *int `json:"ttl,omitempty"`
}

type Model struct {
	Guid      string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TcpRoute struct {
	RouterGroupGuid string `json:"router_group_guid"`
	ExternalPort    uint16 `gorm:"not null; unique_index:idx_tcp_route; type: int" json:"port"`
}

func NewTcpRouteMapping(routerGroupGuid string, externalPort uint16, hostIP string, hostPort uint16, ttl int) TcpRouteMapping {
	guid, _ := uuid.NewV4()

	m := Model{Guid: guid.String()}
	mapping := TcpRouteMappingEntity{
		TcpRoute: TcpRoute{RouterGroupGuid: routerGroupGuid, ExternalPort: externalPort},
		HostPort: hostPort,
		HostIP:   hostIP,
		TTL:      &ttl,
	}
	return TcpRouteMapping{
		ExpiresAt: time.Now().Add(time.Duration(ttl) * time.Second),
		Model:     m,
		TcpRouteMappingEntity: mapping,
	}
}

func (m TcpRouteMapping) String() string {
	return fmt.Sprintf("%s:%d<->%s:%d", m.RouterGroupGuid, m.ExternalPort, m.HostIP, m.HostPort)
}

func (m TcpRouteMapping) Matches(other TcpRouteMapping) bool {
	return m.RouterGroupGuid == other.RouterGroupGuid &&
		m.ExternalPort == other.ExternalPort &&
		m.HostIP == other.HostIP &&
		m.HostPort == other.HostPort &&
		*m.TTL == *other.TTL
}

func (t *TcpRouteMapping) SetDefaults(maxTTL int) {
	// default ttl if not present
	// TTL is a pointer to a uint16 so that we can
	// detect if it's present or not (i.e. nil or 0)
	if t.TTL == nil {
		t.TTL = &maxTTL
	}
}
