import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DetailSettingsComponent } from './detail-settings.component';

describe('DetailSettingsComponent', () => {
  let component: DetailSettingsComponent;
  let fixture: ComponentFixture<DetailSettingsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DetailSettingsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DetailSettingsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
