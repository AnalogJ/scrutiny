export interface ZfsPoolModel {
    pool_guid: string;
    name: string;
    host_id: string;
    state: string; // ONLINE, DEGRADED, FAULTED, OFFLINE, etc.
    txg: string;
    spa_version: string;
    zpl_version: string;
    status: string;
    action: string;
    error_count: string;
    
    // Timestamps
    created_at?: string;
    updated_at?: string;
    
    // Space information (from zpool status)
    alloc_space: string;
    total_space: string;
    def_space: string;
    
    // Pool properties (from zpool list)
    size?: string;
    allocated?: string;
    free?: string;
    fragmentation?: string;
    capacity_percent?: string;
    dedupratio?: string;
    
    // Error counters
    read_errors: string;
    write_errors: string;
    checksum_errors: string;
    
    // Scan information (scrub/resilver)
    scan_function?: string;
    scan_state?: string;
    scan_start_time?: string;
    scan_end_time?: string;
    scan_to_examine?: string;
    scan_examined?: string;
    scan_processed?: string;
    scan_errors?: string;
    scan_issued?: string;
    
    // Related vdevs
    vdevs?: ZfsVdevModel[];
}

export interface ZfsVdevModel {
    id: number;
    pool_guid: string;
    parent_id?: number;
    guid: string;
    name: string;
    
    // Vdev information
    vdev_type: string; // root, mirror, raidz, raidz2, raidz3, disk, file, etc.
    class: string; // normal, spare, cache, etc.
    state: string; // ONLINE, DEGRADED, FAULTED, OFFLINE, etc.
    
    // Space information
    alloc_space?: string;
    total_space?: string;
    def_space?: string;
    rep_dev_size?: string;
    phys_space?: string;
    
    // Error counters
    read_errors: string;
    write_errors: string;
    checksum_errors: string;
    slow_ios?: string;
    
    // Device path information (for leaf devices)
    path?: string;
    phys_path?: string;
    devid?: string;
    
    // Scan information
    scan_processed?: string;
    
    // Child vdevs for hierarchical structure
    children?: ZfsVdevModel[];
}

export interface ZfsPoolResponseWrapper {
    success: boolean;
    errors?: string[];
    data: ZfsPoolModel[];
}

export interface ZfsPoolDetailResponseWrapper {
    success: boolean;
    errors?: string[];
    data: ZfsPoolModel;
}

export interface ZfsDatasetModel {
    id: number;
    name: string;
    host_id: string;
    type: string; // FILESYSTEM, VOLUME, SNAPSHOT
    pool: string;
    createtxg: string;
    
    // Timestamps
    created_at?: string;
    updated_at?: string;
    
    // Space information
    used: string;
    available: string;
    referenced: string;
    mountpoint: string;
}

export interface ZfsDatasetResponseWrapper {
    success: boolean;
    errors?: string[];
    data: ZfsDatasetModel[];
}