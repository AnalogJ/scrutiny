import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DashboardDeviceDeleteDialogComponent } from './dashboard-device-delete-dialog.component';

describe('DashboardDeviceDeleteDialogComponent', () => {
  let component: DashboardDeviceDeleteDialogComponent;
  let fixture: ComponentFixture<DashboardDeviceDeleteDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DashboardDeviceDeleteDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DashboardDeviceDeleteDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
