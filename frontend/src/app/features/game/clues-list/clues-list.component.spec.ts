import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CluesListComponent } from './clues-list.component';

describe('CluesListComponent', () => {
    let component: CluesListComponent;
    let fixture: ComponentFixture<CluesListComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            imports: [CluesListComponent],
        }).compileComponents();

        fixture = TestBed.createComponent(CluesListComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
