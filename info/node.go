package info

type Node struct {
	tableName struct{} `sql:"node_metrics,alias:t" pg:",discard_unknown_columns"`

	NodeID             string    `sql:"node_id,notnull" json:"node_id"`
	NodeIP             string    `sql:"node_ip,notnull" json:"node_ip"`
	CpuInfo            string    `sql:"cpu_info" json:"cpu_info"`
	BtfsVersion        string    `sql:"btfs_version" json:"btfs_version"`
	OsType             string    `sql:"os_type" json:"os_type"`
	ArchType           string    `sql:"arch_type" json:"arch_type"`
	UpTime             uint64    `sql:"up_time,notnull" json:"up_time"`
	StorageUsed        uint64    `sql:"storage_used" json:"storage_used"`
	StorageCap         uint64    `sql:"storage_volume_cap" json:"storage_volume_cap"`
	MemoryUsed         uint64    `sql:"memory_used" json:"memory_used"`
	CpuUsed            float64   `sql:"cpu_used" json:"cpu_used"`
	TimeCreated        time.Time `sql:"time_created" json:"time_created"`
	Upload             uint64    `sql:"upload" json:"upload"`
	Download           uint64    `sql:"download" json:"download"`
	TotalUp            uint64    `sql:"total_upload,notnull" json:"total_upload"`
	TotalDown          uint64    `sql:"total_download,notnull" json:"total_download"`
	BlocksUp           uint64    `sql:"blocks_up,notnull" json:"blocks_up"`
	BlocksDown         uint64    `sql:"blocks_down,notnull" json:"blocks_down"`
	NumPeers           uint64    `sql:"peers_connected,notnull" json:"peers_connected"`
	Reputation         float64   `sql:"reputation,notnull" json:"reputation"`
	StoragePriceDeal   uint64    `sql:"storage_price_deal" json:"storage_price_deal"`
	StoragePriceAsk    uint64    `sql:"storage_price_ask" json:"storage_price_ask"`
	BandwidthPriceDeal uint64    `sql:"bandwidth_price_deal" json:"bandwidth_price_deal"`
	BandwidthPriceAsk  uint64    `sql:"bandwidth_price_ask" json:"bandwidth_price_ask"`
	StorageTimeMin     uint64    `sql:"storage_time_min" json:"storage_time_min"`
	BandwidthLimit     uint64    `sql:"bandwidth_limit" json:"bandwidth_limit"`
	CollateralStake    uint64    `sql:"collateral_stake" json:"collateral_stake"`
}
