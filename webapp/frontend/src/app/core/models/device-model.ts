// maps to webapp/backend/pkg/models/device.go
export interface DeviceModel {
    wwn: string;
    device_name: string;
    device_uuid: string;
    device_serial_id: string;
    device_label: string;

    manufacturer: string;
    model_name: string;
    interface_type: string;
    interface_speed: string;
    serial_number: string;
    firmware: string;
    rotational_speed: number;
    capacity: number;
    form_factor: string;
    smart_support: boolean;
    device_protocol: string;
    device_type: string;

    label: string;
    host_id: string;

    device_status: number;
}
