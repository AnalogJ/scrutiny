export const sda = {
    "data": {
        "device": {
            "CreatedAt": "2021-06-24T21:17:31.301226-07:00",
            "UpdatedAt": "2021-10-24T16:37:56.981833-07:00",
            "DeletedAt": null,
            "wwn": "0x5002538e40a22954",
            "device_name": "sda",
            "manufacturer": "ATA",
            "model_name": "Samsung_SSD_860_EVO_500GB",
            "interface_type": "SCSI",
            "interface_speed": "",
            "serial_number": "S3YZNB0KBXXXXXX",
            "firmware": "002C",
            "rotational_speed": 0,
            "capacity": 500107862016,
            "form_factor": "",
            "smart_support": false,
            "device_protocol": "NVMe",
            "device_type": "",
            "label": "",
            "host_id": "",
            "device_status": 0
        },
        "smart_results": [{
            "date": "2021-10-24T23:20:44Z",
            "device_wwn": "0x5002538e40a22954",
            "device_protocol": "NVMe",
            "temp": 36,
            "power_on_hours": 2401,
            "power_cycle_count": 266,
            "attrs": {
                "available_spare": {
                    "attribute_id": "available_spare",
                    "value": 100,
                    "thresh": 10,
                    "transformed_value": 0,
                    "status": 0
                },
                "controller_busy_time": {
                    "attribute_id": "controller_busy_time",
                    "value": 3060,
                    "thresh": -1,
                    "transformed_value": 0,
                    "status": 0
                },
                "critical_comp_time": {
                    "attribute_id": "critical_comp_time",
                    "value": 0,
                    "thresh": -1,
                    "transformed_value": 0,
                    "status": 0
                },
                "critical_warning": {
                    "attribute_id": "critical_warning",
                    "value": 0,
                    "thresh": 0,
                    "transformed_value": 0,
                    "status": 0
                },
                "data_units_read": {
                    "attribute_id": "data_units_read",
                    "value": 9511859,
                    "thresh": -1,
                    "transformed_value": 0,
                    "status": 0
                },
                "data_units_written": {
                    "attribute_id": "data_units_written",
                    "value": 7773431,
                    "thresh": -1,
                    "transformed_value": 0,
                    "status": 0
                },
                "host_reads": {
                    "attribute_id": "host_reads",
                    "value": 111303174,
                    "thresh": -1,
                    "transformed_value": 0,
                    "status": 0
                },
                "host_writes": {
                    "attribute_id": "host_writes",
                    "value": 83170961,
                    "thresh": -1,
                    "transformed_value": 0,
                    "status": 0
                },
                "media_errors": {
                    "attribute_id": "media_errors",
                    "value": 0,
                    "thresh": 0,
                    "transformed_value": 0,
                    "status": 0
                },
                "num_err_log_entries": {
                    "attribute_id": "num_err_log_entries",
                    "value": 0,
                    "thresh": 0,
                    "transformed_value": 0,
                    "status": 0
                },
                "percentage_used": {
                    "attribute_id": "percentage_used",
                    "value": 0,
                    "thresh": 100,
                    "transformed_value": 0,
                    "status": 0
                },
                "power_cycles": {
                    "attribute_id": "power_cycles",
                    "value": 266,
                    "thresh": -1,
                    "transformed_value": 0,
                    "status": 0
                },
                "power_on_hours": {
                    "attribute_id": "power_on_hours",
                    "value": 2401,
                    "thresh": -1,
                    "transformed_value": 0,
                    "status": 0
                },
                "temperature": {
                    "attribute_id": "temperature",
                    "value": 36,
                    "thresh": -1,
                    "transformed_value": 0,
                    "status": 0
                },
                "unsafe_shutdowns": {
                    "attribute_id": "unsafe_shutdowns",
                    "value": 43,
                    "thresh": -1,
                    "transformed_value": 0,
                    "status": 0
                },
                "warning_temp_time": {
                    "attribute_id": "warning_temp_time",
                    "value": 0,
                    "thresh": -1,
                    "transformed_value": 0,
                    "status": 0
                }
            },
            "Status": 0
        }]
    },
    "metadata": {
        "available_spare": {
            "display_name": "Available Spare",
            "ideal": "high",
            "critical": true,
            "description": "Contains a normalized percentage (0 to 100%) of the remaining spare capacity available.",
            "display_type": ""
        },
        "controller_busy_time": {
            "display_name": "Controller Busy Time",
            "ideal": "",
            "critical": false,
            "description": "Contains the amount of time the controller is busy with I/O commands. The controller is busy when there is a command outstanding to an I/O Queue (specifically, a command was issued via an I/O Submission Queue Tail doorbell write and the corresponding completion queue entry has not been posted yet to the associated I/O Completion Queue). This value is reported in minutes.",
            "display_type": ""
        },
        "critical_comp_time": {
            "display_name": "Critical CompTime",
            "ideal": "",
            "critical": false,
            "description": "Contains the amount of time in minutes that the controller is operational and the Composite Temperature is greater the Critical Composite Temperature Threshold (CCTEMP) field in the Identify Controller data structure.",
            "display_type": ""
        },
        "critical_warning": {
            "display_name": "Critical Warning",
            "ideal": "low",
            "critical": true,
            "description": "This field indicates critical warnings for the state of the controller. Each bit corresponds to a critical warning type; multiple bits may be set. If a bit is cleared to ‘0’, then that critical warning does not apply. Critical warnings may result in an asynchronous event notification to the host. Bits in this field represent the current associated state and are not persistent.",
            "display_type": ""
        },
        "data_units_read": {
            "display_name": "Data Units Read",
            "ideal": "",
            "critical": false,
            "description": "Contains the number of 512 byte data units the host has read from the controller; this value does not include metadata. This value is reported in thousands (i.e., a value of 1 corresponds to 1000 units of 512 bytes read) and is rounded up. When the LBA size is a value other than 512 bytes, the controller shall convert the amount of data read to 512 byte units.",
            "display_type": ""
        },
        "data_units_written": {
            "display_name": "Data Units Written",
            "ideal": "",
            "critical": false,
            "description": "Contains the number of 512 byte data units the host has written to the controller; this value does not include metadata. This value is reported in thousands (i.e., a value of 1 corresponds to 1000 units of 512 bytes written) and is rounded up. When the LBA size is a value other than 512 bytes, the controller shall convert the amount of data written to 512 byte units.",
            "display_type": ""
        },
        "host_reads": {
            "display_name": "Host Reads",
            "ideal": "",
            "critical": false,
            "description": "Contains the number of read commands completed by the controller",
            "display_type": ""
        },
        "host_writes": {
            "display_name": "Host Writes",
            "ideal": "",
            "critical": false,
            "description": "Contains the number of write commands completed by the controller",
            "display_type": ""
        },
        "media_errors": {
            "display_name": "Media Errors",
            "ideal": "low",
            "critical": true,
            "description": "Contains the number of occurrences where the controller detected an unrecovered data integrity error. Errors such as uncorrectable ECC, CRC checksum failure, or LBA tag mismatch are included in this field.",
            "display_type": ""
        },
        "num_err_log_entries": {
            "display_name": "Numb Err Log Entries",
            "ideal": "low",
            "critical": true,
            "description": "Contains the number of Error Information log entries over the life of the controller.",
            "display_type": ""
        },
        "percentage_used": {
            "display_name": "Percentage Used",
            "ideal": "low",
            "critical": true,
            "description": "Contains a vendor specific estimate of the percentage of NVM subsystem life used based on the actual usage and the manufacturer’s prediction of NVM life. A value of 100 indicates that the estimated endurance of the NVM in the NVM subsystem has been consumed, but may not indicate an NVM subsystem failure. The value is allowed to exceed 100. Percentages greater than 254 shall be represented as 255. This value shall be updated once per power-on hour (when the controller is not in a sleep state).",
            "display_type": ""
        },
        "power_cycles": {
            "display_name": "Power Cycles",
            "ideal": "",
            "critical": false,
            "description": "Contains the number of power cycles.",
            "display_type": ""
        },
        "power_on_hours": {
            "display_name": "Power on Hours",
            "ideal": "",
            "critical": false,
            "description": "Contains the number of power-on hours. Power on hours is always logging, even when in low power mode.",
            "display_type": ""
        },
        "temperature": {
            "display_name": "Temperature",
            "ideal": "",
            "critical": false,
            "description": "",
            "display_type": ""
        },
        "unsafe_shutdowns": {
            "display_name": "Unsafe Shutdowns",
            "ideal": "",
            "critical": false,
            "description": "Contains the number of unsafe shutdowns. This count is incremented when a shutdown notification (CC.SHN) is not received prior to loss of power.",
            "display_type": ""
        },
        "warning_temp_time": {
            "display_name": "Warning Temp Time",
            "ideal": "",
            "critical": false,
            "description": "Contains the amount of time in minutes that the controller is operational and the Composite Temperature is greater than or equal to the Warning Composite Temperature Threshold (WCTEMP) field and less than the Critical Composite Temperature Threshold (CCTEMP) field in the Identify Controller data structure.",
            "display_type": ""
        }
    },
    "success": true
}
