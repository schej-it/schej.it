import 'dart:collection';
import 'dart:convert';
import 'dart:io';

import 'package:flutter/foundation.dart';
import 'package:flutter_app/models/calendar_event.dart';
import 'package:flutter_app/models/friend_request.dart';
import 'package:flutter_app/models/status.dart';
import 'package:flutter_app/models/user.dart';
import 'package:flutter_app/utils.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:http/http.dart';
import 'package:property_change_notifier/property_change_notifier.dart';
import 'package:requests/requests.dart';

enum ApiServiceProperties {
  authUser,
  friends,
  friendRequests,
  friendsStatus,
}

// Api is used to keep track of all the data retrieved from the API as well as
// to make API requests
class ApiService extends PropertyChangeNotifier {
  // [serverAddress] is the correct api address to use
  static String get serverAddress {
    if (kReleaseMode) {
      return 'https://schej.it/api';
    }

    // return 'http://localhost:3002';
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
      refreshFriendRequestsList(),
      refreshFriendsList(),
    ]);
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
    notifyListeners(ApiServiceProperties.authUser);
  }

  // Returns a list of all calendar events for the given user between the given
  // dates
  Future<List<CalendarEvent>> getCalendarEvents({
    required String id,
    required DateTime startDate,
    required DateTime endDate,
  }) async {
    final calendarEvents = <CalendarEvent>[];

    final timeMin =
        getLocalDateWithTime(startDate, 0).toUtc().toIso8601String();
    final timeMax = getLocalDateWithTime(
      endDate,
      23.99,
    ).toUtc().toIso8601String();

    String urlPrefix = '';
    if (id == authUser!.id) {
      urlPrefix = '/user/calendar';
    } else {
      urlPrefix = '/friends/$id/calendar';
    }

    final jsonEvents =
        await get('$urlPrefix?timeMin=$timeMin&timeMax=$timeMax');
    for (final event in jsonEvents) {
      calendarEvents.add(CalendarEvent.fromJson(event));
    }
    return calendarEvents;
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

  final Map<String, Status> _friendsStatus = <String, Status>{};
  Map<String, Status> get friendsStatus => _friendsStatus;

  // Gets a user's friends and sets [_friends] to it.
  Future<void> refreshFriendsList() async {
    _friends.clear();
    final result = await get('/friends');
    for (var friend in result) {
      final f = User.fromJson(friend);
      _friends[f.id] = f;
    }
    refreshFriendsStatus();
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

  // Gets a user's friend requests and sets [_friendRequests] to it.
  Future<void> refreshFriendsStatus() async {
    // NOTE: no _friendsStatus.clear() because no statuses are displayed during retrieval.
    for (var friend in friendsList) {
      final result = await getFriendStatus(friend.id);
      _friendsStatus[friend.id] = result;
    }
    notifyListeners(ApiServiceProperties.friendsStatus);
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
    // refreshFriendRequestsList();
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

  Future<void> deleteFriend(String id) async {
    await delete('/friends/$id');
    refreshFriendsList();
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

  // Get status of friend with [id].
  Future<Status> getFriendStatus(String id) async {
    final result = await get('/friends/$id/status');
    return Status.fromJson(result);
  }

  ///////////////////////////////////////////
  // Users
  ///////////////////////////////////////////

  // Returns a filtered list of users using query.
  Future<List<User>> getUserSearchResults(String query) async {
    final List<User> userSearchResults = <User>[];
    final result = await get('/users?query=$query');
    for (var user in result) {
      final u = User.fromJson(user);
      userSearchResults.add(u);
    }
    return userSearchResults;
  }

  ////////////////////////////////////////////
  // Auth
  ////////////////////////////////////////////

  // Returns whether user is authenticated on server based on session cookie.
  Future<bool> isSignedIn() async {
    try {
      await get('/auth/status');
      try {
        await init();
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
        await init();
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
    dynamic json;
    if (response.isNotEmpty) {
      json = jsonDecode(response);
    } else {
      json = null;
    }

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
