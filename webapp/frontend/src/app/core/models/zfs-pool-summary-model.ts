// maps to API response for ZFS pool summary
import { ZFSPoolModel } from './zfs-pool-model';

export interface ZFSPoolSummaryModel {
    pool: ZFSPoolModel;
}

export interface ZFSPoolSummaryResponseWrapper {
    success: boolean;
    data: {
        summary: { [guid: string]: ZFSPoolSummaryModel };
    };
}

export interface ZFSPoolDetailsResponseWrapper {
    success: boolean;
    data: {
        pool: ZFSPoolModel;
        metrics_history: ZFSPoolMetricsHistoryModel[];
    };
}

export interface ZFSPoolMetricsHistoryModel {
    date: string;
    size: number;
    allocated: number;
    free: number;
    capacity_percent: number;
    fragmentation: number;
    status: string;
}
