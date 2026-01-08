import {
    ChangeDetectionStrategy,
    Component,
    OnDestroy,
    OnInit,
    ViewEncapsulation
} from '@angular/core';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { ZFSPoolsService } from 'app/modules/zfs-pools/zfs-pools.service';
import { AppConfig } from 'app/core/config/app.config';
import { ScrutinyConfigService } from 'app/core/config/scrutiny-config.service';
import { Router } from '@angular/router';
import { ZFSPoolSummaryModel } from 'app/core/models/zfs-pool-summary-model';

@Component({
    selector: 'zfs-pools',
    templateUrl: './zfs-pools.component.html',
    styleUrls: ['./zfs-pools.component.scss'],
    encapsulation: ViewEncapsulation.None,
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class ZFSPoolsComponent implements OnInit, OnDestroy {
    summaryData: { [guid: string]: ZFSPoolSummaryModel };
    hostGroups: { [hostId: string]: string[] } = {};
    config: AppConfig;
    showArchived: boolean;

    private _unsubscribeAll: Subject<void>;

    constructor(
        private _zfsPoolsService: ZFSPoolsService,
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
                const oldConfig = JSON.stringify(this.config);
                const newConfig = JSON.stringify(config);

                if (oldConfig !== newConfig) {
                    this.config = config;
                    if (oldConfig) {
                        this.refreshComponent();
                    }
                }
            });

        // Get the data
        this._zfsPoolsService.data$
            .pipe(takeUntil(this._unsubscribeAll))
            .subscribe((data) => {
                this.summaryData = data;

                // Generate group data by host
                this.hostGroups = {};
                for (const guid in this.summaryData) {
                    const hostId = this.summaryData[guid].pool.host_id;
                    const hostPoolList = this.hostGroups[hostId] || [];
                    hostPoolList.push(guid);
                    this.hostGroups[hostId] = hostPoolList;
                }
            });
    }

    ngOnDestroy(): void {
        this._unsubscribeAll.next();
        this._unsubscribeAll.complete();
    }

    private refreshComponent(): void {
        const currentUrl = this.router.url;
        this.router.routeReuseStrategy.shouldReuseRoute = () => false;
        this.router.onSameUrlNavigation = 'reload';
        this.router.navigate([currentUrl]);
    }

    poolSummariesForHostGroup(hostGroupGUIDs: string[]): ZFSPoolSummaryModel[] {
        const poolSummaries: ZFSPoolSummaryModel[] = [];
        for (const guid of hostGroupGUIDs) {
            if (this.summaryData[guid]) {
                poolSummaries.push(this.summaryData[guid]);
            }
        }
        return poolSummaries;
    }

    onPoolDeleted(guid: string): void {
        delete this.summaryData[guid];
    }

    onPoolArchived(guid: string): void {
        this.summaryData[guid].pool.archived = true;
    }

    onPoolUnarchived(guid: string): void {
        this.summaryData[guid].pool.archived = false;
    }

    trackByFn(index: number, item: any): any {
        return item.id || index;
    }
}
