import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FoxTrackComponent } from './fox-track.component';

describe('FoxTrackComponent', () => {
    let component: FoxTrackComponent;
    let fixture: ComponentFixture<FoxTrackComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            imports: [FoxTrackComponent],
        }).compileComponents();

        fixture = TestBed.createComponent(FoxTrackComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
