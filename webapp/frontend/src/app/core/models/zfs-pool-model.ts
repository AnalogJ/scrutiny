// maps to webapp/backend/pkg/models/zfs_pool.go
export interface ZFSPoolModel {
    guid: string;
    name: string;
    host_id: string;
    label: string;
    archived: boolean;
    muted: boolean;

    status: ZFSPoolStatus;
    size: number;
    allocated: number;
    free: number;
    fragmentation: number;
    capacity_percent: number;

    scrub_state: ZFSScrubState;
    scrub_start: string;
    scrub_end: string;
    scrub_percent: number;
    scrub_errors: number;

    total_read_errors: number;
    total_write_errors: number;
    total_checksum_errors: number;

    vdevs?: ZFSVdevModel[];

    created_at: string;
    updated_at: string;
}

export type ZFSPoolStatus = 'ONLINE' | 'DEGRADED' | 'FAULTED' | 'OFFLINE' | 'REMOVED' | 'UNAVAIL';
export type ZFSScrubState = 'none' | 'scanning' | 'finished' | 'canceled';

export interface ZFSVdevModel {
    id: number;
    pool_guid: string;
    parent_id?: number;

    name: string;
    type: ZFSVdevType;
    status: string;
    path: string;

    read_errors: number;
    write_errors: number;
    checksum_errors: number;

    children?: ZFSVdevModel[];

    created_at: string;
    updated_at: string;
}

export type ZFSVdevType = 'disk' | 'file' | 'mirror' | 'raidz1' | 'raidz2' | 'raidz3' | 'spare' | 'log' | 'cache' | 'special' | 'dedup';
