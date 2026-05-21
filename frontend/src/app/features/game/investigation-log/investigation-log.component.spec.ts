import { ComponentFixture, TestBed } from '@angular/core/testing';

import { InvestigationLogComponent } from './investigation-log.component';

describe('InvestigationLogComponent', () => {
    let component: InvestigationLogComponent;
    let fixture: ComponentFixture<InvestigationLogComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            imports: [InvestigationLogComponent],
        }).compileComponents();

        fixture = TestBed.createComponent(InvestigationLogComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
