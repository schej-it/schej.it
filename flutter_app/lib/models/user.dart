import 'package:json_annotation/json_annotation.dart';

part 'user.g.dart';

@JsonSerializable()
class User {
  @JsonKey(name: '_id')
  final String id;

  final String email;
  final String firstName;
  final String lastName;
  final String picture;
  final int visibility;

  const User({
    required this.id,
    required this.email,
    required this.firstName,
    required this.lastName,
    required this.picture,
    required this.visibility,
  });

  String get fullName => '$firstName $lastName';

  @override
  String toString() {
    return '{User id:"$id" email:"$email" firstName:"$firstName" lastName:"$lastName" picture:"$picture" visibility:"$visibility"}';
  }

  factory User.fromJson(Map<String, dynamic> json) => _$UserFromJson(json);
  Map<String, dynamic> toJson() => _$UserToJson(this);
}
