export const sdf = {
    "data": {
        "device": {
            "CreatedAt": "2020-09-29T03:17:30.813487336Z",
            "UpdatedAt": "2020-09-29T03:17:31.04293448Z",
            "DeletedAt": null,
            "wwn": "0x5000c500c0558f02",
            "device_name": "sdc",
            "manufacturer": "",
            "model_name": "ST8000VN004-2M2101",
            "interface_type": "",
            "interface_speed": "6.0 Gb/s",
            "serial_number": "WKD01Y5S",
            "firmware": "SC60",
            "rotational_speed": 7200,
            "capacity": 8001563222016,
            "form_factor": "3.5 inches",
            "smart_support": false,
            "device_protocol": "ATA",
            "device_type": "scsi",
            "device_status": 2,
        },
        "smart_results": [{
            "ID": 1,
            "CreatedAt": "2020-09-29T03:17:31.063859162Z",
            "UpdatedAt": "2020-09-29T03:17:31.063859162Z",
            "DeletedAt": null,
            "device_wwn": "0x5000c500c0558f02",
            "date": "2020-09-29T03:17:30Z",
            "smart_status": "passed",
            "temp": 39,
            "power_on_hours": 9499,
            "power_cycle_count": 22,
            "ata_attributes": [{
                "ID": 1,
                "CreatedAt": "2020-09-29T03:17:31.064174997Z",
                "UpdatedAt": "2020-09-29T03:17:31.064174997Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 1,
                "name": "Read Error Rate",
                "value": 79,
                "worst": 64,
                "thresh": 44,
                "raw_value": 78022392,
                "raw_string": "78022392",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 2,
                "CreatedAt": "2020-09-29T03:17:31.064379412Z",
                "UpdatedAt": "2020-09-29T03:17:31.064379412Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 3,
                "name": "Spin-Up Time",
                "value": 95,
                "worst": 86,
                "thresh": 0,
                "raw_value": 0,
                "raw_string": "0",
                "when_failed": "",
                "transformed_value": 0,
                "status": "warn",
                "status_reason": "Observed Failure Rate for Attribute is greater than 10%",
                "failure_rate": 0.11452195377351217
            }, {
                "ID": 3,
                "CreatedAt": "2020-09-29T03:17:31.064506775Z",
                "UpdatedAt": "2020-09-29T03:17:31.064506775Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 4,
                "name": "Start/Stop Count",
                "value": 100,
                "worst": 100,
                "thresh": 20,
                "raw_value": 220,
                "raw_string": "220",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 4,
                "CreatedAt": "2020-09-29T03:17:31.064621265Z",
                "UpdatedAt": "2020-09-29T03:17:31.064621265Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 5,
                "name": "Reallocated Sectors Count",
                "value": 100,
                "worst": 100,
                "thresh": 10,
                "raw_value": 0,
                "raw_string": "0",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed",
                "failure_rate": 0.025169175350572493
            }, {
                "ID": 5,
                "CreatedAt": "2020-09-29T03:17:31.064742769Z",
                "UpdatedAt": "2020-09-29T03:17:31.064742769Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 7,
                "name": "Seek Error Rate",
                "value": 86,
                "worst": 60,
                "thresh": 45,
                "raw_value": 4679507461,
                "raw_string": "4679507461",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed",
                "failure_rate": 0.08725919610118257
            }, {
                "ID": 6,
                "CreatedAt": "2020-09-29T03:17:31.064850152Z",
                "UpdatedAt": "2020-09-29T03:17:31.064850152Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 9,
                "name": "Power-On Hours",
                "value": 90,
                "worst": 90,
                "thresh": 0,
                "raw_value": 9499,
                "raw_string": "9499",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 7,
                "CreatedAt": "2020-09-29T03:17:31.064970318Z",
                "UpdatedAt": "2020-09-29T03:17:31.064970318Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 10,
                "name": "Spin Retry Count",
                "value": 100,
                "worst": 100,
                "thresh": 97,
                "raw_value": 0,
                "raw_string": "0",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed",
                "failure_rate": 0.05459827163896099
            }, {
                "ID": 8,
                "CreatedAt": "2020-09-29T03:17:31.065072898Z",
                "UpdatedAt": "2020-09-29T03:17:31.065072898Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 12,
                "name": "Power Cycle Count",
                "value": 100,
                "worst": 100,
                "thresh": 20,
                "raw_value": 22,
                "raw_string": "22",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed",
                "failure_rate": 0.038210930067894826
            }, {
                "ID": 9,
                "CreatedAt": "2020-09-29T03:17:31.065366547Z",
                "UpdatedAt": "2020-09-29T03:17:31.065366547Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 18,
                "name": "Unknown_Attribute",
                "value": 100,
                "worst": 100,
                "thresh": 50,
                "raw_value": 0,
                "raw_string": "0",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 10,
                "CreatedAt": "2020-09-29T03:17:31.065539403Z",
                "UpdatedAt": "2020-09-29T03:17:31.065539403Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 187,
                "name": "Reported Uncorrectable Errors",
                "value": 100,
                "worst": 100,
                "thresh": 0,
                "raw_value": 0,
                "raw_string": "0",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed",
                "failure_rate": 0.028130798308190524
            }, {
                "ID": 11,
                "CreatedAt": "2020-09-29T03:17:31.065665808Z",
                "UpdatedAt": "2020-09-29T03:17:31.065665808Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 188,
                "name": "Command Timeout",
                "value": 100,
                "worst": 100,
                "thresh": 0,
                "raw_value": 0,
                "raw_string": "0",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed",
                "failure_rate": 0.024893587674442153
            }, {
                "ID": 12,
                "CreatedAt": "2020-09-29T03:17:31.065780779Z",
                "UpdatedAt": "2020-09-29T03:17:31.065780779Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 190,
                "name": "Temperature Difference",
                "value": 61,
                "worst": 46,
                "thresh": 40,
                "raw_value": 740294695,
                "raw_string": "39 (Min/Max 32/44)",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 13,
                "CreatedAt": "2020-09-29T03:17:31.065931499Z",
                "UpdatedAt": "2020-09-29T03:17:31.065931499Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 192,
                "name": "Power-off Retract Count",
                "value": 100,
                "worst": 100,
                "thresh": 0,
                "raw_value": 4,
                "raw_string": "4",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed",
                "failure_rate": 0.0738571777154862
            }, {
                "ID": 14,
                "CreatedAt": "2020-09-29T03:17:31.066043189Z",
                "UpdatedAt": "2020-09-29T03:17:31.066043189Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 193,
                "name": "Load Cycle Count",
                "value": 100,
                "worst": 100,
                "thresh": 0,
                "raw_value": 1680,
                "raw_string": "1680",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 15,
                "CreatedAt": "2020-09-29T03:17:31.066162028Z",
                "UpdatedAt": "2020-09-29T03:17:31.066162028Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 194,
                "name": "Temperature",
                "value": 39,
                "worst": 54,
                "thresh": 0,
                "raw_value": 81604378663,
                "raw_string": "39 (0 19 0 0 0)",
                "when_failed": "",
                "transformed_value": 39,
                "status": "passed"
            }, {
                "ID": 16,
                "CreatedAt": "2020-09-29T03:17:31.066270283Z",
                "UpdatedAt": "2020-09-29T03:17:31.066270283Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 195,
                "name": "Hardware ECC Recovered",
                "value": 79,
                "worst": 64,
                "thresh": 0,
                "raw_value": 78022392,
                "raw_string": "78022392",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 17,
                "CreatedAt": "2020-09-29T03:17:31.066381938Z",
                "UpdatedAt": "2020-09-29T03:17:31.066381938Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 197,
                "name": "Current Pending Sector Count",
                "value": 100,
                "worst": 100,
                "thresh": 0,
                "raw_value": 0,
                "raw_string": "0",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed",
                "failure_rate": 0.025540791394761345
            }, {
                "ID": 18,
                "CreatedAt": "2020-09-29T03:17:31.066497143Z",
                "UpdatedAt": "2020-09-29T03:17:31.066497143Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 198,
                "name": "(Offline) Uncorrectable Sector Count",
                "value": 100,
                "worst": 100,
                "thresh": 0,
                "raw_value": 0,
                "raw_string": "0",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed",
                "failure_rate": 0.028675322159886437
            }, {
                "ID": 19,
                "CreatedAt": "2020-09-29T03:17:31.066632808Z",
                "UpdatedAt": "2020-09-29T03:17:31.066632808Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 199,
                "name": "UltraDMA CRC Error Count",
                "value": 200,
                "worst": 200,
                "thresh": 0,
                "raw_value": 0,
                "raw_string": "0",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 20,
                "CreatedAt": "2020-09-29T03:17:31.066757646Z",
                "UpdatedAt": "2020-09-29T03:17:31.066757646Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 240,
                "name": "Head Flying Hours",
                "value": 100,
                "worst": 253,
                "thresh": 0,
                "raw_value": 21887153349770,
                "raw_string": "9354 (19 232 0)",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 21,
                "CreatedAt": "2020-09-29T03:17:31.06687199Z",
                "UpdatedAt": "2020-09-29T03:17:31.06687199Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 241,
                "name": "Total LBAs Written",
                "value": 100,
                "worst": 253,
                "thresh": 0,
                "raw_value": 78115437520,
                "raw_string": "78115437520",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 22,
                "CreatedAt": "2020-09-29T03:17:31.066984948Z",
                "UpdatedAt": "2020-09-29T03:17:31.066984948Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 242,
                "name": "Total LBAs Read",
                "value": 100,
                "worst": 253,
                "thresh": 0,
                "raw_value": 226428055205,
                "raw_string": "226428055205",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }],
            "nvme_attributes": [],
            "scsi_attributes": []
        }, {
            "ID": 1,
            "CreatedAt": "2020-09-29T03:17:31.063859162Z",
            "UpdatedAt": "2020-09-29T03:17:31.063859162Z",
            "DeletedAt": null,
            "device_wwn": "0x5000c500c0558f02",
            "date": "2020-09-29T03:17:30Z",
            "smart_status": "passed",
            "temp": 39,
            "power_on_hours": 9499,
            "power_cycle_count": 22,
            "ata_attributes": [{
                "ID": 1,
                "CreatedAt": "2020-09-29T03:17:31.064174997Z",
                "UpdatedAt": "2020-09-29T03:17:31.064174997Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 1,
                "name": "Read Error Rate",
                "value": 79,
                "worst": 64,
                "thresh": 44,
                "raw_value": 78022392,
                "raw_string": "78022392",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 2,
                "CreatedAt": "2020-09-29T03:17:31.064379412Z",
                "UpdatedAt": "2020-09-29T03:17:31.064379412Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 3,
                "name": "Spin-Up Time",
                "value": 95,
                "worst": 86,
                "thresh": 0,
                "raw_value": 0,
                "raw_string": "0",
                "when_failed": "",
                "transformed_value": 0,
                "status": "warn",
                "status_reason": "Observed Failure Rate for Attribute is greater than 10%",
                "failure_rate": 0.11452195377351217
            }, {
                "ID": 3,
                "CreatedAt": "2020-09-29T03:17:31.064506775Z",
                "UpdatedAt": "2020-09-29T03:17:31.064506775Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 4,
                "name": "Start/Stop Count",
                "value": 100,
                "worst": 100,
                "thresh": 20,
                "raw_value": 220,
                "raw_string": "220",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 4,
                "CreatedAt": "2020-09-29T03:17:31.064621265Z",
                "UpdatedAt": "2020-09-29T03:17:31.064621265Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 5,
                "name": "Reallocated Sectors Count",
                "value": 100,
                "worst": 100,
                "thresh": 10,
                "raw_value": 0,
                "raw_string": "0",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed",
                "failure_rate": 0.025169175350572493
            }, {
                "ID": 5,
                "CreatedAt": "2020-09-29T03:17:31.064742769Z",
                "UpdatedAt": "2020-09-29T03:17:31.064742769Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 7,
                "name": "Seek Error Rate",
                "value": 86,
                "worst": 60,
                "thresh": 45,
                "raw_value": 4679507461,
                "raw_string": "4679507461",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed",
                "failure_rate": 0.08725919610118257
            }, {
                "ID": 6,
                "CreatedAt": "2020-09-29T03:17:31.064850152Z",
                "UpdatedAt": "2020-09-29T03:17:31.064850152Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 9,
                "name": "Power-On Hours",
                "value": 90,
                "worst": 90,
                "thresh": 0,
                "raw_value": 9499,
                "raw_string": "9499",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 7,
                "CreatedAt": "2020-09-29T03:17:31.064970318Z",
                "UpdatedAt": "2020-09-29T03:17:31.064970318Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 10,
                "name": "Spin Retry Count",
                "value": 100,
                "worst": 100,
                "thresh": 97,
                "raw_value": 0,
                "raw_string": "0",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed",
                "failure_rate": 0.05459827163896099
            }, {
                "ID": 8,
                "CreatedAt": "2020-09-29T03:17:31.065072898Z",
                "UpdatedAt": "2020-09-29T03:17:31.065072898Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 12,
                "name": "Power Cycle Count",
                "value": 100,
                "worst": 100,
                "thresh": 20,
                "raw_value": 22,
                "raw_string": "22",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed",
                "failure_rate": 0.038210930067894826
            }, {
                "ID": 9,
                "CreatedAt": "2020-09-29T03:17:31.065366547Z",
                "UpdatedAt": "2020-09-29T03:17:31.065366547Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 18,
                "name": "Unknown_Attribute",
                "value": 100,
                "worst": 100,
                "thresh": 50,
                "raw_value": 0,
                "raw_string": "0",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 10,
                "CreatedAt": "2020-09-29T03:17:31.065539403Z",
                "UpdatedAt": "2020-09-29T03:17:31.065539403Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 187,
                "name": "Reported Uncorrectable Errors",
                "value": 100,
                "worst": 100,
                "thresh": 0,
                "raw_value": 0,
                "raw_string": "0",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed",
                "failure_rate": 0.028130798308190524
            }, {
                "ID": 11,
                "CreatedAt": "2020-09-29T03:17:31.065665808Z",
                "UpdatedAt": "2020-09-29T03:17:31.065665808Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 188,
                "name": "Command Timeout",
                "value": 100,
                "worst": 100,
                "thresh": 0,
                "raw_value": 0,
                "raw_string": "0",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed",
                "failure_rate": 0.024893587674442153
            }, {
                "ID": 12,
                "CreatedAt": "2020-09-29T03:17:31.065780779Z",
                "UpdatedAt": "2020-09-29T03:17:31.065780779Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 190,
                "name": "Temperature Difference",
                "value": 61,
                "worst": 46,
                "thresh": 40,
                "raw_value": 740294695,
                "raw_string": "39 (Min/Max 32/44)",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 13,
                "CreatedAt": "2020-09-29T03:17:31.065931499Z",
                "UpdatedAt": "2020-09-29T03:17:31.065931499Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 192,
                "name": "Power-off Retract Count",
                "value": 100,
                "worst": 100,
                "thresh": 0,
                "raw_value": 4,
                "raw_string": "4",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed",
                "failure_rate": 0.0738571777154862
            }, {
                "ID": 14,
                "CreatedAt": "2020-09-29T03:17:31.066043189Z",
                "UpdatedAt": "2020-09-29T03:17:31.066043189Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 193,
                "name": "Load Cycle Count",
                "value": 100,
                "worst": 100,
                "thresh": 0,
                "raw_value": 1680,
                "raw_string": "1680",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 15,
                "CreatedAt": "2020-09-29T03:17:31.066162028Z",
                "UpdatedAt": "2020-09-29T03:17:31.066162028Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 194,
                "name": "Temperature",
                "value": 39,
                "worst": 54,
                "thresh": 0,
                "raw_value": 81604378663,
                "raw_string": "39 (0 19 0 0 0)",
                "when_failed": "",
                "transformed_value": 39,
                "status": "passed"
            }, {
                "ID": 16,
                "CreatedAt": "2020-09-29T03:17:31.066270283Z",
                "UpdatedAt": "2020-09-29T03:17:31.066270283Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 195,
                "name": "Hardware ECC Recovered",
                "value": 79,
                "worst": 64,
                "thresh": 0,
                "raw_value": 78022392,
                "raw_string": "78022392",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 17,
                "CreatedAt": "2020-09-29T03:17:31.066381938Z",
                "UpdatedAt": "2020-09-29T03:17:31.066381938Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 197,
                "name": "Current Pending Sector Count",
                "value": 100,
                "worst": 100,
                "thresh": 0,
                "raw_value": 0,
                "raw_string": "0",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed",
                "failure_rate": 0.025540791394761345
            }, {
                "ID": 18,
                "CreatedAt": "2020-09-29T03:17:31.066497143Z",
                "UpdatedAt": "2020-09-29T03:17:31.066497143Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 198,
                "name": "(Offline) Uncorrectable Sector Count",
                "value": 100,
                "worst": 100,
                "thresh": 0,
                "raw_value": 0,
                "raw_string": "0",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed",
                "failure_rate": 0.028675322159886437
            }, {
                "ID": 19,
                "CreatedAt": "2020-09-29T03:17:31.066632808Z",
                "UpdatedAt": "2020-09-29T03:17:31.066632808Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 199,
                "name": "UltraDMA CRC Error Count",
                "value": 200,
                "worst": 200,
                "thresh": 0,
                "raw_value": 0,
                "raw_string": "0",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 20,
                "CreatedAt": "2020-09-29T03:17:31.066757646Z",
                "UpdatedAt": "2020-09-29T03:17:31.066757646Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 240,
                "name": "Head Flying Hours",
                "value": 100,
                "worst": 253,
                "thresh": 0,
                "raw_value": 21887153349770,
                "raw_string": "9354 (19 232 0)",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 21,
                "CreatedAt": "2020-09-29T03:17:31.06687199Z",
                "UpdatedAt": "2020-09-29T03:17:31.06687199Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 241,
                "name": "Total LBAs Written",
                "value": 100,
                "worst": 253,
                "thresh": 0,
                "raw_value": 78115437520,
                "raw_string": "78115437520",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }, {
                "ID": 22,
                "CreatedAt": "2020-09-29T03:17:31.066984948Z",
                "UpdatedAt": "2020-09-29T03:17:31.066984948Z",
                "DeletedAt": null,
                "smart_id": 1,
                "attribute_id": 242,
                "name": "Total LBAs Read",
                "value": 100,
                "worst": 253,
                "thresh": 0,
                "raw_value": 226428055205,
                "raw_string": "226428055205",
                "when_failed": "",
                "transformed_value": 0,
                "status": "passed"
            }],
            "nvme_attributes": [],
            "scsi_attributes": []
        }]
    }, "metadata": {
        "1": {
            "ideal": "low",
            "critical": false,
            "description": "(Vendor specific raw value.) Stores data related to the rate of hardware read errors that occurred when reading data from a disk surface. The raw value has different structure for different vendors and is often not meaningful as a decimal number.",
            "observed_thresholds": [{
                "low": 80,
                "high": 95,
                "annual_failure_rate": 0.8879749768303985,
                "error_interval": [0.682344353388663, 1.136105732920724]
            }, {
                "low": 95,
                "high": 110,
                "annual_failure_rate": 0.034155719633986996,
                "error_interval": [0.030188482024981093, 0.038499386872354435]
            }, {
                "low": 110,
                "high": 125,
                "annual_failure_rate": 0.06390002135229157,
                "error_interval": [0.05852004676110847, 0.06964160930553712]
            }, {"low": 125, "high": 140, "annual_failure_rate": 0, "error_interval": [0, 0]}, {
                "low": 140,
                "high": 155,
                "annual_failure_rate": 0,
                "error_interval": [0, 0]
            }, {"low": 155, "high": 170, "annual_failure_rate": 0, "error_interval": [0, 0]}, {
                "low": 170,
                "high": 185,
                "annual_failure_rate": 0,
                "error_interval": [0, 0]
            }, {
                "low": 185,
                "high": 200,
                "annual_failure_rate": 0.044823775021490854,
                "error_interval": [0.032022762038723306, 0.06103725943096589]
            }],
            "display_type": "normalized"
        },
        "10": {
            "ideal": "low",
            "critical": true,
            "description": "Count of retry of spin start attempts. This attribute stores a total count of the spin start attempts to reach the fully operational speed (under the condition that the first attempt was unsuccessful). An increase of this attribute value is a sign of problems in the hard disk mechanical subsystem.",
            "observed_thresholds": [{
                "low": 0,
                "high": 0,
                "annual_failure_rate": 0.05459827163896099,
                "error_interval": [0.05113785787727033, 0.05823122757702782]
            }, {
                "low": 0,
                "high": 80,
                "annual_failure_rate": 0.5555555555555556,
                "error_interval": [0.014065448880161053, 3.095357439410498]
            }],
            "display_type": "raw"
        },
        "11": {
            "ideal": "low",
            "critical": false,
            "description": "This attribute indicates the count that recalibration was requested (under the condition that the first attempt was unsuccessful). An increase of this attribute value is a sign of problems in the hard disk mechanical subsystem.",
            "observed_thresholds": [{
                "low": 0,
                "high": 0,
                "annual_failure_rate": 0.04658866433672694,
                "error_interval": [0.03357701137320878, 0.06297433993055492]
            }, {
                "low": 0,
                "high": 80,
                "annual_failure_rate": 0.5555555555555556,
                "error_interval": [0.014065448880161053, 3.095357439410498]
            }, {"low": 80, "high": 160, "annual_failure_rate": 0, "error_interval": [0, 0]}, {
                "low": 160,
                "high": 240,
                "annual_failure_rate": 0,
                "error_interval": [0, 0]
            }, {"low": 240, "high": 320, "annual_failure_rate": 0, "error_interval": [0, 0]}, {
                "low": 320,
                "high": 400,
                "annual_failure_rate": 0,
                "error_interval": [0, 0]
            }, {"low": 400, "high": 480, "annual_failure_rate": 0, "error_interval": [0, 0]}, {
                "low": 480,
                "high": 560,
                "annual_failure_rate": 0,
                "error_interval": [0, 0]
            }],
            "display_type": "raw"
        },
        "12": {
            "ideal": "low",
            "critical": false,
            "description": "This attribute indicates the count of full hard disk power on/off cycles.",
            "observed_thresholds": [{
                "low": 0,
                "high": 13,
                "annual_failure_rate": 0.019835987118930823,
                "error_interval": [0.016560870164523494, 0.023569242386797896]
            }, {
                "low": 13,
                "high": 26,
                "annual_failure_rate": 0.038210930067894826,
                "error_interval": [0.03353859179329295, 0.0433520775718649]
            }, {
                "low": 26,
                "high": 39,
                "annual_failure_rate": 0.11053528307302571,
                "error_interval": [0.09671061589521368, 0.1257816678419765]
            }, {
                "low": 39,
                "high": 52,
                "annual_failure_rate": 0.16831189443375036,
                "error_interval": [0.1440976510675928, 0.19543066007594895]
            }, {
                "low": 52,
                "high": 65,
                "annual_failure_rate": 0.20630344262550107,
                "error_interval": [0.1693965932069108, 0.2488633537247856]
            }, {
                "low": 65,
                "high": 78,
                "annual_failure_rate": 0.1030972634140512,
                "error_interval": [0.06734655535304743, 0.15106137807407605]
            }, {
                "low": 78,
                "high": 91,
                "annual_failure_rate": 0.12354840389522469,
                "error_interval": [0.06578432170016109, 0.21127153335749593]
            }],
            "display_type": "raw"
        },
        "13": {
            "ideal": "low",
            "critical": false,
            "description": "Uncorrected read errors reported to the operating system.",
            "display_type": "normalized"
        },
        "170": {"ideal": "", "critical": false, "description": "See attribute E8.", "display_type": "normalized"},
        "171": {
            "ideal": "",
            "critical": false,
            "description": "(Kingston) The total number of flash program operation failures since the drive was deployed.[33] Identical to attribute 181.",
            "display_type": "normalized"
        },
        "172": {
            "ideal": "",
            "critical": false,
            "description": "(Kingston) Counts the number of flash erase failures. This attribute returns the total number of Flash erase operation failures since the drive was deployed. This attribute is identical to attribute 182.",
            "display_type": "normalized"
        },
        "173": {
            "ideal": "",
            "critical": false,
            "description": "Counts the maximum worst erase count on any block.",
            "display_type": "normalized"
        },
        "174": {
            "ideal": "",
            "critical": false,
            "description": "Also known as Power-off Retract Count per conventional HDD terminology. Raw value reports the number of unclean shutdowns, cumulative over the life of an SSD, where an unclean shutdown is the removal of power without STANDBY IMMEDIATE as the last command (regardless of PLI activity using capacitor power). Normalized value is always 100.",
            "display_type": ""
        },
        "175": {
            "ideal": "",
            "critical": false,
            "description": "Last test result as microseconds to discharge cap, saturated at its maximum value. Also logs minutes since last test and lifetime number of tests. Raw value contains the following data:     Bytes 0-1: Last test result as microseconds to discharge cap, saturates at max value. Test result expected in range 25 \\u003c= result \\u003c= 5000000, lower indicates specific error code. Bytes 2-3: Minutes since last test, saturates at max value.Bytes 4-5: Lifetime number of tests, not incremented on power cycle, saturates at max value. Normalized value is set to one on test failure or 11 if the capacitor has been tested in an excessive temperature condition, otherwise 100.",
            "display_type": "normalized"
        },
        "176": {
            "ideal": "",
            "critical": false,
            "description": "S.M.A.R.T. parameter indicates a number of flash erase command failures.",
            "display_type": "normalized"
        },
        "177": {
            "ideal": "",
            "critical": false,
            "description": "Delta between most-worn and least-worn Flash blocks. It describes how good/bad the wearleveling of the SSD works on a more technical way. ",
            "display_type": "normalized"
        },
        "179": {
            "ideal": "",
            "critical": false,
            "description": "Pre-Fail attribute used at least in Samsung devices.",
            "display_type": "normalized"
        },
        "180": {
            "ideal": "",
            "critical": false,
            "description": "Pre-Fail attribute used at least in HP devices. ",
            "display_type": "normalized"
        },
        "181": {
            "ideal": "",
            "critical": false,
            "description": "Total number of Flash program operation failures since the drive was deployed.",
            "display_type": "normalized"
        },
        "182": {
            "ideal": "",
            "critical": false,
            "description": "Pre-Fail Attribute used at least in Samsung devices.",
            "display_type": "normalized"
        },
        "183": {
            "ideal": "low",
            "critical": false,
            "description": "Western Digital, Samsung or Seagate attribute: Either the number of downshifts of link speed (e.g. from 6Gbit/s to 3Gbit/s) or the total number of data blocks with detected, uncorrectable errors encountered during normal operation. Although degradation of this parameter can be an indicator of drive aging and/or potential electromechanical problems, it does not directly indicate imminent drive failure.",
            "observed_thresholds": [{
                "low": 0,
                "high": 0,
                "annual_failure_rate": 0.09084549203210031,
                "error_interval": [0.08344373475686712, 0.09872777224842152]
            }, {
                "low": 1,
                "high": 2,
                "annual_failure_rate": 0.05756065656498585,
                "error_interval": [0.04657000847949464, 0.07036491775108872]
            }, {
                "low": 2,
                "high": 4,
                "annual_failure_rate": 0.6193088626208925,
                "error_interval": [0.41784508895529787, 0.8841019099092139]
            }, {
                "low": 4,
                "high": 8,
                "annual_failure_rate": 0.5533447034299792,
                "error_interval": [0.31628430884775033, 0.8985971312402635]
            }, {
                "low": 8,
                "high": 16,
                "annual_failure_rate": 0.3882388694727245,
                "error_interval": [0.21225380267814295, 0.6513988534774338]
            }, {
                "low": 16,
                "high": 35,
                "annual_failure_rate": 0.37116708385481856,
                "error_interval": [0.19763084005134446, 0.6347070173754686]
            }, {
                "low": 35,
                "high": 70,
                "annual_failure_rate": 0.2561146752205292,
                "error_interval": [0.10297138269895259, 0.5276941165819332]
            }, {
                "low": 70,
                "high": 130,
                "annual_failure_rate": 0.40299684542586756,
                "error_interval": [0.16202563309223209, 0.8303275247667772]
            }, {"low": 130, "high": 260, "annual_failure_rate": 0, "error_interval": [0, 0]}],
            "display_type": "raw"
        },
        "184": {
            "ideal": "low",
            "critical": true,
            "description": "This attribute is a part of Hewlett-Packards SMART IV technology, as well as part of other vendors IO Error Detection and Correction schemas, and it contains a count of parity errors which occur in the data path to the media via the drives cache RAM",
            "observed_thresholds": [{
                "low": 93,
                "high": 94,
                "annual_failure_rate": 1.631212012870933,
                "error_interval": [1.055634407303844, 2.407990716767714]
            }, {"low": 94, "high": 95, "annual_failure_rate": 0, "error_interval": [0, 0]}, {
                "low": 95,
                "high": 96,
                "annual_failure_rate": 0,
                "error_interval": [0, 0]
            }, {"low": 96, "high": 97, "annual_failure_rate": 0, "error_interval": [0, 0]}, {
                "low": 97,
                "high": 97,
                "annual_failure_rate": 0,
                "error_interval": [0, 0]
            }, {
                "low": 97,
                "high": 98,
                "annual_failure_rate": 1.8069306930693072,
                "error_interval": [0.04574752432804858, 10.067573453924245]
            }, {
                "low": 98,
                "high": 99,
                "annual_failure_rate": 0.8371559633027523,
                "error_interval": [0.10138347095016888, 3.0240951820174824]
            }, {
                "low": 99,
                "high": 100,
                "annual_failure_rate": 0.09334816849865138,
                "error_interval": [0.08689499010435861, 0.10015372448181788]
            }],
            "display_type": "normalized"
        },
        "185": {
            "ideal": "",
            "critical": false,
            "description": "Western Digital attribute.",
            "display_type": "normalized"
        },
        "186": {
            "ideal": "",
            "critical": false,
            "description": "Western Digital attribute.",
            "display_type": "normalized"
        },
        "187": {
            "ideal": "low",
            "critical": true,
            "description": "The count of errors that could not be recovered using hardware ECC (see attribute 195).",
            "observed_thresholds": [{
                "low": 0,
                "high": 0,
                "annual_failure_rate": 0.028130798308190524,
                "error_interval": [0.024487830609364304, 0.032162944988161336]
            }, {
                "low": 1,
                "high": 1,
                "annual_failure_rate": 0.33877621175661743,
                "error_interval": [0.22325565823630591, 0.4929016016666955]
            }, {
                "low": 1,
                "high": 3,
                "annual_failure_rate": 0.24064820598237213,
                "error_interval": [0.14488594021076606, 0.3758019832614595]
            }, {
                "low": 3,
                "high": 6,
                "annual_failure_rate": 0.5014425058387142,
                "error_interval": [0.3062941096766342, 0.7744372808405151]
            }, {
                "low": 6,
                "high": 11,
                "annual_failure_rate": 0.38007108544136836,
                "error_interval": [0.2989500188963677, 0.4764223967570595]
            }, {
                "low": 11,
                "high": 20,
                "annual_failure_rate": 0.5346094598348444,
                "error_interval": [0.40595137663302483, 0.6911066985735377]
            }, {
                "low": 20,
                "high": 35,
                "annual_failure_rate": 0.8428063943161636,
                "error_interval": [0.6504601819243522, 1.0742259350903411]
            }, {
                "low": 35,
                "high": 65,
                "annual_failure_rate": 1.4429071005017484,
                "error_interval": [1.1405581860945952, 1.8008133631629157]
            }, {
                "low": 65,
                "high": 120,
                "annual_failure_rate": 1.6190935390549661,
                "error_interval": [1.0263664163011208, 2.4294352761068576]
            }],
            "display_type": "raw"
        },
        "188": {
            "ideal": "low",
            "critical": true,
            "description": "The count of aborted operations due to HDD timeout. Normally this attribute value should be equal to zero.",
            "observed_thresholds": [{
                "low": 0,
                "high": 0,
                "annual_failure_rate": 0.024893587674442153,
                "error_interval": [0.020857343769186413, 0.0294830350167543]
            }, {
                "low": 0,
                "high": 13,
                "annual_failure_rate": 0.10044174089362015,
                "error_interval": [0.0812633664077498, 0.1227848196758574]
            }, {
                "low": 13,
                "high": 26,
                "annual_failure_rate": 0.334030592234279,
                "error_interval": [0.2523231196342665, 0.4337665082489293]
            }, {
                "low": 26,
                "high": 39,
                "annual_failure_rate": 0.36724705400842445,
                "error_interval": [0.30398009356575617, 0.4397986538328568]
            }, {
                "low": 39,
                "high": 52,
                "annual_failure_rate": 0.29848155926978354,
                "error_interval": [0.2509254838615984, 0.35242890006477073]
            }, {
                "low": 52,
                "high": 65,
                "annual_failure_rate": 0.2203079701535098,
                "error_interval": [0.18366082845676174, 0.26212468677179274]
            }, {
                "low": 65,
                "high": 78,
                "annual_failure_rate": 0.3018169948863018,
                "error_interval": [0.23779746376787655, 0.37776897542831006]
            }, {
                "low": 78,
                "high": 91,
                "annual_failure_rate": 0.32854928239235887,
                "error_interval": [0.2301118782147336, 0.4548506948185028]
            }, {
                "low": 91,
                "high": 104,
                "annual_failure_rate": 0.28488916640649387,
                "error_interval": [0.1366154288236293, 0.5239213202729072]
            }],
            "display_type": "raw"
        },
        "189": {
            "ideal": "low",
            "critical": false,
            "description": "HDD manufacturers implement a flying height sensor that attempts to provide additional protections for write operations by detecting when a recording head is flying outside its normal operating range. If an unsafe fly height condition is encountered, the write process is stopped, and the information is rewritten or reallocated to a safe region of the hard drive. This attribute indicates the count of these errors detected over the lifetime of the drive.",
            "observed_thresholds": [{
                "low": 0,
                "high": 0,
                "annual_failure_rate": 0.09070551401946862,
                "error_interval": [0.08018892683853401, 0.10221801211956287]
            }, {
                "low": 1,
                "high": 2,
                "annual_failure_rate": 0.0844336097370013,
                "error_interval": [0.07299813695315267, 0.09715235540340669]
            }, {
                "low": 2,
                "high": 5,
                "annual_failure_rate": 0.07943219628781906,
                "error_interval": [0.06552176680630226, 0.09542233189887633]
            }, {
                "low": 5,
                "high": 13,
                "annual_failure_rate": 0.09208847603893404,
                "error_interval": [0.07385765060838133, 0.11345557807163456]
            }, {
                "low": 13,
                "high": 30,
                "annual_failure_rate": 0.18161161650924224,
                "error_interval": [0.13858879602902988, 0.23377015012749933]
            }, {
                "low": 30,
                "high": 70,
                "annual_failure_rate": 0.2678117886102384,
                "error_interval": [0.19044036194841887, 0.36610753129699186]
            }, {
                "low": 70,
                "high": 150,
                "annual_failure_rate": 0.26126480798826107,
                "error_interval": [0.15958733218826962, 0.4035023060905559]
            }, {
                "low": 150,
                "high": 350,
                "annual_failure_rate": 0.11337164155924832,
                "error_interval": [0.030889956621649995, 0.2902764300762812]
            }],
            "display_type": "raw"
        },
        "190": {
            "ideal": "",
            "critical": false,
            "description": "Value is equal to (100-temp. °C), allowing manufacturer to set a minimum threshold which corresponds to a maximum temperature. This also follows the convention of 100 being a best-case value and lower values being undesirable. However, some older drives may instead report raw Temperature (identical to 0xC2) or Temperature minus 50 here.",
            "display_type": "normalized"
        },
        "191": {
            "ideal": "low",
            "critical": false,
            "description": "The count of errors resulting from externally induced shock and vibration. ",
            "display_type": "normalized"
        },
        "192": {
            "ideal": "low",
            "critical": false,
            "description": "Number of power-off or emergency retract cycles.",
            "observed_thresholds": [{
                "low": 1,
                "high": 2,
                "annual_failure_rate": 0.02861098445412803,
                "error_interval": [0.022345416230915037, 0.036088863823297186]
            }, {
                "low": 2,
                "high": 6,
                "annual_failure_rate": 0.0738571777154862,
                "error_interval": [0.06406927746420421, 0.0847175264009771]
            }, {
                "low": 6,
                "high": 16,
                "annual_failure_rate": 0.11970378206823593,
                "error_interval": [0.10830059875098269, 0.13198105985656441]
            }, {
                "low": 16,
                "high": 40,
                "annual_failure_rate": 0.027266868552620425,
                "error_interval": [0.021131448605713823, 0.03462795920968522]
            }, {
                "low": 40,
                "high": 100,
                "annual_failure_rate": 0.011741682974559688,
                "error_interval": [0.00430899071133239, 0.025556700631152028]
            }, {
                "low": 100,
                "high": 250,
                "annual_failure_rate": 0.012659940134091309,
                "error_interval": [0.00607093338127348, 0.023282080653656938]
            }, {
                "low": 250,
                "high": 650,
                "annual_failure_rate": 0.01634692899031039,
                "error_interval": [0.009522688540043157, 0.026173016865409605]
            }, {
                "low": 650,
                "high": 1600,
                "annual_failure_rate": 0.005190074354440066,
                "error_interval": [0.0025908664180103293, 0.009286476666453648]
            }],
            "display_type": "raw"
        },
        "193": {
            "ideal": "low",
            "critical": false,
            "description": "Count of load/unload cycles into head landing zone position.[45] Some drives use 225 (0xE1) for Load Cycle Count instead.",
            "display_type": "normalized"
        },
        "194": {
            "ideal": "low",
            "critical": false,
            "description": "Indicates the device temperature, if the appropriate sensor is fitted. Lowest byte of the raw value contains the exact temperature value (Celsius degrees).",
            "transform_value_unit": "°C",
            "display_type": "transformed"
        },
        "195": {
            "ideal": "",
            "critical": false,
            "description": "(Vendor-specific raw value.) The raw value has different structure for different vendors and is often not meaningful as a decimal number.",
            "observed_thresholds": [{
                "low": 12,
                "high": 24,
                "annual_failure_rate": 0.31472916829975706,
                "error_interval": [0.15711166685282174, 0.5631374192486645]
            }, {
                "low": 24,
                "high": 36,
                "annual_failure_rate": 0.15250310197260136,
                "error_interval": [0.10497611828070175, 0.21417105521823687]
            }, {
                "low": 36,
                "high": 48,
                "annual_failure_rate": 0.2193119102723874,
                "error_interval": [0.16475385681835103, 0.28615447006525274]
            }, {
                "low": 48,
                "high": 60,
                "annual_failure_rate": 0.05672658497265746,
                "error_interval": [0.043182904776447234, 0.07317316161437043]
            }, {"low": 60, "high": 72, "annual_failure_rate": 0, "error_interval": [0, 0]}, {
                "low": 72,
                "high": 84,
                "annual_failure_rate": 0,
                "error_interval": [0, 0]
            }, {"low": 84, "high": 96, "annual_failure_rate": 0, "error_interval": [0, 0]}, {
                "low": 96,
                "high": 108,
                "annual_failure_rate": 0.04074570216566197,
                "error_interval": [0.001031591863615295, 0.22702052218047528]
            }],
            "display_type": "normalized"
        },
        "196": {
            "ideal": "low",
            "critical": true,
            "description": "Count of remap operations. The raw value of this attribute shows the total count of attempts to transfer data from reallocated sectors to a spare area. Both successful and unsuccessful attempts are counted.",
            "observed_thresholds": [{
                "low": 0,
                "high": 0,
                "annual_failure_rate": 0.007389855800729792,
                "error_interval": [0.005652654139732716, 0.009492578928212054]
            }, {
                "low": 1,
                "high": 1,
                "annual_failure_rate": 0.026558331312151347,
                "error_interval": [0.005476966404484466, 0.07761471429677293]
            }, {
                "low": 1,
                "high": 2,
                "annual_failure_rate": 0.02471894893674658,
                "error_interval": [0.0006258296027540169, 0.13772516847438018]
            }, {
                "low": 2,
                "high": 4,
                "annual_failure_rate": 0.03200912040691046,
                "error_interval": [0.0008104007642081744, 0.17834340416493005]
            }, {
                "low": 4,
                "high": 7,
                "annual_failure_rate": 0.043078012510326925,
                "error_interval": [0.001090640849081295, 0.24001532369794615]
            }, {
                "low": 7,
                "high": 11,
                "annual_failure_rate": 0.033843300880853036,
                "error_interval": [0.0008568381932559863, 0.18856280368036135]
            }, {
                "low": 11,
                "high": 17,
                "annual_failure_rate": 0.16979376647542252,
                "error_interval": [0.035015556653263225, 0.49620943874336304]
            }, {
                "low": 17,
                "high": 27,
                "annual_failure_rate": 0.059042381106438044,
                "error_interval": [0.0014948236677880642, 0.32896309247698113]
            }, {
                "low": 27,
                "high": 45,
                "annual_failure_rate": 0.24701105346266636,
                "error_interval": [0.050939617608142244, 0.721871118983972]
            }],
            "display_type": "raw"
        },
        "197": {
            "ideal": "low",
            "critical": true,
            "description": "Count of unstable sectors (waiting to be remapped, because of unrecoverable read errors). If an unstable sector is subsequently read successfully, the sector is remapped and this value is decreased. Read errors on a sector will not remap the sector immediately (since the correct value cannot be read and so the value to remap is not known, and also it might become readable later); instead, the drive firmware remembers that the sector needs to be remapped, and will remap it the next time its written.",
            "observed_thresholds": [{
                "low": 0,
                "high": 0,
                "annual_failure_rate": 0.025540791394761345,
                "error_interval": [0.023161777231213983, 0.02809784482748174]
            }, {
                "low": 1,
                "high": 2,
                "annual_failure_rate": 0.34196613799103254,
                "error_interval": [0.22723401523750225, 0.4942362818474496]
            }, {
                "low": 2,
                "high": 6,
                "annual_failure_rate": 0.6823772508117681,
                "error_interval": [0.41083568090070416, 1.0656166047061635]
            }, {
                "low": 6,
                "high": 16,
                "annual_failure_rate": 0.6108100007493069,
                "error_interval": [0.47336936083368364, 0.7757071095273286]
            }, {
                "low": 16,
                "high": 40,
                "annual_failure_rate": 0.9564879341127684,
                "error_interval": [0.7701044196378299, 1.174355230793638]
            }, {
                "low": 40,
                "high": 100,
                "annual_failure_rate": 1.6519989942167461,
                "error_interval": [1.328402276482456, 2.0305872327541317]
            }, {
                "low": 100,
                "high": 250,
                "annual_failure_rate": 2.5137741046831956,
                "error_interval": [1.9772427971560862, 3.1510376077891613]
            }, {
                "low": 250,
                "high": 650,
                "annual_failure_rate": 3.3203378817413904,
                "error_interval": [2.5883662702274406, 4.195047163573006]
            }, {
                "low": 650,
                "high": 1600,
                "annual_failure_rate": 3.133047210300429,
                "error_interval": [1.1497731080460096, 6.819324775707182]
            }],
            "display_type": "raw"
        },
        "198": {
            "ideal": "low",
            "critical": true,
            "description": "The total count of uncorrectable errors when reading/writing a sector. A rise in the value of this attribute indicates defects of the disk surface and/or problems in the mechanical subsystem.",
            "observed_thresholds": [{
                "low": 0,
                "high": 0,
                "annual_failure_rate": 0.028675322159886437,
                "error_interval": [0.026159385510707116, 0.03136793218577656]
            }, {
                "low": 0,
                "high": 2,
                "annual_failure_rate": 0.8135764944275583,
                "error_interval": [0.40613445471964466, 1.4557130815309443]
            }, {
                "low": 2,
                "high": 4,
                "annual_failure_rate": 1.1173469387755102,
                "error_interval": [0.5773494680315332, 1.9517802404552516]
            }, {
                "low": 4,
                "high": 6,
                "annual_failure_rate": 1.3558692421991083,
                "error_interval": [0.4402470522980859, 3.1641465148237544]
            }, {
                "low": 6,
                "high": 8,
                "annual_failure_rate": 0.7324414715719062,
                "error_interval": [0.15104704003805655, 2.140504796291604]
            }, {
                "low": 8,
                "high": 10,
                "annual_failure_rate": 0.5777213677766163,
                "error_interval": [0.43275294849366835, 0.7556737733062419]
            }, {
                "low": 10,
                "high": 12,
                "annual_failure_rate": 1.7464114832535886,
                "error_interval": [0.47583835092536914, 4.471507017371231]
            }, {
                "low": 12,
                "high": 14,
                "annual_failure_rate": 2.6449275362318843,
                "error_interval": [0.3203129951758959, 9.554387676519005]
            }, {
                "low": 14,
                "high": 16,
                "annual_failure_rate": 0.796943231441048,
                "error_interval": [0.5519063550198366, 1.113648286331181]
            }],
            "display_type": "raw"
        },
        "199": {
            "ideal": "low",
            "critical": false,
            "description": "The count of errors in data transfer via the interface cable as determined by ICRC (Interface Cyclic Redundancy Check).",
            "observed_thresholds": [{
                "low": 0,
                "high": 1,
                "annual_failure_rate": 0.04068379316116366,
                "error_interval": [0.037534031558106425, 0.04402730201866553]
            }, {
                "low": 1,
                "high": 2,
                "annual_failure_rate": 0.1513481259734218,
                "error_interval": [0.12037165605991791, 0.18786293065527596]
            }, {
                "low": 2,
                "high": 4,
                "annual_failure_rate": 0.16849758722418978,
                "error_interval": [0.12976367397863445, 0.2151676572000481]
            }, {
                "low": 4,
                "high": 8,
                "annual_failure_rate": 0.15385127340491614,
                "error_interval": [0.10887431782430312, 0.21117289306426648]
            }, {
                "low": 8,
                "high": 16,
                "annual_failure_rate": 0.14882894050104387,
                "error_interval": [0.09631424312463635, 0.2197008753522735]
            }, {
                "low": 16,
                "high": 35,
                "annual_failure_rate": 0.20878219917249793,
                "error_interval": [0.14086447304552446, 0.29804957135975]
            }, {
                "low": 35,
                "high": 70,
                "annual_failure_rate": 0.13742940270409038,
                "error_interval": [0.06860426267470295, 0.24589916335290812]
            }, {
                "low": 70,
                "high": 130,
                "annual_failure_rate": 0.22336578581363,
                "error_interval": [0.11150339549604707, 0.39966309081252904]
            }, {
                "low": 130,
                "high": 260,
                "annual_failure_rate": 0.18277416124186283,
                "error_interval": [0.07890890989692058, 0.3601379610272007]
            }],
            "display_type": "raw"
        },
        "2": {
            "ideal": "high",
            "critical": false,
            "description": "Overall (general) throughput performance of a hard disk drive. If the value of this attribute is decreasing there is a high probability that there is a problem with the disk.",
            "display_type": "normalized"
        },
        "200": {
            "ideal": "low",
            "critical": false,
            "description": "The count of errors found when writing a sector. The higher the value, the worse the disks mechanical condition is.",
            "display_type": "normalized"
        },
        "201": {
            "ideal": "low",
            "critical": true,
            "description": "Count indicates the number of uncorrectable software read errors.",
            "display_type": "normalized"
        },
        "202": {
            "ideal": "low",
            "critical": false,
            "description": "Count of Data Address Mark errors (or vendor-specific).",
            "display_type": "normalized"
        },
        "203": {
            "ideal": "low",
            "critical": false,
            "description": "The number of errors caused by incorrect checksum during the error correction.",
            "display_type": "normalized"
        },
        "204": {
            "ideal": "low",
            "critical": false,
            "description": "Count of errors corrected by the internal error correction software.",
            "display_type": ""
        },
        "205": {
            "ideal": "low",
            "critical": false,
            "description": "Count of errors due to high temperature.",
            "display_type": "normalized"
        },
        "206": {
            "ideal": "",
            "critical": false,
            "description": "Height of heads above the disk surface. If too low, head crash is more likely; if too high, read/write errors are more likely.",
            "display_type": "normalized"
        },
        "207": {
            "ideal": "low",
            "critical": false,
            "description": "Amount of surge current used to spin up the drive.",
            "display_type": "normalized"
        },
        "208": {
            "ideal": "",
            "critical": false,
            "description": "Count of buzz routines needed to spin up the drive due to insufficient power.",
            "display_type": "normalized"
        },
        "209": {
            "ideal": "",
            "critical": false,
            "description": "Drives seek performance during its internal tests.",
            "display_type": "normalized"
        },
        "210": {
            "ideal": "",
            "critical": false,
            "description": "Found in Maxtor 6B200M0 200GB and Maxtor 2R015H1 15GB disks.",
            "display_type": "normalized"
        },
        "211": {
            "ideal": "",
            "critical": false,
            "description": "A recording of a vibration encountered during write operations.",
            "display_type": "normalized"
        },
        "212": {
            "ideal": "",
            "critical": false,
            "description": "A recording of shock encountered during write operations.",
            "display_type": "normalized"
        },
        "22": {
            "ideal": "high",
            "critical": false,
            "description": "Specific to He8 drives from HGST. This value measures the helium inside of the drive specific to this manufacturer. It is a pre-fail attribute that trips once the drive detects that the internal environment is out of specification.",
            "display_type": "normalized"
        },
        "220": {
            "ideal": "low",
            "critical": false,
            "description": "Distance the disk has shifted relative to the spindle (usually due to shock or temperature). Unit of measure is unknown.",
            "display_type": "normalized"
        },
        "221": {
            "ideal": "low",
            "critical": false,
            "description": "The count of errors resulting from externally induced shock and vibration.",
            "display_type": "normalized"
        },
        "222": {
            "ideal": "",
            "critical": false,
            "description": "Time spent operating under data load (movement of magnetic head armature).",
            "display_type": "normalized"
        },
        "223": {
            "ideal": "",
            "critical": false,
            "description": "Count of times head changes position.",
            "display_type": "normalized"
        },
        "224": {
            "ideal": "low",
            "critical": false,
            "description": "Resistance caused by friction in mechanical parts while operating.",
            "display_type": "normalized"
        },
        "225": {
            "ideal": "low",
            "critical": false,
            "description": "Total count of load cycles Some drives use 193 (0xC1) for Load Cycle Count instead. See Description for 193 for significance of this number. ",
            "display_type": "normalized"
        },
        "226": {
            "ideal": "",
            "critical": false,
            "description": "Total time of loading on the magnetic heads actuator (time not spent in parking area).",
            "display_type": "normalized"
        },
        "227": {
            "ideal": "low",
            "critical": false,
            "description": "Count of attempts to compensate for platter speed variations.[66]",
            "display_type": ""
        },
        "228": {
            "ideal": "low",
            "critical": false,
            "description": "The number of power-off cycles which are counted whenever there is a retract event and the heads are loaded off of the media such as when the machine is powered down, put to sleep, or is idle.",
            "display_type": ""
        },
        "230": {
            "ideal": "",
            "critical": false,
            "description": "Amplitude of thrashing (repetitive head moving motions between operations).",
            "display_type": "normalized"
        },
        "231": {
            "ideal": "",
            "critical": false,
            "description": "Indicates the approximate SSD life left, in terms of program/erase cycles or available reserved blocks. A normalized value of 100 represents a new drive, with a threshold value at 10 indicating a need for replacement. A value of 0 may mean that the drive is operating in read-only mode to allow data recovery.",
            "display_type": "normalized"
        },
        "232": {
            "ideal": "",
            "critical": false,
            "description": "Number of physical erase cycles completed on the SSD as a percentage of the maximum physical erase cycles the drive is designed to endure.",
            "display_type": "normalized"
        },
        "233": {
            "ideal": "",
            "critical": false,
            "description": "Intel SSDs report a normalized value from 100, a new drive, to a minimum of 1. It decreases while the NAND erase cycles increase from 0 to the maximum-rated cycles.",
            "display_type": "normalized"
        },
        "234": {
            "ideal": "",
            "critical": false,
            "description": "Decoded as: byte 0-1-2 = average erase count (big endian) and byte 3-4-5 = max erase count (big endian).",
            "display_type": "normalized"
        },
        "235": {
            "ideal": "",
            "critical": false,
            "description": "Decoded as: byte 0-1-2 = good block count (big endian) and byte 3-4 = system (free) block count.",
            "display_type": "normalized"
        },
        "240": {
            "ideal": "",
            "critical": false,
            "description": "Time spent during the positioning of the drive heads.[15][71] Some Fujitsu drives report the count of link resets during a data transfer.",
            "display_type": "normalized"
        },
        "241": {
            "ideal": "",
            "critical": false,
            "description": "Total count of LBAs written.",
            "display_type": "normalized"
        },
        "242": {
            "ideal": "",
            "critical": false,
            "description": "Total count of LBAs read.Some S.M.A.R.T. utilities will report a negative number for the raw value since in reality it has 48 bits rather than 32.",
            "display_type": "normalized"
        },
        "243": {
            "ideal": "",
            "critical": false,
            "description": "The upper 5 bytes of the 12-byte total number of LBAs written to the device. The lower 7 byte value is located at attribute 0xF1.",
            "display_type": "normalized"
        },
        "244": {
            "ideal": "",
            "critical": false,
            "description": "The upper 5 bytes of the 12-byte total number of LBAs read from the device. The lower 7 byte value is located at attribute 0xF2.",
            "display_type": "normalized"
        },
        "249": {
            "ideal": "",
            "critical": false,
            "description": "Total NAND Writes. Raw value reports the number of writes to NAND in 1 GB increments.",
            "display_type": "normalized"
        },
        "250": {
            "ideal": "low",
            "critical": false,
            "description": "Count of errors while reading from a disk.",
            "display_type": "normalized"
        },
        "251": {
            "ideal": "",
            "critical": false,
            "description": "The Minimum Spares Remaining attribute indicates the number of remaining spare blocks as a percentage of the total number of spare blocks available.",
            "display_type": "normalized"
        },
        "252": {
            "ideal": "",
            "critical": false,
            "description": "The Newly Added Bad Flash Block attribute indicates the total number of bad flash blocks the drive detected since it was first initialized in manufacturing.",
            "display_type": "normalized"
        },
        "254": {
            "ideal": "low",
            "critical": false,
            "description": "Count of Free Fall Events detected.",
            "display_type": "normalized"
        },
        "3": {
            "ideal": "low",
            "critical": false,
            "description": "Average time of spindle spin up (from zero RPM to fully operational [milliseconds]).",
            "observed_thresholds": [{
                "low": 78,
                "high": 96,
                "annual_failure_rate": 0.11452195377351217,
                "error_interval": [0.10591837762295722, 0.12363823501915781]
            }, {
                "low": 96,
                "high": 114,
                "annual_failure_rate": 0.040274562840558074,
                "error_interval": [0.03465055611002801, 0.046551312468303144]
            }, {
                "low": 114,
                "high": 132,
                "annual_failure_rate": 0.009100406705780476,
                "error_interval": [0.006530608971356785, 0.012345729280075591]
            }, {
                "low": 132,
                "high": 150,
                "annual_failure_rate": 0.008561351734020232,
                "error_interval": [0.004273795939256936, 0.015318623141355509]
            }, {
                "low": 150,
                "high": 168,
                "annual_failure_rate": 0.015780508262068848,
                "error_interval": [0.005123888078524015, 0.03682644215646287]
            }, {
                "low": 168,
                "high": 186,
                "annual_failure_rate": 0.05262688124794024,
                "error_interval": [0.0325768689524594, 0.08044577830285578]
            }, {
                "low": 186,
                "high": 204,
                "annual_failure_rate": 0.01957419424036038,
                "error_interval": [0.0023705257325185624, 0.0707087198669825]
            }, {
                "low": 204,
                "high": 222,
                "annual_failure_rate": 0.026050959960031404,
                "error_interval": [0.0006595532020744994, 0.1451466588889228]
            }],
            "display_type": "normalized"
        },
        "4": {
            "ideal": "",
            "critical": false,
            "description": "A tally of spindle start/stop cycles. The spindle turns on, and hence the count is increased, both when the hard disk is turned on after having before been turned entirely off (disconnected from power source) and when the hard disk returns from having previously been put to sleep mode.",
            "observed_thresholds": [{
                "low": 0,
                "high": 13,
                "annual_failure_rate": 0.01989335424860646,
                "error_interval": [0.016596548909440657, 0.023653263230617408]
            }, {
                "low": 13,
                "high": 26,
                "annual_failure_rate": 0.03776935438256488,
                "error_interval": [0.03310396052098642, 0.04290806173460437]
            }, {
                "low": 26,
                "high": 39,
                "annual_failure_rate": 0.11022223828187004,
                "error_interval": [0.09655110535164119, 0.12528657238811672]
            }, {
                "low": 39,
                "high": 52,
                "annual_failure_rate": 0.16289995457762474,
                "error_interval": [0.13926541653588131, 0.18939614504497515]
            }, {
                "low": 52,
                "high": 65,
                "annual_failure_rate": 0.19358212432279714,
                "error_interval": [0.15864522253849073, 0.23392418181765526]
            }, {
                "low": 65,
                "high": 78,
                "annual_failure_rate": 0.1157094940074447,
                "error_interval": [0.07861898732346269, 0.16424039052527728]
            }, {
                "low": 78,
                "high": 91,
                "annual_failure_rate": 0.12262136155304391,
                "error_interval": [0.0670382394080032, 0.20573780888032978]
            }, {"low": 91, "high": 104, "annual_failure_rate": 0, "error_interval": [0, 0]}],
            "display_type": "raw"
        },
        "5": {
            "ideal": "low",
            "critical": true,
            "description": "Count of reallocated sectors. The raw value represents a count of the bad sectors that have been found and remapped.Thus, the higher the attribute value, the more sectors the drive has had to reallocate. This value is primarily used as a metric of the life expectancy of the drive; a drive which has had any reallocations at all is significantly more likely to fail in the immediate months.",
            "observed_thresholds": [{
                "low": 0,
                "high": 0,
                "annual_failure_rate": 0.025169175350572493,
                "error_interval": [0.022768612038746357, 0.027753988579272894]
            }, {
                "low": 1,
                "high": 4,
                "annual_failure_rate": 0.027432608477803388,
                "error_interval": [0.010067283827589948, 0.05970923963096652]
            }, {
                "low": 4,
                "high": 16,
                "annual_failure_rate": 0.07501976284584981,
                "error_interval": [0.039944864177334186, 0.12828607921150972]
            }, {
                "low": 16,
                "high": 70,
                "annual_failure_rate": 0.23589260654405794,
                "error_interval": [0.1643078435800227, 0.32806951196017664]
            }, {
                "low": 70,
                "high": 260,
                "annual_failure_rate": 0.36193219378600433,
                "error_interval": [0.2608488901774093, 0.4892271827875412]
            }, {
                "low": 260,
                "high": 1100,
                "annual_failure_rate": 0.5676621428968173,
                "error_interval": [0.4527895568499355, 0.702804359408436]
            }, {
                "low": 1100,
                "high": 4500,
                "annual_failure_rate": 1.5028253400346423,
                "error_interval": [1.2681757596263297, 1.768305221795894]
            }, {
                "low": 4500,
                "high": 17000,
                "annual_failure_rate": 2.0659987547404763,
                "error_interval": [1.6809790460512237, 2.512808045182302]
            }, {
                "low": 17000,
                "high": 70000,
                "annual_failure_rate": 1.7755385684503124,
                "error_interval": [1.2796520259849835, 2.400012341226441]
            }],
            "display_type": "raw"
        },
        "6": {
            "ideal": "",
            "critical": false,
            "description": "Margin of a channel while reading data. The function of this attribute is not specified.",
            "display_type": "normalized"
        },
        "7": {
            "ideal": "",
            "critical": false,
            "description": "(Vendor specific raw value.) Rate of seek errors of the magnetic heads. If there is a partial failure in the mechanical positioning system, then seek errors will arise. Such a failure may be due to numerous factors, such as damage to a servo, or thermal widening of the hard disk. The raw value has different structure for different vendors and is often not meaningful as a decimal number.",
            "observed_thresholds": [{
                "low": 58,
                "high": 76,
                "annual_failure_rate": 0.2040131025936549,
                "error_interval": [0.17032852883286412, 0.2424096283327138]
            }, {
                "low": 76,
                "high": 94,
                "annual_failure_rate": 0.08725919610118257,
                "error_interval": [0.08077138510999876, 0.09412943212007528]
            }, {
                "low": 94,
                "high": 112,
                "annual_failure_rate": 0.01087335627722523,
                "error_interval": [0.008732197944943352, 0.013380600544561905]
            }, {"low": 112, "high": 130, "annual_failure_rate": 0, "error_interval": [0, 0]}, {
                "low": 130,
                "high": 148,
                "annual_failure_rate": 0,
                "error_interval": [0, 0]
            }, {"low": 148, "high": 166, "annual_failure_rate": 0, "error_interval": [0, 0]}, {
                "low": 166,
                "high": 184,
                "annual_failure_rate": 0,
                "error_interval": [0, 0]
            }, {
                "low": 184,
                "high": 202,
                "annual_failure_rate": 0.05316285755900475,
                "error_interval": [0.03370069132942804, 0.07977038905848267]
            }],
            "display_type": "normalized"
        },
        "8": {
            "ideal": "high",
            "critical": false,
            "description": "Average performance of seek operations of the magnetic heads. If this attribute is decreasing, it is a sign of problems in the mechanical subsystem.",
            "display_type": "normalized"
        },
        "9": {
            "ideal": "",
            "critical": false,
            "description": "Count of hours in power-on state. The raw value of this attribute shows total count of hours (or minutes, or seconds, depending on manufacturer) in power-on state. By default, the total expected lifetime of a hard disk in perfect condition is defined as 5 years (running every day and night on all days). This is equal to 1825 days in 24/7 mode or 43800 hours. On some pre-2005 drives, this raw value may advance erratically and/or wrap around (reset to zero periodically).",
            "display_type": "normalized"
        }
    }, "success": true
}
