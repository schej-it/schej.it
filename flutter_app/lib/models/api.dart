import 'dart:convert';
import 'dart:io';

import 'package:flutter/foundation.dart';
import 'package:flutter_app/models/calendar_event.dart';
import 'package:flutter_app/models/user.dart';
import 'package:flutter_app/utils.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:http/http.dart';
import 'package:requests/requests.dart';

// Api is used to keep track of all the data retrieved from the API as well as
// to make API requests
class ApiService extends ChangeNotifier {
  // [serverAddress] is the correct api address to use
  static String get serverAddress {
    if (kReleaseMode) {
      return 'https://schej.it/api';
    }

    if (Platform.isAndroid) {
      return 'http://10.0.2.2:3002';
    } else {
      return 'http://localhost:3002';
    }
  }

  ///////////////////////////////////////////
  // Current user
  ///////////////////////////////////////////

  User? _authUser;
  User? get authUser => _authUser;

  // Gets the user's profile and sets [_authUser] to it
  Future<void> refreshAuthUserProfile() async {
    final userMap = await get('/user/profile');
    _authUser = User.fromJson(userMap);
  }

  final CalendarEvents _authUserSchedule = CalendarEvents(
    events: [
      CalendarEvent(
        title: 'Event',
        startDate: getDateWithTime(DateTime.now(), 9.5),
        endDate: getDateWithTime(DateTime.now(), 12),
      ),
      CalendarEvent(
        title: 'Introduction to Failure Analysis',
        startDate: getDateWithTime(DateTime.now(), 13),
        endDate: getDateWithTime(DateTime.now(), 14.5),
      ),
      CalendarEvent(
        title: 'War',
        startDate:
            getDateWithTime(DateTime.now().add(const Duration(days: 1)), 15),
        endDate:
            getDateWithTime(DateTime.now().add(const Duration(days: 1)), 20),
      ),
    ],
  );
  CalendarEvents get authUserSchedule => _authUserSchedule;

  ///////////////////////////////////////////
  // Friends
  ///////////////////////////////////////////

  final Map<String, User> _friends = <String, User>{
    '123': const User(
      id: '123',
      email: 'liu.z.jonathan@gmail.com',
      firstName: 'Jonathan',
      lastName: 'Liu',
      picture:
          'https://lh3.googleusercontent.com/a-/AFdZucrz7tSsASL-GwauN8bw3wMswC_Kiuo6Ut8ZGvRtnO4=s96-c',
    ),
    '321': const User(
      id: '321',
      email: 'tonyxin@berkeley.edu',
      firstName: 'Tony',
      lastName: 'Xin',
      picture:
          'https://lh3.googleusercontent.com/a-/AFdZucowznIWn8H4iYmZ1SYTMdRKvBOOgO8sBYTOhfp_3Q=s64-p-k-rw-no',
    ),
    'lol': const User(
      id: 'lol',
      email: 'lesleym@usc.edu',
      firstName: 'Lesley',
      lastName: 'Moon',
      picture:
          'https://lh3.googleusercontent.com/a-/AFdZucrm6jLuiTfc8e-wKD3KZsFLfLhocVKUYoSRaLHfBQ=s64-p-k-rw-no',
    ),
  };
  Map<String, User> get friends => _friends;
  List<User> get friendsList {
    List<User> list = <User>[];
    for (User friend in _friends.values) {
      list.add(friend);
    }
    return list;
  }

  User? getFriendById(String id) {
    return friends[id];
  }

  List<User> getFriendsByQuery(String query) {
    final queryRegex = RegExp('.*$query.*', caseSensitive: false);
    List<User> filteredFriends = <User>[];
    _friends.forEach((_, friend) {
      if (queryRegex.hasMatch(friend.fullName)) {
        filteredFriends.add(friend);
      }
    });
    return filteredFriends;
  }

  final Map<String, CalendarEvents> _friendSchedules = {
    '123': CalendarEvents(
      events: [
        CalendarEvent(
          title: 'Hang out',
          startDate: getDateWithTime(DateTime.now(), 7),
          endDate: getDateWithTime(DateTime.now(), 11),
        ),
        CalendarEvent(
          title: 'hehe xd',
          startDate: getDateWithTime(DateTime.now(), 14),
          endDate: getDateWithTime(DateTime.now(), 16),
        ),
        CalendarEvent(
          title: 'Idk man you decide',
          startDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 1)), 11),
          endDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 1)), 13),
        ),
        CalendarEvent(
          title: 'nice',
          startDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 2)), 17),
          endDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 2)), 20),
        ),
        CalendarEvent(
          title: 'okay then',
          startDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 2)), 10),
          endDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 2)), 12),
        ),
      ],
    ),
    '321': CalendarEvents(
      events: [
        CalendarEvent(
          title: 'Hang out',
          startDate: getDateWithTime(DateTime.now(), 12),
          endDate: getDateWithTime(DateTime.now(), 13),
        ),
        CalendarEvent(
          title: 'hehe xd',
          startDate: getDateWithTime(DateTime.now(), 22),
          endDate: getDateWithTime(DateTime.now(), 23),
        ),
        CalendarEvent(
          title: 'Idk man you decide',
          startDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 1)), 10),
          endDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 1)), 11),
        ),
        CalendarEvent(
          title: 'nice',
          startDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 2)), 22),
          endDate: getDateWithTime(
              DateTime.now().add(const Duration(days: 2)), 22.5),
        ),
        CalendarEvent(
          title: 'okay then',
          startDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 2)), 13),
          endDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 2)), 15),
        ),
      ],
    ),
    'lol': CalendarEvents(
      events: [
        CalendarEvent(
          title: 'Hang out',
          startDate: getDateWithTime(DateTime.now(), 15),
          endDate: getDateWithTime(DateTime.now(), 16),
        ),
        CalendarEvent(
          title: 'hehe xd',
          startDate: getDateWithTime(DateTime.now(), 18),
          endDate: getDateWithTime(DateTime.now(), 21),
        ),
        CalendarEvent(
          title: 'Idk man you decide',
          startDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 1)), 17),
          endDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 1)), 19),
        ),
        CalendarEvent(
          title: 'nice',
          startDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 2)), 9),
          endDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 2)), 15),
        ),
        CalendarEvent(
          title: 'okay then',
          startDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 2)), 19),
          endDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 2)), 20),
        ),
      ],
    )
  };

  CalendarEvents? getFriendScheduleById(String id) {
    return _friendSchedules[id];
  }

  ////////////////////////////////////////////
  // Auth
  ////////////////////////////////////////////

  // Returns whether user is authenticated on server based on session cookie
  Future<bool> isSignedIn() async {
    try {
      await get('/auth/status');
      await refreshAuthUserProfile();
      return true;
    } catch (e) {
      return false;
    }
  }

  // Signs user in on the API and returns true if succeeded
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
      await refreshAuthUserProfile();
      return true;
    } catch (e) {
      // TODO: show dialog that sign in failed
      if (kDebugMode) {
        print('Sign in failed! $e');
      }
      return false;
    }
  }

  // Signs user out on the API
  Future<void> signOut() async {
    await post('/auth/sign-out', {});
  }

  /////////////////////////////////////////////
  // Methods to send http requests to the API
  /////////////////////////////////////////////

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
