/* SPDX-License-Identifier: GPL-2.0-only */

#include <device/pci_ops.h>
#include <ec/dell/mec5055/mec5055.h>
#include <southbridge/intel/lynxpoint/pch.h>

void mainboard_config_superio(void)
{
	/* port 0x910 and 0x911 is needed to initialize the EC */
	pci_write_config32(PCH_LPC_DEV, LPC_GEN2_DEC, 0x007c0901);
	mec5055_early_init();
	pci_write_config32(PCH_LPC_DEV, LPC_GEN2_DEC, 0x005c0921);
}
