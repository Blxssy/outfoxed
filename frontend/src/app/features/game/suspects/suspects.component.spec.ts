import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SuspectsComponent } from './suspects.component';

describe('SuspectsComponent', () => {
    let component: SuspectsComponent;
    let fixture: ComponentFixture<SuspectsComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            imports: [SuspectsComponent],
        }).compileComponents();

        fixture = TestBed.createComponent(SuspectsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
