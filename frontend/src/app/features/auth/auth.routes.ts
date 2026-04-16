import { Routes } from '@angular/router';
import { AuthComponent } from './auth.component';

export const routes: Routes = [
    {
        path: '',
        component: AuthComponent,
        children: [
            {
                path: 'register',
                loadComponent: () =>
                    import('./pages/register/register.component').then(
                        (m) => m.RegisterComponent,
                    ),
            },
            {
                path: 'login',
                loadComponent: () =>
                    import('./pages/login/login.component').then(
                        (m) => m.LoginComponent,
                    ),
            },
            {
                path: '**',
                redirectTo: 'register',
                pathMatch: 'full',
            },
        ],
    },
];
