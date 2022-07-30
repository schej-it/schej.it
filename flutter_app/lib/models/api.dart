import 'dart:convert';

import 'package:flutter/foundation.dart';
import 'package:http/http.dart';
import 'package:requests/requests.dart';

// Api is used to keep track of all the data retrieved from the API as well as
// to make API requests
class ApiService extends ChangeNotifier {
  static const serverAddress =
      kReleaseMode ? 'https://schej.it/api' : 'http://localhost:3002';

  Future<void> signIn(String code) async {
    final int timezoneOffset = DateTime.now().timeZoneOffset.inMinutes;
    final json = await post(
      '/auth/sign-in',
      {'code': code, 'timezoneOffset': timezoneOffset},
    );
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

    if (r.hasError) {
      throw '[ERROR ${r.statusCode}] ${r.content()}';
    }
    String response = r.content();
    Map<String, dynamic> json = jsonDecode(response);
    return json;
  }
}
