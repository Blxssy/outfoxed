import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FinishModalComponent } from './finish-modal.component';

describe('FinishModalComponent', () => {
    let component: FinishModalComponent;
    let fixture: ComponentFixture<FinishModalComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            imports: [FinishModalComponent],
        }).compileComponents();

        fixture = TestBed.createComponent(FinishModalComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
