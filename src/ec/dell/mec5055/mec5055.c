/* SPDX-License-Identifier: GPL-2.0-only */

#include <arch/io.h>
#include "mec5055.h"

static const u16 EC_CTRL = 0x910;
static const u16 EC_DATA = 0x911;

static void wait_ec(void)
{
	u8 v;
	do {
		outb(0, EC_CTRL);
		v = inb(EC_DATA);
	} while (v != 0);
}

static int read_ec_regs(u8 start, u8 *sdata, u8 count)
{
	if (start >= 0x20 || start + count > 0x20)
		return -1;

	if (count <= 0)
		return 0;

	while (count--) {
		outb(start + 0x10, EC_CTRL);
		*sdata = inb(EC_DATA);
		sdata++;
		start++;
	}

	return 0;
}

static int write_ec_regs(u8 start, const u8 *sdata, u8 count)
{
	if (start >= 0x20 || start + count > 0x20)
		return -1;

	if (count <= 0)
		return 0;

	while (count--) {
		outb(start + 0x10, EC_CTRL);
		outb(*sdata, EC_DATA);
		sdata++;
		start++;
	}

	return 0;
}

/* Steps to communicate with the EC:
 * 1. send command:
 *   - write arguments to EC[2:]
 *   - send 0 to CTRL then cmd to DATA
 * 2. wait EC
 * 3. receive result: result may starts at EC[0] or EC[1]
 */

/* send command to EC and wait, command at cmd[0], args in cmd[2:] */
static int ec_send_command(const u8 *cmd, int argc)
{
	int result = 0;

	if (argc > 0) {
		if (argc > 0x1e)
			result = -1;
		else
			result = write_ec_regs(2, cmd + 2, argc);
	}

	outb(0, EC_CTRL);
	outb(cmd[0], EC_DATA);
	wait_ec();
	return result;
}

/* receive result from EC[0:] */
int mec5055_ec_command_0(u8 *buf, u8 cmd, int argc)
{
	buf[0] = cmd;
	int res = ec_send_command(buf, argc);
	if (res < 0)
		return res;

	return read_ec_regs(0, buf, 0x10);
}

/* receive result from EC[1:] */
int mec5055_ec_command_1(u8 cmd, const u8 *cmd_buf, int argc, u8 *res_buf, int res_size)
{
	if (argc >= 0x1e)
		argc = 0x1e;

	if (res_size >= 0x1f)
		res_size = 0x1f;

	int t = write_ec_regs(2, cmd_buf, argc);
	if (t < 0)
		return t;

	outb(0, EC_CTRL);
	outb(cmd, EC_DATA);
	wait_ec();

	return read_ec_regs(1, res_buf, res_size);
}
