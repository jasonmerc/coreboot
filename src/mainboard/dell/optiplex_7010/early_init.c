/* SPDX-License-Identifier: GPL-2.0-only */

#include <device/pci_ops.h>
#include <bootblock_common.h>
#include <northbridge/intel/sandybridge/raminit_native.h>
#include <southbridge/intel/bd82x6x/pch.h>

const struct southbridge_usb_port mainboard_usb_ports[] = {
	{ 1, 0, 0 },
	{ 1, 0, 0 },
	{ 1, 1, 1 },
	{ 1, 1, 2 },
	{ 1, 1, 3 },
	{ 1, 1, 3 },
	{ 1, 0, -1 },
	{ 1, 0, -1 },
	{ 1, 0, 4 },
	{ 1, 0, 4 },
	{ 1, 0, 5 },
	{ 1, 1, 5 },
	{ 1, 1, 6 },
	{ 1, 0, 6 },
};

void bootblock_mainboard_early_init(void)
{
	pci_write_config16(PCI_DEV(0, 0x1f, 0), 0x82, 0x3c01);
	pci_write_config16(PCI_DEV(0, 0x1f, 0), 0x80, 0x0000);
}

/* FIXME: Put proper SPD map here. */
void mainboard_get_spd(spd_raw_data *spd, bool id_only)
{
	read_spd(&spd[0], 0x50, id_only);
	read_spd(&spd[1], 0x51, id_only);
	read_spd(&spd[2], 0x52, id_only);
	read_spd(&spd[3], 0x53, id_only);
}
