import { HttpClient } from '@angular/common/http';
import {DetailService} from './detail.service';
import {of} from 'rxjs';
import {sda} from 'app/data/mock/device/details/sda'
import {DeviceDetailsResponseWrapper} from 'app/core/models/device-details-response-wrapper';

describe('DetailService', () => {
    describe('#getData', () => {
        let service: DetailService;
        let httpClientSpy: jasmine.SpyObj<HttpClient>;

        beforeEach(() => {
            httpClientSpy = jasmine.createSpyObj('HttpClient', ['get']);
            service = new DetailService(httpClientSpy);
        });
        it('should return getData() (HttpClient called once)', (done: DoneFn) => {
            httpClientSpy.get.and.returnValue(of(sda));

            service.getData('test').subscribe(value => {
                expect(value).toBe(sda as DeviceDetailsResponseWrapper);
                done();
            });
            expect(httpClientSpy.get.calls.count())
                .withContext('one call')
                .toBe(1);
        });
    })
});
