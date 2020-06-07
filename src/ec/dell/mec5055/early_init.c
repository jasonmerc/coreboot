/* SPDX-License-Identifier: GPL-2.0-only */

#include "mec5055.h"

void mec5055_early_init(void)
{
	u8 buf[32], c;

	mec5055_ec_command_0(buf, 0xc2, 0);
	mec5055_ec_command_1(0xab, NULL, 0, &c, 1);
}
