import { Component, Input, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ZfsPoolModel } from '../../../core/models/zfs-pool-model';
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
                        displayName: this.zfsService.getVdevDisplayName(disk),
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
                    displayName: this.zfsService.getVdevDisplayName(disk),
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
                    displayName: this.zfsService.getVdevDisplayName(disk),
                    isGroup: false
                });
            });
        }
        
        return treeNodes;
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
        return this.zfsService.getDataAge(this.pool?.updated_at || '');
    }
    
    /**
     * Get timestamp color class based on data freshness
     */
    getTimestampColorClass(): string {
        return this.zfsService.getDataAgeColorClass(this.pool?.updated_at || '');
    }

    /**
     * Get utilization percentage from zpool list capacity
     */
    getUtilizationPercentage(): number {
        if (this.pool?.capacity_percent) {
            const match = this.pool.capacity_percent.match(/^(\d+)%?$/);
            if (match) {
                return parseInt(match[1], 10);
            }
        }
        return 0;
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