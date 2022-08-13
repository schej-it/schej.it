// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'friend_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

FriendRequest _$FriendRequestFromJson(Map<String, dynamic> json) =>
    FriendRequest(
      id: json['_id'] as String,
      from: json['from'] as String,
      fromUser: User.fromJson(json['fromUser'] as Map<String, dynamic>),
      to: json['to'] as String,
      toUser: User.fromJson(json['toUser'] as Map<String, dynamic>),
      createdAt: DateTime.parse(json['createdAt'] as String),
    );

Map<String, dynamic> _$FriendRequestToJson(FriendRequest instance) =>
    <String, dynamic>{
      '_id': instance.id,
      'from': instance.from,
      'fromUser': instance.fromUser,
      'to': instance.to,
      'toUser': instance.toUser,
      'createdAt': instance.createdAt.toIso8601String(),
    };
