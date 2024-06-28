#include <lyn_smi.h>
#include <stdio.h>

void main() {
    lynDeviceProperties_t prop;
    lynGetDeviceProperties(0, &prop);
    printf("%s\n", prop.boardBrand);
}

// gcc test_utils/test_smi.c -lLYNSMICLIENTCOMM -o test_utils/test_smi && ./test_utils/test_smi