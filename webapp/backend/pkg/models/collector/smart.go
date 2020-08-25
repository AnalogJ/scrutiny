package collector

type SmartInfo struct {
	JSONFormatVersion []int `json:"json_format_version"`
	Smartctl          struct {
		Version      []int    `json:"version"`
		SvnRevision  string   `json:"svn_revision"`
		PlatformInfo string   `json:"platform_info"`
		BuildInfo    string   `json:"build_info"`
		Argv         []string `json:"argv"`
		ExitStatus   int      `json:"exit_status"`
		Messages     []struct {
			String   string `json:"string"`
			Severity string `json:"severity"`
		} `json:"messages"`
	} `json:"smartctl"`
	Device struct {
		Name     string `json:"name"`
		InfoName string `json:"info_name"`
		Type     string `json:"type"`
		Protocol string `json:"protocol"`
	} `json:"device"`
	ModelName    string `json:"model_name"`
	SerialNumber string `json:"serial_number"`
	Wwn          struct {
		Naa int   `json:"naa"`
		Oui int   `json:"oui"`
		ID  int64 `json:"id"`
	} `json:"wwn"`
	FirmwareVersion string `json:"firmware_version"`
	UserCapacity    struct {
		Blocks int64 `json:"blocks"`
		Bytes  int64 `json:"bytes"`
	} `json:"user_capacity"`
	LogicalBlockSize  int `json:"logical_block_size"`
	PhysicalBlockSize int `json:"physical_block_size"`
	RotationRate      int `json:"rotation_rate"`
	FormFactor        struct {
		AtaValue int    `json:"ata_value"`
		Name     string `json:"name"`
	} `json:"form_factor"`
	InSmartctlDatabase bool `json:"in_smartctl_database"`
	AtaVersion         struct {
		String     string `json:"string"`
		MajorValue int    `json:"major_value"`
		MinorValue int    `json:"minor_value"`
	} `json:"ata_version"`
	SataVersion struct {
		String string `json:"string"`
		Value  int    `json:"value"`
	} `json:"sata_version"`
	InterfaceSpeed struct {
		Max struct {
			SataValue      int    `json:"sata_value"`
			String         string `json:"string"`
			UnitsPerSecond int    `json:"units_per_second"`
			BitsPerUnit    int    `json:"bits_per_unit"`
		} `json:"max"`
		Current struct {
			SataValue      int    `json:"sata_value"`
			String         string `json:"string"`
			UnitsPerSecond int    `json:"units_per_second"`
			BitsPerUnit    int    `json:"bits_per_unit"`
		} `json:"current"`
	} `json:"interface_speed"`
	LocalTime struct {
		TimeT   int64  `json:"time_t"`
		Asctime string `json:"asctime"`
	} `json:"local_time"`
	SmartStatus struct {
		Passed bool `json:"passed"`
	} `json:"smart_status"`

	PowerOnTime struct {
		Hours int64 `json:"hours"`
	} `json:"power_on_time"`
	PowerCycleCount int64 `json:"power_cycle_count"`
	Temperature     struct {
		Current int64 `json:"current"`
	} `json:"temperature"`

	// ATA Protocol Specific Fields
	AtaSmartData struct {
		OfflineDataCollection struct {
			Status struct {
				Value  int    `json:"value"`
				String string `json:"string"`
				Passed bool   `json:"passed"`
			} `json:"status"`
			CompletionSeconds int `json:"completion_seconds"`
		} `json:"offline_data_collection"`
		SelfTest struct {
			Status struct {
				Value            int    `json:"value"`
				String           string `json:"string"`
				RemainingPercent int    `json:"remaining_percent"`
			} `json:"status"`
			PollingMinutes struct {
				Short    int `json:"short"`
				Extended int `json:"extended"`
			} `json:"polling_minutes"`
		} `json:"self_test"`
		Capabilities struct {
			Values                        []int `json:"values"`
			ExecOfflineImmediateSupported bool  `json:"exec_offline_immediate_supported"`
			OfflineIsAbortedUponNewCmd    bool  `json:"offline_is_aborted_upon_new_cmd"`
			OfflineSurfaceScanSupported   bool  `json:"offline_surface_scan_supported"`
			SelfTestsSupported            bool  `json:"self_tests_supported"`
			ConveyanceSelfTestSupported   bool  `json:"conveyance_self_test_supported"`
			SelectiveSelfTestSupported    bool  `json:"selective_self_test_supported"`
			AttributeAutosaveEnabled      bool  `json:"attribute_autosave_enabled"`
			ErrorLoggingSupported         bool  `json:"error_logging_supported"`
			GpLoggingSupported            bool  `json:"gp_logging_supported"`
		} `json:"capabilities"`
	} `json:"ata_smart_data"`
	AtaSctCapabilities struct {
		Value                         int  `json:"value"`
		ErrorRecoveryControlSupported bool `json:"error_recovery_control_supported"`
		FeatureControlSupported       bool `json:"feature_control_supported"`
		DataTableSupported            bool `json:"data_table_supported"`
	} `json:"ata_sct_capabilities"`
	AtaSmartAttributes struct {
		Revision int `json:"revision"`
		Table    []struct {
			ID         int    `json:"id"`
			Name       string `json:"name"`
			Value      int    `json:"value"`
			Worst      int    `json:"worst"`
			Thresh     int    `json:"thresh"`
			WhenFailed string `json:"when_failed"`
			Flags      struct {
				Value         int    `json:"value"`
				String        string `json:"string"`
				Prefailure    bool   `json:"prefailure"`
				UpdatedOnline bool   `json:"updated_online"`
				Performance   bool   `json:"performance"`
				ErrorRate     bool   `json:"error_rate"`
				EventCount    bool   `json:"event_count"`
				AutoKeep      bool   `json:"auto_keep"`
			} `json:"flags"`
			Raw struct {
				Value  int64  `json:"value"`
				String string `json:"string"`
			} `json:"raw"`
		} `json:"table"`
	} `json:"ata_smart_attributes"`
	AtaSmartErrorLog struct {
		Summary struct {
			Revision    int `json:"revision"`
			Count       int `json:"count"`
			LoggedCount int `json:"logged_count"`
			Table       []struct {
				ErrorNumber         int `json:"error_number"`
				LifetimeHours       int `json:"lifetime_hours"`
				CompletionRegisters struct {
					Error  int `json:"error"`
					Status int `json:"status"`
					Count  int `json:"count"`
					Lba    int `json:"lba"`
					Device int `json:"device"`
				} `json:"completion_registers"`
				ErrorDescription string `json:"error_description"`
				PreviousCommands []struct {
					Registers struct {
						Command       int `json:"command"`
						Features      int `json:"features"`
						Count         int `json:"count"`
						Lba           int `json:"lba"`
						Device        int `json:"device"`
						DeviceControl int `json:"device_control"`
					} `json:"registers"`
					PowerupMilliseconds int    `json:"powerup_milliseconds"`
					CommandName         string `json:"command_name"`
				} `json:"previous_commands"`
			} `json:"table"`
		} `json:"summary"`
	} `json:"ata_smart_error_log"`
	AtaSmartSelfTestLog struct {
		Standard struct {
			Revision int `json:"revision"`
			Table    []struct {
				Type struct {
					Value  int    `json:"value"`
					String string `json:"string"`
				} `json:"type"`
				Status struct {
					Value  int    `json:"value"`
					String string `json:"string"`
					Passed bool   `json:"passed"`
				} `json:"status"`
				LifetimeHours int `json:"lifetime_hours"`
			} `json:"table"`
			Count              int `json:"count"`
			ErrorCountTotal    int `json:"error_count_total"`
			ErrorCountOutdated int `json:"error_count_outdated"`
		} `json:"standard"`
	} `json:"ata_smart_self_test_log"`
	AtaSmartSelectiveSelfTestLog struct {
		Revision int `json:"revision"`
		Table    []struct {
			LbaMin int `json:"lba_min"`
			LbaMax int `json:"lba_max"`
			Status struct {
				Value  int    `json:"value"`
				String string `json:"string"`
			} `json:"status"`
		} `json:"table"`
		Flags struct {
			Value                int  `json:"value"`
			RemainderScanEnabled bool `json:"remainder_scan_enabled"`
		} `json:"flags"`
		PowerUpScanResumeMinutes int `json:"power_up_scan_resume_minutes"`
	} `json:"ata_smart_selective_self_test_log"`

	// NVME Protocol Specific Fields
	NvmePciVendor struct {
		ID          int `json:"id"`
		SubsystemID int `json:"subsystem_id"`
	} `json:"nvme_pci_vendor"`
	NvmeIeeeOuiIdentifier  int `json:"nvme_ieee_oui_identifier"`
	NvmeControllerID       int `json:"nvme_controller_id"`
	NvmeNumberOfNamespaces int `json:"nvme_number_of_namespaces"`
	NvmeNamespaces         []struct {
		ID   int `json:"id"`
		Size struct {
			Blocks int   `json:"blocks"`
			Bytes  int64 `json:"bytes"`
		} `json:"size"`
		Capacity struct {
			Blocks int   `json:"blocks"`
			Bytes  int64 `json:"bytes"`
		} `json:"capacity"`
		Utilization struct {
			Blocks int   `json:"blocks"`
			Bytes  int64 `json:"bytes"`
		} `json:"utilization"`
		FormattedLbaSize int `json:"formatted_lba_size"`
	} `json:"nvme_namespaces"`
	NvmeSmartHealthInformationLog struct {
		CriticalWarning         int `json:"critical_warning"`
		Temperature             int `json:"temperature"`
		AvailableSpare          int `json:"available_spare"`
		AvailableSpareThreshold int `json:"available_spare_threshold"`
		PercentageUsed          int `json:"percentage_used"`
		DataUnitsRead           int `json:"data_units_read"`
		DataUnitsWritten        int `json:"data_units_written"`
		HostReads               int `json:"host_reads"`
		HostWrites              int `json:"host_writes"`
		ControllerBusyTime      int `json:"controller_busy_time"`
		PowerCycles             int `json:"power_cycles"`
		PowerOnHours            int `json:"power_on_hours"`
		UnsafeShutdowns         int `json:"unsafe_shutdowns"`
		MediaErrors             int `json:"media_errors"`
		NumErrLogEntries        int `json:"num_err_log_entries"`
		WarningTempTime         int `json:"warning_temp_time"`
		CriticalCompTime        int `json:"critical_comp_time"`
	} `json:"nvme_smart_health_information_log"`

	// SCSI Protocol Specific Fields
	Vendor              string `json:"vendor"`
	Product             string `json:"product"`
	ScsiVersion         string `json:"scsi_version"`
	ScsiGrownDefectList int    `json:"scsi_grown_defect_list"`
	ScsiErrorCounterLog struct {
		Read struct {
			ErrorsCorrectedByEccfast         int    `json:"errors_corrected_by_eccfast"`
			ErrorsCorrectedByEccdelayed      int    `json:"errors_corrected_by_eccdelayed"`
			ErrorsCorrectedByRereadsRewrites int    `json:"errors_corrected_by_rereads_rewrites"`
			TotalErrorsCorrected             int    `json:"total_errors_corrected"`
			CorrectionAlgorithmInvocations   int    `json:"correction_algorithm_invocations"`
			GigabytesProcessed               string `json:"gigabytes_processed"`
			TotalUncorrectedErrors           int    `json:"total_uncorrected_errors"`
		} `json:"read"`
		Write struct {
			ErrorsCorrectedByEccfast         int    `json:"errors_corrected_by_eccfast"`
			ErrorsCorrectedByEccdelayed      int    `json:"errors_corrected_by_eccdelayed"`
			ErrorsCorrectedByRereadsRewrites int    `json:"errors_corrected_by_rereads_rewrites"`
			TotalErrorsCorrected             int    `json:"total_errors_corrected"`
			CorrectionAlgorithmInvocations   int    `json:"correction_algorithm_invocations"`
			GigabytesProcessed               string `json:"gigabytes_processed"`
			TotalUncorrectedErrors           int    `json:"total_uncorrected_errors"`
		} `json:"write"`
	} `json:"scsi_error_counter_log"`
}
