import {ComponentFixture, TestBed, waitForAsync} from '@angular/core/testing';
import {FormsModule} from '@angular/forms';
import {HttpClientTestingModule, HttpTestingController} from '@angular/common/http/testing';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material/dialog';
import {NO_ERRORS_SCHEMA} from '@angular/core';

import {DetailSettingsComponent} from './detail-settings.component';

describe('DetailSettingsComponent', () => {
  let component: DetailSettingsComponent;
  let fixture: ComponentFixture<DetailSettingsComponent>;
  let http: HttpTestingController;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ DetailSettingsComponent ],
      imports: [ FormsModule, HttpClientTestingModule ],
      providers: [
        { provide: MAT_DIALOG_DATA, useValue: { scrutiny_uuid: 'abc', friendly_name: 'Backup Disk' } },
        { provide: MatDialogRef, useValue: { close: () => {} } }
      ],
      schemas: [ NO_ERRORS_SCHEMA ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DetailSettingsComponent);
    component = fixture.componentInstance;
    http = TestBed.inject(HttpTestingController);
    fixture.detectChanges();
  });

  it('should prefill the friendly name from the device without a request', () => {
    expect(component.friendlyName).toBe('Backup Disk');
    http.verify();
  });

  it('should save the friendly name', () => {
    component.friendlyName = 'Data Disk';
    component.saveSettings();
    const req = http.expectOne(r => r.url.endsWith('/api/device/abc/friendly_name'));
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual({ friendly_name: 'Data Disk' });
  });
});
