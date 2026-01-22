import { ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { DetailComponent } from './detail.component';
import { DetailService } from './detail.service';
import { ScrutinyConfigService } from 'app/core/config/scrutiny-config.service';
import { MatDialog } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule } from '@angular/material/sort';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatInputModule } from '@angular/material/input';
import { FormsModule } from '@angular/forms';
import { LOCALE_ID } from '@angular/core';
import { of, Subject } from 'rxjs';
import { SmartModel } from 'app/core/models/measurements/smart-model';
import { SmartAttributeModel } from 'app/core/models/measurements/smart-attribute-model';
import { DeviceModel } from 'app/core/models/device-model';
import { AppConfig } from 'app/core/config/app.config';
import { SharedModule } from 'app/shared/shared.module';
import { DetailSettingsModule } from 'app/layout/common/detail-settings/detail-settings.module';

describe('DetailComponent', () => {
  let component: DetailComponent;
  let fixture: ComponentFixture<DetailComponent>;
  let mockDetailService: jasmine.SpyObj<DetailService>;
  let mockConfigService: jasmine.SpyObj<ScrutinyConfigService>;
  let mockDialog: jasmine.SpyObj<MatDialog>;
  let configSubject: Subject<AppConfig>;
  let dataSubject: Subject<any>;

  beforeEach(() => {
    configSubject = new Subject<AppConfig>();
    dataSubject = new Subject<any>();

    const defaultConfig: AppConfig = {
      dashboard_display: 'name',
      file_size_si_units: false,
      powered_on_hours_unit: 'humanize',
      temperature_unit: 'celsius',
      metrics: { status_threshold: 3 },
      theme: 'light'
    };

    mockDetailService = jasmine.createSpyObj('DetailService', ['getData'], {
      data$: dataSubject.asObservable()
    });
    mockConfigService = jasmine.createSpyObj('ScrutinyConfigService', [], {
      config$: of(defaultConfig)
    });
    mockDialog = jasmine.createSpyObj('MatDialog', ['open']);

    TestBed.configureTestingModule({
      declarations: [DetailComponent],
      imports: [
        SharedModule,
        FormsModule,
        MatButtonModule,
        MatIconModule,
        MatMenuModule,
        MatTableModule,
        MatSortModule,
        MatTooltipModule,
        MatDialogModule,
        MatFormFieldModule,
        MatSelectModule,
        MatInputModule,
        DetailSettingsModule
      ],
      providers: [
        { provide: DetailService, useValue: mockDetailService },
        { provide: ScrutinyConfigService, useValue: mockConfigService },
        { provide: MatDialog, useValue: mockDialog },
        { provide: LOCALE_ID, useValue: 'en-US' }
      ],
      schemas: [NO_ERRORS_SCHEMA]
    });

    fixture = TestBed.createComponent(DetailComponent);
    component = fixture.componentInstance;
    
    // Initialize required properties before ngOnInit runs
    component.device = { device_protocol: 'ATA' } as DeviceModel;
    component.smart_results = [];
    
    // Don't call detectChanges to avoid template rendering issues
    // We're only testing methods, not the template
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('getSSDPercentageUsed', () => {
    it('should return null when smart_results is empty', () => {
      component.smart_results = [];
      expect(component.getSSDPercentageUsed()).toBeNull();
    });

    it('should return null when smart_results is null', () => {
      component.smart_results = null as any;
      expect(component.getSSDPercentageUsed()).toBeNull();
    });

    it('should return null when attrs is missing', () => {
      component.smart_results = [{
        date: '2024-01-01T00:00:00Z',
        device_wwn: 'test-wwn',
        device_protocol: 'ATA',
        temp: 40,
        power_on_hours: 1000,
        power_cycle_count: 10,
        attrs: {}
      } as SmartModel];
      expect(component.getSSDPercentageUsed()).toBeNull();
    });

    it('should return null when attrs is null', () => {
      component.smart_results = [{
        date: '2024-01-01T00:00:00Z',
        device_wwn: 'test-wwn',
        device_protocol: 'ATA',
        temp: 40,
        power_on_hours: 1000,
        power_cycle_count: 10,
        attrs: null as any
      } as SmartModel];
      expect(component.getSSDPercentageUsed()).toBeNull();
    });

    it('should return percentage_used value for NVMe drives', () => {
      component.smart_results = [{
        date: '2024-01-01T00:00:00Z',
        device_wwn: 'test-wwn',
        device_protocol: 'NVMe',
        temp: 40,
        power_on_hours: 1000,
        power_cycle_count: 10,
        attrs: {
          'percentage_used': { 
            attribute_id: 'percentage_used',
            value: 25, 
            raw_value: undefined,
            thresh: 100,
            transformed_value: 25,
            status: 0
          } as SmartAttributeModel
        }
      } as SmartModel];
      expect(component.getSSDPercentageUsed()).toBe(25);
    });

    it('should return devstat_7_8 raw_value for ATA drives when raw_value is available', () => {
      component.smart_results = [{
        date: '2024-01-01T00:00:00Z',
        device_wwn: 'test-wwn',
        device_protocol: 'ATA',
        temp: 40,
        power_on_hours: 1000,
        power_cycle_count: 10,
        attrs: {
          'devstat_7_8': { 
            attribute_id: 'devstat_7_8',
            value: 50, 
            raw_value: 42,
            thresh: 100,
            transformed_value: 42,
            status: 0
          } as SmartAttributeModel
        }
      } as SmartModel];
      expect(component.getSSDPercentageUsed()).toBe(42);
    });

    it('should fallback to value when raw_value is undefined for devstat_7_8', () => {
      component.smart_results = [{
        date: '2024-01-01T00:00:00Z',
        device_wwn: 'test-wwn',
        device_protocol: 'ATA',
        temp: 40,
        power_on_hours: 1000,
        power_cycle_count: 10,
        attrs: {
          'devstat_7_8': { 
            attribute_id: 'devstat_7_8',
            value: 50, 
            raw_value: undefined,
            thresh: 100,
            transformed_value: 50,
            status: 0
          } as SmartAttributeModel
        }
      } as SmartModel];
      expect(component.getSSDPercentageUsed()).toBe(50);
    });

    it('should prioritize percentage_used over devstat_7_8 when both are present', () => {
      component.smart_results = [{
        date: '2024-01-01T00:00:00Z',
        device_wwn: 'test-wwn',
        device_protocol: 'NVMe',
        temp: 40,
        power_on_hours: 1000,
        power_cycle_count: 10,
        attrs: {
          'percentage_used': { 
            attribute_id: 'percentage_used',
            value: 30, 
            raw_value: undefined,
            thresh: 100,
            transformed_value: 30,
            status: 0
          } as SmartAttributeModel,
          'devstat_7_8': { 
            attribute_id: 'devstat_7_8',
            value: 50, 
            raw_value: 42,
            thresh: 100,
            transformed_value: 42,
            status: 0
          } as SmartAttributeModel
        }
      } as SmartModel];
      expect(component.getSSDPercentageUsed()).toBe(30);
    });

    it('should handle zero values correctly', () => {
      component.smart_results = [{
        date: '2024-01-01T00:00:00Z',
        device_wwn: 'test-wwn',
        device_protocol: 'NVMe',
        temp: 40,
        power_on_hours: 1000,
        power_cycle_count: 10,
        attrs: {
          'percentage_used': { 
            attribute_id: 'percentage_used',
            value: 0, 
            raw_value: undefined,
            thresh: 100,
            transformed_value: 0,
            status: 0
          } as SmartAttributeModel
        }
      } as SmartModel];
      expect(component.getSSDPercentageUsed()).toBe(0);
    });
  });

  describe('getSSDWearoutValue', () => {
    it('should return null when smart_results is empty', () => {
      component.smart_results = [];
      expect(component.getSSDWearoutValue()).toBeNull();
    });

    it('should return null when smart_results is null', () => {
      component.smart_results = null as any;
      expect(component.getSSDWearoutValue()).toBeNull();
    });

    it('should return null when attrs is missing', () => {
      component.smart_results = [{
        date: '2024-01-01T00:00:00Z',
        device_wwn: 'test-wwn',
        device_protocol: 'ATA',
        temp: 40,
        power_on_hours: 1000,
        power_cycle_count: 10,
        attrs: {}
      } as SmartModel];
      expect(component.getSSDWearoutValue()).toBeNull();
    });

    it('should return null when attrs is null', () => {
      component.smart_results = [{
        date: '2024-01-01T00:00:00Z',
        device_wwn: 'test-wwn',
        device_protocol: 'ATA',
        temp: 40,
        power_on_hours: 1000,
        power_cycle_count: 10,
        attrs: null as any
      } as SmartModel];
      expect(component.getSSDWearoutValue()).toBeNull();
    });

    it('should return attribute 177 value (Samsung/Crucial)', () => {
      component.smart_results = [{
        date: '2024-01-01T00:00:00Z',
        device_wwn: 'test-wwn',
        device_protocol: 'ATA',
        temp: 40,
        power_on_hours: 1000,
        power_cycle_count: 10,
        attrs: {
          '177': { 
            attribute_id: 177,
            value: 95,
            thresh: 10,
            transformed_value: 95,
            status: 0
          } as SmartAttributeModel
        }
      } as SmartModel];
      expect(component.getSSDWearoutValue()).toBe(95);
    });

    it('should return attribute 233 value (Intel) when 177 is not available', () => {
      component.smart_results = [{
        date: '2024-01-01T00:00:00Z',
        device_wwn: 'test-wwn',
        device_protocol: 'ATA',
        temp: 40,
        power_on_hours: 1000,
        power_cycle_count: 10,
        attrs: {
          '233': { 
            attribute_id: 233,
            value: 87,
            thresh: 1,
            transformed_value: 87,
            status: 0
          } as SmartAttributeModel
        }
      } as SmartModel];
      expect(component.getSSDWearoutValue()).toBe(87);
    });

    it('should prioritize 177 over 233 when both are available', () => {
      component.smart_results = [{
        date: '2024-01-01T00:00:00Z',
        device_wwn: 'test-wwn',
        device_protocol: 'ATA',
        temp: 40,
        power_on_hours: 1000,
        power_cycle_count: 10,
        attrs: {
          '177': { 
            attribute_id: 177,
            value: 95,
            thresh: 10,
            transformed_value: 95,
            status: 0
          } as SmartAttributeModel,
          '233': { 
            attribute_id: 233,
            value: 87,
            thresh: 1,
            transformed_value: 87,
            status: 0
          } as SmartAttributeModel
        }
      } as SmartModel];
      expect(component.getSSDWearoutValue()).toBe(95);
    });

    it('should check all fallback attributes in correct order (177, 233, 231, 232)', () => {
      component.smart_results = [{
        date: '2024-01-01T00:00:00Z',
        device_wwn: 'test-wwn',
        device_protocol: 'ATA',
        temp: 40,
        power_on_hours: 1000,
        power_cycle_count: 10,
        attrs: {
          '231': { 
            attribute_id: 231,
            value: 80,
            thresh: 10,
            transformed_value: 80,
            status: 0
          } as SmartAttributeModel,
          '232': { 
            attribute_id: 232,
            value: 75,
            thresh: 0,
            transformed_value: 75,
            status: 0
          } as SmartAttributeModel
        }
      } as SmartModel];
      // Should return 231 (Life Left) before 232 (Endurance Remaining)
      expect(component.getSSDWearoutValue()).toBe(80);
    });

    it('should return 232 when only 232 is available', () => {
      component.smart_results = [{
        date: '2024-01-01T00:00:00Z',
        device_wwn: 'test-wwn',
        device_protocol: 'ATA',
        temp: 40,
        power_on_hours: 1000,
        power_cycle_count: 10,
        attrs: {
          '232': { 
            attribute_id: 232,
            value: 75,
            thresh: 0,
            transformed_value: 75,
            status: 0
          } as SmartAttributeModel
        }
      } as SmartModel];
      expect(component.getSSDWearoutValue()).toBe(75);
    });

    it('should handle zero values correctly', () => {
      component.smart_results = [{
        date: '2024-01-01T00:00:00Z',
        device_wwn: 'test-wwn',
        device_protocol: 'ATA',
        temp: 40,
        power_on_hours: 1000,
        power_cycle_count: 10,
        attrs: {
          '177': { 
            attribute_id: 177,
            value: 0,
            thresh: 10,
            transformed_value: 0,
            status: 0
          } as SmartAttributeModel
        }
      } as SmartModel];
      expect(component.getSSDWearoutValue()).toBe(0);
    });

    it('should handle low values correctly (near end of life)', () => {
      component.smart_results = [{
        date: '2024-01-01T00:00:00Z',
        device_wwn: 'test-wwn',
        device_protocol: 'ATA',
        temp: 40,
        power_on_hours: 1000,
        power_cycle_count: 10,
        attrs: {
          '233': { 
            attribute_id: 233,
            value: 1,
            thresh: 1,
            transformed_value: 1,
            status: 0
          } as SmartAttributeModel
        }
      } as SmartModel];
      expect(component.getSSDWearoutValue()).toBe(1);
    });
  });

  describe('Integration: Priority between Percentage Used and Wearout Health', () => {
    it('should return percentage_used when both percentage_used and wearout attributes exist', () => {
      component.smart_results = [{
        date: '2024-01-01T00:00:00Z',
        device_wwn: 'test-wwn',
        device_protocol: 'NVMe',
        temp: 40,
        power_on_hours: 1000,
        power_cycle_count: 10,
        attrs: {
          'percentage_used': { 
            attribute_id: 'percentage_used',
            value: 25, 
            raw_value: undefined,
            thresh: 100,
            transformed_value: 25,
            status: 0
          } as SmartAttributeModel,
          '177': { 
            attribute_id: 177,
            value: 95,
            thresh: 10,
            transformed_value: 95,
            status: 0
          } as SmartAttributeModel
        }
      } as SmartModel];
      
      // Both methods should work independently
      expect(component.getSSDPercentageUsed()).toBe(25);
      expect(component.getSSDWearoutValue()).toBe(95);
      
      // In the template, percentage_used takes priority
      // This is handled by the @if/@else if logic in the template
    });

    it('should return devstat_7_8 when both devstat_7_8 and wearout attributes exist', () => {
      component.smart_results = [{
        date: '2024-01-01T00:00:00Z',
        device_wwn: 'test-wwn',
        device_protocol: 'ATA',
        temp: 40,
        power_on_hours: 1000,
        power_cycle_count: 10,
        attrs: {
          'devstat_7_8': { 
            attribute_id: 'devstat_7_8',
            value: 50, 
            raw_value: 42,
            thresh: 100,
            transformed_value: 42,
            status: 0
          } as SmartAttributeModel,
          '233': { 
            attribute_id: 233,
            value: 87,
            thresh: 1,
            transformed_value: 87,
            status: 0
          } as SmartAttributeModel
        }
      } as SmartModel];
      
      expect(component.getSSDPercentageUsed()).toBe(42);
      expect(component.getSSDWearoutValue()).toBe(87);
    });
  });
});
