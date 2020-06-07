/* SPDX-License-Identifier: GPL-2.0-only */

#ifndef _EC_DELL_MEC5055_H_
#define _EC_DELL_MEC5055_H_

#include <stddef.h>
#include <stdint.h>

int mec5055_ec_command_0(u8 *buf, u8 cmd, int argc);
int mec5055_ec_command_1(u8 cmd, const u8 *cmd_buf, int argc, u8 *res_buf, int res_size);
void mec5055_early_init(void);

#endif
