import { render } from '@testing-library/angular';
import { AppComponent } from './app.component';

describe('AppComponent', () => {
    it('should create app component', async () => {
        const { fixture } = await render(AppComponent);

        expect(fixture.componentInstance).toBeTruthy();
    });

    it('should render router outlet', async () => {
        const { container } = await render(AppComponent);

        expect(container.querySelector('router-outlet')).toBeTruthy();
    });
});
