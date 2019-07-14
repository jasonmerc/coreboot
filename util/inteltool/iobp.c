/* inteltool - dump all registers on an Intel CPU + chipset based system */
/* SPDX-License-Identifier: GPL-2.0-only */

#include <stdio.h>
#include <stdlib.h>
#include "inteltool.h"

#define RCBA8(x) (*((volatile u8 *)(rcba + x)))
#define RCBA16(x) (*((volatile u16 *)(rcba + x)))
#define RCBA32(x) (*((volatile u32 *)(rcba + x)))
#define RCBA64(x) (*((volatile u64 *)(rcba + x)))
/* IO Buffer Programming */
#define IOBPIRI 0x2330
#define IOBPD 0x2334
#define IOBPS 0x2338
#define  IOBPS_READY 0x0001
#define  IOBPS_TX_MASK 0x0006
#define  IOBPS_MASK 0xff00
#define  IOBPS_READ 0x0600
#define  IOBPS_WRITE 0x0700
#define IOBPU 0x233a
#define  IOBPU_MAGIC 0xf000

#define IOBP_RETRY 1000
static inline int iobp_poll(volatile uint8_t *rcba)
{
	for (int try = IOBP_RETRY; try > 0; try --) {
		u16 status = RCBA16(IOBPS);
		if ((status & IOBPS_READY) == 0)
			return 1;
		// udelay(10);
	}

	printf("IOBP: timeout waiting for transaction to complete\n");
	return 0;
}

static u32 pch_iobp_read(volatile uint8_t *rcba, u32 address)
{
	u16 status;

	if (!iobp_poll(rcba))
		return 0;

	/* Set the address */
	RCBA32(IOBPIRI) = address;

	/* READ OPCODE */
	status = RCBA16(IOBPS);
	status &= ~IOBPS_MASK;
	status |= IOBPS_READ;
	RCBA16(IOBPS) = status;

	/* Undocumented magic */
	RCBA16(IOBPU) = IOBPU_MAGIC;

	/* Set ready bit */
	status = RCBA16(IOBPS);
	status |= IOBPS_READY;
	RCBA16(IOBPS) = status;

	if (!iobp_poll(rcba))
		return 0;

	/* Check for successful transaction */
	status = RCBA16(IOBPS);
	if (status & IOBPS_TX_MASK) {
		printf("IOBP: read 0x%08x failed\n", address);
		return 0;
	}

	/* Read IOBP data */
	return RCBA32(IOBPD);
}

typedef struct {
	u32 addr;
	char *name;
} iobp_register_t;

static const iobp_register_t lynxpoint_iobp_registers[] = {
	/* SATA IOBP */
	{0xea000151, "SATA_IOBP_SP0G3IR"},
	{0xea000051, "SATA_IOBP_SP1G3IR"},
	{0xea002550, "SATA_IOBP_SP0DTLE_DATA"},
	{0xea002554, "SATA_IOBP_SP0DTLE_EDGE"},
	{0xea002750, "SATA_IOBP_SP1DTLE_DATA"},
	{0xea002754, "SATA_IOBP_SP1DTLE_EDGE"},
	/* USB port IOBP */
	{0xe5004100, "USBP01"},
	{0xe5004200, "USBP02"},
	{0xe5004300, "USBP03"},
	{0xe5004400, "USBP04"},
	{0xe5004500, "USBP05"},
	{0xe5004600, "USBP06"},
	{0xe5004700, "USBP07"},
	{0xe5004800, "USBP08"},
	/* LynxPoint-LP doesn't have port 9-14, the IOBP values are zero */
	{0xe5004900, "USBP09"},
	{0xe5004a00, "USBP10"},
	{0xe5004b00, "USBP11"},
	{0xe5004c00, "USBP12"},
	{0xe5004d00, "USBP13"},
	{0xe5004e00, "USBP14"},
};

void print_iobp(struct pci_dev *sb, volatile uint8_t *rcba)
{
	const iobp_register_t *iobp_registers = NULL;
	size_t iobp_size = 0;

	switch (sb->device_id) {
	case PCI_DEVICE_ID_INTEL_LYNXPOINT_LP_FULL:
	case PCI_DEVICE_ID_INTEL_LYNXPOINT_LP_PREM:
	case PCI_DEVICE_ID_INTEL_LYNXPOINT_LP_BASE:
	case PCI_DEVICE_ID_INTEL_C8_MOBILE:
	case PCI_DEVICE_ID_INTEL_C8_DESKTOP:
	case PCI_DEVICE_ID_INTEL_Z87:
	case PCI_DEVICE_ID_INTEL_Z85:
	case PCI_DEVICE_ID_INTEL_HM86:
	case PCI_DEVICE_ID_INTEL_H87:
	case PCI_DEVICE_ID_INTEL_HM87:
	case PCI_DEVICE_ID_INTEL_Q85:
	case PCI_DEVICE_ID_INTEL_Q87:
	case PCI_DEVICE_ID_INTEL_QM87:
	case PCI_DEVICE_ID_INTEL_B85:
	case PCI_DEVICE_ID_INTEL_C222:
	case PCI_DEVICE_ID_INTEL_C224:
	case PCI_DEVICE_ID_INTEL_C226:
	case PCI_DEVICE_ID_INTEL_H81:
		iobp_registers = lynxpoint_iobp_registers;
		iobp_size = ARRAY_SIZE(lynxpoint_iobp_registers);
		break;
	default:
		break;
	}

	if (iobp_size == 0)
		return;

	printf("\n============= IOBP ==============\n\n");

	for (size_t i = 0; i < iobp_size; i++) {
		u32 address = iobp_registers[i].addr;
		const char *name = iobp_registers[i].name;
		u32 v = pch_iobp_read(rcba, address);
		printf("0x%08x: 0x%08x (%s)\n", address, v, name);
	}
}
