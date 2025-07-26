import { Component, Input, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ZfsPoolModel, ZfsVdevModel } from '../../../core/models/zfs-pool-model';
import { ZfsService } from '../../../modules/dashboard/zfs.service';

@Component({
    selector: 'app-dashboard-zfs-pool',
    templateUrl: './dashboard-zfs-pool.component.html'
})
export class DashboardZfsPoolComponent implements OnInit {
    @Input() pool: ZfsPoolModel;

    constructor(
        public zfsService: ZfsService,
        private _router: Router
    ) {}

    ngOnInit(): void {
    }

    /**
     * Check if the pool has any errors
     */
    hasErrors(): boolean {
        return (
            (this.pool.read_errors && this.pool.read_errors !== '0') ||
            (this.pool.write_errors && this.pool.write_errors !== '0') ||
            (this.pool.checksum_errors && this.pool.checksum_errors !== '0')
        );
    }

    /**
     * Get top-level vdevs (RAID groups and individual disks)
     */
    getTopLevelVdevs(): ZfsVdevModel[] {
        if (!this.pool.vdevs) return [];
        
        // Find root vdev first
        const rootVdev = this.pool.vdevs.find(vdev => vdev.vdev_type === 'root');
        
        if (rootVdev) {
            // Return direct children of root
            return this.pool.vdevs.filter(vdev => 
                vdev.parent_id === rootVdev.id && vdev.vdev_type !== 'root'
            );
        } else {
            // Fallback: return vdevs without parent_id (top-level)
            return this.pool.vdevs.filter(vdev => 
                vdev.vdev_type !== 'root' && 
                (!vdev.parent_id || vdev.parent_id === null || vdev.parent_id === 0)
            );
        }
    }

    /**
     * Get single disk vdevs (for pools without raidz/mirror)
     */
    getSingleDiskVdevs(): ZfsVdevModel[] {
        if (!this.pool.vdevs) return [];
        
        // Find root vdev first
        const rootVdev = this.pool.vdevs.find(vdev => vdev.vdev_type === 'root');
        
        if (rootVdev) {
            // Return direct disk children of root
            return this.pool.vdevs.filter(vdev => 
                vdev.vdev_type === 'disk' && vdev.parent_id === rootVdev.id
            );
        } else {
            // Fallback: return top-level disks
            return this.pool.vdevs.filter(vdev => 
                vdev.vdev_type === 'disk' && 
                (!vdev.parent_id || vdev.parent_id === null || vdev.parent_id === 0)
            );
        }
    }

    /**
     * Get child vdevs for a given parent
     */
    getChildVdevs(parentId: number): ZfsVdevModel[] {
        if (!this.pool.vdevs) return [];
        
        return this.pool.vdevs.filter(vdev => vdev.parent_id === parentId);
    }

    /**
     * Get device name with fallback
     */
    getDeviceName(vdev: ZfsVdevModel): string {
        if (vdev.path && vdev.path.trim()) {
            // Extract just the device name from the path
            const parts = vdev.path.split('/');
            return parts[parts.length - 1];
        }
        return vdev.name || vdev.guid;
    }

    /**
     * Get state class for compact display (dark mode compatible)
     */
    getStateClass(state: string): string {
        switch (state?.toUpperCase()) {
            case 'ONLINE':
                return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300';
            case 'DEGRADED':
                return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300';
            case 'FAULTED':
            case 'OFFLINE':
            case 'UNAVAIL':
                return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300';
            default:
                return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300';
        }
    }

    /**
     * Get vdev icon based on type
     */
    getVdevIcon(vdevType: string): string {
        switch (vdevType?.toLowerCase()) {
            case 'raidz':
            case 'raidz2':
            case 'raidz3':
                return 'heroicons_outline:server-stack';
            case 'mirror':
                return 'heroicons_outline:rectangle-stack';
            case 'disk':
                return 'heroicons_outline:circle-stack';
            default:
                return 'heroicons_outline:server';
        }
    }

    /**
     * Get formatted scan status
     */
    getScanStatus(): string {
        if (!this.pool.scan_function) return 'No scan';
        
        const scanType = this.pool.scan_function.charAt(0).toUpperCase() + this.pool.scan_function.slice(1).toLowerCase();
        const scanState = this.pool.scan_state ? this.pool.scan_state.toLowerCase() : '';
        
        if (scanState === 'finished') {
            return `${scanType} completed`;
        } else if (scanState === 'scanning') {
            return `${scanType} in progress`;
        } else if (scanState === 'canceled') {
            return `${scanType} canceled`;
        }
        
        return `${scanType}`;
    }
    
    /**
     * Get scan status color class
     */
    getScanStatusClass(): string {
        if (!this.pool.scan_state) return '';
        
        switch (this.pool.scan_state.toLowerCase()) {
            case 'finished':
                return 'text-green-600 dark:text-green-400';
            case 'scanning':
                return 'text-blue-600 dark:text-blue-400';
            case 'canceled':
                return 'text-yellow-600 dark:text-yellow-400';
            default:
                return 'text-gray-600 dark:text-gray-400';
        }
    }
    
    /**
     * Get formatted data age string
     */
    getDataAge(): string {
        if (!this.pool?.updated_at) return 'Unknown';
        
        const updatedTime = new Date(this.pool.updated_at);
        const now = new Date();
        const diffMs = now.getTime() - updatedTime.getTime();
        
        const minutes = Math.floor(diffMs / (1000 * 60));
        const hours = Math.floor(diffMs / (1000 * 60 * 60));
        const days = Math.floor(diffMs / (1000 * 60 * 60 * 24));
        
        if (minutes < 1) return 'just now';
        if (minutes < 60) return `${minutes}m ago`;
        if (hours < 24) return `${hours}h ago`;
        return `${days}d ago`;
    }
    
    /**
     * Get timestamp color class based on data freshness (matching disk component style)
     */
    getTimestampColorClass(): string {
        if (!this.pool?.updated_at) return 'text-gray-600 dark:text-gray-400';
        
        const updatedTime = new Date(this.pool.updated_at);
        const now = new Date();
        const diffMs = now.getTime() - updatedTime.getTime();
        const minutes = Math.floor(diffMs / (1000 * 60));
        
        // Match the styling logic from the disk component
        if (minutes < 10) return 'text-gray-600 dark:text-gray-400'; // Fresh data
        if (minutes < 30) return 'text-yellow-600 dark:text-yellow-400'; // Moderately stale
        return 'text-red-600 dark:text-red-400'; // Stale data
    }

    /**
     * Get utilization percentage
     */
    getUtilizationPercentage(): number {
        if (!this.pool.total_space || parseInt(this.pool.total_space) === 0) return 0;
        return Math.round((parseInt(this.pool.alloc_space) / parseInt(this.pool.total_space)) * 100);
    }

    /**
     * Get utilization color class based on percentage
     */
    getUtilizationClass(): string {
        const percentage = this.getUtilizationPercentage();
        
        if (percentage < 70) {
            return 'text-green-600 dark:text-green-400';
        } else if (percentage < 85) {
            return 'text-yellow-600 dark:text-yellow-400';
        } else {
            return 'text-red-600 dark:text-red-400';
        }
    }

    /**
     * Navigate to pool details page
     */
    viewPoolDetails(): void {
        this._router.navigate(['/zfs-pool', this.pool.pool_guid]);
    }
}