import 'dart:collection';
import 'dart:convert';
import 'dart:io';

import 'package:flutter/foundation.dart';
import 'package:flutter_app/models/calendar_event.dart';
import 'package:flutter_app/models/friend_request.dart';
import 'package:flutter_app/models/user.dart';
import 'package:flutter_app/utils.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:http/http.dart';
import 'package:property_change_notifier/property_change_notifier.dart';
import 'package:requests/requests.dart';

enum ApiServiceProperties {
  authUser,
  authUserSchedule,
  friends,
  friendRequests,
}

// Api is used to keep track of all the data retrieved from the API as well as
// to make API requests
class ApiService extends PropertyChangeNotifier {
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

  // Initialize everything
  Future<void> init() async {
    await Future.wait([
      refreshAuthUserProfile(),
      refreshAuthUserSchedule(),
      refreshFriendRequestsList(),
      refreshFriendsList(),
    ]);
  }

  ///////////////////////////////////////////
  // Current user
  ///////////////////////////////////////////

  User? _authUser;
  User? get authUser => _authUser;

  CalendarEvents _authUserSchedule = CalendarEvents(events: []);
  CalendarEvents get authUserSchedule => _authUserSchedule;

  // Gets the user's profile and sets [_authUser] to it
  Future<void> refreshAuthUserProfile() async {
    final userMap = await get('/user/profile');
    _authUser = User.fromJson(userMap);
    notifyListeners(ApiServiceProperties.authUser);
  }

  // Gets the user's schedule and sets [_authUserSchedule] to it
  Future<void> refreshAuthUserSchedule() async {
    final calendarEvents = <CalendarEvent>[];

    final timeMin =
        getLocalDateWithTime(DateTime.now(), 0).toUtc().toIso8601String();
    final timeMax = getLocalDateWithTime(
      DateTime.now().add(const Duration(days: 7)),
      23.99,
    ).toUtc().toIso8601String();

    final jsonEvents =
        await get('/user/calendar?timeMin=$timeMin&timeMax=$timeMax');
    for (final event in jsonEvents) {
      calendarEvents.add(CalendarEvent.fromJson(event));
    }
    _authUserSchedule = CalendarEvents(events: calendarEvents);
    notifyListeners(ApiServiceProperties.authUserSchedule);
  }

  // Updates a user's visibility
  Future<void> updateUserVisibility(int visibility) async {
    await post('/user/visibility', {'visibility': visibility});
  }

  ///////////////////////////////////////////
  // Friends
  ///////////////////////////////////////////

  final Map<String, User> _friends = <String, User>{};
  Map<String, User> get friends => _friends;

  final List<FriendRequest> _friendRequests = <FriendRequest>[];
  List<FriendRequest> get friendRequests => _friendRequests;

  // Gets a user's friends and sets [_friends] to it.
  Future<void> refreshFriendsList() async {
    _friends.clear();
    final result = await get('/friends');
    for (var friend in result) {
      final f = User.fromJson(friend);
      _friends[f.id] = f;
    }
    notifyListeners(ApiServiceProperties.friends);
  }

  // Gets a user's friend requests and sets [_friendRequests] to it.
  Future<void> refreshFriendRequestsList() async {
    _friendRequests.clear();
    final result = await get('/friends/requests');
    for (var request in result) {
      final r = FriendRequest.fromJson(request);
      _friendRequests.add(r);
    }
    notifyListeners(ApiServiceProperties.friendRequests);
  }

  // Refreshes every friend related variable.
  Future<void> refreshFriends() async {
    refreshFriendsList();
    refreshFriendRequestsList();
  }

  // Returns the incoming friend requests.
  List<FriendRequest> getIncomingFriendRequests() {
    List<FriendRequest> list = <FriendRequest>[];
    for (FriendRequest r in _friendRequests) {
      if (r.to == _authUser?.id) list.add(r);
    }
    return list;
  }

  // Returns the outgoing friend requests.
  List<FriendRequest> getOutgoingFriendRequests() {
    List<FriendRequest> list = <FriendRequest>[];
    for (FriendRequest r in _friendRequests) {
      if (r.from == _authUser?.id) list.add(r);
    }
    return list;
  }

  // Generates an array of friends that you've sent a friend request to.
  HashSet<String> getOutgoingFriendRequestsUserIds() {
    final HashSet<String> result = HashSet<String>();
    for (FriendRequest r in _friendRequests) {
      if (r.from == _authUser?.id) {
        result.add(r.to);
      }
    }
    return result;
  }

  List<User> get friendsList {
    List<User> list = <User>[];
    for (User friend in _friends.values) {
      list.add(friend);
    }
    return list;
  }

  User? getFriendById(String id) {
    refreshFriendRequestsList();
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

  // Sends a friend request to a user with [id].
  Future<void> sendFriendRequest(String id) async {
    await post('/friends/requests', {
      'to': id,
    });
    refreshFriendRequestsList();
  }

  // Deletes a friend request with [id].
  Future<void> deleteFriendRequest(String id) async {
    await delete('/friends/requests/$id');
  }

  // Accept a friend request with [id].
  Future<void> acceptFriendRequest(String id) async {
    await post('/friends/requests/$id/accept', {});
    refreshFriends();
  }

  // Reject a friend request with [id].
  Future<void> rejectFriendRequest(String id) async {
    await post('/friends/requests/$id/reject', {});
    refreshFriends();
  }

  final Map<String, CalendarEvents> _friendSchedules = {
    '123': CalendarEvents(
      events: [
        CalendarEvent(
          title: 'Hang out',
          startDate: getLocalDateWithTime(DateTime.now(), 7),
          endDate: getLocalDateWithTime(DateTime.now(), 11),
        ),
        CalendarEvent(
          title: 'hehe xd',
          startDate: getLocalDateWithTime(DateTime.now(), 14),
          endDate: getLocalDateWithTime(DateTime.now(), 16),
        ),
        CalendarEvent(
          title: 'Idk man you decide',
          startDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 1)), 11),
          endDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 1)), 13),
        ),
        CalendarEvent(
          title: 'nice',
          startDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 2)), 17),
          endDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 2)), 20),
        ),
        CalendarEvent(
          title: 'okay then',
          startDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 2)), 10),
          endDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 2)), 12),
        ),
      ],
    ),
    '321': CalendarEvents(
      events: [
        CalendarEvent(
          title: 'Hang out',
          startDate: getLocalDateWithTime(DateTime.now(), 12),
          endDate: getLocalDateWithTime(DateTime.now(), 13),
        ),
        CalendarEvent(
          title: 'hehe xd',
          startDate: getLocalDateWithTime(DateTime.now(), 22),
          endDate: getLocalDateWithTime(DateTime.now(), 23),
        ),
        CalendarEvent(
          title: 'Idk man you decide',
          startDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 1)), 10),
          endDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 1)), 11),
        ),
        CalendarEvent(
          title: 'Cooliooo',
          startDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 1)), 15),
          endDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 1)), 18),
        ),
        CalendarEvent(
          title: 'nice',
          startDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 2)), 22),
          endDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 2)), 22.5),
        ),
        CalendarEvent(
          title: 'okay then',
          startDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 2)), 13),
          endDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 2)), 15),
        ),
      ],
    ),
    'lol': CalendarEvents(
      events: [
        CalendarEvent(
          title: 'Hang out',
          startDate: getLocalDateWithTime(DateTime.now(), 15),
          endDate: getLocalDateWithTime(DateTime.now(), 16),
        ),
        CalendarEvent(
          title: 'hehe xd',
          startDate: getLocalDateWithTime(DateTime.now(), 18),
          endDate: getLocalDateWithTime(DateTime.now(), 21),
        ),
        CalendarEvent(
          title: 'Idk man you decide',
          startDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 1)), 17),
          endDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 1)), 19),
        ),
        CalendarEvent(
          title: 'nice',
          startDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 2)), 9),
          endDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 2)), 15),
        ),
        CalendarEvent(
          title: 'okay then',
          startDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 2)), 19),
          endDate: getLocalDateWithTime(
              DateTime.now().add(const Duration(days: 2)), 20),
        ),
      ],
    )
  };

  CalendarEvents? getFriendScheduleById(String id) {
    return _friendSchedules[id];
  }

  ///////////////////////////////////////////
  // Users
  ///////////////////////////////////////////

  final List<User> _userSearchResults = <User>[];
  List<User> get userSearchResults => _userSearchResults;

  // Gets search results and sets [_searchResults] to it.
  Future<void> refreshUserSearchResults(String query) async {
    _userSearchResults.clear();
    final result = await get('/users?query=$query');
    for (var user in result) {
      final u = User.fromJson(user);
      _userSearchResults.add(u);
    }
    notifyListeners();
  }

  ////////////////////////////////////////////
  // Auth
  ////////////////////////////////////////////

  // Returns whether user is authenticated on server based on session cookie.
  Future<bool> isSignedIn() async {
    try {
      await get('/auth/status');
      try {
        init();
      } catch (e) {
        if (kDebugMode) {
          print(e);
        }
      }
      return true;
    } catch (e) {
      return false;
    }
  }

  // Signs user in on the API and returns true if succeeded
  Future<bool> signIn(String accessToken, int expiresIn, String idToken,
      String refreshToken) async {
    final int timezoneOffset = -1 * DateTime.now().timeZoneOffset.inMinutes;
    try {
      await post(
        '/auth/sign-in-mobile',
        {
          'timezoneOffset': timezoneOffset,
          'accessToken': accessToken,
          'expiresIn': expiresIn,
          'idToken': idToken,
          'refreshToken': refreshToken,
          'tokenOrigin': Platform.isAndroid ? 'android' : 'ios',
        },
      );

      try {
        init();
      } catch (e) {
        if (kDebugMode) {
          print(e);
        }
      }
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

  static Future<dynamic> get(String path) async {
    return await requestMethod(HttpMethod.GET, path);
  }

  static Future<dynamic> post(String path, dynamic body) async {
    return await requestMethod(HttpMethod.POST, path, body: body);
  }

  static Future<dynamic> delete(String path) async {
    return await requestMethod(HttpMethod.DELETE, path);
  }

  static Future<dynamic> patch(String path, dynamic body) async {
    return await requestMethod(HttpMethod.PATCH, path, body: body);
  }

  // Send the given request with the specified method
  static Future<dynamic> requestMethod(HttpMethod method, String path,
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
      case HttpMethod.DELETE:
        r = await Requests.delete(url);
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
    dynamic json = jsonDecode(response);

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
