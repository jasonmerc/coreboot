/* SPDX-License-Identifier: GPL-2.0-only */

#include <device/azalia_device.h>

const u32 cim_verb_data[] = {
	0x10ec0283,	/* Codec Vendor / Device ID:  */
	0x17aa3806,	/* Subsystem ID */
	4,		/* Number of 4 dword sets */
	AZALIA_SUBVENDOR(0, 0x10ec0283),
	AZALIA_PIN_CFG(0,0x12,0x90a60140),
	AZALIA_PIN_CFG(0,0x14,0x90170110),
	AZALIA_PIN_CFG(0,0x17,0x40000000),
	AZALIA_PIN_CFG(0,0x18,0x411111f0),
	AZALIA_PIN_CFG(0,0x19,0x04a11030),
	AZALIA_PIN_CFG(0,0x1a,0x411111f0),
	AZALIA_PIN_CFG(0,0x1b,0x411111f0),
	AZALIA_PIN_CFG(0,0x1d,0x40e00001),
	AZALIA_PIN_CFG(0,0x1e,0x04452120),
	AZALIA_PIN_CFG(0,0x21,0x0421101f),

};

const u32 pc_beep_verbs[0] = {};

AZALIA_ARRAY_SIZES;
