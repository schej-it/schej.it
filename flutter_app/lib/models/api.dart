import 'dart:convert';
import 'dart:io';

import 'package:flutter/foundation.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:http/http.dart';
import 'package:requests/requests.dart';

// Api is used to keep track of all the data retrieved from the API as well as
// to make API requests
class ApiService extends ChangeNotifier {
  static final serverAddress = getServerAddress();

  static String getServerAddress() {
    if (kReleaseMode) {
      return 'https://schej.it/api';
    }

    if (Platform.isAndroid) {
      return 'http://10.0.2.2:3002';
    } else {
      return 'http://localhost:3002';
    }
  }

  Future<bool> isSignedIn() async {
    try {
      await get('/auth/status');
      return true;
    } catch (e) {
      return false;
    }
  }

  Future<bool> signIn(String accessToken, int expiresIn, String idToken,
      String refreshToken) async {
    final int timezoneOffset = DateTime.now().timeZoneOffset.inMinutes;
    try {
      await post(
        '/auth/sign-in-mobile',
        {
          'timezoneOffset': timezoneOffset,
          'accessToken': accessToken,
          'expiresIn': expiresIn,
          'idToken': idToken,
          'refreshToken': refreshToken,
        },
      );
      return true;
    } catch (e) {
      // TODO: show dialog that sign in failed
      print('Sign in failed! $e');
      return false;
    }
  }

  Future<void> signOut() async {
    await post('/auth/sign-out', {});
  }

  //
  // Methods to send http requests to the API
  //
  static Future<Map<String, dynamic>> get(String path) async {
    return await requestMethod(HttpMethod.GET, path);
  }

  static Future<Map<String, dynamic>> post(String path, dynamic body) async {
    return await requestMethod(HttpMethod.POST, path, body: body);
  }

  static Future<Map<String, dynamic>> patch(String path, dynamic body) async {
    return await requestMethod(HttpMethod.PATCH, path, body: body);
  }

  // Send the given request with the specified method
  static Future<Map<String, dynamic>> requestMethod(
      HttpMethod method, String path,
      {dynamic body = ''}) async {
    String url = '$serverAddress$path';
    Response r;

    switch (method) {
      case HttpMethod.GET:
        r = await Requests.get(url);
        break;
      case HttpMethod.POST:
        r = await Requests.post(
          url,
          json: body,
        );
        break;
      case HttpMethod.PATCH:
        r = await Requests.patch(
          url,
          json: body,
        );
        break;
      default:
        return {};
    }

    // Throw error if there was one
    if (r.hasError) {
      throw '[ERROR ${r.statusCode}] ${r.content()}';
    }

    // Get json object to return
    String response = r.content();
    Map<String, dynamic> json = jsonDecode(response);

    // Write sessionCookie to secure storage if it was set on this request
    final cookieJar = Requests.extractResponseCookies(r.headers);
    if (cookieJar.isNotEmpty) {
      if (cookieJar['session'] != null) {
        String sessionCookieString = cookieJar['session'].toString();
        sessionCookieString =
            sessionCookieString.substring(sessionCookieString.indexOf(':') + 1);
        const storage = FlutterSecureStorage();
        await storage.write(key: 'sessionCookie', value: sessionCookieString);
      }
    }

    return json;
  }
}
