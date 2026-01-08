import {
    ChangeDetectionStrategy,
    Component,
    OnDestroy,
    OnInit,
    ViewEncapsulation
} from '@angular/core';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { ApexOptions } from 'ng-apexcharts';
import { ZFSPoolDetailService } from 'app/modules/zfs-pool-detail/zfs-pool-detail.service';
import { AppConfig } from 'app/core/config/app.config';
import { ScrutinyConfigService } from 'app/core/config/scrutiny-config.service';
import { Router } from '@angular/router';
import { ZFSPoolModel, ZFSPoolStatus, ZFSVdevModel } from 'app/core/models/zfs-pool-model';
import { ZFSPoolMetricsHistoryModel } from 'app/core/models/zfs-pool-summary-model';

@Component({
    selector: 'zfs-pool-detail',
    templateUrl: './zfs-pool-detail.component.html',
    styleUrls: ['./zfs-pool-detail.component.scss'],
    encapsulation: ViewEncapsulation.None,
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class ZFSPoolDetailComponent implements OnInit, OnDestroy {
    pool: ZFSPoolModel;
    metricsHistory: ZFSPoolMetricsHistoryModel[];
    capacityOptions: ApexOptions;
    config: AppConfig;

    private _unsubscribeAll: Subject<void>;

    constructor(
        private _zfsPoolDetailService: ZFSPoolDetailService,
        private _configService: ScrutinyConfigService,
        private router: Router,
    ) {
        this._unsubscribeAll = new Subject();
    }

    ngOnInit(): void {
        // Subscribe to config changes
        this._configService.config$
            .pipe(takeUntil(this._unsubscribeAll))
            .subscribe((config: AppConfig) => {
                this.config = config;
            });

        // Get the data
        this._zfsPoolDetailService.data$
            .pipe(takeUntil(this._unsubscribeAll))
            .subscribe((data) => {
                if (data) {
                    this.pool = data.data.pool;
                    this.metricsHistory = data.data.metrics_history;
                    this._prepareChartData();
                }
            });
    }

    ngOnDestroy(): void {
        this._unsubscribeAll.next();
        this._unsubscribeAll.complete();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Private methods
    // -----------------------------------------------------------------------------------------------------

    private _prepareChartData(): void {
        if (!this.metricsHistory || this.metricsHistory.length === 0) {
            return;
        }

        const capacityData = this.metricsHistory.map(m => ({
            x: new Date(m.date),
            y: m.capacity_percent
        }));

        this.capacityOptions = {
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
            colors: ['#667eea'],
            fill: {
                colors: ['#b2bef4'],
                opacity: 0.5,
                type: 'gradient'
            },
            series: [{
                name: 'Capacity',
                data: capacityData
            }],
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
                    formatter: (value) => `${value}%`
                }
            },
            xaxis: {
                type: 'datetime',
                labels: {
                    datetimeUTC: false
                }
            },
            yaxis: {
                min: 0,
                max: 100
            }
        };
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    getPoolTitle(): string {
        if (this.pool?.label) {
            return this.pool.label;
        }
        return this.pool?.name || 'Unknown Pool';
    }

    getStatusColorClass(status: ZFSPoolStatus): string {
        switch (status) {
            case 'ONLINE':
                return 'text-green';
            case 'DEGRADED':
                return 'text-yellow';
            case 'FAULTED':
            case 'UNAVAIL':
            case 'OFFLINE':
            case 'REMOVED':
                return 'text-red';
            default:
                return '';
        }
    }

    getVdevIcon(type: string): string {
        switch (type) {
            case 'disk':
            case 'file':
                return 'heroicons_outline:server';
            case 'mirror':
                return 'heroicons_outline:document-duplicate';
            case 'raidz1':
            case 'raidz2':
            case 'raidz3':
                return 'heroicons_outline:database';
            case 'spare':
                return 'heroicons_outline:shield-check';
            case 'log':
                return 'heroicons_outline:document-text';
            case 'cache':
                return 'heroicons_outline:lightning-bolt';
            default:
                return 'heroicons_outline:server';
        }
    }

    getVdevStatusClass(status: string): string {
        switch (status) {
            case 'ONLINE':
                return 'text-green';
            case 'DEGRADED':
                return 'text-yellow';
            case 'FAULTED':
            case 'UNAVAIL':
            case 'OFFLINE':
            case 'REMOVED':
                return 'text-red';
            default:
                return '';
        }
    }

    hasErrors(vdev: ZFSVdevModel): boolean {
        return vdev.read_errors > 0 || vdev.write_errors > 0 || vdev.checksum_errors > 0;
    }

    toggleMuted(): void {
        const newMutedState = !this.pool.muted;
        this._zfsPoolDetailService.setMuted(this.pool.guid, newMutedState).subscribe(() => {
            this.pool.muted = newMutedState;
        });
    }

    goBack(): void {
        this.router.navigate(['/zfs-pools']);
    }

    trackByFn(index: number, item: any): any {
        return item.id || index;
    }
}
