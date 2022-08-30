// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'status.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Status _$StatusFromJson(Map<String, dynamic> json) => Status(
      eventName: json['eventName'] as String,
      status: json['status'] as String,
    );

Map<String, dynamic> _$StatusToJson(Status instance) => <String, dynamic>{
      'eventName': instance.eventName,
      'status': instance.status,
    };
