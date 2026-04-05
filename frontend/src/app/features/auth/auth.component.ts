import { Component } from '@angular/core';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';

@Component({
    selector: 'app-auth',
    imports: [RouterOutlet, RouterLink, RouterLinkActive],
    templateUrl: './auth.component.html',
    styleUrl: './auth.component.scss',
})
export class AuthComponent {}
