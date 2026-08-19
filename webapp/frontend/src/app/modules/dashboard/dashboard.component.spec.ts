import {MatDialog} from '@angular/material/dialog';
import {Router} from '@angular/router';
import {DashboardComponent} from './dashboard.component';
import {DashboardService} from './dashboard.service';
import {ScrutinyConfigService} from 'app/core/config/scrutiny-config.service';
import {DeviceSummaryModel} from 'app/core/models/device-summary-model';

describe('DashboardComponent', () => {
    let component: DashboardComponent;

    const deviceSummary = (scrutinyUuid: string, hostId: string, archived: boolean): DeviceSummaryModel => ({
        device: {
            scrutiny_uuid: scrutinyUuid,
            host_id: hostId,
            archived,
        },
    } as DeviceSummaryModel)

    beforeEach(() => {
        component = new DashboardComponent(
            jasmine.createSpyObj('DashboardService', ['getSummaryData', 'getSummaryTempData']) as DashboardService,
            jasmine.createSpyObj('ScrutinyConfigService', ['config']) as ScrutinyConfigService,
            jasmine.createSpyObj('MatDialog', ['open']) as MatDialog,
            jasmine.createSpyObj('Router', ['navigate']) as Router,
        );
    });

    describe('#hostGroupHasVisibleDevices()', () => {

        it('should be false when every device on the host is archived', () => {
            component.summaryData = {
                'uuid-1': deviceSummary('uuid-1', 'my-host', true),
                'uuid-2': deviceSummary('uuid-2', 'my-host', true),
            }

            expect(component.hostGroupHasVisibleDevices(['uuid-1', 'uuid-2'])).toBeFalse();
        });

        it('should be true when at least one device on the host is not archived', () => {
            component.summaryData = {
                'uuid-1': deviceSummary('uuid-1', 'my-host', true),
                'uuid-2': deviceSummary('uuid-2', 'my-host', false),
            }

            expect(component.hostGroupHasVisibleDevices(['uuid-1', 'uuid-2'])).toBeTrue();
        });

        it('should be true when every device is archived but archived devices are shown', () => {
            component.summaryData = {
                'uuid-1': deviceSummary('uuid-1', 'my-host', true),
            }
            component.showArchived = true;

            expect(component.hostGroupHasVisibleDevices(['uuid-1'])).toBeTrue();
        });

        it('should be false when every device on the host has been deleted', () => {
            component.summaryData = {
                'uuid-1': deviceSummary('uuid-1', 'my-host', false),
            }

            component.onDeviceDeleted('uuid-1')

            expect(component.hostGroupHasVisibleDevices(['uuid-1'])).toBeFalse();
        });
    });
});
