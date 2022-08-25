// **************************************************************************
// AutoRouteGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND

// **************************************************************************
// AutoRouteGenerator
// **************************************************************************
//
// ignore_for_file: type=lint

import 'package:auto_route/auto_route.dart' as _i13;
import 'package:flutter/material.dart' as _i14;

import '../pages/events/event.dart' as _i9;
import '../pages/events/events.dart' as _i8;
import '../pages/events/events_tab.dart' as _i4;
import '../pages/friends/compare_schej.dart' as _i11;
import '../pages/friends/friends.dart' as _i10;
import '../pages/friends/friends_tab.dart' as _i5;
import '../pages/home.dart' as _i2;
import '../pages/my_schej/my_schej.dart' as _i7;
import '../pages/my_schej/my_schej_tab.dart' as _i3;
import '../pages/profile/profile.dart' as _i12;
import '../pages/profile/profile_tab.dart' as _i6;
import '../pages/sign_in.dart' as _i1;
import 'auth_guard.dart' as _i15;

class AppRouter extends _i13.RootStackRouter {
  AppRouter(
      {_i14.GlobalKey<_i14.NavigatorState>? navigatorKey,
      required this.authGuard})
      : super(navigatorKey);

  final _i15.AuthGuard authGuard;

  @override
  final Map<String, _i13.PageFactory> pagesMap = {
    SignInPageRoute.name: (routeData) {
      final args = routeData.argsAs<SignInPageRouteArgs>();
      return _i13.MaterialPageX<dynamic>(
          routeData: routeData,
          child: _i1.SignInPage(key: args.key, onSignIn: args.onSignIn));
    },
    HomePageRoute.name: (routeData) {
      return _i13.MaterialPageX<dynamic>(
          routeData: routeData, child: const _i2.HomePage());
    },
    MySchejTabRoute.name: (routeData) {
      return _i13.MaterialPageX<dynamic>(
          routeData: routeData, child: const _i3.MySchejTab());
    },
    EventsTabRoute.name: (routeData) {
      return _i13.MaterialPageX<dynamic>(
          routeData: routeData, child: const _i4.EventsTab());
    },
    FriendsTabRoute.name: (routeData) {
      return _i13.MaterialPageX<dynamic>(
          routeData: routeData, child: const _i5.FriendsTab());
    },
    ProfileTabRoute.name: (routeData) {
      return _i13.MaterialPageX<dynamic>(
          routeData: routeData, child: const _i6.ProfileTab());
    },
    MySchejPageRoute.name: (routeData) {
      return _i13.MaterialPageX<dynamic>(
          routeData: routeData, child: const _i7.MySchejPage());
    },
    EventsPageRoute.name: (routeData) {
      return _i13.MaterialPageX<dynamic>(
          routeData: routeData, child: const _i8.EventsPage());
    },
    EventPageRoute.name: (routeData) {
      final pathParams = routeData.inheritedPathParams;
      final args = routeData.argsAs<EventPageRouteArgs>(
          orElse: () =>
              EventPageRouteArgs(eventId: pathParams.getString('id')));
      return _i13.MaterialPageX<dynamic>(
          routeData: routeData,
          child: _i9.EventPage(key: args.key, eventId: args.eventId));
    },
    FriendsPageRoute.name: (routeData) {
      return _i13.MaterialPageX<dynamic>(
          routeData: routeData, child: const _i10.FriendsPage());
    },
    CompareSchejPageRoute.name: (routeData) {
      final args = routeData.argsAs<CompareSchejPageRouteArgs>();
      return _i13.MaterialPageX<dynamic>(
          routeData: routeData,
          child: _i11.CompareSchejPage(
              key: args.key,
              friendId: args.friendId,
              initialIncludeSelf: args.initialIncludeSelf));
    },
    ProfilePageRoute.name: (routeData) {
      return _i13.MaterialPageX<dynamic>(
          routeData: routeData, child: const _i12.ProfilePage());
    }
  };

  @override
  List<_i13.RouteConfig> get routes => [
        _i13.RouteConfig(SignInPageRoute.name, path: '/sign-in'),
        _i13.RouteConfig(HomePageRoute.name, path: '/', guards: [
          authGuard
        ], children: [
          _i13.RouteConfig('#redirect',
              path: '',
              parent: HomePageRoute.name,
              redirectTo: 'my-schej',
              fullMatch: true),
          _i13.RouteConfig(MySchejTabRoute.name,
              path: 'my-schej',
              parent: HomePageRoute.name,
              children: [
                _i13.RouteConfig(MySchejPageRoute.name,
                    path: '', parent: MySchejTabRoute.name)
              ]),
          _i13.RouteConfig(EventsTabRoute.name,
              path: 'events',
              parent: HomePageRoute.name,
              children: [
                _i13.RouteConfig(EventsPageRoute.name,
                    path: '', parent: EventsTabRoute.name),
                _i13.RouteConfig(EventPageRoute.name,
                    path: ':id', parent: EventsTabRoute.name)
              ]),
          _i13.RouteConfig(FriendsTabRoute.name,
              path: 'friends',
              parent: HomePageRoute.name,
              children: [
                _i13.RouteConfig(FriendsPageRoute.name,
                    path: '', parent: FriendsTabRoute.name),
                _i13.RouteConfig(CompareSchejPageRoute.name,
                    path: 'schej', parent: FriendsTabRoute.name)
              ]),
          _i13.RouteConfig(ProfileTabRoute.name,
              path: 'profile',
              parent: HomePageRoute.name,
              children: [
                _i13.RouteConfig(ProfilePageRoute.name,
                    path: '', parent: ProfileTabRoute.name)
              ])
        ])
      ];
}

/// generated route for
/// [_i1.SignInPage]
class SignInPageRoute extends _i13.PageRouteInfo<SignInPageRouteArgs> {
  SignInPageRoute({_i14.Key? key, required void Function() onSignIn})
      : super(SignInPageRoute.name,
            path: '/sign-in',
            args: SignInPageRouteArgs(key: key, onSignIn: onSignIn));

  static const String name = 'SignInPageRoute';
}

class SignInPageRouteArgs {
  const SignInPageRouteArgs({this.key, required this.onSignIn});

  final _i14.Key? key;

  final void Function() onSignIn;

  @override
  String toString() {
    return 'SignInPageRouteArgs{key: $key, onSignIn: $onSignIn}';
  }
}

/// generated route for
/// [_i2.HomePage]
class HomePageRoute extends _i13.PageRouteInfo<void> {
  const HomePageRoute({List<_i13.PageRouteInfo>? children})
      : super(HomePageRoute.name, path: '/', initialChildren: children);

  static const String name = 'HomePageRoute';
}

/// generated route for
/// [_i3.MySchejTab]
class MySchejTabRoute extends _i13.PageRouteInfo<void> {
  const MySchejTabRoute({List<_i13.PageRouteInfo>? children})
      : super(MySchejTabRoute.name,
            path: 'my-schej', initialChildren: children);

  static const String name = 'MySchejTabRoute';
}

/// generated route for
/// [_i4.EventsTab]
class EventsTabRoute extends _i13.PageRouteInfo<void> {
  const EventsTabRoute({List<_i13.PageRouteInfo>? children})
      : super(EventsTabRoute.name, path: 'events', initialChildren: children);

  static const String name = 'EventsTabRoute';
}

/// generated route for
/// [_i5.FriendsTab]
class FriendsTabRoute extends _i13.PageRouteInfo<void> {
  const FriendsTabRoute({List<_i13.PageRouteInfo>? children})
      : super(FriendsTabRoute.name, path: 'friends', initialChildren: children);

  static const String name = 'FriendsTabRoute';
}

/// generated route for
/// [_i6.ProfileTab]
class ProfileTabRoute extends _i13.PageRouteInfo<void> {
  const ProfileTabRoute({List<_i13.PageRouteInfo>? children})
      : super(ProfileTabRoute.name, path: 'profile', initialChildren: children);

  static const String name = 'ProfileTabRoute';
}

/// generated route for
/// [_i7.MySchejPage]
class MySchejPageRoute extends _i13.PageRouteInfo<void> {
  const MySchejPageRoute() : super(MySchejPageRoute.name, path: '');

  static const String name = 'MySchejPageRoute';
}

/// generated route for
/// [_i8.EventsPage]
class EventsPageRoute extends _i13.PageRouteInfo<void> {
  const EventsPageRoute() : super(EventsPageRoute.name, path: '');

  static const String name = 'EventsPageRoute';
}

/// generated route for
/// [_i9.EventPage]
class EventPageRoute extends _i13.PageRouteInfo<EventPageRouteArgs> {
  EventPageRoute({_i14.Key? key, required String eventId})
      : super(EventPageRoute.name,
            path: ':id',
            args: EventPageRouteArgs(key: key, eventId: eventId),
            rawPathParams: {'id': eventId});

  static const String name = 'EventPageRoute';
}

class EventPageRouteArgs {
  const EventPageRouteArgs({this.key, required this.eventId});

  final _i14.Key? key;

  final String eventId;

  @override
  String toString() {
    return 'EventPageRouteArgs{key: $key, eventId: $eventId}';
  }
}

/// generated route for
/// [_i10.FriendsPage]
class FriendsPageRoute extends _i13.PageRouteInfo<void> {
  const FriendsPageRoute() : super(FriendsPageRoute.name, path: '');

  static const String name = 'FriendsPageRoute';
}

/// generated route for
/// [_i11.CompareSchejPage]
class CompareSchejPageRoute
    extends _i13.PageRouteInfo<CompareSchejPageRouteArgs> {
  CompareSchejPageRoute(
      {_i14.Key? key,
      required String friendId,
      bool initialIncludeSelf = false})
      : super(CompareSchejPageRoute.name,
            path: 'schej',
            args: CompareSchejPageRouteArgs(
                key: key,
                friendId: friendId,
                initialIncludeSelf: initialIncludeSelf));

  static const String name = 'CompareSchejPageRoute';
}

class CompareSchejPageRouteArgs {
  const CompareSchejPageRouteArgs(
      {this.key, required this.friendId, this.initialIncludeSelf = false});

  final _i14.Key? key;

  final String friendId;

  final bool initialIncludeSelf;

  @override
  String toString() {
    return 'CompareSchejPageRouteArgs{key: $key, friendId: $friendId, initialIncludeSelf: $initialIncludeSelf}';
  }
}

/// generated route for
/// [_i12.ProfilePage]
class ProfilePageRoute extends _i13.PageRouteInfo<void> {
  const ProfilePageRoute() : super(ProfilePageRoute.name, path: '');

  static const String name = 'ProfilePageRoute';
}
