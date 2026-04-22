import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EnterCodeModalComponent } from './enter-code-modal.component';

describe('EnterCodeModalComponent', () => {
    let component: EnterCodeModalComponent;
    let fixture: ComponentFixture<EnterCodeModalComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            imports: [EnterCodeModalComponent],
        }).compileComponents();

        fixture = TestBed.createComponent(EnterCodeModalComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
