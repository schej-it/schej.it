import 'package:flutter/material.dart';
import 'package:google_sign_in/google_sign_in.dart';

// AuthService is used to keep track of the auth state of the user
class AuthService extends ChangeNotifier {
  final GoogleSignIn _googleSignIn = GoogleSignIn(
    scopes: [
      'openid',
      'profile',
      'email',
      'https://www.googleapis.com/auth/calendar.calendarlist.readonly',
      'https://www.googleapis.com/auth/calendar.events.readonly',
    ],
    serverClientId:
        '523323684219-jfakov2bgsleeb6den4ktpohq4lcnae2.apps.googleusercontent.com',
  );

  GoogleSignInAccount? get currentUser => _googleSignIn.currentUser;

  Future<void> signInSilently() async {
    await _googleSignIn.signInSilently();
    // print('user: ${_googleSignIn.currentUser}');
    notifyListeners();
  }

  Future<bool> signIn() async {
    GoogleSignInAccount? user = await _googleSignIn.signIn();
    // print('user: ${_googleSignIn.currentUser}');
    notifyListeners();
    return user != null;
  }

  Future<void> signOut() async {
    await _googleSignIn.signOut();
    notifyListeners();
  }
}
