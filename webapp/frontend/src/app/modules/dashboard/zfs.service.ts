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

    /**
     * Parse bytes from string format (e.g., "12.8T", "65.5G")
     */
    parseBytes(bytesStr: string): number {
        if (!bytesStr || bytesStr === "0" || bytesStr === "") return 0;

        const match = bytesStr.match(/^(\d+\.?\d*)\s*([KMGTPE]?)B?$/i);
        if (!match) return 0;

        const value = parseFloat(match[1]);
        const unit = match[2].toUpperCase();

        const multipliers: { [key: string]: number } = {
            "": 1,
            K: 1024,
            M: 1024 ** 2,
            G: 1024 ** 3,
            T: 1024 ** 4,
            P: 1024 ** 5,
            E: 1024 ** 6,
        };

        return value * (multipliers[unit] || 1);
    }

    /**
     * Get vdev display name (use path if available, otherwise name)
     */
    getVdevDisplayName(vdev: any): string {
        if (vdev.path && vdev.path.trim()) {
            return vdev.path;
        }
        return vdev.name || vdev.guid;
    }

    /**
     * Get state background class for badges (dark mode compatible)
     */
    getStateBackgroundClass(state: string): string {
        switch (state?.toUpperCase()) {
            case "ONLINE":
                return "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300";
            case "DEGRADED":
                return "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300";
            case "FAULTED":
            case "OFFLINE":
            case "UNAVAIL":
                return "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300";
            default:
                return "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300";
        }
    }

    /**
     * Get formatted data age string
     */
    getDataAge(updatedAt: string): string {
        if (!updatedAt) return "Unknown";

        const updatedTime = new Date(updatedAt);
        const now = new Date();
        const diffMs = now.getTime() - updatedTime.getTime();

        const minutes = Math.floor(diffMs / (1000 * 60));
        const hours = Math.floor(diffMs / (1000 * 60 * 60));
        const days = Math.floor(diffMs / (1000 * 60 * 60 * 24));

        if (minutes < 1) return "Just now";
        if (minutes < 60)
            return `${minutes} minute${minutes !== 1 ? "s" : ""} ago`;
        if (hours < 24) return `${hours} hour${hours !== 1 ? "s" : ""} ago`;
        return `${days} day${days !== 1 ? "s" : ""} ago`;
    }

    /**
     * Get data age color class based on staleness
     */
    getDataAgeColorClass(updatedAt: string): string {
        if (!updatedAt) return "text-gray-500";

        const updatedTime = new Date(updatedAt);
        const now = new Date();
        const diffMs = now.getTime() - updatedTime.getTime();
        const minutes = Math.floor(diffMs / (1000 * 60));

        if (minutes < 10) return "text-green-600 dark:text-green-400";
        if (minutes < 30) return "text-yellow-600 dark:text-yellow-400";
        return "text-red-600 dark:text-red-400";
    }
}