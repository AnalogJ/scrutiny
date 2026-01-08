import {ComponentFixture, TestBed} from '@angular/core/testing';

import {DashboardDeviceComponent} from './dashboard-device.component';
import {MatDialog as MatDialog} from '@angular/material/dialog';
import {MatButtonModule as MatButtonModule} from '@angular/material/button';
import {MatIconModule} from '@angular/material/icon';
import {SharedModule} from 'app/shared/shared.module';
import {MatMenuModule as MatMenuModule} from '@angular/material/menu';
import {TREO_APP_CONFIG} from '@treo/services/config/config.constants';
import {DeviceSummaryModel} from 'app/core/models/device-summary-model';
import moment from 'moment';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { HttpClient, provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import {ScrutinyConfigService} from 'app/core/config/scrutiny-config.service';
import {of} from 'rxjs';
import {MetricsStatusThreshold} from 'app/core/config/app.config';

describe('DashboardDeviceComponent', () => {
    let component: DashboardDeviceComponent;
    let fixture: ComponentFixture<DashboardDeviceComponent>;

    const matDialogSpy = jasmine.createSpyObj('MatDialog', ['open']);
    // const configServiceSpy = jasmine.createSpyObj('ScrutinyConfigService', ['config$']);
    let configService: ScrutinyConfigService;
    let httpClientSpy: jasmine.SpyObj<HttpClient>;

    beforeEach(() => {

        httpClientSpy = jasmine.createSpyObj('HttpClient', ['get']);
        configService = new ScrutinyConfigService(httpClientSpy, {});

        TestBed.configureTestingModule({
    declarations: [DashboardDeviceComponent],
    imports: [MatButtonModule,
        MatIconModule,
        MatMenuModule,
        SharedModule],
    providers: [
        { provide: MatDialog, useValue: matDialogSpy },
        { provide: TREO_APP_CONFIG, useValue: { dashboard_display: 'name', metrics: { status_threshold: 3 } } },
        { provide: ScrutinyConfigService, useValue: configService },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting()
    ]
})
            .compileComponents();
    });

    beforeEach(() => {
        // configServiceSpy.config$.and.returnValue(of({'success': true});
        fixture = TestBed.createComponent(DashboardDeviceComponent);
        component = fixture.componentInstance;
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });

    describe('#classDeviceLastUpdatedOn()', () => {

        it('if non-zero device status, should be red', () => {
            httpClientSpy.get.and.returnValue(of({
                settings: {
                    metrics: {
                        status_threshold: MetricsStatusThreshold.Both,
                    }
                }
            }));
            component.ngOnInit()
            // component.deviceSummary = summary.data.summary['0x5000c500673e6b5f'] as DeviceSummaryModel
            expect(component.classDeviceLastUpdatedOn({
                device: {
                    device_status: 2,
                },
                smart: {
                    collector_date: moment().subtract(13, 'days').toISOString()
                },
            } as DeviceSummaryModel)).toBe('text-red')
        });

        it('if non-zero device status, should be red', () => {
            httpClientSpy.get.and.returnValue(of({
                settings: {
                    metrics: {
                        status_threshold: MetricsStatusThreshold.Both,
                    }
                }
            }));
            component.ngOnInit()
            expect(component.classDeviceLastUpdatedOn({
                device: {
                    device_status: 2
                },
                smart: {
                    collector_date: moment().subtract(13, 'days').toISOString()
                },
            } as DeviceSummaryModel)).toBe('text-red')
        });

        it('if healthy device status and updated in the last two weeks, should be green', () => {
            httpClientSpy.get.and.returnValue(of({
                settings: {
                    metrics: {
                        status_threshold: MetricsStatusThreshold.Both,
                    }
                }
            }));
            component.ngOnInit()
            expect(component.classDeviceLastUpdatedOn({
                device: {
                    device_status: 0
                },
                smart: {
                    collector_date: moment().subtract(13, 'days').toISOString()
                }
            } as DeviceSummaryModel)).toBe('text-green')
        });

        it('if healthy device status and updated more than two weeks ago, but less than 1 month, should be yellow', () => {
            httpClientSpy.get.and.returnValue(of({
                settings: {
                    metrics: {
                        status_threshold: MetricsStatusThreshold.Both,
                    }
                }
            }));
            component.ngOnInit()
            expect(component.classDeviceLastUpdatedOn({
                device: {
                    device_status: 0
                },
                smart: {
                    collector_date: moment().subtract(3, 'weeks').toISOString()
                }
            } as DeviceSummaryModel)).toBe('text-yellow')
        });

        it('if healthy device status and updated more 1 month ago, should be red', () => {
            httpClientSpy.get.and.returnValue(of({
                settings: {
                    metrics: {
                        status_threshold: MetricsStatusThreshold.Both,
                    }
                }
            }));
            component.ngOnInit()
            expect(component.classDeviceLastUpdatedOn({
                device: {
                    device_status: 0
                },
                smart: {
                    collector_date: moment().subtract(5, 'weeks').toISOString()
                }
            } as DeviceSummaryModel)).toBe('text-red')
        });
    })
});
