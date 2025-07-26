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
     * Get hierarchical vdev tree for dashboard display
     */
    getVdevTree(): any[] {
        if (!this.pool.vdevs) return [];
        
        const treeNodes: any[] = [];
        
        // Find RAID groups (raidz, mirror) - these are the top level containers
        const raidGroups = this.pool.vdevs.filter(vdev => 
            vdev.vdev_type.includes('raidz') || vdev.vdev_type === 'mirror'
        );
        
        // Sort by name to ensure consistent ordering
        raidGroups.sort((a, b) => (a.name || '').localeCompare(b.name || ''));
        
        // Process each RAID group
        raidGroups.forEach(group => {
            // Add the RAID group
            treeNodes.push({
                ...group,
                indent: 1,
                displayName: group.name || group.vdev_type,
                isGroup: true
            });
            
            // Find replacing vdevs (these are special containers for disk replacement)
            const replacingVdevs = this.pool.vdevs.filter(vdev => vdev.vdev_type === 'replacing');
            
            // Find regular disks (not in replacing operations)
            const regularDisks = this.pool.vdevs.filter(vdev => 
                vdev.vdev_type === 'disk' || vdev.vdev_type === 'file'
            );
            
            // Add replacing vdevs under the RAID group
            replacingVdevs.forEach(replacing => {
                treeNodes.push({
                    ...replacing,
                    indent: 2,
                    displayName: replacing.name || 'replacing',
                    isGroup: true
                });
                
                // Add disks that are part of the replacing operation
                // In your data, these would be the file and disk that are being replaced
                const replacingDisks = regularDisks.filter(disk => {
                    // The offline file and the online disk that's replacing it
                    return disk.state === 'OFFLINE' || 
                           (disk.state === 'ONLINE' && disk.scan_processed);
                });
                
                replacingDisks.forEach(disk => {
                    treeNodes.push({
                        ...disk,
                        indent: 3,
                        displayName: this.getVdevDisplayName(disk),
                        isGroup: false
                    });
                });
            });
            
            // Add regular disks that are not part of replacing operations
            // Only exclude disks if there's actually a replacing operation happening
            const hasReplacing = replacingVdevs.length > 0;
            const nonReplacingDisks = regularDisks.filter(disk => {
                if (!hasReplacing) {
                    // No replacing operations, show all regular disks
                    return true;
                }
                // If there are replacing operations, exclude offline disks and actively replacing disks
                return disk.state !== 'OFFLINE' && 
                       !(disk.state === 'ONLINE' && disk.scan_processed && 
                         replacingVdevs.some(r => r.scan_processed));
            });
            
            nonReplacingDisks.forEach(disk => {
                treeNodes.push({
                    ...disk,
                    indent: 2,
                    displayName: this.getVdevDisplayName(disk),
                    isGroup: false
                });
            });
        });
        
        // Handle pools without RAID groups (simple disk pools)
        if (raidGroups.length === 0) {
            const standaloneDisks = this.pool.vdevs.filter(vdev => 
                vdev.vdev_type === 'disk' || vdev.vdev_type === 'file'
            );
            
            standaloneDisks.forEach(disk => {
                treeNodes.push({
                    ...disk,
                    indent: 1,
                    displayName: this.getVdevDisplayName(disk),
                    isGroup: false
                });
            });
        }
        
        return treeNodes;
    }

    /**
     * Get vdev display name (use path if available, otherwise name)
     */
    private getVdevDisplayName(vdev: ZfsVdevModel): string {
        if (vdev.path && vdev.path.trim()) {
            return vdev.path;
        }
        return vdev.name || vdev.guid;
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