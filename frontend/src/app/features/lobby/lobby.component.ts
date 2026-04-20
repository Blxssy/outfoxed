import { Component } from '@angular/core';
import { ButtonComponent } from '@fox/ui-kit/button';
import { CardComponent } from '@fox/ui-kit/card/card.component';

@Component({
    selector: 'app-lobby',
    imports: [ButtonComponent, CardComponent],
    templateUrl: './lobby.component.html',
    styleUrl: './lobby.component.scss',
})
export class LobbyComponent {}
