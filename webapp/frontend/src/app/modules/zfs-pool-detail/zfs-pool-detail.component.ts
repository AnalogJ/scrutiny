import { Component, OnInit, OnDestroy } from "@angular/core";
import { ActivatedRoute, Router } from "@angular/router";
import { Subject } from "rxjs";
import { takeUntil } from "rxjs/operators";
import { ZfsPoolModel, ZfsVdevModel, ZfsDatasetModel } from "../../core/models/zfs-pool-model";
import { ZfsService } from "../dashboard/zfs.service";

interface VdevTreeNode extends ZfsVdevModel {
    indent: number;
    displayName: string;
}

@Component({
    selector: "zfs-pool-detail",
    templateUrl: "./zfs-pool-detail.component.html",
    styleUrls: ["./zfs-pool-detail.component.scss"],
})
export class ZfsPoolDetailComponent implements OnInit, OnDestroy {
    pool: ZfsPoolModel;
    datasets: ZfsDatasetModel[] = [];

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
                    this.loadDatasets();
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
        this._router.navigate(["/dashboard"]);
    }

    /**
     * Load ZFS datasets for this pool
     */
    loadDatasets(): void {
        if (!this.pool?.name) return;

        this.zfsService.getZfsDatasets(this.pool.name)
            .pipe(takeUntil(this._unsubscribeAll))
            .subscribe(
                (response) => {
                    if (response.success) {
                        this.datasets = response.data;
                    }
                },
                (error) => {
                    console.error('Error loading datasets:', error);
                }
            );
    }

    /**
     * Get background color class for pool/vdev state (dark mode compatible)
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
     * Get capacity percentage from zpool list data
     */
    getCapacityPercentage(): number {
        // Use zpool list capacity percentage (most accurate)
        if (this.pool?.capacity_percent) {
            const match = this.pool.capacity_percent.match(/^(\d+)%?$/);
            if (match) {
                return parseInt(match[1], 10);
            }
        }

        // Calculate from zpool list allocated/size
        if (this.pool?.allocated && this.pool?.size) {
            const allocBytes = this.parseBytes(this.pool.allocated);
            const sizeBytes = this.parseBytes(this.pool.size);
            if (sizeBytes === 0) return 0;
            return Math.round((allocBytes / sizeBytes) * 100);
        }

        return 0;
    }


    /**
     * Get formatted used space from zpool list data
     */
    getUsedSpaceFromPool(): string {
        return this.pool?.allocated || "0B";
    }

    /**
     * Get formatted available space from zpool list data
     */
    getAvailableSpaceFromPool(): string {
        return this.pool?.free || "0B";
    }

    /**
     * Get formatted total space from zpool list data
     */
    getTotalSpaceFromPool(): string {
        return this.pool?.size || "0B";
    }

    /**
     * Get capacity color class based on usage percentage
     */
    getCapacityColorClass(): string {
        const percentage = this.getCapacityPercentage();
        if (percentage >= 95) return "bg-red-500";
        if (percentage >= 85) return "bg-orange-500";
        if (percentage >= 70) return "bg-yellow-500";
        return "bg-green-500";
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
            case "raidz":
            case "raidz2":
            case "raidz3":
                return "bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300";
            case "mirror":
                return "bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300";
            case "disk":
                return "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300";
            default:
                return "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300";
        }
    }


    /**
     * Get indent lines for tree structure
     */
    getIndentLines(indent: number): string[] {
        const lines: string[] = [];
        for (let i = 0; i < indent - 1; i++) {
            lines.push("line");
        }
        if (indent > 0) {
            lines.push("corner");
        }
        return lines;
    }

    /**
     * Parse bytes from string format (e.g., "12.8T", "65.5G")
     */
    private parseBytes(bytesStr: string): number {
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
     * Format scan information similar to zpool status
     */
    formatScanInfo(): string {
        if (!this.pool.scan_function || !this.pool.scan_state) return "";

        let scanInfo = `${this.pool.scan_function}`;

        if (this.pool.scan_processed) {
            scanInfo += ` ${this.pool.scan_processed}`;
        }

        if (this.pool.scan_end_time) {
            scanInfo += ` on ${this.pool.scan_end_time}`;
        }

        if (this.pool.scan_errors && this.pool.scan_errors !== "0") {
            scanInfo += ` with ${this.pool.scan_errors} errors`;
        } else {
            scanInfo += " with 0 errors";
        }

        return scanInfo;
    }

    /**
     * Build hierarchical vdev tree for display: pool → RAID type → disks/replacing → individual disks
     */
    getVdevTree(): VdevTreeNode[] {
        if (!this.pool.vdevs) return [];

        const treeNodes: VdevTreeNode[] = [];

        // Find RAID groups (raidz, mirror) - these are the top level containers
        const raidGroups = this.pool.vdevs.filter(
            (vdev) =>
                vdev.vdev_type.includes("raidz") || vdev.vdev_type === "mirror"
        );

        // Sort by name to ensure consistent ordering
        raidGroups.sort((a, b) => (a.name || "").localeCompare(b.name || ""));

        // Process each RAID group
        raidGroups.forEach((group) => {
            // Add the RAID group
            treeNodes.push({
                ...group,
                indent: 1,
                displayName: group.name || group.vdev_type,
            });

            // Find replacing vdevs (these are special containers for disk replacement)
            const replacingVdevs = this.pool.vdevs.filter(
                (vdev) => vdev.vdev_type === "replacing"
            );

            // Find regular disks (not in replacing operations)
            const regularDisks = this.pool.vdevs.filter(
                (vdev) => vdev.vdev_type === "disk" || vdev.vdev_type === "file"
            );

            // Add replacing vdevs under the RAID group
            replacingVdevs.forEach((replacing) => {
                treeNodes.push({
                    ...replacing,
                    indent: 2,
                    displayName: replacing.name || "replacing",
                });

                // Add disks that are part of the replacing operation
                // In your data, these would be the file and disk that are being replaced
                const replacingDisks = regularDisks.filter((disk) => {
                    // The offline file and the online disk that's replacing it
                    return (
                        disk.state === "OFFLINE" ||
                        (disk.state === "ONLINE" && disk.scan_processed)
                    );
                });

                replacingDisks.forEach((disk) => {
                    treeNodes.push({
                        ...disk,
                        indent: 3,
                        displayName: this.getVdevDisplayName(disk),
                    });
                });
            });

            // Add regular disks that are not part of replacing operations
            // Only exclude disks if there's actually a replacing operation happening
            const hasReplacing = replacingVdevs.length > 0;
            const nonReplacingDisks = regularDisks.filter((disk) => {
                if (!hasReplacing) {
                    // No replacing operations, show all regular disks
                    return true;
                }
                // If there are replacing operations, exclude offline disks and actively replacing disks
                return (
                    disk.state !== "OFFLINE" &&
                    !(
                        disk.state === "ONLINE" &&
                        disk.scan_processed &&
                        replacingVdevs.some((r) => r.scan_processed)
                    )
                );
            });

            nonReplacingDisks.forEach((disk) => {
                treeNodes.push({
                    ...disk,
                    indent: 2,
                    displayName: this.getVdevDisplayName(disk),
                });
            });
        });

        // Handle pools without RAID groups (simple disk pools)
        if (raidGroups.length === 0) {
            const standaloneDisks = this.pool.vdevs.filter(
                (vdev) => vdev.vdev_type === "disk" || vdev.vdev_type === "file"
            );

            standaloneDisks.forEach((disk) => {
                treeNodes.push({
                    ...disk,
                    indent: 1,
                    displayName: this.getVdevDisplayName(disk),
                });
            });
        }

        return treeNodes;
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
        if (!this.pool?.updated_at) return "Unknown";

        const updatedTime = new Date(this.pool.updated_at);
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
    getDataAgeColorClass(): string {
        if (!this.pool?.updated_at) return "text-gray-500";

        const updatedTime = new Date(this.pool.updated_at);
        const now = new Date();
        const diffMs = now.getTime() - updatedTime.getTime();
        const minutes = Math.floor(diffMs / (1000 * 60));

        if (minutes < 10) return "text-green-600 dark:text-green-400";
        if (minutes < 30) return "text-yellow-600 dark:text-yellow-400";
        return "text-red-600 dark:text-red-400";
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
            (this.pool.read_errors && this.pool.read_errors !== "0") ||
            (this.pool.write_errors && this.pool.write_errors !== "0") ||
            (this.pool.checksum_errors && this.pool.checksum_errors !== "0") ||
            (this.pool.error_count && this.pool.error_count !== "0")
        );
    }

    /**
     * Get error summary text
     */
    getErrorSummary(): string {
        const errors = [];

        if (this.pool.read_errors && this.pool.read_errors !== "0") {
            errors.push(`${this.pool.read_errors} read errors`);
        }

        if (this.pool.write_errors && this.pool.write_errors !== "0") {
            errors.push(`${this.pool.write_errors} write errors`);
        }

        if (this.pool.checksum_errors && this.pool.checksum_errors !== "0") {
            errors.push(`${this.pool.checksum_errors} checksum errors`);
        }

        return errors.length > 0 ? errors.join(", ") : "No known data errors";
    }

    /**
     * Get CSS classes for the required action message
     */
    getActionMessageClass(): string {
        const isNoAction = this.pool?.action
            ?.toLowerCase()
            .includes("no action required");

        if (isNoAction) {
            return "text-green-600 dark:text-green-300 bg-green-50 dark:bg-green-900 border-green-200 dark:border-green-700";
        } else {
            return "text-orange-600 dark:text-orange-300 bg-orange-50 dark:bg-orange-900 border-orange-200 dark:border-orange-700";
        }
    }

}
