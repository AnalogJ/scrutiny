import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import moment from 'moment';
import { Subject } from 'rxjs';
import { MatDialog } from '@angular/material/dialog';
import { ZFSPoolModel, ZFSPoolStatus } from 'app/core/models/zfs-pool-model';
import { AppConfig } from 'app/core/config/app.config';
import { ZFSPoolsService } from 'app/modules/zfs-pools/zfs-pools.service';

@Component({
    selector: 'app-zfs-pool-card',
    templateUrl: './zfs-pool-card.component.html',
    styleUrls: ['./zfs-pool-card.component.scss'],
    standalone: false,
})
export class ZFSPoolCardComponent implements OnInit {
    constructor(private _zfsPoolsService: ZFSPoolsService, public dialog: MatDialog) {
        this._unsubscribeAll = new Subject();
    }

    @Input() poolSummary: ZFSPoolModel;
    @Input() config: AppConfig;
    @Output() poolArchived = new EventEmitter<string>();
    @Output() poolUnarchived = new EventEmitter<string>();
    @Output() poolDeleted = new EventEmitter<string>();

    private _unsubscribeAll: Subject<void>;

    ngOnInit(): void {}

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    getPoolStatus(pool: ZFSPoolModel): 'passed' | 'failed' | 'unknown' {
        if (!pool) {
            return 'unknown';
        }
        switch (pool.status) {
            case 'ONLINE':
                return 'passed';
            case 'DEGRADED':
            case 'FAULTED':
                return 'failed';
            default:
                return 'unknown';
        }
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

    classPoolLastUpdatedOn(pool: ZFSPoolModel): string {
        const poolStatus = this.getPoolStatus(pool);
        if (poolStatus === 'failed') {
            return 'text-red';
        } else if (poolStatus === 'passed') {
            if (moment().subtract(14, 'days').isBefore(pool.updated_at)) {
                return 'text-green';
            } else if (moment().subtract(1, 'months').isBefore(pool.updated_at)) {
                return 'text-yellow';
            } else {
                return 'text-red';
            }
        } else {
            return '';
        }
    }

    getPoolTitle(pool: ZFSPoolModel): string {
        if (pool.label) {
            return pool.label;
        }
        return pool.name;
    }

    getCapacityPercentClass(percent: number): string {
        if (percent >= 90) {
            return 'bg-red-500';
        } else if (percent >= 80) {
            return 'bg-yellow-500';
        } else {
            return 'bg-green-500';
        }
    }

    getScrubStatusText(pool: ZFSPoolModel): string {
        switch (pool.scrub_state) {
            case 'none':
                return 'Never';
            case 'scanning':
                return `In Progress (${pool.scrub_percent}%)`;
            case 'finished':
                return moment(pool.scrub_end).fromNow();
            case 'canceled':
                return 'Canceled';
            default:
                return 'Unknown';
        }
    }

    archivePool(): void {
        if (this.poolSummary.archived) {
            this._zfsPoolsService.unarchivePool(this.poolSummary.guid).subscribe(() => {
                this.poolUnarchived.emit(this.poolSummary.guid);
            });
        } else {
            this._zfsPoolsService.archivePool(this.poolSummary.guid).subscribe(() => {
                this.poolArchived.emit(this.poolSummary.guid);
            });
        }
    }

    deletePool(): void {
        if (confirm(`Are you sure you want to delete pool "${this.getPoolTitle(this.poolSummary)}"?`)) {
            this._zfsPoolsService.deletePool(this.poolSummary.guid).subscribe(() => {
                this.poolDeleted.emit(this.poolSummary.guid);
            });
        }
    }

    getTotalErrors(pool: ZFSPoolModel): number {
        return pool.total_read_errors + pool.total_write_errors + pool.total_checksum_errors;
    }
}
