// **************************************************************************
// AutoRouteGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND

// **************************************************************************
// AutoRouteGenerator
// **************************************************************************
//
// ignore_for_file: type=lint

import 'package:auto_route/auto_route.dart' as _i12;
import 'package:flutter/material.dart' as _i13;

import '../pages/events/event.dart' as _i9;
import '../pages/events/events.dart' as _i8;
import '../pages/events/events_tab.dart' as _i4;
import '../pages/friends/friends.dart' as _i10;
import '../pages/friends/friends_tab.dart' as _i5;
import '../pages/home.dart' as _i2;
import '../pages/my_schej/my_schej.dart' as _i7;
import '../pages/my_schej/my_schej_tab.dart' as _i3;
import '../pages/profile/profile.dart' as _i11;
import '../pages/profile/profile_tab.dart' as _i6;
import '../pages/sign_in.dart' as _i1;
import 'auth_guard.dart' as _i14;

class AppRouter extends _i12.RootStackRouter {
  AppRouter(
      {_i13.GlobalKey<_i13.NavigatorState>? navigatorKey,
      required this.authGuard})
      : super(navigatorKey);

  final _i14.AuthGuard authGuard;

  @override
  final Map<String, _i12.PageFactory> pagesMap = {
    SignInPageRoute.name: (routeData) {
      final args = routeData.argsAs<SignInPageRouteArgs>();
      return _i12.MaterialPageX<dynamic>(
          routeData: routeData,
          child: _i1.SignInPage(key: args.key, onSignIn: args.onSignIn));
    },
    HomePageRoute.name: (routeData) {
      return _i12.MaterialPageX<dynamic>(
          routeData: routeData, child: const _i2.HomePage());
    },
    MySchejTabRoute.name: (routeData) {
      return _i12.MaterialPageX<dynamic>(
          routeData: routeData, child: const _i3.MySchejTab());
    },
    EventsTabRoute.name: (routeData) {
      return _i12.MaterialPageX<dynamic>(
          routeData: routeData, child: const _i4.EventsTab());
    },
    FriendsTabRoute.name: (routeData) {
      return _i12.MaterialPageX<dynamic>(
          routeData: routeData, child: const _i5.FriendsTab());
    },
    ProfileTabRoute.name: (routeData) {
      return _i12.MaterialPageX<dynamic>(
          routeData: routeData, child: const _i6.ProfileTab());
    },
    MySchejPageRoute.name: (routeData) {
      return _i12.MaterialPageX<dynamic>(
          routeData: routeData, child: const _i7.MySchejPage());
    },
    EventsPageRoute.name: (routeData) {
      return _i12.MaterialPageX<dynamic>(
          routeData: routeData, child: const _i8.EventsPage());
    },
    EventPageRoute.name: (routeData) {
      final pathParams = routeData.inheritedPathParams;
      final args = routeData.argsAs<EventPageRouteArgs>(
          orElse: () =>
              EventPageRouteArgs(eventId: pathParams.getString('id')));
      return _i12.MaterialPageX<dynamic>(
          routeData: routeData,
          child: _i9.EventPage(key: args.key, eventId: args.eventId));
    },
    FriendsPageRoute.name: (routeData) {
      return _i12.MaterialPageX<dynamic>(
          routeData: routeData, child: const _i10.FriendsPage());
    },
    ProfilePageRoute.name: (routeData) {
      return _i12.MaterialPageX<dynamic>(
          routeData: routeData, child: const _i11.ProfilePage());
    }
  };

  @override
  List<_i12.RouteConfig> get routes => [
        _i12.RouteConfig(SignInPageRoute.name, path: '/sign-in'),
        _i12.RouteConfig(HomePageRoute.name, path: '/', guards: [
          authGuard
        ], children: [
          _i12.RouteConfig('#redirect',
              path: '',
              parent: HomePageRoute.name,
              redirectTo: 'my-schej',
              fullMatch: true),
          _i12.RouteConfig(MySchejTabRoute.name,
              path: 'my-schej',
              parent: HomePageRoute.name,
              children: [
                _i12.RouteConfig(MySchejPageRoute.name,
                    path: '', parent: MySchejTabRoute.name)
              ]),
          _i12.RouteConfig(EventsTabRoute.name,
              path: 'events',
              parent: HomePageRoute.name,
              children: [
                _i12.RouteConfig(EventsPageRoute.name,
                    path: '', parent: EventsTabRoute.name),
                _i12.RouteConfig(EventPageRoute.name,
                    path: ':id', parent: EventsTabRoute.name)
              ]),
          _i12.RouteConfig(FriendsTabRoute.name,
              path: 'friends',
              parent: HomePageRoute.name,
              children: [
                _i12.RouteConfig(FriendsPageRoute.name,
                    path: '', parent: FriendsTabRoute.name)
              ]),
          _i12.RouteConfig(ProfileTabRoute.name,
              path: 'profile',
              parent: HomePageRoute.name,
              children: [
                _i12.RouteConfig(ProfilePageRoute.name,
                    path: '', parent: ProfileTabRoute.name)
              ])
        ])
      ];
}

/// generated route for
/// [_i1.SignInPage]
class SignInPageRoute extends _i12.PageRouteInfo<SignInPageRouteArgs> {
  SignInPageRoute({_i13.Key? key, required void Function() onSignIn})
      : super(SignInPageRoute.name,
            path: '/sign-in',
            args: SignInPageRouteArgs(key: key, onSignIn: onSignIn));

  static const String name = 'SignInPageRoute';
}

class SignInPageRouteArgs {
  const SignInPageRouteArgs({this.key, required this.onSignIn});

  final _i13.Key? key;

  final void Function() onSignIn;

  @override
  String toString() {
    return 'SignInPageRouteArgs{key: $key, onSignIn: $onSignIn}';
  }
}

/// generated route for
/// [_i2.HomePage]
class HomePageRoute extends _i12.PageRouteInfo<void> {
  const HomePageRoute({List<_i12.PageRouteInfo>? children})
      : super(HomePageRoute.name, path: '/', initialChildren: children);

  static const String name = 'HomePageRoute';
}

/// generated route for
/// [_i3.MySchejTab]
class MySchejTabRoute extends _i12.PageRouteInfo<void> {
  const MySchejTabRoute({List<_i12.PageRouteInfo>? children})
      : super(MySchejTabRoute.name,
            path: 'my-schej', initialChildren: children);

  static const String name = 'MySchejTabRoute';
}

/// generated route for
/// [_i4.EventsTab]
class EventsTabRoute extends _i12.PageRouteInfo<void> {
  const EventsTabRoute({List<_i12.PageRouteInfo>? children})
      : super(EventsTabRoute.name, path: 'events', initialChildren: children);

  static const String name = 'EventsTabRoute';
}

/// generated route for
/// [_i5.FriendsTab]
class FriendsTabRoute extends _i12.PageRouteInfo<void> {
  const FriendsTabRoute({List<_i12.PageRouteInfo>? children})
      : super(FriendsTabRoute.name, path: 'friends', initialChildren: children);

  static const String name = 'FriendsTabRoute';
}

/// generated route for
/// [_i6.ProfileTab]
class ProfileTabRoute extends _i12.PageRouteInfo<void> {
  const ProfileTabRoute({List<_i12.PageRouteInfo>? children})
      : super(ProfileTabRoute.name, path: 'profile', initialChildren: children);

  static const String name = 'ProfileTabRoute';
}

/// generated route for
/// [_i7.MySchejPage]
class MySchejPageRoute extends _i12.PageRouteInfo<void> {
  const MySchejPageRoute() : super(MySchejPageRoute.name, path: '');

  static const String name = 'MySchejPageRoute';
}

/// generated route for
/// [_i8.EventsPage]
class EventsPageRoute extends _i12.PageRouteInfo<void> {
  const EventsPageRoute() : super(EventsPageRoute.name, path: '');

  static const String name = 'EventsPageRoute';
}

/// generated route for
/// [_i9.EventPage]
class EventPageRoute extends _i12.PageRouteInfo<EventPageRouteArgs> {
  EventPageRoute({_i13.Key? key, required String eventId})
      : super(EventPageRoute.name,
            path: ':id',
            args: EventPageRouteArgs(key: key, eventId: eventId),
            rawPathParams: {'id': eventId});

  static const String name = 'EventPageRoute';
}

class EventPageRouteArgs {
  const EventPageRouteArgs({this.key, required this.eventId});

  final _i13.Key? key;

  final String eventId;

  @override
  String toString() {
    return 'EventPageRouteArgs{key: $key, eventId: $eventId}';
  }
}

/// generated route for
/// [_i10.FriendsPage]
class FriendsPageRoute extends _i12.PageRouteInfo<void> {
  const FriendsPageRoute() : super(FriendsPageRoute.name, path: '');

  static const String name = 'FriendsPageRoute';
}

/// generated route for
/// [_i11.ProfilePage]
class ProfilePageRoute extends _i12.PageRouteInfo<void> {
  const ProfilePageRoute() : super(ProfilePageRoute.name, path: '');

  static const String name = 'ProfilePageRoute';
}
