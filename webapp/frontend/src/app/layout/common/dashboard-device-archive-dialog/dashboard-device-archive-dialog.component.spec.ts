import {ComponentFixture, TestBed} from '@angular/core/testing';

import {DashboardDeviceArchiveDialogComponent} from './dashboard-device-archive-dialog.component';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import {MAT_DIALOG_DATA, MatDialogModule as MatDialogModule, MatDialogRef as MatDialogRef} from '@angular/material/dialog';
import {MatButtonModule as MatButtonModule} from '@angular/material/button';
import {MatIconModule} from '@angular/material/icon';
import {SharedModule} from '../../../shared/shared.module';
import {DashboardDeviceArchiveDialogService} from './dashboard-device-archive-dialog.service';
import {of} from 'rxjs';


describe('DashboardDeviceArchiveDialogComponent', () => {
    let component: DashboardDeviceArchiveDialogComponent;
    let fixture: ComponentFixture<DashboardDeviceArchiveDialogComponent>;

    const matDialogRefSpy = jasmine.createSpyObj('MatDialogRef', ['closeDialog', 'close']);
    const dashboardDeviceArchiveDialogServiceSpy = jasmine.createSpyObj('DashboardDeviceArchiveDialogService', ['archiveDevice']);

    beforeEach(() => {
        TestBed.configureTestingModule({
    declarations: [DashboardDeviceArchiveDialogComponent],
    imports: [MatDialogModule,
        MatButtonModule,
        MatIconModule,
        SharedModule],
    providers: [
        { provide: MatDialogRef, useValue: matDialogRefSpy },
        { provide: MAT_DIALOG_DATA, useValue: { wwn: 'test-wwn', title: 'my-test-device-title' } },
        { provide: DashboardDeviceArchiveDialogService, useValue: dashboardDeviceArchiveDialogServiceSpy },
        provideHttpClient(withInterceptorsFromDi())
    ]
})
            .compileComponents()
    });

    beforeEach(() => {
        fixture = TestBed.createComponent(DashboardDeviceArchiveDialogComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });

    it('should close the component if cancel is clicked', () => {
        matDialogRefSpy.closeDialog.calls.reset();
        matDialogRefSpy.closeDialog()
        expect(matDialogRefSpy.closeDialog).toHaveBeenCalled();
    });

    it('should attempt to archive device if archive is clicked', () => {
        dashboardDeviceArchiveDialogServiceSpy.archiveDevice.and.returnValue(of({'success': true}));

        component.onArchiveClick()
        expect(dashboardDeviceArchiveDialogServiceSpy.archiveDevice).toHaveBeenCalledWith('test-wwn');
        expect(dashboardDeviceArchiveDialogServiceSpy.archiveDevice.calls.count())
            .withContext('one call')
            .toBe(1);
    });
});
