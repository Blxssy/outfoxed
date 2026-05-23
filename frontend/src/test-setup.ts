import '@angular/compiler';

import '@testing-library/jest-dom';

import '@analogjs/vitest-angular/setup-zone';

import {
    BrowserTestingModule,
    platformBrowserTesting,
} from '@angular/platform-browser/testing';

import { getTestBed } from '@angular/core/testing';

getTestBed().initTestEnvironment(
    BrowserTestingModule,
    platformBrowserTesting(),
);
