#include <stdlib.h>

extern void termPutImage(unsigned char *, int, int);
extern void termInit();
extern void termFlush();
extern void termClose();
extern void termWait();

int main() {
  unsigned char *data = (unsigned char *)malloc(sizeof(char) * 3 * 128 * 128);
  for (int i = 0; i < 3 * 128 * 128; i += 3) {
    if (i < 128 * 3 * 24) {
      data[i] = (unsigned char)(i % 255);
      data[i + 1] = (unsigned char)0;
      data[i + 2] = (unsigned char)0;
    } else if (i > 128 * 24 * 3 && i < 128 * 48 * 3) {
      data[i] = (unsigned char)(i % 255);
      data[i + 1] = (unsigned char)(i % 255);
      data[i + 2] = (unsigned char)0;
    } else if (i > 128 * 48 * 3 && i < 128 * 72 * 3) {
      data[i] = (unsigned char)0;
      data[i + 1] = (unsigned char)(i % 255);
      data[i + 2] = (unsigned char)(i % 255);
    } else if (i > 128 * 72 * 3 && i < 128 * 96 * 3) {
      data[i] = (unsigned char)0;
      data[i + 1] = (unsigned char)0;
      data[i + 2] = (unsigned char)(i % 255);
    } else {
      data[i] = (unsigned char)(i % 255);
      data[i + 1] = (unsigned char)0;
      data[i + 2] = (unsigned char)(i % 255);
    }
  }
  termInit();
  termPutImage(data, 128, 128);
  termFlush();
  termWait();
  termClose();
  return 0;
}