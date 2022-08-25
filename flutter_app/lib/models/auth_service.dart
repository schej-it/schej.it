import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_app/models/api.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:oauth2_client/access_token_response.dart';
import 'package:oauth2_client/google_oauth2_client.dart';

// AuthService is used to keep track of the auth state of the user
class AuthService extends ChangeNotifier {
  final ApiService apiService;

  AuthService({
    required this.apiService,
  });

  final GoogleOAuth2Client _client = Platform.isAndroid
      ? GoogleOAuth2Client(
          redirectUri: 'it.schej.app:/oauth2redirect',
          customUriScheme: 'it.schej.app',
        )
      : GoogleOAuth2Client(
          redirectUri:
              'com.googleusercontent.apps.523323684219-vntbcabt43u6ago35a8s9mkjlhrserdg:/oauth2redirect',
          customUriScheme:
              'com.googleusercontent.apps.523323684219-vntbcabt43u6ago35a8s9mkjlhrserdg',
        );

  bool _authenticated = false;

  bool get authenticated => _authenticated;

  Future<void> signInSilently() async {
    bool signedIn = await apiService.isSignedIn();
    _authenticated = signedIn;
    notifyListeners();
  }

  Future<bool> signIn() async {
    // Show sign in dialog
    AccessTokenResponse res = await _client.getTokenWithAuthCodeFlow(
      clientId: Platform.isAndroid
          ? '523323684219-rltj890jhj1ncb0dk1532vq4k8t3a6bp.apps.googleusercontent.com'
          : '523323684219-vntbcabt43u6ago35a8s9mkjlhrserdg.apps.googleusercontent.com',
      scopes: [
        'openid',
        'profile',
        'email',
        'https://www.googleapis.com/auth/calendar.calendarlist.readonly',
        'https://www.googleapis.com/auth/calendar.events.readonly',
      ],
    );

    // Check if user cancelled sign in or access token invalid
    if (!res.isValid()) {
      return false;
    }

    // Sign user in with the api
    String idToken = res.getRespField('id_token');
    bool success = await apiService.signIn(
        res.accessToken!, res.expiresIn!, idToken, res.refreshToken!);

    return success;
  }

  Future<void> signOut() async {
    await apiService.signOut();

    const storage = FlutterSecureStorage();
    await storage.deleteAll();

    _authenticated = false;
    notifyListeners();
  }
}
