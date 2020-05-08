#!/bin/sh

i=0
while true
do
   echo "$i"
   echo "Line $i " >>${MOUNT_POINT}/testfile.txt  
   (( i+=1 ))
   sleep 10
done
