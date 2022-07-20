import 'package:flutter/material.dart';

// AuthService is used to keep track of the auth state of the user
class AuthService extends ChangeNotifier {
  bool _authenticated = false;

  bool get authenticated => _authenticated;

  set authenticated(bool value) {
    _authenticated = value;
    notifyListeners();
  }
}
