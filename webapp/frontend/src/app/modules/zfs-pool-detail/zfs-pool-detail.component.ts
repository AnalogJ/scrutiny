import { Component, OnInit, OnDestroy, ViewChild } from "@angular/core";
import { ActivatedRoute, Router } from "@angular/router";
import { Subject } from "rxjs";
import { takeUntil } from "rxjs/operators";
import { ApexOptions, ChartComponent } from 'ng-apexcharts';
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
    datasetUtilizationOptions: ApexOptions;

    private _unsubscribeAll: Subject<void>;
    @ViewChild('datasetChart', { static: false }) datasetChart: ChartComponent;

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
                    this._prepareDatasetUtilizationChart();
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
                        this._prepareDatasetUtilizationChart();
                    }
                },
                (error) => {
                    console.error('Error loading datasets:', error);
                }
            );
    }


    /**
     * Get capacity percentage from zpool list data
     */
    getCapacityPercentage(): number | null {
        if (this.pool?.capacity_percent) {
            const match = this.pool.capacity_percent.match(/^(\d+)%?$/);
            if (match) {
                return parseInt(match[1], 10);
            }
        }
        return null;
    }

    /**
     * Get formatted capacity percentage with fallback
     */
    getFormattedCapacityPercentage(): string {
        const percentage = this.getCapacityPercentage();
        return percentage !== null ? `${percentage}%` : "N/A";
    }


    /**
     * Get formatted used space from zpool list data
     */
    getUsedSpaceFromPool(): string {
        return this.pool?.allocated || "N/A";
    }

    /**
     * Get formatted available space from zpool list data
     */
    getAvailableSpaceFromPool(): string {
        return this.pool?.free || "N/A";
    }

    /**
     * Get formatted total space from zpool list data
     */
    getTotalSpaceFromPool(): string {
        return this.pool?.size || "N/A";
    }

    /**
     * Get formatted fragmentation with fallback
     */
    getFormattedFragmentation(): string {
        return this.pool?.fragmentation || "N/A";
    }


    /**
     * Get formatted deduplication ratio with fallback
     */
    getFormattedDedupratio(): string {
        return this.pool?.dedupratio || "N/A";
    }

    /**
     * Get capacity color class based on usage percentage
     */
    getCapacityColorClass(): string {
        const percentage = this.getCapacityPercentage();
        if (percentage === null) return "bg-gray-400";
        if (percentage >= 95) return "bg-red-500";
        if (percentage >= 85) return "bg-orange-500";
        if (percentage >= 70) return "bg-yellow-500";
        return "bg-green-500";
    }

    /**
     * Get vdev capacity percentage
     */
    getVdevCapacityPercentage(vdev: ZfsVdevModel): string {
        if (!vdev.alloc_space || !vdev.total_space) return "N/A";

        const allocBytes = this.zfsService.parseBytes(vdev.alloc_space);
        const totalBytes = this.zfsService.parseBytes(vdev.total_space);

        if (totalBytes === 0) return "N/A";
        const percentage = Math.round((allocBytes / totalBytes) * 100);
        return `${percentage}%`;
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
                        displayName: this.zfsService.getVdevDisplayName(disk),
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
                    displayName: this.zfsService.getVdevDisplayName(disk),
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
                    displayName: this.zfsService.getVdevDisplayName(disk),
                });
            });
        }

        return treeNodes;
    }



    /**
     * Get formatted data age string
     */
    getDataAge(): string {
        return this.zfsService.getDataAge(this.pool?.updated_at || '');
    }

    /**
     * Get data age color class based on staleness
     */
    getDataAgeColorClass(): string {
        return this.zfsService.getDataAgeColorClass(this.pool?.updated_at || '');
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

    /**
     * Prepare dataset utilization chart data
     */
    private _prepareDatasetUtilizationChart(): void {
        if (!this.datasets || this.datasets.length === 0) {
            return;
        }

        const series = this._generateDatasetUtilizationSeries();

        this.datasetUtilizationOptions = {
            chart: {
                animations: {
                    speed: 400,
                    animateGradually: {
                        enabled: false
                    }
                },
                fontFamily: 'inherit',
                foreColor: 'inherit',
                width: '100%',
                height: '100%',
                type: 'area',
                sparkline: {
                    enabled: true
                }
            },
            colors: ['#667eea', '#9066ea', '#66c0ea', '#66ead2', '#d266ea', '#66ea90', '#ea6690', '#90ea66'],
            fill: {
                colors: ['#b2bef4', '#c7b2f4', '#b2dff4', '#b2f4e8', '#e8b2f4', '#b2f4c7', '#f4b2c7', '#c7f4b2'],
                opacity: 0.5,
                type: 'gradient'
            },
            series: series,
            stroke: {
                curve: 'smooth',
                width: 2
            },
            tooltip: {
                theme: 'dark',
                x: {
                    format: 'MMM dd, yyyy HH:mm:ss'
                },
                y: {
                    formatter: (value: number) => {
                        return this.zfsService.formatBytes(value.toString());
                    }
                }
            },
            xaxis: {
                type: 'datetime'
            },
            yaxis: {
                labels: {
                    formatter: (value: number) => {
                        return this.zfsService.formatBytes(value.toString());
                    }
                }
            }
        };
    }

    /**
     * Generate dataset utilization series data
     * For now, this creates mock historical data based on current usage
     * In a real implementation, this would fetch historical data from the backend
     */
    private _generateDatasetUtilizationSeries(): any[] {
        if (!this.datasets || this.datasets.length === 0) {
            return [];
        }

        const series = [];
        const now = new Date();
        
        // Generate data for the last 30 days
        const dataPoints = 30;
        
        this.datasets.forEach((dataset) => {
            if (!dataset.used || dataset.used === 'N/A') return;
            
            const currentUsed = this.zfsService.parseBytes(dataset.used);
            if (currentUsed === 0) return;
            
            const datasetSeries = {
                name: dataset.name,
                data: []
            };
            
            // Generate mock historical data points
            for (let i = dataPoints - 1; i >= 0; i--) {
                const date = new Date(now.getTime() - (i * 24 * 60 * 60 * 1000));
                
                // Generate realistic usage progression (gradual increase with some variation)
                const progressFactor = (dataPoints - i) / dataPoints;
                const variationFactor = 0.9 + (Math.random() * 0.2); // ±10% variation
                const historicalUsed = Math.floor(currentUsed * progressFactor * variationFactor);
                
                datasetSeries.data.push({
                    x: date.getTime(),
                    y: historicalUsed
                });
            }
            
            series.push(datasetSeries);
        });
        
        return series;
    }

}
