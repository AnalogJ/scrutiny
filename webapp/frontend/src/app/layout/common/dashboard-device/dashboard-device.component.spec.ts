import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DashboardDeviceComponent } from './dashboard-device.component';

describe('DashboardDeviceComponent', () => {
  let component: DashboardDeviceComponent;
  let fixture: ComponentFixture<DashboardDeviceComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DashboardDeviceComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DashboardDeviceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
