import {Layout} from 'app/layout/layout.types';

// Theme type
export type Theme = 'light' | 'dark' | 'system';

// Device title to display on the dashboard
export type DashboardDisplay = 'name' | 'serial_id' | 'uuid' | 'label'

export type DashboardSort =
    | 'status' | 'status_asc' | 'status_desc'
    | 'title' | 'title_asc' | 'title_desc'
    | 'age' | 'age_asc' | 'age_desc'
    | 'capacity_asc' | 'capacity_desc'
    | 'temperature_asc' | 'temperature_desc'

export type TemperatureUnit = 'celsius' | 'fahrenheit'

export type LineStroke = 'smooth' | 'straight' | 'stepline'

export type DevicePoweredOnUnit = 'humanize' | 'device_hours'


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

    file_size_si_units?: boolean;

    powered_on_hours_unit?: DevicePoweredOnUnit;

    line_stroke?: LineStroke;

    // Settings from Scrutiny API
    
    collector?: {
        retrieve_sct_temperature_history?: boolean
    }

    metrics?: {
        notify_level?: MetricsNotifyLevel
        status_filter_attributes?: MetricsStatusFilterAttributes
        status_threshold?: MetricsStatusThreshold
        repeat_notifications?: boolean
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
    file_size_si_units: false,
    powered_on_hours_unit: 'humanize',

    line_stroke: 'smooth',
    
    collector: {
        retrieve_sct_temperature_history : true,
    },

    metrics: {
        notify_level: MetricsNotifyLevel.Fail,
        status_filter_attributes: MetricsStatusFilterAttributes.All,
        status_threshold: MetricsStatusThreshold.Both,
        repeat_notifications: true
    }
};

