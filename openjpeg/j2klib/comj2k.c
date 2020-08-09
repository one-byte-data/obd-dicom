
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <assert.h>
#include <string.h>
#include "openjpeg.h"

#define J2K_CFMT 0

typedef  signed char         int8_t;
typedef  signed short        int16_t;
typedef  signed int          int32_t;
typedef  unsigned char       uint8_t;
typedef  unsigned short      uint16_t;
typedef  unsigned int        uint32_t;

/*
 * Divide an integer by a power of 2 and round upwards.
 *
 * a divided by 2^b
 */
extern int int_ceildivpow2(int a, int b);

void rawtoimg_fill8(int8_t *inputbuffer, int w, int h, int numcomps, opj_image_t *image) {
  int8_t *p = inputbuffer;
  for (int i = 0; i < w * h; i++)
    {
    for(int compno = 0; compno < numcomps; compno++)
      {
      /* compno : 0 = GREY, (0, 1, 2) = (R, G, B) */
      image->comps[compno].data[i] = *p;
      ++p;
      }
    }
}

void rawtoimg_fillu8(uint8_t *inputbuffer, int w, int h, int numcomps, opj_image_t *image) {
  uint8_t *p = inputbuffer;
  for (int i = 0; i < w * h; i++)
    {
    for(int compno = 0; compno < numcomps; compno++)
      {
      /* compno : 0 = GREY, (0, 1, 2) = (R, G, B) */
      image->comps[compno].data[i] = *p;
      ++p;
      }
    }
}

void rawtoimg_fill16(int16_t *inputbuffer, int w, int h, int numcomps, opj_image_t *image)
{
  int16_t *p = inputbuffer;
  for (int i = 0; i < w * h; i++)
    {
    for(int compno = 0; compno < numcomps; compno++)
      {
      /* compno : 0 = GREY, (0, 1, 2) = (R, G, B) */
      image->comps[compno].data[i] = *p;
      ++p;
      }
    }
}

void rawtoimg_fillu16(uint16_t *inputbuffer, int w, int h, int numcomps, opj_image_t *image)
{
  uint16_t *p = inputbuffer;
  for (int i = 0; i < w * h; i++)
    {
    for(int compno = 0; compno < numcomps; compno++)
      {
      /* compno : 0 = GREY, (0, 1, 2) = (R, G, B) */
      image->comps[compno].data[i] = *p;
      ++p;
      }
    }
}

void rawtoimg_fill32(int32_t *inputbuffer, int w, int h, int numcomps, opj_image_t *image)
{
  int32_t *p = inputbuffer;
  for (int i = 0; i < w * h; i++)
    {
    for(int compno = 0; compno < numcomps; compno++)
      {
      /* compno : 0 = GREY, (0, 1, 2) = (R, G, B) */
      image->comps[compno].data[i] = *p;
      ++p;
      }
    }
}

void rawtoimg_fillu32(uint32_t *inputbuffer, int w, int h, int numcomps, opj_image_t *image)
{
  uint32_t *p = inputbuffer;
  for (int i = 0; i < w * h; i++)
    {
    for(int compno = 0; compno < numcomps; compno++)
      {
      /* compno : 0 = GREY, (0, 1, 2) = (R, G, B) */
      image->comps[compno].data[i] = *p;
      ++p;
      }
    }
}

opj_image_t* rawtoimage(char *inputbuffer, opj_cparameters_t *parameters, int image_width, int image_height,
 int sample_pixel, int bitsallocated, int sign)
{
  int w, h;
  int numcomps;
  OPJ_COLOR_SPACE color_space;
  opj_image_cmptparm_t cmptparm[3]; /* maximum of 3 components */
  opj_image_t * image = NULL;

  assert( sample_pixel == 1 || sample_pixel == 3 );
  if( sample_pixel == 1 )
    {
    numcomps = 1;
    color_space = CLRSPC_GRAY;
    }
  else // sample_pixel == 3
    {
    numcomps = 3;
    color_space = CLRSPC_SRGB;
    }
  int subsampling_dx = parameters->subsampling_dx;
  int subsampling_dy = parameters->subsampling_dy;

  // FIXME
  w = image_width;
  h = image_height;

  /* initialize image components */
  memset(&cmptparm[0], 0, 3 * sizeof(opj_image_cmptparm_t));
  //assert( bitsallocated == 8 );
  for(int i = 0; i < numcomps; i++) {
    cmptparm[i].prec = bitsallocated;
    cmptparm[i].bpp = bitsallocated;
    cmptparm[i].sgnd = sign;
    cmptparm[i].dx = subsampling_dx;
    cmptparm[i].dy = subsampling_dy;
    cmptparm[i].w = w;
    cmptparm[i].h = h;
  }

  /* create the image */
  image = opj_image_create(numcomps, &cmptparm[0], color_space);
  if(!image) {
    return NULL;
  }
  /* set image offset and reference grid */
  image->x0 = parameters->image_offset_x0;
  image->y0 = parameters->image_offset_y0;
  image->x1 = parameters->image_offset_x0 + (w - 1) * subsampling_dx + 1;
  image->y1 = parameters->image_offset_y0 + (h - 1) * subsampling_dy + 1;

  /* set image data */

  //assert( fragment_size == numcomps*w*h*(bitsallocated/8) );
  if (bitsallocated <= 8)
    {
    if( sign )
      {
      rawtoimg_fill8((int8_t*)inputbuffer,w,h,numcomps,image);
      }
    else
      {
      rawtoimg_fillu8((uint8_t*)inputbuffer,w,h,numcomps,image);
      }
    }
  else if (bitsallocated <= 16)
    {
    if( sign )
      {
      rawtoimg_fill16((int16_t*)inputbuffer,w,h,numcomps,image);
      }
    else
      {
      rawtoimg_fillu16((uint16_t*)inputbuffer,w,h,numcomps,image);
      }
    }
  else if (bitsallocated <= 32)
    {
    if( sign )
      {
      rawtoimg_fill32((int32_t*)inputbuffer,w,h,numcomps,image);
      }
    else
      {
      rawtoimg_fillu32((uint32_t*)inputbuffer,w,h,numcomps,image);
      }
    }
  else
    {
    abort();
    }

  return image;
}

/*
 * The following function was copy paste from image_to_j2k.c with part from convert.c
 */

bool J2KEncode(char *raw_data, int image_width, int image_height, int sample_pixel, int bitsallocated, char **jpeg_data, int *encodedlength, int ratio)
{
//// input_buffer is ONE image
//// fragment_size is the size of this image (fragment)
  bool bSuccess;
  opj_cparameters_t parameters;  /* compression parameters */
  opj_event_mgr_t event_mgr;    /* event manager */
  opj_image_t *image = NULL;
  //quality = 100;

  /*
  configure the event callbacks (not required)
  setting of each callback is optionnal
  */
  memset(&event_mgr, 0, sizeof(opj_event_mgr_t));
//  event_mgr.error_handler = error_callback;
//  event_mgr.warning_handler = warning_callback;
//  event_mgr.info_handler = info_callback;

  /* set encoding parameters to default values */
  memset(&parameters, 0, sizeof(parameters));
  opj_set_default_encoder_parameters(&parameters);

   parameters.tcp_rates[0] = ratio;
  parameters.tcp_numlayers = 1;
  parameters.cp_disto_alloc = 1;

  if(parameters.cp_comment == NULL) {
    const char comment[] = "Created by OpenJPEG version 1.5";
    parameters.cp_comment = (char*)malloc(strlen(comment) + 1);
    strcpy(parameters.cp_comment, comment);
    /* no need to delete parameters.cp_comment on exit */
    //delete_comment = false;
  }

  /* decode the source image */
  /* ----------------------- */

  image = rawtoimage(raw_data, &parameters, image_width, image_height, sample_pixel, bitsallocated, 0);
  if (!image) {
    return false;
  }

    /* encode the destination image */
  /* ---------------------------- */
   parameters.cod_format = J2K_CFMT; /* J2K format output */
    int codestream_length;
    opj_cio_t *cio = NULL;

    /* get a J2K compressor handle */
    opj_cinfo_t* cinfo = opj_create_compress(CODEC_J2K);

    /* catch events using our callbacks and give a local context */
//    opj_set_event_mgr((opj_common_ptr)cinfo, &event_mgr, stderr);

    /* setup the encoder parameters using the current image and using user parameters */
    opj_setup_encoder(cinfo, &parameters, image);

    /* open a byte stream for writing */
    /* allocate memory for all tiles */
    cio = opj_cio_open((opj_common_ptr)cinfo, NULL, 0);

    /* encode the image */
    bSuccess = opj_encode(cinfo, cio, image, parameters.index);
    if (!bSuccess) {
      opj_cio_close(cio);
      return false;
    }
    codestream_length = cio_tell(cio);

    *jpeg_data = (char *)malloc(codestream_length);
    *encodedlength=codestream_length;
    printf("Encoded Length %d\r\n", codestream_length);
    memcpy(*jpeg_data, (char *)(cio->buffer), codestream_length);
//    fwrite((char*)(cio->buffer), codestream_length,1, fp);

    /* close and free the byte stream */
    opj_cio_close(cio);

    /* free remaining compression structures */
    opj_destroy_compress(cinfo);

      /* free user parameters structure */
  //if(delete_comment) {
    if(parameters.cp_comment) free(parameters.cp_comment);
  //}
  if(parameters.cp_matrice) free(parameters.cp_matrice);

  /* free image data */
  opj_image_destroy(image);
  return true;
}

/*
int main() {
char *jpeg_data;
unsigned char *img;
int jpeg_size, size;
FILE *fp;

puts("INFO, Starting test");
if((fp=fopen("test.raw", "rb"))==NULL) {
  puts("ERROR, can't open test.raw");
	return -1;
  }
size = 1576*1134*3; 
img = (unsigned char *) malloc(size);
size=fread(img, 1, size, fp);
fclose(fp);

if(J2Kencode((char *)img, 1576, 1134, 3, 8, &jpeg_data, &jpeg_size, 0)==true){
  if((fp=fopen("out.j2k", "wb"))==NULL) {
    puts("ERROR, can't write out.j2k");
    return -1;
    }
  fwrite(jpeg_data, 1, jpeg_size, fp);
  fclose(fp);
  }
else {
  puts("ERROR, encoding JPEG data");
}
free(jpeg_data);

free(img);
puts("INFO, Finish test");
return 0;
}
*/
