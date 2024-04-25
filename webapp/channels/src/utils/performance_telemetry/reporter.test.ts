// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {Client4} from '@mattermost/client';

import configureStore from 'store';

import {initializePerformanceMocks, waitForObservations} from 'tests/helpers/performance_mocks';

import PerformanceReporter from './reporter';

import {measureAndReport} from '.';

jest.mock('web-vitals');

initializePerformanceMocks();

const sendBeacon = jest.fn().mockReturnValue(true);
navigator.sendBeacon = sendBeacon;

describe('PerformanceReporter', () => {
    test('should report measurements to the server periodically', async () => {
        const reporter = newTestReporter();
        reporter.observe();

        expect(sendBeacon).not.toHaveBeenCalled();

        const testMarkA = performance.mark('testMarkA');
        const testMarkB = performance.mark('testMarkB');
        measureAndReport('testMeasure', 'testMarkA', 'testMarkB');

        await waitForObservations();

        expect(reporter.handleMeasures).toHaveBeenCalled();

        await waitForReport();

        expect(sendBeacon).toHaveBeenCalledTimes(1);
        expect(sendBeacon.mock.calls[0][0]).toEqual('/api/v4/metrics');
        const report = JSON.parse(sendBeacon.mock.calls[0][1]);
        expect(report).toMatchObject({
            measures: [
                {
                    name: 'testMeasure',
                    value: testMarkB.startTime - testMarkA.startTime,
                },
            ],
        });
    });
});

class TestPerformanceReporter extends PerformanceReporter {
    public reportPeriodBase = 10;
    public reportPeriodJitter = 0;

    public handleMeasures = jest.fn(super.handleMeasures);
}

function newTestReporter(telemetryEnabled = true) {
    return new TestPerformanceReporter(new Client4(), configureStore({
        entities: {
            general: {
                config: {
                    DiagnosticsEnabled: String(telemetryEnabled),
                },
            },
        },
    }));
}

function waitForReport() {
    // Reports are set every 10ms by default
    return new Promise((resolve) => setTimeout(resolve, 10));
}
