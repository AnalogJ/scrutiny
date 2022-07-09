import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {DashboardDeviceDeleteDialogComponent} from './dashboard-device-delete-dialog.component';
import {HttpClientModule} from '@angular/common/http';
import {MAT_DIALOG_DATA, MatDialogModule, MatDialogRef} from '@angular/material/dialog';
import {MatButtonModule} from '@angular/material/button';
import {MatIconModule} from '@angular/material/icon';
import {SharedModule} from '../../../shared/shared.module';
import {DashboardDeviceDeleteDialogService} from './dashboard-device-delete-dialog.service';
import {of} from 'rxjs';


describe('DashboardDeviceDeleteDialogComponent', () => {
    let component: DashboardDeviceDeleteDialogComponent;
    let fixture: ComponentFixture<DashboardDeviceDeleteDialogComponent>;

    const matDialogRefSpy = jasmine.createSpyObj('MatDialogRef', ['closeDialog', 'close']);
    const dashboardDeviceDeleteDialogServiceSpy = jasmine.createSpyObj('DashboardDeviceDeleteDialogService', ['deleteDevice']);

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            imports: [
                HttpClientModule,
                MatDialogModule,
                MatButtonModule,
                MatIconModule,
                SharedModule,
            ],
            providers: [
                {provide: MatDialogRef, useValue: matDialogRefSpy},
                {provide: MAT_DIALOG_DATA, useValue: {wwn: 'test-wwn', title: 'my-test-device-title'}},
                {provide: DashboardDeviceDeleteDialogService, useValue: dashboardDeviceDeleteDialogServiceSpy}
            ],
            declarations: [DashboardDeviceDeleteDialogComponent]
        })
            .compileComponents()
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(DashboardDeviceDeleteDialogComponent);
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

    it('should attempt to delete device if delete is clicked', () => {
        dashboardDeviceDeleteDialogServiceSpy.deleteDevice.and.returnValue(of({'success': true}));

        component.onDeleteClick()
        expect(dashboardDeviceDeleteDialogServiceSpy.deleteDevice).toHaveBeenCalledWith('test-wwn');
        expect(dashboardDeviceDeleteDialogServiceSpy.deleteDevice.calls.count())
            .withContext('one call')
            .toBe(1);
    });
});
