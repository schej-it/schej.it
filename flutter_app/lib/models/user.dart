class User {
  final String id;
  final String email;
  final String firstName;
  final String lastName;
  final String picture;

  const User({
    required this.id,
    required this.email,
    required this.firstName,
    required this.lastName,
    required this.picture,
  });

  String get fullName => '$firstName $lastName';

  @override
  String toString() {
    return '{User id:"$id" email:"$email" firstName:"$firstName" lastName:"$lastName" picture:"$picture"}';
  }
}
