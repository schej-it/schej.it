import 'package:flutter/material.dart';
import 'package:flutter_app/models/api.dart';
import 'package:google_sign_in/google_sign_in.dart';
import 'package:requests/requests.dart';

// AuthService is used to keep track of the auth state of the user
class AuthService extends ChangeNotifier {
  AuthService({
    required this.apiService,
  });

  final ApiService apiService;
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
    GoogleSignInAccount? user = await _googleSignIn.signInSilently();
    print('user: ${_googleSignIn.currentUser}');
    if (user != null && user.serverAuthCode != null) {
      apiService.signIn(user.serverAuthCode!);
    }
    notifyListeners();
  }

  Future<bool> signIn() async {
    GoogleSignInAccount? user = await _googleSignIn.signIn();
    // print('user: ${_googleSignIn.currentUser}');
    if (user != null && user.serverAuthCode != null) {
      apiService.signIn(user.serverAuthCode!);
    }
    notifyListeners();
    return user != null;
  }

  Future<void> signOut() async {
    await _googleSignIn.signOut();
    notifyListeners();
  }

  Future<void> _signInServer(String authCode) async {
    await Requests.post('');
  }
}
