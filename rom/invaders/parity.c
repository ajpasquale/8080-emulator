#include <stdio.h>

sdad

int parity(int x, int size)
{
	int i;
	int p = 0;
	x = (x & ((1<<size)-1));
	for (i=0; i<size; i++)
	{
		if (x & 0x1) p++;
		x = x >> 1;
	}
	return (0 == (p & 0x1));
}


int main(void){
  for(int i=0; i< 100;i++)
{
  printf("{%d, %d},\n" ,i, parity(i, 8));
}
  return 0;
}
