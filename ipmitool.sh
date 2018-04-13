#!/usr/bin/env bash
ipmitool -H 118.184.72.242 -I lanplus -U root -P u4ymHar7 -v sdr list|while read line
do
if [[ $line =~ 'Sensor ID' ]];then
sensor_id=$(echo $line|awk -F':' '{print $2}')
fi
if [[ $line =~ 'Entity ID' ]];then
entity_id=$(echo $line|awk -F':' '{print $2}')
fi
if [[ $line =~ 'Sensor Type' ]];then
sensor_type=$(echo $line|awk -F':' '{print $2}')
fi
if [[ $line =~ 'Sensor Reading' ]];then
sensor_reading=$(echo $line|awk -F':' '{print $2}')
fi
if [[ $line = '' ]];then
    echo $entity_id\|$sensor_id\|$sensor_type\|$sensor_reading
fi
#if [[ -n $entity_id ]] && [[ -n $sensor_id ]] && [[ -n $sensor_type ]] && [[ -n $sensor_reading ]];then
#    echo $entity_id\|$sensor_id\|$sensor_type\|$sensor_reading
#fi
done
