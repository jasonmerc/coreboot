## SPDX-License-Identifier: GPL-2.0-only

ifeq ($(CONFIG_EC_DELL_MEC5055),y)

bootblock-y += mec5055.c early_init.c
romstage-y += mec5055.c early_init.c

endif
