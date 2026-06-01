import 'zone.js';
import 'zone.js/testing';
import { getTestBed } from '@angular/core/testing';
import {
    BrowserTestingModule,
    platformBrowserTesting,
} from '@angular/platform-browser/testing';
import { destroyPlatform } from '@angular/core';

destroyPlatform();

getTestBed().initTestEnvironment(
    BrowserTestingModule,
    platformBrowserTesting(),
    { teardown: { destroyAfterEach: true } },
);
