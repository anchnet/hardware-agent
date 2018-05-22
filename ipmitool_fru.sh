#!/usr/bin/env bash
HOST=118.184.72.242
USER=root
PWD=u4ymHar7
ipmitool -H $HOST -I lanplus -U $USER -P $PWD fru list|while read line
do
if [[ $line =~ 'Board Mfg' ]];then
brand=$(echo $line|awk -F':' '{print $2}')
fi
if [[ $line =~ 'Board Product' ]];then
model=$(echo $line|awk -F':' '{print $2}')
echo $brand\|$model
break
fi
#if [[ -n $entity_id ]] && [[ -n $sensor_id ]] && [[ -n $sensor_type ]] && [[ -n $sensor_reading ]];then
#    echo $entity_id\|$sensor_id\|$sensor_type\|$sensor_reading
#fi
done
