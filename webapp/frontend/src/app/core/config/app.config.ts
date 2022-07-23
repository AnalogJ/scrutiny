import {Layout} from 'app/layout/layout.types';

// Theme type
export type Theme = 'light' | 'dark' | 'system';

// Device title to display on the dashboard
export type DashboardDisplay = 'name' | 'serial_id' | 'uuid' | 'label'

export type DashboardSort = 'status' | 'title' | 'age'

export type TemperatureUnit = 'celsius' | 'fahrenheit'


export enum MetricsNotifyLevel {
    Warn = 1,
    Fail = 2
}

export enum MetricsStatusFilterAttributes {
    All = 0,
    Critical = 1
}

export enum MetricsStatusThreshold {
    Smart = 1,
    Scrutiny = 2,

    // shortcut
    Both = 3
}

/**
 * AppConfig interface. Update this interface to strictly type your config
 * object.
 */
export interface AppConfig {
    theme?: Theme;
    layout?: Layout;

    // Dashboard options
    dashboard_display?: DashboardDisplay;
    dashboard_sort?: DashboardSort;

    temperature_unit?: TemperatureUnit;

    // Settings from Scrutiny API

    metrics?: {
        notify_level?: MetricsNotifyLevel
        status_filter_attributes?: MetricsStatusFilterAttributes
        status_threshold?: MetricsStatusThreshold
    }

}

/**
 * Default configuration for the entire application. This object is used by
 * "ConfigService" to set the default configuration.
 *
 * If you need to store global configuration for your app, you can use this
 * object to set the defaults. To access, update and reset the config, use
 * "ConfigService".
 */
export const appConfig: AppConfig = {
    theme: 'light',
    layout: 'material',

    dashboard_display: 'name',
    dashboard_sort: 'status',

    temperature_unit: 'celsius',
    metrics: {
        notify_level: MetricsNotifyLevel.Fail,
        status_filter_attributes: MetricsStatusFilterAttributes.All,
        status_threshold: MetricsStatusThreshold.Both
    }
};

