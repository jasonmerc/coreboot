/* SPDX-License-Identifier: GPL-2.0-only */

#include <stdint.h>
#include <northbridge/intel/haswell/haswell.h>
#include <northbridge/intel/haswell/raminit.h>
#include <southbridge/intel/lynxpoint/pch.h>

void mainboard_config_rcba(void)
{
	RCBA16(D31IR) = DIR_ROUTE(PIRQA, PIRQD, PIRQC, PIRQA);
	RCBA16(D29IR) = DIR_ROUTE(PIRQF, PIRQD, PIRQA, PIRQC);
	RCBA16(D28IR) = DIR_ROUTE(PIRQA, PIRQA, PIRQA, PIRQA);
	RCBA16(D27IR) = DIR_ROUTE(PIRQG, PIRQB, PIRQC, PIRQD);
	RCBA16(D26IR) = DIR_ROUTE(PIRQA, PIRQF, PIRQC, PIRQD);
	RCBA16(D25IR) = DIR_ROUTE(PIRQE, PIRQF, PIRQG, PIRQH);
	RCBA16(D22IR) = DIR_ROUTE(PIRQA, PIRQD, PIRQC, PIRQB);
	RCBA16(D20IR) = DIR_ROUTE(PIRQA, PIRQB, PIRQC, PIRQD);
}

/* FIXME: called after romstage_common, remove it if not used */
void mb_late_romstage_setup(void)
{
}

void mb_get_spd_map(uint8_t spd_map[4])
{
	/* FIXME: check this */
	/* note that SPD addresses are left-shifted by 1. */
	spd_map[0] = 0xa0;
	spd_map[1] = 0xa2;
	spd_map[2] = 0xa4;
	spd_map[3] = 0xa6;
}

void mainboard_fill_pei_data(struct pei_data *pei_data)
{
	struct usb2_port_setting usb2_ports[MAX_USB2_PORTS] = {
		/* Length, Enable, OCn#, Location */
		{ 0x0040, 1, 0, USB_PORT_BACK_PANEL },
		{ 0x0040, 1, 1, USB_PORT_BACK_PANEL },
		{ 0x0080, 1, USB_OC_PIN_SKIP, USB_PORT_DOCK },
		{ 0x0040, 1, USB_OC_PIN_SKIP, USB_PORT_BACK_PANEL },
		{ 0x0040, 1, USB_OC_PIN_SKIP, USB_PORT_BACK_PANEL },
		{ 0x0110, 1, 2, USB_PORT_MINI_PCIE },
		{ 0x0040, 1, USB_OC_PIN_SKIP, USB_PORT_DOCK },
		{ 0x0040, 1, USB_OC_PIN_SKIP, USB_PORT_BACK_PANEL },
	};
	struct usb3_port_setting usb3_ports[MAX_USB3_PORTS] = {
		{ 1, 1 },
		{ 1, 0 },
		{ 1, 2 },
		{ 1, USB_OC_PIN_SKIP },
		{ 1, USB_OC_PIN_SKIP },
		{ 1, USB_OC_PIN_SKIP },
	};
	memcpy(pei_data->usb2_ports, usb2_ports, sizeof(usb2_ports));
	memcpy(pei_data->usb3_ports, usb3_ports, sizeof(usb3_ports));
}