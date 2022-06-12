export const sdd = {
    'data': {
        'device': {
            'CreatedAt': '2021-06-24T21:17:31.30374-07:00',
            'UpdatedAt': '2021-10-24T16:37:57.013758-07:00',
            'DeletedAt': null,
            'wwn': '0x5000cca252c859cc',
            'device_name': 'sdd',
            'manufacturer': 'ATA',
            'model_name': 'WDC_WD80EFAX-68LHPN0',
            'interface_type': 'SCSI',
            'interface_speed': '',
            'serial_number': '7SGLXXXXX',
            'firmware': '',
            'rotational_speed': 0,
            'capacity': 8001563222016,
            'form_factor': '',
            'smart_support': false,
            'device_protocol': 'SCSI',
            'device_type': '',
            'label': '',
            'host_id': '',
            'device_status': 0
        },
        'smart_results': [{
            'date': '2021-10-24T23:20:44Z',
            'device_wwn': '0x5000cca252c859cc',
            'device_protocol': 'SCSI',
            'temp': 34,
            'power_on_hours': 43549,
            'power_cycle_count': 0,
            'attrs': {
                'read_correction_algorithm_invocations': {
                    'attribute_id': 'read_correction_algorithm_invocations',
                    'value': 0,
                    'thresh': -1,
                    'transformed_value': 0,
                    'status': 0
                },
                'read_errors_corrected_by_eccdelayed': {
                    'attribute_id': 'read_errors_corrected_by_eccdelayed',
                    'value': 0,
                    'thresh': -1,
                    'transformed_value': 0,
                    'status': 0
                },
                'read_errors_corrected_by_eccfast': {
                    'attribute_id': 'read_errors_corrected_by_eccfast',
                    'value': 300357663,
                    'thresh': -1,
                    'transformed_value': 0,
                    'status': 0
                },
                'read_errors_corrected_by_rereads_rewrites': {
                    'attribute_id': 'read_errors_corrected_by_rereads_rewrites',
                    'value': 0,
                    'thresh': 0,
                    'transformed_value': 0,
                    'status': 0
                },
                'read_total_errors_corrected': {
                    'attribute_id': 'read_total_errors_corrected',
                    'value': 300357663,
                    'thresh': -1,
                    'transformed_value': 0,
                    'status': 0
                },
                'read_total_uncorrected_errors': {
                    'attribute_id': 'read_total_uncorrected_errors',
                    'value': 0,
                    'thresh': 0,
                    'transformed_value': 0,
                    'status': 0
                },
                'scsi_grown_defect_list': {
                    'attribute_id': 'scsi_grown_defect_list',
                    'value': 56,
                    'thresh': 0,
                    'transformed_value': 0,
                    'status': 0
                },
                'write_correction_algorithm_invocations': {
                    'attribute_id': 'write_correction_algorithm_invocations',
                    'value': 0,
                    'thresh': -1,
                    'transformed_value': 0,
                    'status': 0
                },
                'write_errors_corrected_by_eccdelayed': {
                    'attribute_id': 'write_errors_corrected_by_eccdelayed',
                    'value': 0,
                    'thresh': -1,
                    'transformed_value': 0,
                    'status': 0
                },
                'write_errors_corrected_by_eccfast': {
                    'attribute_id': 'write_errors_corrected_by_eccfast',
                    'value': 0,
                    'thresh': -1,
                    'transformed_value': 0,
                    'status': 0
                },
                'write_errors_corrected_by_rereads_rewrites': {
                    'attribute_id': 'write_errors_corrected_by_rereads_rewrites',
                    'value': 0,
                    'thresh': 0,
                    'transformed_value': 0,
                    'status': 0
                },
                'write_total_errors_corrected': {
                    'attribute_id': 'write_total_errors_corrected',
                    'value': 0,
                    'thresh': -1,
                    'transformed_value': 0,
                    'status': 0
                },
                'write_total_uncorrected_errors': {
                    'attribute_id': 'write_total_uncorrected_errors',
                    'value': 0,
                    'thresh': 0,
                    'transformed_value': 0,
                    'status': 0
                }
            },
            'Status': 0
        }]
    },
    'metadata': {
        'read_correction_algorithm_invocations': {
            'display_name': 'Read Correction Algorithm Invocations',
            'ideal': '',
            'critical': false,
            'description': '',
            'display_type': ''
        },
        'read_errors_corrected_by_eccdelayed': {
            'display_name': 'Read Errors Corrected by ECC Delayed',
            'ideal': '',
            'critical': false,
            'description': '',
            'display_type': ''
        },
        'read_errors_corrected_by_eccfast': {
            'display_name': 'Read Errors Corrected by ECC Fast',
            'ideal': '',
            'critical': false,
            'description': '',
            'display_type': ''
        },
        'read_errors_corrected_by_rereads_rewrites': {
            'display_name': 'Read Errors Corrected by ReReads/ReWrites',
            'ideal': 'low',
            'critical': true,
            'description': '',
            'display_type': ''
        },
        'read_total_errors_corrected': {
            'display_name': 'Read Total Errors Corrected',
            'ideal': '',
            'critical': false,
            'description': '',
            'display_type': ''
        },
        'read_total_uncorrected_errors': {
            'display_name': 'Read Total Uncorrected Errors',
            'ideal': 'low',
            'critical': true,
            'description': '',
            'display_type': ''
        },
        'scsi_grown_defect_list': {
            'display_name': 'Grown Defect List',
            'ideal': 'low',
            'critical': true,
            'description': '',
            'display_type': ''
        },
        'write_correction_algorithm_invocations': {
            'display_name': 'Write Correction Algorithm Invocations',
            'ideal': '',
            'critical': false,
            'description': '',
            'display_type': ''
        },
        'write_errors_corrected_by_eccdelayed': {
            'display_name': 'Write Errors Corrected by ECC Delayed',
            'ideal': '',
            'critical': false,
            'description': '',
            'display_type': ''
        },
        'write_errors_corrected_by_eccfast': {
            'display_name': 'Write Errors Corrected by ECC Fast',
            'ideal': '',
            'critical': false,
            'description': '',
            'display_type': ''
        },
        'write_errors_corrected_by_rereads_rewrites': {
            'display_name': 'Write Errors Corrected by ReReads/ReWrites',
            'ideal': 'low',
            'critical': true,
            'description': '',
            'display_type': ''
        },
        'write_total_errors_corrected': {
            'display_name': 'Write Total Errors Corrected',
            'ideal': '',
            'critical': false,
            'description': '',
            'display_type': ''
        },
        'write_total_uncorrected_errors': {
            'display_name': 'Write Total Uncorrected Errors',
            'ideal': 'low',
            'critical': true,
            'description': '',
            'display_type': ''
        }
    },
    'success': true
}
