import 'package:json_annotation/json_annotation.dart';

part 'status.g.dart';

@JsonSerializable()
class Status {
  final String eventName;
  final String status;

  const Status({
    required this.eventName,
    required this.status,
  });

  @override
  String toString() {
    return '{Status eventName:"$eventName" status:"$status"}';
  }

  factory Status.fromJson(Map<String, dynamic> json) => _$StatusFromJson(json);
  Map<String, dynamic> toJson() => _$StatusToJson(this);
}
