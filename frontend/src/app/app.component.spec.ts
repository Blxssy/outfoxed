import { render } from '@testing-library/angular';
import { AppComponent } from './app.component';
import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { provideHttpClient } from '@angular/common/http';

describe('AppComponent', () => {
    beforeEach(async () => {
        await TestBed.configureTestingModule({
            providers: [provideRouter([]), provideHttpClient()],
        }).compileComponents();
    });
    it('should create app component', async () => {
        const { fixture } = await render(AppComponent);

        expect(fixture.componentInstance).toBeTruthy();
    });

    it('should render router outlet', async () => {
        const { container } = await render(AppComponent);

        expect(container.querySelector('router-outlet')).toBeTruthy();
    });
});
