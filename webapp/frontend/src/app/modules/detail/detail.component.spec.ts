import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientModule } from '@angular/common/http';
import { DetailComponent } from './detail.component';

import {TreoConfigService} from '@treo/services/config';
import { TREO_APP_CONFIG } from '@treo/services/config/config.constants';
const TREO_APP_CONFIG_PROVIDER = [ { provide: TREO_APP_CONFIG, useValue: TreoConfigService } ];
import { MatDialogModule } from '@angular/material/dialog';
import {DeviceTitlePipe} from 'app/shared/device-title.pipe';


describe('DetailComponent', () => {
  let component: DetailComponent;
  let fixture: ComponentFixture<DetailComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
        imports: [
            HttpClientModule,
            MatDialogModule

        ],
        declarations: [ DetailComponent, DeviceTitlePipe ],
        providers: [ TREO_APP_CONFIG_PROVIDER ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
