package core

import "strings"

const StarterRebuildBlueprint = "bp_starter_rebuild"

var catalogPartIDs = map[string]struct{}{
	"part_chassis_starter_prototype": {},
	"part_engine_ice_street_i":       {},
	"part_engine_torque_offroad_i":   {},
	"part_drive_close_ratio_i":       {},
	"part_drive_launch_i":            {},
	"part_suspension_angle_i":        {},
	"part_suspension_travel_i":       {},
	"part_brakes_track_i":            {},
	"part_tires_street_i":            {},
	"part_tires_gravel_i":            {},
	"part_ecu_standard_i":            {},
	"part_ecu_legacy_zero":           {},
	"part_aero_low_drag_i":           {},
	"part_aero_downforce_i":          {},
	"part_interface_delivery_i":      {},
	"part_sensor_salvage_i":          {},
	"part_sensor_race_telemetry_i":   {},
	"part_module_recovery_winch_i":   {},
	"part_module_cargo_lock_i":       {},
	"part_module_story_zero_bridge":  {},
}

func IsCatalogPart(partID string) bool {
	_, ok := catalogPartIDs[strings.TrimSpace(partID)]
	return ok
}
