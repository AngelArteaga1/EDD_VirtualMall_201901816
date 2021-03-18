import { TestBed } from '@angular/core/testing';

import { GetTiendasService } from './get-tiendas.service';

describe('GetTiendasService', () => {
  let service: GetTiendasService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(GetTiendasService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
