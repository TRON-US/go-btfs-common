package info

import (
	"time"
)

// We're deprecating this package to switch using protobuf
// TODO: Remove this package after it's fully deprecated.

type Node struct {
	tableName struct{} `pg:"node_metrics,alias:t,discard_unknown_columns"`

	NodeID             string    `pg:"node_id,notnull" json:"node_id"`
	NodeIP             string    `pg:"node_ip,notnull" json:"node_ip"`
	CpuInfo            string    `pg:"cpu_info" json:"cpu_info"`
	BtfsVersion        string    `pg:"btfs_version" json:"btfs_version"`
	OsType             string    `pg:"os_type" json:"os_type"`
	ArchType           string    `pg:"arch_type" json:"arch_type"`
	UpTime             uint64    `pg:"up_time,notnull" json:"up_time"`
	StorageUsed        uint64    `pg:"storage_used" json:"storage_used"`
	StorageCap         uint64    `pg:"storage_volume_cap" json:"storage_volume_cap"`
	MemoryUsed         uint64    `pg:"memory_used" json:"memory_used"`
	CpuUsed            float64   `pg:"cpu_used" json:"cpu_used"`
	TimeCreated        time.Time `pg:"time_created" json:"time_created"`
	Upload             uint64    `pg:"upload" json:"upload"`
	Download           uint64    `pg:"download" json:"download"`
	TotalUp            uint64    `pg:"total_upload,notnull" json:"total_upload"`
	TotalDown          uint64    `pg:"total_download,notnull" json:"total_download"`
	BlocksUp           uint64    `pg:"blocks_up,notnull" json:"blocks_up"`
	BlocksDown         uint64    `pg:"blocks_down,notnull" json:"blocks_down"`
	NumPeers           uint64    `pg:"peers_connected,notnull" json:"peers_connected"`
	Reputation         float64   `pg:"reputation,notnull" json:"reputation"`
	StoragePriceDeal   uint64    `pg:"storage_price_deal" json:"storage_price_deal"`
	BandwidthPriceDeal uint64    `pg:"bandwidth_price_deal" json:"bandwidth_price_deal"`

	NodeStorage
}

// Host storage publishable information
type NodeStorage struct {
	StoragePriceAsk   uint64  `pg:"storage_price_ask" json:"storage_price_ask"`
	BandwidthPriceAsk uint64  `pg:"bandwidth_price_ask" json:"bandwidth_price_ask"`
	StorageTimeMin    uint64  `pg:"storage_time_min" json:"storage_time_min"`
	BandwidthLimit    float64 `pg:"bandwidth_limit" json:"bandwidth_limit"`
	CollateralStake   uint64  `pg:"collateral_stake" json:"collateral_stake"`
}
