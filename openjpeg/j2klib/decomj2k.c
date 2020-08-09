#include <stdio.h>
#include <stdlib.h>
#include <stddef.h>
#include <stdbool.h>
#include <memory.h>
#include <assert.h>

#include "openjpeg.h"

#define J2K_CFMT 0
#define JP2_CFMT 1
#define JPT_CFMT 2
#define MJ2_CFMT 3
#define PXM_DFMT 0
#define PGX_DFMT 1
#define BMP_DFMT 2
#define YUV_DFMT 3

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
int int_ceildivpow2(int a, int b) {
  return (a + (1 << b) - 1) >> b;
}

/*
 * The following function was copy paste from j2k_to_image.c with part from convert.c
 */
bool J2KDecode(char *inputdata, int inputlength, char *raw){
  opj_dparameters_t parameters;  /* decompression parameters */
  opj_event_mgr_t event_mgr;    /* event manager */
  opj_image_t *image;
  opj_dinfo_t* dinfo;  /* handle to a decompressor */
  opj_cio_t *cio;
  unsigned char *src = (unsigned char*)inputdata;
  int file_length = inputlength;

  /* configure the event callbacks (not required) */
  memset(&event_mgr, 0, sizeof(opj_event_mgr_t));
//  event_mgr.error_handler = error_callback;
//  event_mgr.warning_handler = warning_callback;
//  event_mgr.info_handler = info_callback;

  /* set decoding parameters to default values */
  opj_set_default_decoder_parameters(&parameters);
 
   // default blindly copied
   parameters.cp_layer=0;
   parameters.cp_reduce=0;
//   parameters.decod_format=-1;
//   parameters.cod_format=-1;

     /* JPEG-2000 codestream */
     parameters.decod_format = J2K_CFMT;
     assert(parameters.decod_format == J2K_CFMT);
      /* get a decoder handle */
      dinfo = opj_create_decompress(CODEC_J2K);
  parameters.cod_format = PGX_DFMT;
  assert(parameters.cod_format == PGX_DFMT);


      /* catch events using our callbacks and give a local context */
      opj_set_event_mgr((opj_common_ptr)dinfo, &event_mgr, NULL);

      /* setup the decoder decoding parameters using user parameters */
      opj_setup_decoder(dinfo, &parameters);

      /* open a byte stream */
      cio = opj_cio_open((opj_common_ptr)dinfo, src, file_length);

      /* decode the stream and fill the image structure */
      image = opj_decode(dinfo, cio);
      if(!image) {
        opj_destroy_decompress(dinfo);
        opj_cio_close(cio);
        return 1;
      }
      
      /* close the byte stream */
      opj_cio_close(cio);

   // Copy buffer
   for (int compno = 0; compno < image->numcomps; compno++)
   {
      opj_image_comp_t *comp = &image->comps[compno];

      int w = image->comps[compno].w;
      int wr = int_ceildivpow2(image->comps[compno].w, image->comps[compno].factor);

      //int h = image.comps[compno].h;
      int hr = int_ceildivpow2(image->comps[compno].h, image->comps[compno].factor);

      if (comp->prec <= 8)
      {
         uint8_t *data8 = (uint8_t*)raw + compno;
         for (int i = 0; i < wr * hr; i++)
         {
            int v = image->comps[compno].data[i / wr * w + i % wr];
            *data8 = (uint8_t)v;
            data8 += image->numcomps;
         }
      }
      else if (comp->prec <= 16)
      {
         uint16_t *data16 = (uint16_t*)raw + compno;
         for (int i = 0; i < wr * hr; i++)
         {
            int v = image->comps[compno].data[i / wr * w + i % wr];
            *data16 = (int16_t)v;
            data16 += image->numcomps;
         }
      }
      else
      {
         uint32_t *data32 = (uint32_t*)raw + compno;
         for (int i = 0; i < wr * hr; i++)
         {
            int v = image->comps[compno].data[i / wr * w + i % wr];
            *data32 = (uint32_t)v;
            data32 += image->numcomps;
         }
      }
      //free(image.comps[compno].data);
   }


  /* free remaining structures */
  if(dinfo) {
    opj_destroy_decompress(dinfo);
  }

  /* free image data structure */
  opj_image_destroy(image);

  return true;
}
/*
int FileSize(FILE *fp){
  int size;

  fseek(fp, 0L, SEEK_END);
  size = ftell(fp);
  fseek(fp, 0L, SEEK_SET);
return(size);
}

int main(){
  unsigned char *jpeg_data;
  unsigned char *output_data;
  int jpeg_size, output_size;
  FILE *fp;

  if((fp=fopen("test.j2k", "rb"))==NULL) {
    puts("ERROR, can't open test.j2k");
    return -1;
    }
  jpeg_size = FileSize(fp);
  jpeg_data = malloc(jpeg_size);
  fread(jpeg_data, 1, jpeg_size, fp);
  fclose(fp);
  printf("INFO, J2KDecode: %d\r\n", jpeg_size);
  
  output_size = 1576*1134*3;
  output_data = malloc(output_size);
  if(J2Kdecode(jpeg_data, jpeg_size, output_data)==true){
    puts("INFO, decode was succesfull");
    if((fp=fopen("test.raw", "wb"))==NULL) {
      puts("ERROR, can't open test.raw");
      return -1;
      }
    fwrite(output_data, 1, output_size, fp);
    fclose(fp);
  } else {
    puts("ERROR, decode failed!");  
  }
  free(jpeg_data);
  free(output_data);
  return 0;
}
*/
