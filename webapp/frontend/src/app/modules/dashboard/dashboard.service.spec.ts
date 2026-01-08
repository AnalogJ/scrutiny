import { HttpClient } from '@angular/common/http';
import {DashboardService} from './dashboard.service';
import {of} from 'rxjs';
import {summary} from 'app/data/mock/summary/data'
import {temp_history} from 'app/data/mock/summary/temp_history'
import {DeviceSummaryModel} from 'app/core/models/device-summary-model';
import {SmartTemperatureModel} from 'app/core/models/measurements/smart-temperature-model';

describe('DashboardService', () => {
    let service: DashboardService;
    let httpClientSpy: jasmine.SpyObj<HttpClient>;

    beforeEach(() => {
        httpClientSpy = jasmine.createSpyObj('HttpClient', ['get']);
        service = new DashboardService(httpClientSpy);
    });

    it('should unwrap and return getSummaryData() (HttpClient called once)', (done: DoneFn) => {
        httpClientSpy.get.and.returnValue(of(summary));

        service.getSummaryData().subscribe(value => {
            expect(value).toBe(summary.data.summary as { [key: string]: DeviceSummaryModel });
            done();
        });
        expect(httpClientSpy.get.calls.count())
            .withContext('one call')
            .toBe(1);
    });

    it('should unwrap and return getSummaryTempData() (HttpClient called once)', (done: DoneFn) => {
        // const expectedHeroes: any[] =
        //     [{ id: 1, name: 'A' }, { id: 2, name: 'B' }];

        httpClientSpy.get.and.returnValue(of(temp_history));

        service.getSummaryTempData('weekly').subscribe(value => {
            expect(value).toBe(temp_history.data.temp_history as { [key: string]: SmartTemperatureModel[] });
            done();
        });
        expect(httpClientSpy.get.calls.count())
            .withContext('one call')
            .toBe(1);
    });
});
