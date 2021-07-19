#include <webp/encode.h>

// A custom writer function for intercepting WebP data
int writeTo(uint8_t* data, size_t length, WebPPicture* picture);