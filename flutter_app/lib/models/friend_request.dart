import 'package:flutter_app/models/user.dart';
import 'package:json_annotation/json_annotation.dart';

part 'friend_request.g.dart';

@JsonSerializable()
class FriendRequest {
  @JsonKey(name: '_id')
  final String id;
  final String from;
  final User fromUser;
  final String to;
  final User toUser;
  final DateTime createdAt;

  const FriendRequest({
    required this.id,
    required this.from,
    required this.fromUser,
    required this.to,
    required this.toUser,
    required this.createdAt,
  });

  @override
  String toString() {
    return '{FriendRequest id:"$id" from:"$from" to:"$to" createdAt:"$createdAt"}';
  }

  factory FriendRequest.fromJson(Map<String, dynamic> json) =>
      _$FriendRequestFromJson(json);
  Map<String, dynamic> toJson() => _$FriendRequestToJson(this);
}
