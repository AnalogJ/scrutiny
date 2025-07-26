import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ZfsPoolResponseWrapper, ZfsPoolDetailResponseWrapper, ZfsPoolModel, ZfsDatasetResponseWrapper, ZfsDatasetModel } from '../../core/models/zfs-pool-model';
import { getBasePath } from '../../app.routing';

@Injectable({
    providedIn: 'root'
})
export class ZfsService {

    constructor(private _httpClient: HttpClient) {
    }

    /**
     * Get all ZFS pools
     */
    getZfsPools(hostId?: string): Observable<ZfsPoolResponseWrapper> {
        const url = `${getBasePath()}/api/zfs/pools`;
        const params = hostId ? { host_id: hostId } : {};
        return this._httpClient.get<ZfsPoolResponseWrapper>(url, { params });
    }

    /**
     * Get specific ZFS pool details
     */
    getZfsPoolDetails(poolGuid: string): Observable<ZfsPoolDetailResponseWrapper> {
        const url = `${getBasePath()}/api/zfs/pool/${poolGuid}`;
        return this._httpClient.get<ZfsPoolDetailResponseWrapper>(url);
    }

    /**
     * Get ZFS pool status for display
     */
    getPoolStatusClass(state: string): string {
        switch (state?.toUpperCase()) {
            case 'ONLINE':
                return 'text-green-600 dark:text-green-400';
            case 'DEGRADED':
                return 'text-yellow-600 dark:text-yellow-400';
            case 'FAULTED':
            case 'OFFLINE':
            case 'UNAVAIL':
                return 'text-red-600 dark:text-red-400';
            default:
                return 'text-gray-600 dark:text-gray-400';
        }
    }

    /**
     * Get ZFS pool status icon
     */
    getPoolStatusIcon(state: string): string {
        switch (state?.toUpperCase()) {
            case 'ONLINE':
                return 'heroicons_outline:check-circle';
            case 'DEGRADED':
                return 'heroicons_outline:exclamation-circle';
            case 'FAULTED':
            case 'OFFLINE':
            case 'UNAVAIL':
                return 'heroicons_outline:exclamation-circle';
            default:
                return 'heroicons_outline:question-mark-circle';
        }
    }

    /**
     * Format bytes to human readable format
     */
    formatBytes(bytes: string): string {
        if (!bytes || bytes === '0' || bytes === '') return '0 B';
        
        // Handle values like "12.8T", "65.5T", etc.
        const match = bytes.match(/^(\d+\.?\d*)\s*([KMGTPE]?)B?$/i);
        if (match) {
            const value = parseFloat(match[1]);
            const unit = match[2].toUpperCase();
            
            const units = ['', 'K', 'M', 'G', 'T', 'P', 'E'];
            const index = units.indexOf(unit);
            
            if (index >= 0) {
                return `${value}${unit}B`;
            }
        }
        
        return bytes;
    }

    /**
     * Get vdev type display name
     */
    getVdevTypeDisplayName(vdevType: string): string {
        switch (vdevType?.toLowerCase()) {
            case 'raidz':
                return 'RAID-Z1';
            case 'raidz2':
                return 'RAID-Z2';
            case 'raidz3':
                return 'RAID-Z3';
            case 'mirror':
                return 'Mirror';
            case 'disk':
                return 'Single Disk';
            case 'file':
                return 'File';
            case 'root':
                return 'Root';
            default:
                return vdevType || 'Unknown';
        }
    }

    /**
     * Get vdev icon based on type
     */
    getVdevIcon(vdevType: string): string {
        switch (vdevType?.toLowerCase()) {
            case 'replacing':
                return 'heroicons_outline:exclamation-circle';
            default:
                return 'heroicons_outline:check-circle';
        }
    }

    /**
     * Get ZFS datasets for a pool
     */
    getZfsDatasets(poolName?: string, hostId?: string): Observable<ZfsDatasetResponseWrapper> {
        const url = `${getBasePath()}/api/zfs/datasets`;
        const params: any = {};
        if (poolName) params.pool = poolName;
        if (hostId) params.host_id = hostId;
        return this._httpClient.get<ZfsDatasetResponseWrapper>(url, { params });
    }
}