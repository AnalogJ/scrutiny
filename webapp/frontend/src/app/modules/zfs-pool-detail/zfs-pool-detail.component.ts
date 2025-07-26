import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { ZfsPoolModel, ZfsVdevModel } from '../../core/models/zfs-pool-model';
import { ZfsService } from '../dashboard/zfs.service';

interface VdevTreeNode extends ZfsVdevModel {
    indent: number;
    displayName: string;
}

@Component({
    selector: 'zfs-pool-detail',
    templateUrl: './zfs-pool-detail.component.html',
    styleUrls: ['./zfs-pool-detail.component.scss']
})

export class ZfsPoolDetailComponent implements OnInit, OnDestroy {
    pool: ZfsPoolModel;

    private _unsubscribeAll: Subject<void>;

    constructor(
        private _activatedRoute: ActivatedRoute,
        private _router: Router,
        public zfsService: ZfsService
    ) {
        this._unsubscribeAll = new Subject<void>();
    }

    ngOnInit(): void {
        // Get the pool details from resolver
        this._activatedRoute.data
            .pipe(takeUntil(this._unsubscribeAll))
            .subscribe((data) => {
                if (data.poolDetail && data.poolDetail.success) {
                    this.pool = data.poolDetail.data;
                }
            });
    }

    ngOnDestroy(): void {
        this._unsubscribeAll.next();
        this._unsubscribeAll.complete();
    }

    /**
     * Navigate back to dashboard
     */
    goBack(): void {
        this._router.navigate(['/dashboard']);
    }

    /**
     * Get background color class for pool/vdev state (dark mode compatible)
     */
    getStateBackgroundClass(state: string): string {
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
     * Get capacity percentage (accounts for RAID overhead to show effective utilization)
     */
    getCapacityPercentage(): number {
        if (!this.pool?.alloc_space || !this.pool?.total_space || !this.pool?.def_space) return 0;
        
        const allocBytes = this.parseBytes(this.pool.alloc_space);
        const totalBytes = this.parseBytes(this.pool.total_space);
        const freeBytes = this.parseBytes(this.pool.def_space);
        
        if (totalBytes === 0) return 0;
        
        // Calculate the usable capacity by determining the effective space
        // The usable space is approximately: total - free = used, but we need to account for parity
        const usableCapacity = this.calculateUsableCapacity(totalBytes);
        const usableAllocated = this.calculateUsableSpace(allocBytes);
        
        if (usableCapacity === 0) return 0;
        return Math.round((usableAllocated / usableCapacity) * 100);
    }
    
    /**
     * Calculate usable capacity accounting for RAID overhead
     */
    private calculateUsableCapacity(totalBytes: number): number {
        const raidConfig = this.getRaidConfiguration();
        
        if (!raidConfig) {
            // No RAID, return total capacity
            return totalBytes;
        }
        
        // Calculate usable capacity based on RAID type and disk count
        switch (raidConfig.type) {
            case 'raidz':
                // RAID-Z1: (n-1)/n usable where n is number of disks
                return totalBytes * (raidConfig.diskCount - 1) / raidConfig.diskCount;
            case 'raidz2':
                // RAID-Z2: (n-2)/n usable where n is number of disks
                return totalBytes * (raidConfig.diskCount - 2) / raidConfig.diskCount;
            case 'raidz3':
                // RAID-Z3: (n-3)/n usable where n is number of disks
                return totalBytes * (raidConfig.diskCount - 3) / raidConfig.diskCount;
            case 'mirror':
                // Mirror: 50% usable
                return totalBytes * 0.5;
            default:
                return totalBytes;
        }
    }
    
    /**
     * Calculate usable allocated space accounting for RAID overhead
     */
    private calculateUsableSpace(allocatedBytes: number): number {
        const raidConfig = this.getRaidConfiguration();
        
        if (!raidConfig) {
            return allocatedBytes;
        }
        
        // For allocated space, we apply the same ratio as usable capacity
        switch (raidConfig.type) {
            case 'raidz':
                return allocatedBytes * (raidConfig.diskCount - 1) / raidConfig.diskCount;
            case 'raidz2':
                return allocatedBytes * (raidConfig.diskCount - 2) / raidConfig.diskCount;
            case 'raidz3':
                return allocatedBytes * (raidConfig.diskCount - 3) / raidConfig.diskCount;
            case 'mirror':
                return allocatedBytes * 0.5;
            default:
                return allocatedBytes;
        }
    }
    
    /**
     * Get RAID configuration from vdevs
     */
    private getRaidConfiguration(): { type: string; diskCount: number } | null {
        if (!this.pool?.vdevs) return null;
        
        // Find the first RAID vdev configuration
        for (const vdev of this.pool.vdevs) {
            if (vdev.vdev_type.includes('raidz') || vdev.vdev_type === 'mirror') {
                // Count child disks for this vdev
                const childDisks = this.pool.vdevs.filter(child => child.parent_id === vdev.id);
                const diskCount = childDisks.length || 2; // Default to 2 if no children found
                
                return {
                    type: vdev.vdev_type,
                    diskCount: diskCount
                };
            }
        }
        
        return null; // No RAID configuration found
    }

    /**
     * Get capacity color class based on usage percentage
     */
    getCapacityColorClass(): string {
        const percentage = this.getCapacityPercentage();
        if (percentage >= 95) return 'bg-red-500';
        if (percentage >= 85) return 'bg-orange-500';
        if (percentage >= 70) return 'bg-yellow-500';
        return 'bg-green-500';
    }

    /**
     * Get vdev capacity percentage
     */
    getVdevCapacityPercentage(vdev: ZfsVdevModel): number {
        if (!vdev.alloc_space || !vdev.total_space) return 0;
        
        const allocBytes = this.parseBytes(vdev.alloc_space);
        const totalBytes = this.parseBytes(vdev.total_space);
        
        if (totalBytes === 0) return 0;
        return Math.round((allocBytes / totalBytes) * 100);
    }

    /**
     * Get vdev type background class (dark mode compatible)
     */
    getVdevTypeClass(vdevType: string): string {
        switch (vdevType?.toLowerCase()) {
            case 'raidz':
            case 'raidz2':
            case 'raidz3':
                return 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300';
            case 'mirror':
                return 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300';
            case 'disk':
                return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300';
            default:
                return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300';
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
            case 'replacing':
                return 'heroicons_outline:arrow-path';
            case 'disk':
                return 'heroicons_outline:circle-stack';
            default:
                return 'heroicons_outline:server';
        }
    }

    /**
     * Get indent lines for tree structure
     */
    getIndentLines(indent: number): string[] {
        const lines: string[] = [];
        for (let i = 0; i < indent - 1; i++) {
            lines.push('line');
        }
        if (indent > 0) {
            lines.push('corner');
        }
        return lines;
    }

    /**
     * Parse bytes from string format (e.g., "12.8T", "65.5G")
     */
    private parseBytes(bytesStr: string): number {
        if (!bytesStr || bytesStr === '0' || bytesStr === '') return 0;
        
        const match = bytesStr.match(/^(\d+\.?\d*)\s*([KMGTPE]?)B?$/i);
        if (!match) return 0;
        
        const value = parseFloat(match[1]);
        const unit = match[2].toUpperCase();
        
        const multipliers: { [key: string]: number } = {
            '': 1,
            'K': 1024,
            'M': 1024 ** 2,
            'G': 1024 ** 3,
            'T': 1024 ** 4,
            'P': 1024 ** 5,
            'E': 1024 ** 6
        };
        
        return value * (multipliers[unit] || 1);
    }

    /**
     * Format scan information similar to zpool status
     */
    formatScanInfo(): string {
        if (!this.pool.scan_function || !this.pool.scan_state) return '';
        
        let scanInfo = `${this.pool.scan_function}`;
        
        if (this.pool.scan_processed) {
            scanInfo += ` ${this.pool.scan_processed}`;
        }
        
        if (this.pool.scan_end_time) {
            scanInfo += ` on ${this.pool.scan_end_time}`;
        }
        
        if (this.pool.scan_errors && this.pool.scan_errors !== '0') {
            scanInfo += ` with ${this.pool.scan_errors} errors`;
        } else {
            scanInfo += ' with 0 errors';
        }
        
        return scanInfo;
    }

    /**
     * Build hierarchical vdev tree for display: pool → RAID type → disks/replacing → individual disks
     */
    getVdevTree(): VdevTreeNode[] {
        if (!this.pool.vdevs) return [];
        
        const treeNodes: VdevTreeNode[] = [];
        
        // Create a map for easier lookup
        const vdevMap = new Map<string, ZfsVdevModel>();
        this.pool.vdevs.forEach(vdev => {
            vdevMap.set(vdev.guid, vdev);
        });
        
        // Find top-level vdevs (raidz, mirror, or standalone disks)
        const topLevelVdevs = this.pool.vdevs.filter(vdev => 
            vdev.vdev_type !== 'root' && 
            (vdev.vdev_type.includes('raidz') || vdev.vdev_type === 'mirror' || 
             vdev.vdev_type === 'disk')
        );
        
        // Sort by name to ensure consistent ordering
        topLevelVdevs.sort((a, b) => (a.name || '').localeCompare(b.name || ''));
        
        // Group vdevs: find the main group vdevs first
        const mainGroups = topLevelVdevs.filter(vdev => 
            vdev.vdev_type.includes('raidz') || vdev.vdev_type === 'mirror'
        );
        
        const standaloneDisks = topLevelVdevs.filter(vdev => 
            vdev.vdev_type === 'disk' && 
            !vdev.parent_id // Top-level disks have no parent
        );
        
        // Process main groups (raidz, mirror)
        mainGroups.forEach(group => {
            // Add the main group
            treeNodes.push({
                ...group,
                indent: 1,
                displayName: group.name || `${group.vdev_type}`
            });
            
            // Find all disks and special vdevs that belong to this group (by parent_id)
            const groupMembers = this.pool.vdevs.filter(vdev => 
                vdev.parent_id === group.id && (
                    vdev.vdev_type === 'disk' || 
                    vdev.vdev_type === 'replacing' || 
                    vdev.vdev_type === 'spare'
                )
            );
            
            // Sort group members by name
            groupMembers.sort((a, b) => (a.name || '').localeCompare(b.name || ''));
            
            groupMembers.forEach(member => {
                if (member.vdev_type === 'replacing') {
                    // Add replacing vdev
                    treeNodes.push({
                        ...member,
                        indent: 2,
                        displayName: member.name || 'replacing'
                    });
                    
                    // Find disks that belong to this replacing vdev (by parent_id)
                    const replacingDisks = this.pool.vdevs.filter(vdev => 
                        vdev.parent_id === member.id && vdev.vdev_type === 'disk'
                    );
                    
                    // Show disks under this replacing vdev
                    replacingDisks.forEach(disk => {
                        treeNodes.push({
                            ...disk,
                            indent: 3,
                            displayName: this.getVdevDisplayName(disk)
                        });
                    });
                    
                } else if (member.vdev_type === 'disk') {
                    // Regular disk member of the group
                    treeNodes.push({
                        ...member,
                        indent: 2,
                        displayName: this.getVdevDisplayName(member)
                    });
                }
            });
        });
        
        // Add standalone disks (for simple pools without raidz/mirror)
        standaloneDisks.forEach(disk => {
            treeNodes.push({
                ...disk,
                indent: 1,
                displayName: this.getVdevDisplayName(disk)
            });
        });
        
        return treeNodes;
    }
    
    /**
     * Check if a disk is a child of any of the groups
     */
    private isChildOfGroup(vdev: ZfsVdevModel, groups: ZfsVdevModel[]): boolean {
        // This is a simplified check - in reality, ZFS relationships are complex
        // For now, we assume disks with specific names are part of groups
        return false; // Let the main logic handle grouping
    }

    /**
     * Get display name for vdev (use path if available, otherwise name)
     */
    private getVdevDisplayName(vdev: ZfsVdevModel): string {
        if (vdev.path && vdev.path.trim()) {
            return vdev.path;
        }
        return vdev.name || vdev.guid;
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
        
        if (minutes < 1) return 'Just now';
        if (minutes < 60) return `${minutes} minute${minutes !== 1 ? 's' : ''} ago`;
        if (hours < 24) return `${hours} hour${hours !== 1 ? 's' : ''} ago`;
        return `${days} day${days !== 1 ? 's' : ''} ago`;
    }
    
    /**
     * Get data age color class based on staleness
     */
    getDataAgeColorClass(): string {
        if (!this.pool?.updated_at) return 'text-gray-500';
        
        const updatedTime = new Date(this.pool.updated_at);
        const now = new Date();
        const diffMs = now.getTime() - updatedTime.getTime();
        const minutes = Math.floor(diffMs / (1000 * 60));
        
        if (minutes < 10) return 'text-green-600 dark:text-green-400';
        if (minutes < 30) return 'text-yellow-600 dark:text-yellow-400';
        return 'text-red-600 dark:text-red-400';
    }
    
    /**
     * Check if data is stale (older than 30 minutes)
     */
    isDataStale(): boolean {
        if (!this.pool?.updated_at) return true;
        
        const updatedTime = new Date(this.pool.updated_at);
        const now = new Date();
        const diffMs = now.getTime() - updatedTime.getTime();
        const minutes = Math.floor(diffMs / (1000 * 60));
        
        return minutes > 30;
    }

    /**
     * Check if the pool has any errors
     */
    hasErrors(): boolean {
        return (
            (this.pool.read_errors && this.pool.read_errors !== '0') ||
            (this.pool.write_errors && this.pool.write_errors !== '0') ||
            (this.pool.checksum_errors && this.pool.checksum_errors !== '0') ||
            (this.pool.error_count && this.pool.error_count !== '0')
        );
    }

    /**
     * Get error summary text
     */
    getErrorSummary(): string {
        const errors = [];
        
        if (this.pool.read_errors && this.pool.read_errors !== '0') {
            errors.push(`${this.pool.read_errors} read errors`);
        }
        
        if (this.pool.write_errors && this.pool.write_errors !== '0') {
            errors.push(`${this.pool.write_errors} write errors`);
        }
        
        if (this.pool.checksum_errors && this.pool.checksum_errors !== '0') {
            errors.push(`${this.pool.checksum_errors} checksum errors`);
        }
        
        return errors.length > 0 ? errors.join(', ') : 'No known data errors';
    }
}